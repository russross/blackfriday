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
	"strconv"
)

const (
	HTML_SKIP_HTML = 1 << iota
	HTML_SKIP_STYLE
	HTML_SKIP_IMAGES
	HTML_SKIP_LINKS
	HTML_SAFELINK
	HTML_TOC
	HTML_OMIT_CONTENTS
	HTML_COMPLETE_PAGE
	HTML_GITHUB_BLOCKCODE
	HTML_USE_XHTML
	HTML_USE_SMARTYPANTS
	HTML_SMARTYPANTS_FRACTIONS
	HTML_SMARTYPANTS_LATEX_DASHES
)

type htmlOptions struct {
	flags    int    // HTML_* options
	closeTag string // how to end singleton tags: either " />\n" or ">\n"
	title    string // document title
	css      string // optional css file url (used with HTML_COMPLETE_PAGE)

	// table of contents data
	tocMarker    int
	headerCount  int
	currentLevel int
	toc          *bytes.Buffer

	smartypants *SmartypantsRenderer
}

var xhtmlClose = " />\n"
var htmlClose = ">\n"

func HtmlRenderer(flags int, title string, css string) *Renderer {
	// configure the rendering engine
	r := new(Renderer)

	r.BlockCode = htmlBlockCode
	r.BlockQuote = htmlBlockQuote
	r.BlockHtml = htmlBlockHtml
	r.Header = htmlHeader
	r.HRule = htmlHRule
	r.List = htmlList
	r.ListItem = htmlListItem
	r.Paragraph = htmlParagraph
	r.Table = htmlTable
	r.TableRow = htmlTableRow
	r.TableCell = htmlTableCell

	r.AutoLink = htmlAutoLink
	r.CodeSpan = htmlCodeSpan
	r.DoubleEmphasis = htmlDoubleEmphasis
	r.Emphasis = htmlEmphasis
	r.Image = htmlImage
	r.LineBreak = htmlLineBreak
	r.Link = htmlLink
	r.RawHtmlTag = htmlRawTag
	r.TripleEmphasis = htmlTripleEmphasis
	r.StrikeThrough = htmlStrikeThrough

	r.Entity = htmlEntity
	r.NormalText = htmlNormalText

	r.DocumentHeader = htmlDocumentHeader
	r.DocumentFooter = htmlDocumentFooter

	closeTag := htmlClose
	if flags&HTML_USE_XHTML != 0 {
		closeTag = xhtmlClose
	}

	r.Opaque = &htmlOptions{
		flags:    flags,
		closeTag: closeTag,
		title:    title,
		css:      css,

		headerCount:  0,
		currentLevel: 0,
		toc:          new(bytes.Buffer),

		smartypants: Smartypants(flags),
	}
	return r
}

