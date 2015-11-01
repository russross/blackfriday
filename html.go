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
	"regexp"
	"strconv"
	"strings"
)

type HtmlFlags int

// Html renderer configuration options.
const (
	HtmlFlagsNone           HtmlFlags = 0
	SkipHTML                HtmlFlags = 1 << iota // Skip preformatted HTML blocks
	SkipStyle                                     // Skip embedded <style> elements
	SkipImages                                    // Skip embedded images
	SkipLinks                                     // Skip all links
	Safelink                                      // Only link to trusted protocols
	NofollowLinks                                 // Only link with rel="nofollow"
	NoreferrerLinks                               // Only link with rel="noreferrer"
	HrefTargetBlank                               // Add a blank target
	Toc                                           // Generate a table of contents
	OmitContents                                  // Skip the main contents (for a standalone table of contents)
	CompletePage                                  // Generate a complete HTML page
	UseXHTML                                      // Generate XHTML output instead of HTML
	UseSmartypants                                // Enable smart punctuation substitutions
	SmartypantsFractions                          // Enable smart fractions (with UseSmartypants)
	SmartypantsDashes                             // Enable smart dashes (with UseSmartypants)
	SmartypantsLatexDashes                        // Enable LaTeX-style dashes (with UseSmartypants)
	SmartypantsAngledQuotes                       // Enable angled double quotes (with UseSmartypants) for double quotes rendering
	FootnoteReturnLinks                           // Generate a link at the end of a footnote to return to the source
)

var (
	alignments = []string{
		"left",
		"right",
		"center",
	}

	// TODO: improve this regexp to catch all possible entities:
	htmlEntity = regexp.MustCompile(`&[a-z]{2,5};`)
)

type HtmlRendererParameters struct {
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
}

// Html is a type that implements the Renderer interface for HTML output.
//
// Do not create this directly, instead use the HtmlRenderer function.
type Html struct {
	flags    HtmlFlags
	closeTag string // how to end singleton tags: either " />" or ">"
	title    string // document title
	css      string // optional css file url (used with HTML_COMPLETE_PAGE)

	parameters HtmlRendererParameters

	// table of contents data
	tocMarker    int
	headerCount  int
	currentLevel int
	toc          *bytes.Buffer

	// Track header IDs to prevent ID collision in a single generation.
	headerIDs map[string]int

	smartypants *smartypantsRenderer
}

const (
	xhtmlClose = " />"
	htmlClose  = ">"
)

// HtmlRenderer creates and configures an Html object, which
// satisfies the Renderer interface.
//
// flags is a set of HtmlFlags ORed together.
// title is the title of the document, and css is a URL for the document's
// stylesheet.
// title and css are only used when HTML_COMPLETE_PAGE is selected.
func HtmlRenderer(flags HtmlFlags, title string, css string) Renderer {
	return HtmlRendererWithParameters(flags, title, css, HtmlRendererParameters{})
}

func HtmlRendererWithParameters(flags HtmlFlags, title string,
	css string, renderParameters HtmlRendererParameters) Renderer {
	// configure the rendering engine
	closeTag := htmlClose
	if flags&UseXHTML != 0 {
		closeTag = xhtmlClose
	}

	if renderParameters.FootnoteReturnLinkContents == "" {
		renderParameters.FootnoteReturnLinkContents = `<sup>[return]</sup>`
	}

	return &Html{
		flags:      flags,
		closeTag:   closeTag,
		title:      title,
		css:        css,
		parameters: renderParameters,

		headerCount:  0,
		currentLevel: 0,
		toc:          new(bytes.Buffer),

		headerIDs: make(map[string]int),

		smartypants: smartypants(flags),
	}
}

