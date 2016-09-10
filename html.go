//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
//
// HTML rendering backend
//
//

package blackfriday

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
)

// HTMLFlags control optional behavior of HTML renderer.
type HTMLFlags int

// HTML renderer configuration options.
const (
	HTMLFlagsNone           HTMLFlags = 0
	SkipHTML                HTMLFlags = 1 << iota // Skip preformatted HTML blocks
	SkipStyle                                     // Skip embedded <style> elements
	SkipImages                                    // Skip embedded images
	SkipLinks                                     // Skip all links
	Safelink                                      // Only link to trusted protocols
	NofollowLinks                                 // Only link with rel="nofollow"
	NoreferrerLinks                               // Only link with rel="noreferrer"
	HrefTargetBlank                               // Add a blank target
	CompletePage                                  // Generate a complete HTML page
	UseXHTML                                      // Generate XHTML output instead of HTML
	FootnoteReturnLinks                           // Generate a link at the end of a footnote to return to the source
	Smartypants                                   // Enable smart punctuation substitutions
	SmartypantsFractions                          // Enable smart fractions (with Smartypants)
	SmartypantsDashes                             // Enable smart dashes (with Smartypants)
	SmartypantsLatexDashes                        // Enable LaTeX-style dashes (with Smartypants)
	SmartypantsAngledQuotes                       // Enable angled double quotes (with Smartypants) for double quotes rendering

	TagName               = "[A-Za-z][A-Za-z0-9-]*"
	AttributeName         = "[a-zA-Z_:][a-zA-Z0-9:._-]*"
	UnquotedValue         = "[^\"'=<>`\\x00-\\x20]+"
	SingleQuotedValue     = "'[^']*'"
	DoubleQuotedValue     = "\"[^\"]*\""
	AttributeValue        = "(?:" + UnquotedValue + "|" + SingleQuotedValue + "|" + DoubleQuotedValue + ")"
	AttributeValueSpec    = "(?:" + "\\s*=" + "\\s*" + AttributeValue + ")"
	Attribute             = "(?:" + "\\s+" + AttributeName + AttributeValueSpec + "?)"
	OpenTag               = "<" + TagName + Attribute + "*" + "\\s*/?>"
	CloseTag              = "</" + TagName + "\\s*[>]"
	HTMLComment           = "<!---->|<!--(?:-?[^>-])(?:-?[^-])*-->"
	ProcessingInstruction = "[<][?].*?[?][>]"
	Declaration           = "<![A-Z]+" + "\\s+[^>]*>"
	CDATA                 = "<!\\[CDATA\\[[\\s\\S]*?\\]\\]>"
	HTMLTag               = "(?:" + OpenTag + "|" + CloseTag + "|" + HTMLComment + "|" +
		ProcessingInstruction + "|" + Declaration + "|" + CDATA + ")"
)

var (
	htmlTagRe = regexp.MustCompile("(?i)^" + HTMLTag)
)

// HTMLRendererParameters is a collection of supplementary parameters tweaking
// the behavior of various parts of HTML renderer.
type HTMLRendererParameters struct {
	// Prepend this text to each relative URL.
	AbsolutePrefix string
	// Add this text to each footnote anchor, to ensure uniqueness.
	FootnoteAnchorPrefix string
	// Show this text inside the <a> tag for a footnote return link, if the
	// HTML_FOOTNOTE_RETURN_LINKS flag is enabled. If blank, the string
	// <sup>[return]</sup> is used.
	FootnoteReturnLinkContents string
	// If set, add this text to the front of each Header ID, to ensure
	// uniqueness.
	HeaderIDPrefix string
	// If set, add this text to the back of each Header ID, to ensure uniqueness.
	HeaderIDSuffix string

	Title string // Document title (used if CompletePage is set)
	CSS   string // Optional CSS file URL (used if CompletePage is set)
	Icon  string // Optional icon file URL (used if CompletePage is set)

	Flags      HTMLFlags  // Flags allow customizing this renderer's behavior
	Extensions Extensions // Extensions give Smartypants and HTML renderer access to Blackfriday's global extensions
}

// HTMLRenderer is a type that implements the Renderer interface for HTML output.
//
// Do not create this directly, instead use the NewHTMLRenderer function.
type HTMLRenderer struct {
	HTMLRendererParameters

	closeTag string // how to end singleton tags: either " />" or ">"

	// Track header IDs to prevent ID collision in a single generation.
	headerIDs map[string]int

	lastOutputLen int
	disableTags   int

	sr *SPRenderer
}