func attrEscape(out *bytes.Buffer, src []byte) {
	org := 0
	for i, ch := range src {
		// using if statements is a bit faster than a switch statement.
		// as the compiler improves, this should be unnecessary
		// this is only worthwhile because attrEscape is the single
		// largest CPU user in normal use
		if ch == '"' {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString("&quot;")
			continue
		}
		if ch == '&' {
			if i > org {
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString("&amp;")
			continue
		}
		if ch == '<' {
			if i > org {
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString("&lt;")
			continue
		}
		if ch == '>' {
			if i > org {
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString("&gt;")
			continue
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}
}

func htmlHeader(out *bytes.Buffer, text func() bool, level int, opaque interface{}) {
	options := opaque.(*htmlOptions)
	marker := out.Len()

	if marker > 0 {
		out.WriteByte('\n')
	}

	if options.flags&HTML_TOC != 0 {
		// headerCount is incremented in htmlTocHeader
		out.WriteString(fmt.Sprintf("<h%d id=\"toc_%d\">", level, options.headerCount))
	} else {
		out.WriteString(fmt.Sprintf("<h%d>", level))
	}

	tocMarker := out.Len()
	if !text() {
		out.Truncate(marker)
		return
	}

	// are we building a table of contents?
	if options.flags&HTML_TOC != 0 {
		htmlTocHeader(out.Bytes()[tocMarker:], level, opaque)
	}

	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func htmlBlockHtml(out *bytes.Buffer, text []byte, opaque interface{}) {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_SKIP_HTML != 0 {
		return
	}

	sz := len(text)
	for sz > 0 && text[sz-1] == '\n' {
		sz--
	}
	org := 0
	for org < sz && text[org] == '\n' {
		org++
	}
	if org >= sz {
		return
	}
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
	out.Write(text[org:sz])
	out.WriteByte('\n')
}

func htmlHRule(out *bytes.Buffer, opaque interface{}) {
	options := opaque.(*htmlOptions)

	if out.Len() > 0 {
		out.WriteByte('\n')
	}
	out.WriteString("<hr")
	out.WriteString(options.closeTag)
}

func htmlBlockCode(out *bytes.Buffer, text []byte, lang string, opaque interface{}) {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_GITHUB_BLOCKCODE != 0 {
		htmlBlockCodeGithub(out, text, lang, opaque)
	} else {
		htmlBlockCodeNormal(out, text, lang, opaque)
	}
}

func htmlBlockCodeNormal(out *bytes.Buffer, text []byte, lang string, opaque interface{}) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}

	if lang != "" {
		out.WriteString("<pre><code class=\"")

		for i, cls := 0, 0; i < len(lang); i, cls = i+1, cls+1 {
			for i < len(lang) && isspace(lang[i]) {
				i++
			}

			if i < len(lang) {
				org := i
				for i < len(lang) && !isspace(lang[i]) {
					i++
				}

				if lang[org] == '.' {
					org++
				}

				if cls > 0 {
					out.WriteByte(' ')
				}
				attrEscape(out, []byte(lang[org:]))
			}
		}

		out.WriteString("\">")
	} else {
		out.WriteString("<pre><code>")
	}

	if len(text) > 0 {
		attrEscape(out, text)
	}

	out.WriteString("</code></pre>\n")
}

/*
 * GitHub style code block:
 *
 *              <pre lang="LANG"><code>
 *              ...
 *              </pre></code>
 *
 * Unlike other parsers, we store the language identifier in the <pre>,
 * and don't let the user generate custom classes.
 *
 * The language identifier in the <pre> block gets postprocessed and all
 * the code inside gets syntax highlighted with Pygments. This is much safer
 * than letting the user specify a CSS class for highlighting.
 *
 * Note that we only generate HTML for the first specifier.
 * E.g.
 *              ~~~~ {.python .numbered}        =>      <pre lang="python"><code>
 */
func htmlBlockCodeGithub(out *bytes.Buffer, text []byte, lang string, opaque interface{}) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}

	if len(lang) > 0 {
		out.WriteString("<pre lang=\"")

		i := 0
		for i < len(lang) && !isspace(lang[i]) {
			i++
		}

		if lang[0] == '.' {
			attrEscape(out, []byte(lang[1:i]))
		} else {
			attrEscape(out, []byte(lang[:i]))
		}

		out.WriteString("\"><code>")
	} else {
		out.WriteString("<pre><code>")
	}

	if len(text) > 0 {
		attrEscape(out, text)
	}

	out.WriteString("</code></pre>\n")
}


func htmlBlockQuote(out *bytes.Buffer, text []byte, opaque interface{}) {
	out.WriteString("<blockquote>\n")
	out.Write(text)
	out.WriteString("</blockquote>")
}

func htmlTable(out *bytes.Buffer, header []byte, body []byte, columnData []int, opaque interface{}) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
	out.WriteString("<table><thead>\n")
	out.Write(header)
	out.WriteString("\n</thead><tbody>\n")
	out.Write(body)
	out.WriteString("\n</tbody></table>")
}

func htmlTableRow(out *bytes.Buffer, text []byte, opaque interface{}) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
	out.WriteString("<tr>\n")
	out.Write(text)
	out.WriteString("\n</tr>")
}

func htmlTableCell(out *bytes.Buffer, text []byte, align int, opaque interface{}) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
	switch align {
	case TABLE_ALIGNMENT_LEFT:
		out.WriteString("<td align=\"left\">")
	case TABLE_ALIGNMENT_RIGHT:
		out.WriteString("<td align=\"right\">")
	case TABLE_ALIGNMENT_CENTER:
		out.WriteString("<td align=\"center\">")
	default:
		out.WriteString("<td>")
	}

	out.Write(text)
	out.WriteString("</td>")
}

func htmlList(out *bytes.Buffer, text func() bool, flags int, opaque interface{}) {
	marker := out.Len()

	if marker > 0 {
		out.WriteByte('\n')
	}
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString("<ol>\n")
	} else {
		out.WriteString("<ul>\n")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString("</ol>\n")
	} else {
		out.WriteString("</ul>\n")
	}
}

