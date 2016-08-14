//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
// Functions to parse inline elements.
//

package blackfriday

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	urlRe    = `((https?|ftp):\/\/|\/)[-A-Za-z0-9+&@#\/%?=~_|!:,.;\(\)]+`
	anchorRe = regexp.MustCompile(`^(<a\shref="` + urlRe + `"(\stitle="[^"<>]+")?\s?>` + urlRe + `<\/a>)`)

	// TODO: improve this regexp to catch all possible entities:
	htmlEntityRe = regexp.MustCompile(`&[a-z]{2,5};`)
)

// Functions to parse text within a block
// Each function returns the number of chars taken care of
// data is the complete block being rendered
// offset is the number of valid chars before the current cursor

func (p *parser) inline(data []byte) {
	// this is called recursively: enforce a maximum depth
	if p.nesting >= p.maxNesting {
		return
	}
	p.nesting++

	i, end := 0, 0
	for i < len(data) {
		// Stop at EOL
		if data[i] == '\n' && i+1 == len(data) {
			break
		}
		// Copy inactive chars into the output, but first check for one quirk:
		// 'h', 'm' and 'f' all might trigger a check for autolink processing
		// and end this run of inactive characters. However, there's one nasty
		// case where breaking this run would be bad: in smartypants fraction
		// detection, we expect things like "1/2th" to be in a single run. So
		// we check here if an 'h' is followed by 't' (from 'http') and if it's
		// not, we short circuit the 'h' into the run of inactive characters.
		//
		// Also, in a similar fashion maybeLineBreak breaks this run of chars,
		// but smartDash processor relies on seeing context around the dashes.
		// Fix this somehow.
		for end < len(data) {
			if data[end] == ' ' {
				consumed, br := maybeLineBreak(p, data, end)
				if consumed > 0 {
					p.currBlock.AppendChild(text(data[i:end]))
					if br {
						p.currBlock.AppendChild(NewNode(Hardbreak))
					}
					i = end
					i += consumed
					end = i
				} else {
					end++
				}
				continue
			}
			if p.inlineCallback[data[end]] != nil {
				if end+1 < len(data) && data[end] == 'h' && data[end+1] != 't' {
					end++
				} else if strings.ContainsRune("hHmMfF", rune(data[end])) && !isAutoLink(p, data, end) {
					// Autolink callback on something that is not a link.
					end++
				} else {
					break
				}
			} else {
				end++
			}
		}

		p.currBlock.AppendChild(text(data[i:end]))

		if end >= len(data) {
			break
		}
		i = end

		// call the trigger
		handler := p.inlineCallback[data[end]]
		if consumed := handler(p, data, i); consumed == 0 {
			// no action from the callback; buffer the byte for later
			end = i + 1
		} else {
			// skip past whatever the callback used
			i += consumed
			end = i
		}
	}

	p.nesting--
}

// single and double emphasis parsing
func emphasis(p *parser, data []byte, offset int) int {
	data = data[offset:]
	c := data[0]
	ret := 0

	if len(data) > 2 && data[1] != c {
		// whitespace cannot follow an opening emphasis;
		// strikethrough only takes two characters '~~'
		if c == '~' || isspace(data[1]) {
			return 0
		}
		if ret = helperEmphasis(p, data[1:], c); ret == 0 {
			return 0
		}

		return ret + 1
	}

	if len(data) > 3 && data[1] == c && data[2] != c {
		if isspace(data[2]) {
			return 0
		}
		if ret = helperDoubleEmphasis(p, data[2:], c); ret == 0 {
			return 0
		}

		return ret + 2
	}

	if len(data) > 4 && data[1] == c && data[2] == c && data[3] != c {
		if c == '~' || isspace(data[3]) {
			return 0
		}
		if ret = helperTripleEmphasis(p, data, 3, c); ret == 0 {
			return 0
		}

		return ret + 3
	}

	return 0
}