// Using if statements is a bit faster than a switch statement. As the compiler
// improves, this should be unnecessary this is only worthwhile because
// attrEscape is the single largest CPU user in normal use.
// Also tried using map, but that gave a ~3x slowdown.
func escapeSingleChar(char byte) (string, bool) {
	if char == '"' {
		return "&quot;", true
	}
	if char == '&' {
		return "&amp;", true
	}
	if char == '<' {
		return "&lt;", true
	}
	if char == '>' {
		return "&gt;", true
	}
	return "", false
}

func attrEscape(src []byte) {
	org := 0
	for i, ch := range src {
		if entity, ok := escapeSingleChar(ch); ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}
}

func entityEscapeWithSkip(src []byte, skipRanges [][]int) {
	end := 0
	for _, rang := range skipRanges {
		attrEscape(out, src[end:rang[0]])
		out.Write(src[rang[0]:rang[1]])
		end = rang[1]
	}
	attrEscape(out, src[end:])
}

func (r *Html) GetFlags() HtmlFlags {
	return r.flags
}

func (r *Html) TitleBlock(text []byte) {
	text = bytes.TrimPrefix(text, []byte("% "))
	text = bytes.Replace(text, []byte("\n% "), []byte("\n"), -1)
	out.WriteString("<h1 class=\"title\">")
	out.Write(text)
	out.WriteString("\n</h1>")
}

func (r *Html) BeginHeader(level int, id string) int {
	doubleSpace(out)

	if id == "" && r.flags&Toc != 0 {
		id = fmt.Sprintf("toc_%d", r.headerCount)
	}

	if id != "" {
		id = r.ensureUniqueHeaderID(id)

		if r.parameters.HeaderIDPrefix != "" {
			id = r.parameters.HeaderIDPrefix + id
		}

		if r.parameters.HeaderIDSuffix != "" {
			id = id + r.parameters.HeaderIDSuffix
		}

		out.WriteString(fmt.Sprintf("<h%d id=\"%s\">", level, id))
	} else {
		out.WriteString(fmt.Sprintf("<h%d>", level))
	}

	return out.Len()
}

func (r *Html) EndHeader(level int, id string, tocMarker int) {
	// are we building a table of contents?
	if r.flags&Toc != 0 {
		r.TocHeaderWithAnchor(out.Bytes()[tocMarker:], level, id)
	}

	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func (r *Html) BlockHtml(text []byte) {
	if r.flags&SkipHTML != 0 {
		return
	}

	doubleSpace(out)
	out.Write(text)
	out.WriteByte('\n')
}

func (r *Html) HRule() {
	doubleSpace(out)
	out.WriteString("<hr")
	out.WriteString(r.closeTag)
	out.WriteByte('\n')
}

func (r *Html) BlockCode(text []byte, lang string) {
	doubleSpace(out)

	// parse out the language names/classes
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		if count == 0 {
			out.WriteString("<pre><code class=\"language-")
		} else {
			out.WriteByte(' ')
		}
		attrEscape(out, []byte(elt))
		count++
	}

	if count == 0 {
		out.WriteString("<pre><code>")
	} else {
		out.WriteString("\">")
	}

	attrEscape(out, text)
	out.WriteString("</code></pre>\n")
}

func (r *Html) BlockQuote(text []byte) {
	doubleSpace(out)
	out.WriteString("<blockquote>\n")
	out.Write(text)
	out.WriteString("</blockquote>\n")
}

func (r *Html) Table(header []byte, body []byte, columnData []int) {
	doubleSpace(out)
	out.WriteString("<table>\n<thead>\n")
	out.Write(header)
	out.WriteString("</thead>\n\n<tbody>\n")
	out.Write(body)
	out.WriteString("</tbody>\n</table>\n")
}

func (r *Html) TableRow(text []byte) {
	doubleSpace(out)
	out.WriteString("<tr>\n")
	out.Write(text)
	out.WriteString("\n</tr>\n")
}

