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
)

// Functions to parse text within a block
// Each function returns the number of chars taken care of
// data is the complete block being rendered
// offset is the number of valid chars before the current cursor

func parseInline(out *bytes.Buffer, rndr *render, data []byte) {
	// this is called recursively: enforce a maximum depth
	if rndr.nesting >= rndr.maxNesting {
		return
	}
	rndr.nesting++

	i, end := 0, 0
	for i < len(data) {
		// copy inactive chars into the output
		for end < len(data) && rndr.inline[data[end]] == nil {
			end++
		}

		if rndr.mk.NormalText != nil {
			rndr.mk.NormalText(out, data[i:end], rndr.mk.Opaque)
		} else {
			out.Write(data[i:end])
		}

		if end >= len(data) {
			break
		}
		i = end

		// call the trigger
		parser := rndr.inline[data[end]]
		if consumed := parser(out, rndr, data, i); consumed == 0 {
			// no action from the callback; buffer the byte for later
			end = i + 1
		} else {
			// skip past whatever the callback used
			i += consumed
			end = i
		}
	}

	rndr.nesting--
}

// single and double emphasis parsing
func inlineEmphasis(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
	data = data[offset:]
	c := data[0]
	ret := 0

	if len(data) > 2 && data[1] != c {
		// whitespace cannot follow an opening emphasis;
		// strikethrough only takes two characters '~~'
		if c == '~' || isspace(data[1]) {
			return 0
		}
		if ret = inlineHelperEmph1(out, rndr, data[1:], c); ret == 0 {
			return 0
		}

		return ret + 1
	}

	if len(data) > 3 && data[1] == c && data[2] != c {
		if isspace(data[2]) {
			return 0
		}
		if ret = inlineHelperEmph2(out, rndr, data[2:], c); ret == 0 {
			return 0
		}

		return ret + 2
	}

	if len(data) > 4 && data[1] == c && data[2] == c && data[3] != c {
		if c == '~' || isspace(data[3]) {
			return 0
		}
		if ret = inlineHelperEmph3(out, rndr, data, 3, c); ret == 0 {
			return 0
		}

		return ret + 3
	}

	return 0
}

func inlineCodeSpan(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
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
	for fBegin < end && (data[fBegin] == ' ' || data[fBegin] == '\t') {
		fBegin++
	}

	fEnd := end - nb
	for fEnd > fBegin && (data[fEnd-1] == ' ' || data[fEnd-1] == '\t') {
		fEnd--
	}

	// render the code span
	if rndr.mk.CodeSpan == nil {
		return 0
	}
	if rndr.mk.CodeSpan(out, data[fBegin:fEnd], rndr.mk.Opaque) == 0 {
		end = 0
	}

	return end

}

// newline preceded by two spaces becomes <br>
// newline without two spaces works when EXTENSION_HARD_LINE_BREAK is enabled
func inlineLineBreak(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
	// remove trailing spaces from out
	outBytes := out.Bytes()
	end := len(outBytes)
	eol := end
	for eol > 0 && (outBytes[eol-1] == ' ' || outBytes[eol-1] == '\t') {
		eol--
	}
	out.Truncate(eol)

	// should there be a hard line break here?
	if rndr.flags&EXTENSION_HARD_LINE_BREAK == 0 && end-eol < 2 {
		return 0
	}

	if rndr.mk.LineBreak == nil {
		return 0
	}
	if rndr.mk.LineBreak(out, rndr.mk.Opaque) > 0 {
		return 1
	} else {
		return 0
	}

	return 0
}