func codeSpan(p *parser, data []byte, offset int) int {
	data = data[offset:]

	nb := 0

	// count the number of backticks in the delimiter
	for nb < len(data) && data[nb] == '`' {
		nb++
	}

	// find the next delimiter
	i, end := 0, 0
	for end = nb; end < len(data) && i < nb; end++ {
		if data[end] == '`' {
			i++
		} else {
			i = 0
		}
	}

	// no matching delimiter?
	if i < nb && end >= len(data) {
		return 0
	}

	// trim outside whitespace
	fBegin := nb
	for fBegin < end && data[fBegin] == ' ' {
		fBegin++
	}

	fEnd := end - nb
	for fEnd > fBegin && data[fEnd-1] == ' ' {
		fEnd--
	}

	// render the code span
	if fBegin != fEnd {
		code := NewNode(Code)
		code.Literal = data[fBegin:fEnd]
		p.currBlock.AppendChild(code)
	}

	return end

}

// newline preceded by two spaces becomes <br>
func maybeLineBreak(p *parser, data []byte, offset int) (int, bool) {
	origOffset := offset
	for offset < len(data) && data[offset] == ' ' {
		offset++
	}
	if offset < len(data) && data[offset] == '\n' {
		if offset-origOffset >= 2 {
			return offset - origOffset + 1, true
		}
		return offset - origOffset, false
	}
	return 0, false
}

// newline without two spaces works when HardLineBreak is enabled
func lineBreak(p *parser, data []byte, offset int) int {
	if p.flags&HardLineBreak != 0 {
		p.currBlock.AppendChild(NewNode(Hardbreak))
		return 1
	}
	return 0
}

type linkType int

const (
	linkNormal linkType = iota
	linkImg
	linkDeferredFootnote
	linkInlineFootnote
)

func isReferenceStyleLink(data []byte, pos int, t linkType) bool {
	if t == linkDeferredFootnote {
		return false
	}
	return pos < len(data)-1 && data[pos] == '[' && data[pos+1] != '^'
}

func maybeImage(p *parser, data []byte, offset int) int {
	if offset < len(data)-1 && data[offset+1] == '[' {
		return link(p, data, offset)
	}
	return 0
}

func maybeInlineFootnote(p *parser, data []byte, offset int) int {
	if offset < len(data)-1 && data[offset+1] == '[' {
		return link(p, data, offset)
	}
	return 0
}