const (
	xhtmlClose = " />"
	htmlClose  = ">"
)

// NewHTMLRenderer creates and configures an HTMLRenderer object, which
// satisfies the Renderer interface.
func NewHTMLRenderer(params HTMLRendererParameters) *HTMLRenderer {
	// configure the rendering engine
	closeTag := htmlClose
	if params.Flags&UseXHTML != 0 {
		closeTag = xhtmlClose
	}

	if params.FootnoteReturnLinkContents == "" {
		params.FootnoteReturnLinkContents = `<sup>[return]</sup>`
	}

	return &HTMLRenderer{
		HTMLRendererParameters: params,

		closeTag:  closeTag,
		headerIDs: make(map[string]int),

		sr: NewSmartypantsRenderer(params.Flags),
	}
}

func isHTMLTag(tag []byte, tagname string) bool {
	found, _ := findHTMLTagPos(tag, tagname)
	return found
}

// Look for a character, but ignore it when it's in any kind of quotes, it
// might be JavaScript
func skipUntilCharIgnoreQuotes(html []byte, start int, char byte) int {
	inSingleQuote := false
	inDoubleQuote := false
	inGraveQuote := false
	i := start
	for i < len(html) {
		switch {
		case html[i] == char && !inSingleQuote && !inDoubleQuote && !inGraveQuote:
			return i
		case html[i] == '\'':
			inSingleQuote = !inSingleQuote
		case html[i] == '"':
			inDoubleQuote = !inDoubleQuote
		case html[i] == '`':
			inGraveQuote = !inGraveQuote
		}
		i++
	}
	return start
}

func findHTMLTagPos(tag []byte, tagname string) (bool, int) {
	i := 0
	if i < len(tag) && tag[0] != '<' {
		return false, -1
	}
	i++
	i = skipSpace(tag, i)

	if i < len(tag) && tag[i] == '/' {
		i++
	}

	i = skipSpace(tag, i)
	j := 0
	for ; i < len(tag); i, j = i+1, j+1 {
		if j >= len(tagname) {
			break
		}

		if strings.ToLower(string(tag[i]))[0] != tagname[j] {
			return false, -1
		}
	}

	if i == len(tag) {
		return false, -1
	}

	rightAngle := skipUntilCharIgnoreQuotes(tag, i, '>')
	if rightAngle >= i {
		return true, rightAngle
	}

	return false, -1
}

func skipSpace(tag []byte, i int) int {
	for i < len(tag) && isspace(tag[i]) {
		i++
	}
	return i
}

func isRelativeLink(link []byte) (yes bool) {
	// a tag begin with '#'
	if link[0] == '#' {
		return true
	}

	// link begin with '/' but not '//', the second maybe a protocol relative link
	if len(link) >= 2 && link[0] == '/' && link[1] != '/' {
		return true
	}

	// only the root '/'
	if len(link) == 1 && link[0] == '/' {
		return true
	}

	// current directory : begin with "./"
	if bytes.HasPrefix(link, []byte("./")) {
		return true
	}

	// parent directory : begin with "../"
	if bytes.HasPrefix(link, []byte("../")) {
		return true
	}

	return false
}

func (r *HTMLRenderer) ensureUniqueHeaderID(id string) string {
	for count, found := r.headerIDs[id]; found; count, found = r.headerIDs[id] {
		tmp := fmt.Sprintf("%s-%d", id, count+1)

		if _, tmpFound := r.headerIDs[tmp]; !tmpFound {
			r.headerIDs[id] = count + 1
			id = tmp
		} else {
			id = id + "-1"
		}
	}

	if _, found := r.headerIDs[id]; !found {
		r.headerIDs[id] = 0
	}

	return id
}

func (r *HTMLRenderer) addAbsPrefix(link []byte) []byte {
	if r.AbsolutePrefix != "" && isRelativeLink(link) && link[0] != '.' {
		newDest := r.AbsolutePrefix
		if link[0] != '/' {
			newDest += "/"
		}
		newDest += string(link)
		return []byte(newDest)
	}
	return link
}

func appendLinkAttrs(attrs []string, flags HTMLFlags, link []byte) []string {
	if isRelativeLink(link) {
		return attrs
	}
	val := []string{}
	if flags&NofollowLinks != 0 {
		val = append(val, "nofollow")
	}
	if flags&NoreferrerLinks != 0 {
		val = append(val, "noreferrer")
	}
	if flags&HrefTargetBlank != 0 {
		attrs = append(attrs, "target=\"_blank\"")
	}
	if len(val) == 0 {
		return attrs
	}
	attr := fmt.Sprintf("rel=%q", strings.Join(val, " "))
	return append(attrs, attr)
}