// '[': parse a link or an image
func inlineLink(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
	// no links allowed inside other links
	if rndr.insideLink {
		return 0
	}

	isImg := offset > 0 && data[offset-1] == '!'

	data = data[offset:]

	i := 1
	var title, link []byte
	textHasNl := false

	// check whether the correct renderer exists
	if (isImg && rndr.mk.Image == nil) || (!isImg && rndr.mk.Link == nil) {
		return 0
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
		for i < len(data) {
			if data[i] == '\\' {
				i += 2
			} else {
				if data[i] == ')' || data[i] == '\'' || data[i] == '"' {
					break
				}
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

			for i < len(data) {
				if data[i] == '\\' {
					i += 2
				} else {
					if data[i] == ')' {
						break
					}
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
	case i < len(data) && data[i] == '[':
		var id []byte

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
			}
		} else {
			id = data[linkB:linkE]
		}

		// find the reference with matching id (ids are case-insensitive)
		key := string(bytes.ToLower(id))
		lr, ok := rndr.refs[key]
		if !ok {
			return 0
		}

		// keep link and title from reference
		link = lr.link
		title = lr.title
		i++

	// shortcut reference style link
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
			id = data[1:txtE]
		}

		// find the reference with matching id
		key := string(bytes.ToLower(id))
		lr, ok := rndr.refs[key]
		if !ok {
			return 0
		}

		// keep link and title from reference
		link = lr.link
		title = lr.title

		// rewind the whitespace
		i = txtE + 1
	}

	// build content: img alt is escaped, link content is parsed
	var content bytes.Buffer
	if txtE > 1 {
		if isImg {
			content.Write(data[1:txtE])
		} else {
			// links cannot contain other links, so turn off link parsing temporarily
			insideLink := rndr.insideLink
			rndr.insideLink = true
			parseInline(&content, rndr, data[1:txtE])
			rndr.insideLink = insideLink
		}
	}

	var uLink []byte
	if len(link) > 0 {
		var uLinkBuf bytes.Buffer
		unescapeText(&uLinkBuf, link)
		uLink = uLinkBuf.Bytes()
	}

	// call the relevant rendering function
	ret := 0
	if isImg {
		outSize := out.Len()
		outBytes := out.Bytes()
		if outSize > 0 && outBytes[outSize-1] == '!' {
			out.Truncate(outSize - 1)
		}

		ret = rndr.mk.Image(out, uLink, title, content.Bytes(), rndr.mk.Opaque)
	} else {
		ret = rndr.mk.Link(out, uLink, title, content.Bytes(), rndr.mk.Opaque)
	}

	if ret > 0 {
		return i
	}
	return 0
}

// '<' when tags or autolinks are allowed
func inlineLAngle(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
	data = data[offset:]
	altype := LINK_TYPE_NOT_AUTOLINK
	end := tagLength(data, &altype)
	ret := 0

	if end > 2 {
		switch {
		case rndr.mk.AutoLink != nil && altype != LINK_TYPE_NOT_AUTOLINK:
			var uLink bytes.Buffer
			unescapeText(&uLink, data[1:end+1-2])
			ret = rndr.mk.AutoLink(out, uLink.Bytes(), altype, rndr.mk.Opaque)
		case rndr.mk.RawHtmlTag != nil:
			ret = rndr.mk.RawHtmlTag(out, data[:end], rndr.mk.Opaque)
		}
	}

	if ret == 0 {
		return 0
	}
	return end
}

// '\\' backslash escape
var escapeChars = []byte("\\`*_{}[]()#+-.!:|&<>")

func inlineEscape(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
	data = data[offset:]

	if len(data) > 1 {
		if bytes.IndexByte(escapeChars, data[1]) < 0 {
			return 0
		}

		if rndr.mk.NormalText != nil {
			rndr.mk.NormalText(out, data[1:2], rndr.mk.Opaque)
		} else {
			out.WriteByte(data[1])
		}
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
func inlineEntity(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
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

	if rndr.mk.Entity != nil {
		rndr.mk.Entity(out, data[:end], rndr.mk.Opaque)
	} else {
		out.Write(data[:end])
	}

	return end
}

func inlineAutoLink(out *bytes.Buffer, rndr *render, data []byte, offset int) int {
	// quick check to rule out most false hits on ':'
	if rndr.insideLink || len(data) < offset+3 || data[offset+1] != '/' || data[offset+2] != '/' {
		return 0
	}

	// scan backward for a word boundary
	rewind := 0
	for offset-rewind > 0 && rewind <= 7 && !isspace(data[offset-rewind-1]) && !isspace(data[offset-rewind-1]) {
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
	for linkEnd < len(data) && !isspace(data[linkEnd]) {
		linkEnd++
	}

	// Skip punctuation at the end of the link
	if (data[linkEnd-1] == '.' || data[linkEnd-1] == ',' || data[linkEnd-1] == ';') && data[linkEnd-2] != '\\' {
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

	// we were triggered on the ':', so we need to rewind the output a bit
	if out.Len() >= rewind {
		out.Truncate(len(out.Bytes()) - rewind)
	}

	if rndr.mk.AutoLink != nil {
		var uLink bytes.Buffer
		unescapeText(&uLink, data[:linkEnd])

		rndr.mk.AutoLink(out, uLink.Bytes(), LINK_TYPE_NORMAL, rndr.mk.Opaque)
	}

	return linkEnd - rewind
}

var validUris = [][]byte{[]byte("http://"), []byte("https://"), []byte("ftp://"), []byte("mailto://")}

func isSafeLink(link []byte) bool {
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
func tagLength(data []byte, autolink *int) int {
	var i, j int

	// a valid tag can't be shorter than 3 chars
	if len(data) < 3 {
		return 0
	}

	// begins with a '<' optionally followed by '/', followed by letter or number
	if data[0] != '<' {
		return 0
	}
	if data[1] == '/' {
		i = 2
	} else {
		i = 1
	}

	if !isalnum(data[i]) {
		return 0
	}

	// scheme test
	*autolink = LINK_TYPE_NOT_AUTOLINK

	// try to find the beginning of an URI
	for i < len(data) && (isalnum(data[i]) || data[i] == '.' || data[i] == '+' || data[i] == '-') {
		i++
	}

	if i > 1 && data[i] == '@' {
		if j = isMailtoAutoLink(data[i:]); j != 0 {
			*autolink = LINK_TYPE_EMAIL
			return i + j
		}
	}

	if i > 2 && data[i] == ':' {
		*autolink = LINK_TYPE_NORMAL
		i++
	}

	// complete autolink test: no whitespace or ' or "
	switch {
	case i >= len(data):
		*autolink = LINK_TYPE_NOT_AUTOLINK
	case *autolink != 0:
		j = i

		for i < len(data) {
			if data[i] == '\\' {
				i += 2
			} else {
				if data[i] == '>' || data[i] == '\'' || data[i] == '"' || isspace(data[i]) {
					break
				} else {
					i++
				}
			}

		}

		if i >= len(data) {
			return 0
		}
		if i > j && data[i] == '>' {
			return i + 1
		}

		// one of the forbidden chars has been found
		*autolink = LINK_TYPE_NOT_AUTOLINK
	}

	// look for something looking like a tag end
	for i < len(data) && data[i] != '>' {
		i++
	}
	if i >= len(data) {
		return 0
	}
	return i + 1
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
			} else {
				return 0
			}
		default:
			return 0
		}
	}

	return 0
}

// look for the next emph char, skipping other constructs
func inlineHelperFindEmphChar(data []byte, c byte) int {
	i := 1

	for i < len(data) {
		for i < len(data) && data[i] != c && data[i] != '`' && data[i] != '[' {
			i++
		}
		if i >= len(data) {
			return 0
		}
		if data[i] == c {
			return i
		}

		// do not count escaped chars
		if i != 0 && data[i-1] == '\\' {
			i++
			continue
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
		} else {
			if data[i] == '[' {
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
				for i < len(data) && (data[i] == ' ' || data[i] == '\t' || data[i] == '\n') {
					i++
				}
				if i >= len(data) {
					return tmpI
				}
				if data[i] != '[' && data[i] != '(' { // not a link
					if tmpI > 0 {
						return tmpI
					} else {
						continue
					}
				}
				cc := data[i]
				i++
				for i < len(data) && data[i] != cc {
					if tmpI == 0 && data[i] == c {
						tmpI = i
					}
					i++
				}
				if i >= len(data) {
					return tmpI
				}
				i++
			}
		}
	}
	return 0
}

func inlineHelperEmph1(out *bytes.Buffer, rndr *render, data []byte, c byte) int {
	i := 0

	if rndr.mk.Emphasis == nil {
		return 0
	}

	// skip one symbol if coming from emph3
	if len(data) > 1 && data[0] == c && data[1] == c {
		i = 1
	}

	for i < len(data) {
		length := inlineHelperFindEmphChar(data[i:], c)
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

			if rndr.flags&EXTENSION_NO_INTRA_EMPHASIS != 0 {
				if !(i+1 == len(data) || isspace(data[i+1]) || ispunct(data[i+1])) {
					continue
				}
			}

			var work bytes.Buffer
			parseInline(&work, rndr, data[:i])
			r := rndr.mk.Emphasis(out, work.Bytes(), rndr.mk.Opaque)
			if r > 0 {
				return i + 1
			} else {
				return 0
			}
		}
	}

	return 0
}

func inlineHelperEmph2(out *bytes.Buffer, rndr *render, data []byte, c byte) int {
	renderMethod := rndr.mk.DoubleEmphasis
	if c == '~' {
		renderMethod = rndr.mk.StrikeThrough
	}

	if renderMethod == nil {
		return 0
	}

	i := 0

	for i < len(data) {
		length := inlineHelperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length

		if i+1 < len(data) && data[i] == c && data[i+1] == c && i > 0 && !isspace(data[i-1]) {
			var work bytes.Buffer
			parseInline(&work, rndr, data[:i])
			r := renderMethod(out, work.Bytes(), rndr.mk.Opaque)
			if r > 0 {
				return i + 2
			} else {
				return 0
			}
		}
		i++
	}
	return 0
}

func inlineHelperEmph3(out *bytes.Buffer, rndr *render, data []byte, offset int, c byte) int {
	i := 0
	origData := data
	data = data[offset:]

	for i < len(data) {
		length := inlineHelperFindEmphChar(data[i:], c)
		if length == 0 {
			return 0
		}
		i += length

		// skip whitespace preceded symbols
		if data[i] != c || isspace(data[i-1]) {
			continue
		}

		switch {
		case (i+2 < len(data) && data[i+1] == c && data[i+2] == c && rndr.mk.TripleEmphasis != nil):
			// triple symbol found
			var work bytes.Buffer

			parseInline(&work, rndr, data[:i])
			r := rndr.mk.TripleEmphasis(out, work.Bytes(), rndr.mk.Opaque)
			if r > 0 {
				return i + 3
			} else {
				return 0
			}
		case (i+1 < len(data) && data[i+1] == c):
			// double symbol found, hand over to emph1
			length = inlineHelperEmph1(out, rndr, origData[offset-2:], c)
			if length == 0 {
				return 0
			} else {
				return length - 2
			}
		default:
			// single symbol found, hand over to emph2
			length = inlineHelperEmph2(out, rndr, origData[offset-1:], c)
			if length == 0 {
				return 0
			} else {
				return length - 1
			}
		}
	}
	return 0
}