// '[': parse a link or an image or a footnote
func link(p *parser, data []byte, offset int) int {
	// no links allowed inside regular links, footnote, and deferred footnotes
	if p.insideLink && (offset > 0 && data[offset-1] == '[' || len(data)-1 > offset && data[offset+1] == '^') {
		return 0
	}

	var t linkType
	switch {
	// special case: ![^text] == deferred footnote (that follows something with
	// an exclamation point)
	case p.flags&Footnotes != 0 && len(data)-1 > offset && data[offset+1] == '^':
		t = linkDeferredFootnote
	// ![alt] == image
	case offset >= 0 && data[offset] == '!':
		t = linkImg
		offset++
	// ^[text] == inline footnote
	// [^refId] == deferred footnote
	case p.flags&Footnotes != 0:
		if offset >= 0 && data[offset] == '^' {
			t = linkInlineFootnote
			offset++
		} else if len(data)-1 > offset && data[offset+1] == '^' {
			t = linkDeferredFootnote
		}
	// [text] == regular link
	default:
		t = linkNormal
	}

	data = data[offset:]

	var (
		i                       = 1
		noteID                  int
		title, link, altContent []byte
		textHasNl               = false
	)

	if t == linkDeferredFootnote {
		i++
	}

	// look for the matching closing bracket
	for level := 1; level > 0 && i < len(data); i++ {
		switch {
		case data[i] == '\n':
			textHasNl = true

		case data[i-1] == '\\':
			continue

		case data[i] == '[':
			level++

		case data[i] == ']':
			level--
			if level <= 0 {
				i-- // compensate for extra i++ in for loop
			}
		}
	}

	if i >= len(data) {
		return 0
	}

	txtE := i
	i++

	// skip any amount of whitespace or newline
	// (this is much more lax than original markdown syntax)
	for i < len(data) && isspace(data[i]) {
		i++
	}

	// inline style link
	switch {
	case i < len(data) && data[i] == '(':
		// skip initial whitespace
		i++

		for i < len(data) && isspace(data[i]) {
			i++
		}

		linkB := i

		// look for link end: ' " )
	findlinkend:
		for i < len(data) {
			switch {
			case data[i] == '\\':
				i += 2

			case data[i] == ')' || data[i] == '\'' || data[i] == '"':
				break findlinkend

			default:
				i++
			}
		}

		if i >= len(data) {
			return 0
		}
		linkE := i

		// look for title end if present
		titleB, titleE := 0, 0
		if data[i] == '\'' || data[i] == '"' {
			i++
			titleB = i

		findtitleend:
			for i < len(data) {
				switch {
				case data[i] == '\\':
					i += 2

				case data[i] == ')':
					break findtitleend

				default:
					i++
				}
			}

			if i >= len(data) {
				return 0
			}

			// skip whitespace after title
			titleE = i - 1
			for titleE > titleB && isspace(data[titleE]) {
				titleE--
			}

			// check for closing quote presence
			if data[titleE] != '\'' && data[titleE] != '"' {
				titleB, titleE = 0, 0
				linkE = i
			}
		}

		// remove whitespace at the end of the link
		for linkE > linkB && isspace(data[linkE-1]) {
			linkE--
		}

		// remove optional angle brackets around the link
		if data[linkB] == '<' {
			linkB++
		}
		if data[linkE-1] == '>' {
			linkE--
		}

		// build escaped link and title
		if linkE > linkB {
			link = data[linkB:linkE]
		}

		if titleE > titleB {
			title = data[titleB:titleE]
		}

		i++

	// reference style link
	case isReferenceStyleLink(data, i, t):
		var id []byte
		altContentConsidered := false

		// look for the id
		i++
		linkB := i
		for i < len(data) && data[i] != ']' {
			i++
		}
		if i >= len(data) {
			return 0
		}
		linkE := i

		// find the reference
		if linkB == linkE {
			if textHasNl {
				var b bytes.Buffer

				for j := 1; j < txtE; j++ {
					switch {
					case data[j] != '\n':
						b.WriteByte(data[j])
					case data[j-1] != ' ':
						b.WriteByte(' ')
					}
				}

				id = b.Bytes()
			} else {
				id = data[1:txtE]
				altContentConsidered = true
			}
		} else {
			id = data[linkB:linkE]
		}

		// find the reference with matching id
		lr, ok := p.getRef(string(id))
		if !ok {
			return 0
		}

		// keep link and title from reference
		link = lr.link
		title = lr.title
		if altContentConsidered {
			altContent = lr.text
		}
		i++

	// shortcut reference style link or reference or inline footnote
	default:
		var id []byte

		// craft the id
		if textHasNl {
			var b bytes.Buffer

			for j := 1; j < txtE; j++ {
				switch {
				case data[j] != '\n':
					b.WriteByte(data[j])
				case data[j-1] != ' ':
					b.WriteByte(' ')
				}
			}

			id = b.Bytes()
		} else {
			if t == linkDeferredFootnote {
				id = data[2:txtE] // get rid of the ^
			} else {
				id = data[1:txtE]
			}
		}

		if t == linkInlineFootnote {
			// create a new reference
			noteID = len(p.notes) + 1

			var fragment []byte
			if len(id) > 0 {
				if len(id) < 16 {
					fragment = make([]byte, len(id))
				} else {
					fragment = make([]byte, 16)
				}
				copy(fragment, slugify(id))
			} else {
				fragment = append([]byte("footnote-"), []byte(strconv.Itoa(noteID))...)
			}

			ref := &reference{
				noteID:   noteID,
				hasBlock: false,
				link:     fragment,
				title:    id,
			}

			p.notes = append(p.notes, ref)

			link = ref.link
			title = ref.title
		} else {
			// find the reference with matching id
			lr, ok := p.getRef(string(id))
			if !ok {
				return 0
			}

			if t == linkDeferredFootnote {
				lr.noteID = len(p.notes) + 1
				p.notes = append(p.notes, lr)
			}

			// keep link and title from reference
			link = lr.link
			// if inline footnote, title == footnote contents
			title = lr.title
			noteID = lr.noteID
		}

		// rewind the whitespace
		i = txtE + 1
	}

	var uLink []byte
	if t == linkNormal || t == linkImg {
		if len(link) > 0 {
			var uLinkBuf bytes.Buffer
			unescapeText(&uLinkBuf, link)
			uLink = uLinkBuf.Bytes()
		}

		// links need something to click on and somewhere to go
		if len(uLink) == 0 || (t == linkNormal && txtE <= 1) {
			return 0
		}
	}

	// call the relevant rendering function
	switch t {
	case linkNormal:
		linkNode := NewNode(Link)
		linkNode.Destination = normalizeURI(uLink)
		linkNode.Title = title
		p.currBlock.AppendChild(linkNode)
		if len(altContent) > 0 {
			linkNode.AppendChild(text(altContent))
		} else {
			// links cannot contain other links, so turn off link parsing
			// temporarily and recurse
			insideLink := p.insideLink
			p.insideLink = true
			tmpNode := p.currBlock
			p.currBlock = linkNode
			p.inline(data[1:txtE])
			p.currBlock = tmpNode
			p.insideLink = insideLink
		}

	case linkImg:
		linkNode := NewNode(Image)
		linkNode.Destination = uLink
		linkNode.Title = title
		p.currBlock.AppendChild(linkNode)
		linkNode.AppendChild(text(data[1:txtE]))
		i++

	case linkInlineFootnote, linkDeferredFootnote:
		linkNode := NewNode(Link)
		linkNode.Destination = link
		linkNode.Title = title
		linkNode.NoteID = noteID
		p.currBlock.AppendChild(linkNode)
		if t == linkInlineFootnote {
			i++
		}

	default:
		return 0
	}

	return i
}