func isMailto(link []byte) bool {
	return bytes.HasPrefix(link, []byte("mailto:"))
}

func needSkipLink(flags HTMLFlags, dest []byte) bool {
	if flags&SkipLinks != 0 {
		return true
	}
	return flags&Safelink != 0 && !isSafeLink(dest) && !isMailto(dest)
}

func isSmartypantable(node *Node) bool {
	pt := node.Parent.Type
	return pt != Link && pt != CodeBlock && pt != Code
}

func appendLanguageAttr(attrs []string, info []byte) []string {
	infoWords := bytes.Split(info, []byte("\t "))
	if len(infoWords) > 0 && len(infoWords[0]) > 0 {
		attrs = append(attrs, fmt.Sprintf("class=\"language-%s\"", infoWords[0]))
	}
	return attrs
}

func tag(name string, attrs []string, selfClosing bool) []byte {
	result := "<" + name
	if attrs != nil && len(attrs) > 0 {
		result += " " + strings.Join(attrs, " ")
	}
	if selfClosing {
		result += " /"
	}
	return []byte(result + ">")
}

func footnoteRef(prefix string, node *Node) []byte {
	urlFrag := prefix + string(slugify(node.Destination))
	anchor := fmt.Sprintf(`<a rel="footnote" href="#fn:%s">%d</a>`, urlFrag, node.NoteID)
	return []byte(fmt.Sprintf(`<sup class="footnote-ref" id="fnref:%s">%s</sup>`, urlFrag, anchor))
}

func footnoteItem(prefix string, slug []byte) []byte {
	return []byte(fmt.Sprintf(`<li id="fn:%s%s">`, prefix, slug))
}

func footnoteReturnLink(prefix, returnLink string, slug []byte) []byte {
	const format = ` <a class="footnote-return" href="#fnref:%s%s">%s</a>`
	return []byte(fmt.Sprintf(format, prefix, slug, returnLink))
}

func itemOpenCR(node *Node) bool {
	if node.Prev == nil {
		return false
	}
	ld := node.Parent.ListData
	return !ld.Tight && ld.ListFlags&ListTypeDefinition == 0
}

func skipParagraphTags(node *Node) bool {
	grandparent := node.Parent.Parent
	if grandparent == nil || grandparent.Type != List {
		return false
	}
	tightOrTerm := grandparent.Tight || node.Parent.ListFlags&ListTypeTerm != 0
	return grandparent.Type == List && tightOrTerm
}

func cellAlignment(align CellAlignFlags) string {
	switch align {
	case TableAlignmentLeft:
		return "left"
	case TableAlignmentRight:
		return "right"
	case TableAlignmentCenter:
		return "center"
	default:
		return ""
	}
}

func esc(text []byte) []byte {
	unesc := []byte(html.UnescapeString(string(text)))
	return escCode(unesc)
}

func escCode(text []byte) []byte {
	e1 := []byte(html.EscapeString(string(text)))
	e2 := bytes.Replace(e1, []byte("&#34;"), []byte("&quot;"), -1)
	return bytes.Replace(e2, []byte("&#39;"), []byte{'\''}, -1)
}

func (r *HTMLRenderer) out(w io.Writer, text []byte) {
	if r.disableTags > 0 {
		w.Write(htmlTagRe.ReplaceAll(text, []byte{}))
	} else {
		w.Write(text)
	}
	r.lastOutputLen = len(text)
}

func (r *HTMLRenderer) cr(w io.Writer) {
	if r.lastOutputLen > 0 {
		r.out(w, []byte{'\n'})
	}
}