func (r *Html) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	doubleSpace(out)
	switch align {
	case TableAlignmentLeft:
		out.WriteString("<th align=\"left\">")
	case TableAlignmentRight:
		out.WriteString("<th align=\"right\">")
	case TableAlignmentCenter:
		out.WriteString("<th align=\"center\">")
	default:
		out.WriteString("<th>")
	}

	out.Write(text)
	out.WriteString("</th>")
}

func (r *Html) TableCell(out *bytes.Buffer, text []byte, align int) {
	doubleSpace(out)
	switch align {
	case TableAlignmentLeft:
		out.WriteString("<td align=\"left\">")
	case TableAlignmentRight:
		out.WriteString("<td align=\"right\">")
	case TableAlignmentCenter:
		out.WriteString("<td align=\"center\">")
	default:
		out.WriteString("<td>")
	}

	out.Write(text)
	out.WriteString("</td>")
}

func (r *Html) BeginFootnotes() {
	out.WriteString("<div class=\"footnotes\">\n")
	r.HRule(out)
	r.BeginList(out, ListTypeOrdered)
}

func (r *Html) EndFootnotes() {
	r.EndList(out, ListTypeOrdered)
	out.WriteString("</div>\n")
}

func (r *Html) FootnoteItem(name, text []byte, flags ListType) {
	if flags&ListItemContainsBlock != 0 || flags&ListItemBeginningOfList != 0 {
		doubleSpace(out)
	}
	slug := slugify(name)
	out.WriteString(`<li id="`)
	out.WriteString(`fn:`)
	out.WriteString(r.parameters.FootnoteAnchorPrefix)
	out.Write(slug)
	out.WriteString(`">`)
	out.Write(text)
	if r.flags&FootnoteReturnLinks != 0 {
		out.WriteString(` <a class="footnote-return" href="#`)
		out.WriteString(`fnref:`)
		out.WriteString(r.parameters.FootnoteAnchorPrefix)
		out.Write(slug)
		out.WriteString(`">`)
		out.WriteString(r.parameters.FootnoteReturnLinkContents)
		out.WriteString(`</a>`)
	}
	out.WriteString("</li>\n")
}

func (r *Html) BeginList(flags ListType) {
	doubleSpace(out)

	if flags&ListTypeDefinition != 0 {
		out.WriteString("<dl>")
	} else if flags&ListTypeOrdered != 0 {
		out.WriteString("<ol>")
	} else {
		out.WriteString("<ul>")
	}
}

func (r *Html) EndList(flags ListType) {
	if flags&ListTypeDefinition != 0 {
		out.WriteString("</dl>\n")
	} else if flags&ListTypeOrdered != 0 {
		out.WriteString("</ol>\n")
	} else {
		out.WriteString("</ul>\n")
	}
}

func (r *Html) ListItem(text []byte, flags ListType) {
	if (flags&ListItemContainsBlock != 0 && flags&ListTypeDefinition == 0) ||
		flags&ListItemBeginningOfList != 0 {
		doubleSpace(out)
	}
	if flags&ListTypeTerm != 0 {
		out.WriteString("<dt>")
	} else if flags&ListTypeDefinition != 0 {
		out.WriteString("<dd>")
	} else {
		out.WriteString("<li>")
	}
	out.Write(text)
	if flags&ListTypeTerm != 0 {
		out.WriteString("</dt>\n")
	} else if flags&ListTypeDefinition != 0 {
		out.WriteString("</dd>\n")
	} else {
		out.WriteString("</li>\n")
	}
}

func (r *Html) BeginParagraph() {
	doubleSpace(out)
	out.WriteString("<p>")
}

func (r *Html) EndParagraph() {
	out.WriteString("</p>\n")
}