func htmlListItem(out *bytes.Buffer, text []byte, flags int, opaque interface{}) {
	out.WriteString("<li>")
	size := len(text)
	for size > 0 && text[size-1] == '\n' {
		size--
	}
	out.Write(text[:size])
	out.WriteString("</li>\n")
}

func htmlParagraph(out *bytes.Buffer, text func() bool, opaque interface{}) {
	marker := out.Len()
	if marker > 0 {
		out.WriteByte('\n')
	}

	out.WriteString("<p>")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</p>\n")
}

func htmlAutoLink(out *bytes.Buffer, link []byte, kind int, opaque interface{}) bool {
	options := opaque.(*htmlOptions)

	if len(link) == 0 {
		return false
	}
	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) && kind != LINK_TYPE_EMAIL {
		return false
	}

	out.WriteString("<a href=\"")
	if kind == LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	attrEscape(out, link)
	out.WriteString("\">")

	/*
	 * Pretty print: if we get an email address as
	 * an actual URI, e.g. `mailto:foo@bar.com`, we don't
	 * want to print the `mailto:` prefix
	 */
	switch {
	case bytes.HasPrefix(link, []byte("mailto://")):
		attrEscape(out, link[9:])
	case bytes.HasPrefix(link, []byte("mailto:")):
		attrEscape(out, link[7:])
	default:
		attrEscape(out, link)
	}

	out.WriteString("</a>")

	return true
}

func htmlCodeSpan(out *bytes.Buffer, text []byte, opaque interface{}) bool {
	out.WriteString("<code>")
	attrEscape(out, text)
	out.WriteString("</code>")
	return true
}

func htmlDoubleEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) bool {
	if len(text) == 0 {
		return false
	}
	out.WriteString("<strong>")
	out.Write(text)
	out.WriteString("</strong>")
	return true
}

func htmlEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) bool {
	if len(text) == 0 {
		return false
	}
	out.WriteString("<em>")
	out.Write(text)
	out.WriteString("</em>")
	return true
}

func htmlImage(out *bytes.Buffer, link []byte, title []byte, alt []byte, opaque interface{}) bool {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_SKIP_IMAGES != 0 {
		return false
	}

	if len(link) == 0 {
		return false
	}
	out.WriteString("<img src=\"")
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
	out.WriteString(options.closeTag)
	return true
}

func htmlLineBreak(out *bytes.Buffer, opaque interface{}) bool {
	options := opaque.(*htmlOptions)
	out.WriteString("<br")
	out.WriteString(options.closeTag)
	return true
}

func htmlLink(out *bytes.Buffer, link []byte, title []byte, content []byte, opaque interface{}) bool {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_SKIP_LINKS != 0 {
		return false
	}

	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) {
		return false
	}

	out.WriteString("<a href=\"")
	attrEscape(out, link)
	if len(title) > 0 {
		out.WriteString("\" title=\"")
		attrEscape(out, title)
	}
	out.WriteString("\">")
	out.Write(content)
	out.WriteString("</a>")
	return true
}

func htmlRawTag(out *bytes.Buffer, text []byte, opaque interface{}) bool {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_SKIP_HTML != 0 {
		return true
	}
	if options.flags&HTML_SKIP_STYLE != 0 && isHtmlTag(text, "style") {
		return true
	}
	if options.flags&HTML_SKIP_LINKS != 0 && isHtmlTag(text, "a") {
		return true
	}
	if options.flags&HTML_SKIP_IMAGES != 0 && isHtmlTag(text, "img") {
		return true
	}
	out.Write(text)
	return true
}

func htmlTripleEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) bool {
	if len(text) == 0 {
		return false
	}
	out.WriteString("<strong><em>")
	out.Write(text)
	out.WriteString("</em></strong>")
	return true
}

func htmlStrikeThrough(out *bytes.Buffer, text []byte, opaque interface{}) bool {
	if len(text) == 0 {
		return false
	}
	out.WriteString("<del>")
	out.Write(text)
	out.WriteString("</del>")
	return true
}

func htmlEntity(out *bytes.Buffer, entity []byte, opaque interface{}) {
	out.Write(entity)
}

func htmlNormalText(out *bytes.Buffer, text []byte, opaque interface{}) {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_USE_SMARTYPANTS != 0 {
		htmlSmartypants(out, text, opaque)
	} else {
		attrEscape(out, text)
	}
}