func (p *parser) inlineHTMLComment(data []byte) int {
	if len(data) < 5 {
		return 0
	}
	if data[0] != '<' || data[1] != '!' || data[2] != '-' || data[3] != '-' {
		return 0
	}
	i := 5
	// scan for an end-of-comment marker, across lines if necessary
	for i < len(data) && !(data[i-2] == '-' && data[i-1] == '-' && data[i] == '>') {
		i++
	}
	// no end-of-comment marker
	if i >= len(data) {
		return 0
	}
	return i + 1
}

func stripMailto(link []byte) []byte {
	if bytes.HasPrefix(link, []byte("mailto://")) {
		return link[9:]
	} else if bytes.HasPrefix(link, []byte("mailto:")) {
		return link[7:]
	} else {
		return link
	}
}

// autolinkType specifies a kind of autolink that gets detected.
type autolinkType int

// These are the possible flag values for the autolink renderer.
const (
	notAutolink autolinkType = iota
	normalAutolink
	emailAutolink
)

// '<' when tags or autolinks are allowed
func leftAngle(p *parser, data []byte, offset int) int {
	data = data[offset:]
	altype, end := tagLength(data)
	if size := p.inlineHTMLComment(data); size > 0 {
		end = size
	}
	if end > 2 {
		if altype != notAutolink {
			var uLink bytes.Buffer
			unescapeText(&uLink, data[1:end+1-2])
			if uLink.Len() > 0 {
				link := uLink.Bytes()
				node := NewNode(Link)
				node.Destination = link
				if altype == emailAutolink {
					node.Destination = append([]byte("mailto:"), link...)
				}
				p.currBlock.AppendChild(node)
				node.AppendChild(text(stripMailto(link)))
			}
		} else {
			htmlTag := NewNode(HTMLSpan)
			htmlTag.Literal = data[:end]
			p.currBlock.AppendChild(htmlTag)
		}
	}

	return end
}

// '\\' backslash escape
var escapeChars = []byte("\\`*_{}[]()#+-.!:|&<>~")

func escape(p *parser, data []byte, offset int) int {
	data = data[offset:]

	if len(data) > 1 {
		if p.flags&BackslashLineBreak != 0 && data[1] == '\n' {
			p.currBlock.AppendChild(NewNode(Hardbreak))
			return 2
		}
		if bytes.IndexByte(escapeChars, data[1]) < 0 {
			return 0
		}

		p.currBlock.AppendChild(text(data[1:2]))
	}

	return 2
}

func unescapeText(ob *bytes.Buffer, src []byte) {
	i := 0
	for i < len(src) {
		org := i
		for i < len(src) && src[i] != '\\' {
			i++
		}

		if i > org {
			ob.Write(src[org:i])
		}

		if i+1 >= len(src) {
			break
		}

		ob.WriteByte(src[i+1])
		i += 2
	}
}