func (r *Html) AutoLink(link []byte, kind LinkType) {
	skipRanges := htmlEntity.FindAllIndex(link, -1)
	if r.flags&Safelink != 0 && !isSafeLink(link) && kind != LinkTypeEmail {
		// mark it but don't link it if it is not a safe link: no smartypants
		out.WriteString("<tt>")
		entityEscapeWithSkip(out, link, skipRanges)
		out.WriteString("</tt>")
		return
	}

	out.WriteString("<a href=\"")
	if kind == LinkTypeEmail {
		out.WriteString("mailto:")
	} else {
		r.maybeWriteAbsolutePrefix(out, link)
	}

	entityEscapeWithSkip(out, link, skipRanges)

	var relAttrs []string
	if r.flags&NofollowLinks != 0 && !isRelativeLink(link) {
		relAttrs = append(relAttrs, "nofollow")
	}
	if r.flags&NoreferrerLinks != 0 && !isRelativeLink(link) {
		relAttrs = append(relAttrs, "noreferrer")
	}
	if len(relAttrs) > 0 {
		out.WriteString(fmt.Sprintf("\" rel=\"%s", strings.Join(relAttrs, " ")))
	}

	// blank target only add to external link
	if r.flags&HrefTargetBlank != 0 && !isRelativeLink(link) {
		out.WriteString("\" target=\"_blank")
	}

	out.WriteString("\">")

	// Pretty print: if we get an email address as
	// an actual URI, e.g. `mailto:foo@bar.com`, we don't
	// want to print the `mailto:` prefix
	switch {
	case bytes.HasPrefix(link, []byte("mailto://")):
		attrEscape(out, link[len("mailto://"):])
	case bytes.HasPrefix(link, []byte("mailto:")):
		attrEscape(out, link[len("mailto:"):])
	default:
		entityEscapeWithSkip(out, link, skipRanges)
	}

	out.WriteString("</a>")
}

func (r *Html) CodeSpan(text []byte) {
	out.WriteString("<code>")
	attrEscape(out, text)
	out.WriteString("</code>")
}

func (r *Html) DoubleEmphasis(text []byte) {
	out.WriteString("<strong>")
	out.Write(text)
	out.WriteString("</strong>")
}

func (r *Html) Emphasis(text []byte) {
	if len(text) == 0 {
		return
	}
	out.WriteString("<em>")
	out.Write(text)
	out.WriteString("</em>")
}

func (r *Html) maybeWriteAbsolutePrefix(link []byte) {
	if r.parameters.AbsolutePrefix != "" && isRelativeLink(link) && link[0] != '.' {
		out.WriteString(r.parameters.AbsolutePrefix)
		if link[0] != '/' {
			out.WriteByte('/')
		}
	}
}

func (r *Html) Image(link []byte, title []byte, alt []byte) {
	if r.flags&SkipImages != 0 {
		return
	}

	out.WriteString("<img src=\"")
	r.maybeWriteAbsolutePrefix(out, link)
	attrEscape(out, link)
	out.WriteString("\" alt=\"")
	if len(alt) > 0 {
		attrEscape(out, alt)
	}
	if len(title) > 0 {
		out.WriteString("\" title=\"")
		attrEscape(out, title)
	}

	out.WriteByte('"')
	out.WriteString(r.closeTag)
}

func (r *Html) LineBreak() {
	out.WriteString("<br")
	out.WriteString(r.closeTag)
	out.WriteByte('\n')
}

func (r *Html) Link(link []byte, title []byte, content []byte) {
	if r.flags&SkipLinks != 0 {
		// write the link text out but don't link it, just mark it with typewriter font
		out.WriteString("<tt>")
		attrEscape(out, content)
		out.WriteString("</tt>")
		return
	}

	if r.flags&Safelink != 0 && !isSafeLink(link) {
		// write the link text out but don't link it, just mark it with typewriter font
		out.WriteString("<tt>")
		attrEscape(out, content)
		out.WriteString("</tt>")
		return
	}

	out.WriteString("<a href=\"")
	r.maybeWriteAbsolutePrefix(out, link)
	attrEscape(out, link)
	if len(title) > 0 {
		out.WriteString("\" title=\"")
		attrEscape(out, title)
	}
	var relAttrs []string
	if r.flags&NofollowLinks != 0 && !isRelativeLink(link) {
		relAttrs = append(relAttrs, "nofollow")
	}
	if r.flags&NoreferrerLinks != 0 && !isRelativeLink(link) {
		relAttrs = append(relAttrs, "noreferrer")
	}
	if len(relAttrs) > 0 {
		out.WriteString(fmt.Sprintf("\" rel=\"%s", strings.Join(relAttrs, " ")))
	}

	// blank target only add to external link
	if r.flags&HrefTargetBlank != 0 && !isRelativeLink(link) {
		out.WriteString("\" target=\"_blank")
	}

	out.WriteString("\">")
	out.Write(content)
	out.WriteString("</a>")
	return
}