// RenderNode is a default renderer of a single node of a syntax tree. For
// block nodes it will be called twice: first time with entering=true, second
// time with entering=false, so that it could know when it's working on an open
// tag and when on close. It writes the result to w.
//
// The return value is a way to tell the calling walker to adjust its walk
// pattern: e.g. it can terminate the traversal by returning Terminate. Or it
// can ask the walker to skip a subtree of this node by returning SkipChildren.
// The typical behavior is to return GoToNext, which asks for the usual
// traversal to the next node.
func (r *HTMLRenderer) RenderNode(w io.Writer, node *Node, entering bool) WalkStatus {
	attrs := []string{}
	switch node.Type {
	case Text:
		node.Literal = esc(node.Literal)
		if r.Flags&Smartypants != 0 {
			node.Literal = r.sr.Process(node.Literal)
		}
		r.out(w, node.Literal)
	case Softbreak:
		r.out(w, []byte{'\n'})
		// TODO: make it configurable via out(renderer.softbreak)
	case Hardbreak:
		r.out(w, tag("br", nil, true))
		r.cr(w)
	case Emph:
		if entering {
			r.out(w, tag("em", nil, false))
		} else {
			r.out(w, tag("/em", nil, false))
		}
	case Strong:
		if entering {
			r.out(w, tag("strong", nil, false))
		} else {
			r.out(w, tag("/strong", nil, false))
		}
	case Del:
		if entering {
			r.out(w, tag("del", nil, false))
		} else {
			r.out(w, tag("/del", nil, false))
		}
	case HTMLSpan:
		if r.Flags&SkipHTML != 0 {
			break
		}
		if r.Flags&SkipStyle != 0 && isHTMLTag(node.Literal, "style") {
			break
		}
		//if options.safe {
		//	out(w, "<!-- raw HTML omitted -->")
		//} else {
		r.out(w, node.Literal)
		//}
	case Link:
		// mark it but don't link it if it is not a safe link: no smartypants
		dest := node.LinkData.Destination
		if needSkipLink(r.Flags, dest) {
			if entering {
				r.out(w, tag("tt", nil, false))
			} else {
				r.out(w, tag("/tt", nil, false))
			}
		} else {
			if entering {
				dest = r.addAbsPrefix(dest)
				//if (!(options.safe && potentiallyUnsafe(node.destination))) {
				attrs = append(attrs, fmt.Sprintf("href=%q", esc(dest)))
				//}
				if node.NoteID != 0 {
					r.out(w, footnoteRef(r.FootnoteAnchorPrefix, node))
					break
				}
				attrs = appendLinkAttrs(attrs, r.Flags, dest)
				if len(node.LinkData.Title) > 0 {
					attrs = append(attrs, fmt.Sprintf("title=%q", esc(node.LinkData.Title)))
				}
				r.out(w, tag("a", attrs, false))
			} else {
				if node.NoteID != 0 {
					break
				}
				r.out(w, tag("/a", nil, false))
			}
		}
	case Image:
		if r.Flags&SkipImages != 0 {
			return SkipChildren
		}
		if entering {
			dest := node.LinkData.Destination
			dest = r.addAbsPrefix(dest)
			if r.disableTags == 0 {
				//if options.safe && potentiallyUnsafe(dest) {
				//out(w, `<img src="" alt="`)
				//} else {
				r.out(w, []byte(fmt.Sprintf(`<img src="%s" alt="`, esc(dest))))
				//}
			}
			r.disableTags++
		} else {
			r.disableTags--
			if r.disableTags == 0 {
				if node.LinkData.Title != nil {
					r.out(w, []byte(`" title="`))
					r.out(w, esc(node.LinkData.Title))
				}
				r.out(w, []byte(`" />`))
			}
		}
	case Code:
		r.out(w, tag("code", nil, false))
		r.out(w, escCode(node.Literal))
		r.out(w, tag("/code", nil, false))
	case Document:
		break
	case Paragraph:
		if skipParagraphTags(node) {
			break
		}
		if entering {
			// TODO: untangle this clusterfuck about when the newlines need
			// to be added and when not.
			if node.Prev != nil {
				switch node.Prev.Type {
				case HTMLBlock, List, Paragraph, Header, CodeBlock, BlockQuote, HorizontalRule:
					r.cr(w)
				}
			}
			if node.Parent.Type == BlockQuote && node.Prev == nil {
				r.cr(w)
			}
			r.out(w, tag("p", attrs, false))
		} else {
			r.out(w, tag("/p", attrs, false))
			if !(node.Parent.Type == Item && node.Next == nil) {
				r.cr(w)
			}
		}
	case BlockQuote:
		if entering {
			r.cr(w)
			r.out(w, tag("blockquote", attrs, false))
		} else {
			r.out(w, tag("/blockquote", nil, false))
			r.cr(w)
		}
	case HTMLBlock:
		if r.Flags&SkipHTML != 0 {
			break
		}
		r.cr(w)
		r.out(w, node.Literal)
		r.cr(w)
	case Header:
		tagname := fmt.Sprintf("h%d", node.Level)
		if entering {
			if node.IsTitleblock {
				attrs = append(attrs, `class="title"`)
			}
			if node.HeaderID != "" {
				id := r.ensureUniqueHeaderID(node.HeaderID)
				if r.HeaderIDPrefix != "" {
					id = r.HeaderIDPrefix + id
				}
				if r.HeaderIDSuffix != "" {
					id = id + r.HeaderIDSuffix
				}
				attrs = append(attrs, fmt.Sprintf(`id="%s"`, id))
			}
			r.cr(w)
			r.out(w, tag(tagname, attrs, false))
		} else {
			r.out(w, tag("/"+tagname, nil, false))
			if !(node.Parent.Type == Item && node.Next == nil) {
				r.cr(w)
			}
		}
	case HorizontalRule:
		r.cr(w)
		r.out(w, tag("hr", attrs, r.Flags&UseXHTML != 0))
		r.cr(w)
	case List:
		tagName := "ul"
		if node.ListFlags&ListTypeOrdered != 0 {
			tagName = "ol"
		}
		if node.ListFlags&ListTypeDefinition != 0 {
			tagName = "dl"
		}
		if entering {
			if node.IsFootnotesList {
				r.out(w, []byte("\n<div class=\"footnotes\">\n\n"))
				r.out(w, tag("hr", attrs, r.Flags&UseXHTML != 0))
				r.cr(w)
			}
			r.cr(w)
			if node.Parent.Type == Item && node.Parent.Parent.Tight {
				r.cr(w)
			}
			r.out(w, tag(tagName, attrs, false))
			r.cr(w)
		} else {
			r.out(w, tag("/"+tagName, nil, false))
			//cr(w)
			//if node.parent.Type != Item {
			//	cr(w)
			//}
			if node.Parent.Type == Item && node.Next != nil {
				r.cr(w)
			}
			if node.Parent.Type == Document || node.Parent.Type == BlockQuote {
				r.cr(w)
			}
			if node.IsFootnotesList {
				r.out(w, []byte("\n</div>\n"))
			}
		}
	case Item:
		tagName := "li"
		if node.ListFlags&ListTypeDefinition != 0 {
			tagName = "dd"
		}
		if node.ListFlags&ListTypeTerm != 0 {
			tagName = "dt"
		}
		if entering {
			if itemOpenCR(node) {
				r.cr(w)
			}
			if node.ListData.RefLink != nil {
				slug := slugify(node.ListData.RefLink)
				r.out(w, footnoteItem(r.FootnoteAnchorPrefix, slug))
				break
			}
			r.out(w, tag(tagName, nil, false))
		} else {
			if node.ListData.RefLink != nil {
				slug := slugify(node.ListData.RefLink)
				if r.Flags&FootnoteReturnLinks != 0 {
					r.out(w, footnoteReturnLink(r.FootnoteAnchorPrefix, r.FootnoteReturnLinkContents, slug))
				}
			}
			r.out(w, tag("/"+tagName, nil, false))
			r.cr(w)
		}
	case CodeBlock:
		attrs = appendLanguageAttr(attrs, node.Info)
		r.cr(w)
		r.out(w, tag("pre", nil, false))
		r.out(w, tag("code", attrs, false))
		r.out(w, escCode(node.Literal))
		r.out(w, tag("/code", nil, false))
		r.out(w, tag("/pre", nil, false))
		if node.Parent.Type != Item {
			r.cr(w)
		}
	case Table:
		if entering {
			r.cr(w)
			r.out(w, tag("table", nil, false))
		} else {
			r.out(w, tag("/table", nil, false))
			r.cr(w)
		}
	case TableCell:
		tagName := "td"
		if node.IsHeader {
			tagName = "th"
		}
		if entering {
			align := cellAlignment(node.Align)
			if align != "" {
				attrs = append(attrs, fmt.Sprintf(`align="%s"`, align))
			}
			if node.Prev == nil {
				r.cr(w)
			}
			r.out(w, tag(tagName, attrs, false))
		} else {
			r.out(w, tag("/"+tagName, nil, false))
			r.cr(w)
		}
	case TableHead:
		if entering {
			r.cr(w)
			r.out(w, tag("thead", nil, false))
		} else {
			r.out(w, tag("/thead", nil, false))
			r.cr(w)
		}
	case TableBody:
		if entering {
			r.cr(w)
			r.out(w, tag("tbody", nil, false))
			// XXX: this is to adhere to a rather silly test. Should fix test.
			if node.FirstChild == nil {
				r.cr(w)
			}
		} else {
			r.out(w, tag("/tbody", nil, false))
			r.cr(w)
		}
	case TableRow:
		if entering {
			r.cr(w)
			r.out(w, tag("tr", nil, false))
		} else {
			r.out(w, tag("/tr", nil, false))
			r.cr(w)
		}
	default:
		panic("Unknown node type " + node.Type.String())
	}
	return GoToNext
}