// '&' escaped when it doesn't belong to an entity
// valid entities are assumed to be anything matching &#?[A-Za-z0-9]+;
func entity(p *parser, data []byte, offset int) int {
	data = data[offset:]

	end := 1

	if end < len(data) && data[end] == '#' {
		end++
	}

	for end < len(data) && isalnum(data[end]) {
		end++
	}

	if end < len(data) && data[end] == ';' {
		end++ // real entity
	} else {
		return 0 // lone '&'
	}

	ent := data[:end]
	// undo &amp; escaping or it will be converted to &amp;amp; by another
	// escaper in the renderer
	if bytes.Equal(ent, []byte("&amp;")) {
		ent = []byte{'&'}
	}
	p.currBlock.AppendChild(text(ent))

	return end
}

func linkEndsWithEntity(data []byte, linkEnd int) bool {
	entityRanges := htmlEntityRe.FindAllIndex(data[:linkEnd], -1)
	return entityRanges != nil && entityRanges[len(entityRanges)-1][1] == linkEnd
}

func isAutoLink(p *parser, data []byte, offset int) bool {
	// quick check to rule out most false hits
	if p.insideLink || len(data) < offset+6 { // 6 is the len() of the shortest prefix below
		return false
	}
	prefixes := []string{
		"http://",
		"https://",
		"ftp://",
		"file://",
		"mailto:",
	}
	for _, prefix := range prefixes {
		endOfHead := offset + 8 // 8 is the len() of the longest prefix
		if endOfHead > len(data) {
			endOfHead = len(data)
		}
		head := bytes.ToLower(data[offset:endOfHead])
		if bytes.HasPrefix(head, []byte(prefix)) {
			return true
		}
	}
	return false
}

func autoLink(p *parser, data []byte, offset int) int {
	// Now a more expensive check to see if we're not inside an anchor element
	anchorStart := offset
	offsetFromAnchor := 0
	for anchorStart > 0 && data[anchorStart] != '<' {
		anchorStart--
		offsetFromAnchor++
	}

	anchorStr := anchorRe.Find(data[anchorStart:])
	if anchorStr != nil {
		anchorClose := NewNode(HTMLSpan)
		anchorClose.Literal = anchorStr[offsetFromAnchor:]
		p.currBlock.AppendChild(anchorClose)
		return len(anchorStr) - offsetFromAnchor
	}

	// scan backward for a word boundary
	rewind := 0
	for offset-rewind > 0 && rewind <= 7 && isletter(data[offset-rewind-1]) {
		rewind++
	}
	if rewind > 6 { // longest supported protocol is "mailto" which has 6 letters
		return 0
	}

	origData := data
	data = data[offset-rewind:]

	if !isSafeLink(data) {
		return 0
	}

	linkEnd := 0
	for linkEnd < len(data) && !isEndOfLink(data[linkEnd]) {
		linkEnd++
	}

	// Skip punctuation at the end of the link
	if (data[linkEnd-1] == '.' || data[linkEnd-1] == ',') && data[linkEnd-2] != '\\' {
		linkEnd--
	}

	// But don't skip semicolon if it's a part of escaped entity:
	if data[linkEnd-1] == ';' && data[linkEnd-2] != '\\' && !linkEndsWithEntity(data, linkEnd) {
		linkEnd--
	}

	// See if the link finishes with a punctuation sign that can be closed.
	var copen byte
	switch data[linkEnd-1] {
	case '"':
		copen = '"'
	case '\'':
		copen = '\''
	case ')':
		copen = '('
	case ']':
		copen = '['
	case '}':
		copen = '{'
	default:
		copen = 0
	}

	if copen != 0 {
		bufEnd := offset - rewind + linkEnd - 2

		openDelim := 1

		/* Try to close the final punctuation sign in this same line;
		 * if we managed to close it outside of the URL, that means that it's
		 * not part of the URL. If it closes inside the URL, that means it
		 * is part of the URL.
		 *
		 * Examples:
		 *
		 *      foo http://www.pokemon.com/Pikachu_(Electric) bar
		 *              => http://www.pokemon.com/Pikachu_(Electric)
		 *
		 *      foo (http://www.pokemon.com/Pikachu_(Electric)) bar
		 *              => http://www.pokemon.com/Pikachu_(Electric)
		 *
		 *      foo http://www.pokemon.com/Pikachu_(Electric)) bar
		 *              => http://www.pokemon.com/Pikachu_(Electric))
		 *
		 *      (foo http://www.pokemon.com/Pikachu_(Electric)) bar
		 *              => foo http://www.pokemon.com/Pikachu_(Electric)
		 */

		for bufEnd >= 0 && origData[bufEnd] != '\n' && openDelim != 0 {
			if origData[bufEnd] == data[linkEnd-1] {
				openDelim++
			}

			if origData[bufEnd] == copen {
				openDelim--
			}

			bufEnd--
		}

		if openDelim == 0 {
			linkEnd--
		}
	}

	var uLink bytes.Buffer
	unescapeText(&uLink, data[:linkEnd])

	if uLink.Len() > 0 {
		node := NewNode(Link)
		node.Destination = uLink.Bytes()
		p.currBlock.AppendChild(node)
		node.AppendChild(text(uLink.Bytes()))
	}

	return linkEnd
}