func (r *Html) RawHtmlTag(text []byte) {
	if r.flags&SkipHTML != 0 {
		return
	}
	if r.flags&SkipStyle != 0 && isHtmlTag(text, "style") {
		return
	}
	if r.flags&SkipLinks != 0 && isHtmlTag(text, "a") {
		return
	}
	if r.flags&SkipImages != 0 && isHtmlTag(text, "img") {
		return
	}
	out.Write(text)
}

func (r *Html) TripleEmphasis(text []byte) {
	out.WriteString("<strong><em>")
	out.Write(text)
	out.WriteString("</em></strong>")
}

func (r *Html) StrikeThrough(text []byte) {
	out.WriteString("<del>")
	out.Write(text)
	out.WriteString("</del>")
}

func (r *Html) FootnoteRef(ref []byte, id int) {
	slug := slugify(ref)
	out.WriteString(`<sup class="footnote-ref" id="`)
	out.WriteString(`fnref:`)
	out.WriteString(r.parameters.FootnoteAnchorPrefix)
	out.Write(slug)
	out.WriteString(`"><a rel="footnote" href="#`)
	out.WriteString(`fn:`)
	out.WriteString(r.parameters.FootnoteAnchorPrefix)
	out.Write(slug)
	out.WriteString(`">`)
	out.WriteString(strconv.Itoa(id))
	out.WriteString(`</a></sup>`)
}

func (r *Html) Entity(entity []byte) {
	out.Write(entity)
}

func (r *Html) NormalText(text []byte) {
	if r.flags&UseSmartypants != 0 {
		r.Smartypants(out, text)
	} else {
		attrEscape(out, text)
	}
}

func (r *Html) Smartypants(text []byte) {
	smrt := smartypantsData{false, false}

	// first do normal entity escaping
	var escaped bytes.Buffer
	attrEscape(&escaped, text)
	text = escaped.Bytes()

	mark := 0
	for i := 0; i < len(text); i++ {
		if action := r.smartypants[text[i]]; action != nil {
			if i > mark {
				out.Write(text[mark:i])
			}

			previousChar := byte(0)
			if i > 0 {
				previousChar = text[i-1]
			}
			i += action(out, &smrt, previousChar, text[i:])
			mark = i + 1
		}
	}

	if mark < len(text) {
		out.Write(text[mark:])
	}
}