func htmlDocumentHeader(out *bytes.Buffer, opaque interface{}) {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_COMPLETE_PAGE == 0 {
		return
	}

	ending := ""
	if options.flags&HTML_USE_XHTML != 0 {
		out.WriteString("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" ")
		out.WriteString("\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
		out.WriteString("<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")
		ending = " /"
	} else {
		out.WriteString("<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.01//EN\" ")
		out.WriteString("\"http://www.w3.org/TR/html4/strict.dtd\">\n")
		out.WriteString("<html>\n")
	}
	out.WriteString("<head>\n")
	out.WriteString("  <title>")
	htmlNormalText(out, []byte(options.title), opaque)
	out.WriteString("</title>\n")
	out.WriteString("  <meta name=\"GENERATOR\" content=\"Blackfriday Markdown Processor v")
	out.WriteString(VERSION)
	out.WriteString("\"")
	out.WriteString(ending)
	out.WriteString(">\n")
	out.WriteString("  <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"")
	out.WriteString(ending)
	out.WriteString(">\n")
	if options.css != "" {
		out.WriteString("  <link rel=\"stylesheet\" type=\"text/css\" href=\"")
		attrEscape(out, []byte(options.css))
		out.WriteString("\"")
		out.WriteString(ending)
		out.WriteString(">\n")
	}
	out.WriteString("</head>\n")
	out.WriteString("<body>\n")

	options.tocMarker = out.Len()
}

func htmlDocumentFooter(out *bytes.Buffer, opaque interface{}) {
	options := opaque.(*htmlOptions)

	// finalize and insert the table of contents
	if options.flags&HTML_TOC != 0 {
		htmlTocFinalize(opaque)

		// now we have to insert the table of contents into the document
		var temp bytes.Buffer

		// start by making a copy of everything after the document header
		temp.Write(out.Bytes()[options.tocMarker:])

		// now clear the copied material from the main output buffer
		out.Truncate(options.tocMarker)

		// insert the table of contents
		out.Write(options.toc.Bytes())

		// write out everything that came after it
		if options.flags&HTML_OMIT_CONTENTS == 0 {
			out.Write(temp.Bytes())
		}
	}

	if options.flags&HTML_COMPLETE_PAGE != 0 {
		out.WriteString("\n</body>\n")
		out.WriteString("</html>\n")
	}

}

func htmlTocHeader(text []byte, level int, opaque interface{}) {
	options := opaque.(*htmlOptions)

	for level > options.currentLevel {
		switch {
		case bytes.HasSuffix(options.toc.Bytes(), []byte("</li>\n")):
			size := options.toc.Len()
			options.toc.Truncate(size - len("</li>\n"))

		case options.currentLevel > 0:
			options.toc.WriteString("<li>")
		}
		options.toc.WriteString("\n<ul>\n")
		options.currentLevel++
	}

	for level < options.currentLevel {
		options.toc.WriteString("</ul>")
		if options.currentLevel > 1 {
			options.toc.WriteString("</li>\n")
		}
		options.currentLevel--
	}

	options.toc.WriteString("<li><a href=\"#toc_")
	options.toc.WriteString(strconv.Itoa(options.headerCount))
	options.toc.WriteString("\">")
	options.headerCount++

	options.toc.Write(text)

	options.toc.WriteString("</a></li>\n")
}

func htmlTocFinalize(opaque interface{}) {
	options := opaque.(*htmlOptions)
	for options.currentLevel > 1 {
		options.toc.WriteString("</ul></li>\n")
		options.currentLevel--
	}

	if options.currentLevel > 0 {
		options.toc.WriteString("</ul>\n")
	}
}

func isHtmlTag(tag []byte, tagname string) bool {
	i := 0
	if i < len(tag) && tag[0] != '<' {
		return false
	}
	i++
	for i < len(tag) && isspace(tag[i]) {
		i++
	}

	if i < len(tag) && tag[i] == '/' {
		i++
	}

	for i < len(tag) && isspace(tag[i]) {
		i++
	}

	j := i
	for ; i < len(tag); i, j = i+1, j+1 {
		if j >= len(tagname) {
			break
		}

		if tag[i] != tagname[j] {
			return false
		}
	}

	if i == len(tag) {
		return false
	}

	return isspace(tag[i]) || tag[i] == '>'
}