func (r *HTMLRenderer) writeDocumentHeader(w *bytes.Buffer) {
	if r.Flags&CompletePage == 0 {
		return
	}
	ending := ""
	if r.Flags&UseXHTML != 0 {
		w.WriteString("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" ")
		w.WriteString("\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
		w.WriteString("<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")
		ending = " /"
	} else {
		w.WriteString("<!DOCTYPE html>\n")
		w.WriteString("<html>\n")
	}
	w.WriteString("<head>\n")
	w.WriteString("  <title>")
	if r.Flags&Smartypants != 0 {
		w.Write(r.sr.Process([]byte(r.Title)))
	} else {
		w.Write(esc([]byte(r.Title)))
	}
	w.WriteString("</title>\n")
	w.WriteString("  <meta name=\"GENERATOR\" content=\"Blackfriday Markdown Processor v")
	w.WriteString(Version)
	w.WriteString("\"")
	w.WriteString(ending)
	w.WriteString(">\n")
	w.WriteString("  <meta charset=\"utf-8\"")
	w.WriteString(ending)
	w.WriteString(">\n")
	if r.CSS != "" {
		w.WriteString("  <link rel=\"stylesheet\" type=\"text/css\" href=\"")
		w.Write(esc([]byte(r.CSS)))
		w.WriteString("\"")
		w.WriteString(ending)
		w.WriteString(">\n")
	}
	if r.Icon != "" {
		w.WriteString("  <link rel=\"icon\" type=\"image/x-icon\" href=\"")
		w.Write(esc([]byte(r.Icon)))
		w.WriteString("\"")
		w.WriteString(ending)
		w.WriteString(">\n")
	}
	w.WriteString("</head>\n")
	w.WriteString("<body>\n\n")
}

func (r *HTMLRenderer) writeTOC(w *bytes.Buffer, ast *Node) {
	buf := bytes.Buffer{}

	inHeader := false
	tocLevel := 0
	headerCount := 0

	ast.Walk(func(node *Node, entering bool) WalkStatus {
		if node.Type == Header && !node.HeaderData.IsTitleblock {
			inHeader = entering
			if entering {
				node.HeaderID = fmt.Sprintf("toc_%d", headerCount)
				if node.Level == tocLevel {
					buf.WriteString("</li>\n\n<li>")
				} else if node.Level < tocLevel {
					for node.Level < tocLevel {
						tocLevel--
						buf.WriteString("</li>\n</ul>")
					}
					buf.WriteString("</li>\n\n<li>")
				} else {
					for node.Level > tocLevel {
						tocLevel++
						buf.WriteString("\n<ul>\n<li>")
					}
				}

				fmt.Fprintf(&buf, `<a href="#toc_%d">`, headerCount)
				headerCount++
			} else {
				buf.WriteString("</a>")
			}
			return GoToNext
		}

		if inHeader {
			return r.RenderNode(&buf, node, entering)
		}

		return GoToNext
	})

	for ; tocLevel > 0; tocLevel-- {
		buf.WriteString("</li>\n</ul>")
	}

	if buf.Len() > 0 {
		w.WriteString("<nav>\n")
		w.Write(buf.Bytes())
		w.WriteString("\n\n</nav>\n")
	}
}

func (r *HTMLRenderer) writeDocumentFooter(w *bytes.Buffer) {
	if r.Flags&CompletePage == 0 {
		return
	}
	w.WriteString("\n</body>\n</html>\n")
}

// Render walks the specified syntax (sub)tree and returns a HTML document.
func (r *HTMLRenderer) Render(ast *Node) []byte {
	//println("render_Blackfriday")
	//dump(ast)
	var buff bytes.Buffer
	r.writeDocumentHeader(&buff)
	if r.Extensions&TOC != 0 || r.Extensions&OmitContents != 0 {
		r.writeTOC(&buff, ast)
		if r.Extensions&OmitContents != 0 {
			return buff.Bytes()
		}
	}
	ast.Walk(func(node *Node, entering bool) WalkStatus {
		return r.RenderNode(&buff, node, entering)
	})
	r.writeDocumentFooter(&buff)
	return buff.Bytes()
}