func (r *Html) DocumentHeader() {
	if r.flags&CompletePage == 0 {
		return
	}

	ending := ""
	if r.flags&UseXHTML != 0 {
		out.WriteString("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" ")
		out.WriteString("\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
		out.WriteString("<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")
		ending = " /"
	} else {
		out.WriteString("<!DOCTYPE html>\n")
		out.WriteString("<html>\n")
	}
	out.WriteString("<head>\n")
	out.WriteString("  <title>")
	r.NormalText(out, []byte(r.title))
	out.WriteString("</title>\n")
	out.WriteString("  <meta name=\"GENERATOR\" content=\"Blackfriday Markdown Processor v")
	out.WriteString(VERSION)
	out.WriteString("\"")
	out.WriteString(ending)
	out.WriteString(">\n")
	out.WriteString("  <meta charset=\"utf-8\"")
	out.WriteString(ending)
	out.WriteString(">\n")
	if r.css != "" {
		out.WriteString("  <link rel=\"stylesheet\" type=\"text/css\" href=\"")
		attrEscape(out, []byte(r.css))
		out.WriteString("\"")
		out.WriteString(ending)
		out.WriteString(">\n")
	}
	out.WriteString("</head>\n")
	out.WriteString("<body>\n")

	r.tocMarker = out.Len()
}

func (r *Html) DocumentFooter() {
	// finalize and insert the table of contents
	if r.flags&Toc != 0 {
		r.TocFinalize()

		// now we have to insert the table of contents into the document
		var temp bytes.Buffer

		// start by making a copy of everything after the document header
		temp.Write(out.Bytes()[r.tocMarker:])

		// now clear the copied material from the main output buffer
		out.Truncate(r.tocMarker)

		// corner case spacing issue
		if r.flags&CompletePage != 0 {
			out.WriteByte('\n')
		}

		// insert the table of contents
		out.WriteString("<nav>\n")
		out.Write(r.toc.Bytes())
		out.WriteString("</nav>\n")

		// corner case spacing issue
		if r.flags&CompletePage == 0 && r.flags&OmitContents == 0 {
			out.WriteByte('\n')
		}

		// write out everything that came after it
		if r.flags&OmitContents == 0 {
			out.Write(temp.Bytes())
		}
	}

	if r.flags&CompletePage != 0 {
		out.WriteString("\n</body>\n")
		out.WriteString("</html>\n")
	}

}

func (r *Html) TocHeaderWithAnchor(text []byte, level int, anchor string) {
	for level > r.currentLevel {
		switch {
		case bytes.HasSuffix(r.toc.Bytes(), []byte("</li>\n")):
			// this sublist can nest underneath a header
			size := r.toc.Len()
			r.toc.Truncate(size - len("</li>\n"))

		case r.currentLevel > 0:
			r.toc.WriteString("<li>")
		}
		if r.toc.Len() > 0 {
			r.toc.WriteByte('\n')
		}
		r.toc.WriteString("<ul>\n")
		r.currentLevel++
	}

	for level < r.currentLevel {
		r.toc.WriteString("</ul>")
		if r.currentLevel > 1 {
			r.toc.WriteString("</li>\n")
		}
		r.currentLevel--
	}

	r.toc.WriteString("<li><a href=\"#")
	if anchor != "" {
		r.toc.WriteString(anchor)
	} else {
		r.toc.WriteString("toc_")
		r.toc.WriteString(strconv.Itoa(r.headerCount))
	}
	r.toc.WriteString("\">")
	r.headerCount++

	r.toc.Write(text)

	r.toc.WriteString("</a></li>\n")
}

func (r *Html) TocHeader(text []byte, level int) {
	r.TocHeaderWithAnchor(text, level, "")
}

func (r *Html) TocFinalize() {
	for r.currentLevel > 1 {
		r.toc.WriteString("</ul></li>\n")
		r.currentLevel--
	}

	if r.currentLevel > 0 {
		r.toc.WriteString("</ul>\n")
	}
}

func isHtmlTag(tag []byte, tagname string) bool {
	found, _ := findHtmlTagPos(tag, tagname)
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

func findHtmlTagPos(tag []byte, tagname string) (bool, int) {
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
	if rightAngle > i {
		return true, rightAngle
	}

	return false, -1
}

func skipUntilChar(text []byte, start int, char byte) int {
	i := start
	for i < len(text) && text[i] != char {
		i++
	}
	return i
}

func skipSpace(tag []byte, i int) int {
	for i < len(tag) && isspace(tag[i]) {
		i++
	}
	return i
}

func skipChar(data []byte, start int, char byte) int {
	i := start
	for i < len(data) && data[i] == char {
		i++
	}
	return i
}

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
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

func (r *Html) ensureUniqueHeaderID(id string) string {
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