func isEndOfLink(char byte) bool {
	return isspace(char) || char == '<'
}

var validUris = [][]byte{[]byte("http://"), []byte("https://"), []byte("ftp://"), []byte("mailto://")}
var validPaths = [][]byte{[]byte("/"), []byte("./"), []byte("../")}

func isSafeLink(link []byte) bool {
	for _, path := range validPaths {
		if len(link) >= len(path) && bytes.Equal(link[:len(path)], path) {
			if len(link) == len(path) {
				return true
			} else if isalnum(link[len(path)]) {
				return true
			}
		}
	}

	for _, prefix := range validUris {
		// TODO: handle unicode here
		// case-insensitive prefix test
		if len(link) > len(prefix) && bytes.Equal(bytes.ToLower(link[:len(prefix)]), prefix) && isalnum(link[len(prefix)]) {
			return true
		}
	}

	return false
}

// return the length of the given tag, or 0 is it's not valid
func tagLength(data []byte) (autolink autolinkType, end int) {
	var i, j int

	// a valid tag can't be shorter than 3 chars
	if len(data) < 3 {
		return notAutolink, 0
	}

	// begins with a '<' optionally followed by '/', followed by letter or number
	if data[0] != '<' {
		return notAutolink, 0
	}
	if data[1] == '/' {
		i = 2
	} else {
		i = 1
	}

	if !isalnum(data[i]) {
		return notAutolink, 0
	}

	// scheme test
	autolink = notAutolink

	// try to find the beginning of an URI
	for i < len(data) && (isalnum(data[i]) || data[i] == '.' || data[i] == '+' || data[i] == '-') {
		i++
	}

	if i > 1 && i < len(data) && data[i] == '@' {
		if j = isMailtoAutoLink(data[i:]); j != 0 {
			return emailAutolink, i + j
		}
	}

	if i > 2 && i < len(data) && data[i] == ':' {
		autolink = normalAutolink
		i++
	}

	// complete autolink test: no whitespace or ' or "
	switch {
	case i >= len(data):
		autolink = notAutolink
	case autolink != notAutolink:
		j = i

		for i < len(data) {
			if data[i] == '\\' {
				i += 2
			} else if data[i] == '>' || data[i] == '\'' || data[i] == '"' || isspace(data[i]) {
				break
			} else {
				i++
			}

		}

		if i >= len(data) {
			return autolink, 0
		}
		if i > j && data[i] == '>' {
			return autolink, i + 1
		}

		// one of the forbidden chars has been found
		autolink = notAutolink
	}
	i += bytes.IndexByte(data[i:], '>')
	if i < 0 {
		return autolink, 0
	}
	return autolink, i + 1
}

// look for the address part of a mail autolink and '>'
// this is less strict than the original markdown e-mail address matching
func isMailtoAutoLink(data []byte) int {
	nb := 0

	// address is assumed to be: [-@._a-zA-Z0-9]+ with exactly one '@'
	for i := 0; i < len(data); i++ {
		if isalnum(data[i]) {
			continue
		}

		switch data[i] {
		case '@':
			nb++

		case '-', '.', '_':
			break

		case '>':
			if nb == 1 {
				return i + 1
			}
			return 0
		default:
			return 0
		}
	}

	return 0
}

// look for the next emph char, skipping other constructs
func helperFindEmphChar(data []byte, c byte) int {
	i := 0

	for i < len(data) {
		for i < len(data) && data[i] != c && data[i] != '`' && data[i] != '[' {
			i++
		}
		if i >= len(data) {
			return 0
		}
		// do not count escaped chars
		if i != 0 && data[i-1] == '\\' {
			i++
			continue
		}
		if data[i] == c {
			return i
		}

		if data[i] == '`' {
			// skip a code span
			tmpI := 0
			i++
			for i < len(data) && data[i] != '`' {
				if tmpI == 0 && data[i] == c {
					tmpI = i
				}
				i++
			}
			if i >= len(data) {
				return tmpI
			}
			i++
		} else if data[i] == '[' {
			// skip a link
			tmpI := 0
			i++
			for i < len(data) && data[i] != ']' {
				if tmpI == 0 && data[i] == c {
					tmpI = i
				}
				i++
			}
			i++
			for i < len(data) && (data[i] == ' ' || data[i] == '\n') {
				i++
			}
			if i >= len(data) {
				return tmpI
			}
			if data[i] != '[' && data[i] != '(' { // not a link
				if tmpI > 0 {
					return tmpI
				}
				continue
			}
			cc := data[i]
			i++
			for i < len(data) && data[i] != cc {
				if tmpI == 0 && data[i] == c {
					return i
				}
				i++
			}
			if i >= len(data) {
				return tmpI
			}
			i++
		}
	}
	return 0
}

func helperEmphasis(p *parser, data []byte, c byte) int {
	i := 0

	// skip one symbol if coming from emph3
	if len(data) > 1 && data[0] == c && data[1] == c {
		i = 1
	}

	for i < len(data) {
		length := helperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length
		if i >= len(data) {
			return 0
		}

		if i+1 < len(data) && data[i+1] == c {
			i++
			continue
		}

		if data[i] == c && !isspace(data[i-1]) {

			if p.flags&NoIntraEmphasis != 0 {
				if !(i+1 == len(data) || isspace(data[i+1]) || ispunct(data[i+1])) {
					continue
				}
			}

			emph := NewNode(Emph)
			p.currBlock.AppendChild(emph)
			tmp := p.currBlock
			p.currBlock = emph
			p.inline(data[:i])
			p.currBlock = tmp
			return i + 1
		}
	}

	return 0
}

func helperDoubleEmphasis(p *parser, data []byte, c byte) int {
	i := 0

	for i < len(data) {
		length := helperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length

		if i+1 < len(data) && data[i] == c && data[i+1] == c && i > 0 && !isspace(data[i-1]) {
			nodeType := Strong
			if c == '~' {
				nodeType = Del
			}
			node := NewNode(nodeType)
			p.currBlock.AppendChild(node)
			tmp := p.currBlock
			p.currBlock = node
			p.inline(data[:i])
			p.currBlock = tmp
			return i + 2
		}
		i++
	}
	return 0
}

func helperTripleEmphasis(p *parser, data []byte, offset int, c byte) int {
	i := 0
	origData := data
	data = data[offset:]

	for i < len(data) {
		length := helperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length

		// skip whitespace preceded symbols
		if data[i] != c || isspace(data[i-1]) {
			continue
		}

		switch {
		case i+2 < len(data) && data[i+1] == c && data[i+2] == c:
			// triple symbol found
			strong := NewNode(Strong)
			em := NewNode(Emph)
			strong.AppendChild(em)
			p.currBlock.AppendChild(strong)
			tmp := p.currBlock
			p.currBlock = em
			p.inline(data[:i])
			p.currBlock = tmp
			return i + 3
		case (i+1 < len(data) && data[i+1] == c):
			// double symbol found, hand over to emph1
			length = helperEmphasis(p, origData[offset-2:], c)
			if length == 0 {
				return 0
			}
			return length - 2
		default:
			// single symbol found, hand over to emph2
			length = helperDoubleEmphasis(p, origData[offset-1:], c)
			if length == 0 {
				return 0
			}
			return length - 1
		}
	}
	return 0
}

func text(s []byte) *Node {
	node := NewNode(Text)
	node.Literal = s
	return node
}

func normalizeURI(s []byte) []byte {
	return s // TODO: implement
}
