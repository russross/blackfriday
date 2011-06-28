//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2011 Russ Ross <russ@russross.com>.
// Licensed under the Simplified BSD License.
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
	HTML_GITHUB_BLOCKCODE
	HTML_USE_XHTML
	HTML_USE_SMARTYPANTS
	HTML_SMARTYPANTS_FRACTIONS
	HTML_SMARTYPANTS_LATEX_DASHES
)

type htmlOptions struct {
	flags    int
	closeTag string // how to end singleton tags: either " />\n" or ">\n"
	tocData  struct {
		headerCount  int
		currentLevel int
	}
	smartypants *SmartypantsRenderer
}

var xhtmlClose = " />\n"
var htmlClose = ">\n"

func HtmlRenderer(flags int) *Renderer {
	// configure the rendering engine
	r := new(Renderer)
	if flags&HTML_GITHUB_BLOCKCODE == 0 {
		r.BlockCode = htmlBlockCode
	} else {
		r.BlockCode = htmlBlockCodeGithub
	}
	r.BlockQuote = htmlBlockQuote
	if flags&HTML_SKIP_HTML == 0 {
		r.BlockHtml = htmlRawBlock
	}
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
	if flags&HTML_SKIP_IMAGES == 0 {
		r.Image = htmlImage
	}
	r.LineBreak = htmlLineBreak
	if flags&HTML_SKIP_LINKS == 0 {
		r.Link = htmlLink
	}
	r.RawHtmlTag = htmlRawTag
	r.TripleEmphasis = htmlTripleEmphasis
	r.StrikeThrough = htmlStrikeThrough

	var cb *SmartypantsRenderer
	if flags&HTML_USE_SMARTYPANTS == 0 {
		r.NormalText = htmlNormalText
	} else {
		cb = Smartypants(flags)
		r.NormalText = htmlSmartypants
	}

	closeTag := htmlClose
	if flags&HTML_USE_XHTML != 0 {
		closeTag = xhtmlClose
	}
	r.Opaque = &htmlOptions{flags: flags, closeTag: closeTag, smartypants: cb}
	return r
}

func HtmlTocRenderer(flags int) *Renderer {
	// configure the rendering engine
	r := new(Renderer)
	r.Header = htmlTocHeader

	r.CodeSpan = htmlCodeSpan
	r.DoubleEmphasis = htmlDoubleEmphasis
	r.Emphasis = htmlEmphasis
	r.TripleEmphasis = htmlTripleEmphasis
	r.StrikeThrough = htmlStrikeThrough

	r.DocumentFooter = htmlTocFinalize

	closeTag := ">\n"
	if flags&HTML_USE_XHTML != 0 {
		closeTag = " />\n"
	}
	r.Opaque = &htmlOptions{flags: flags | HTML_TOC, closeTag: closeTag}
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
		out.WriteString(fmt.Sprintf("<h%d id=\"toc_%d\">", level, options.tocData.headerCount))
		options.tocData.headerCount++
	} else {
		out.WriteString(fmt.Sprintf("<h%d>", level))
	}

	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func htmlRawBlock(out *bytes.Buffer, text []byte, opaque interface{}) {
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

func htmlAutoLink(out *bytes.Buffer, link []byte, kind int, opaque interface{}) int {
	options := opaque.(*htmlOptions)

	if len(link) == 0 {
		return 0
	}
	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) && kind != LINK_TYPE_EMAIL {
		return 0
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

	return 1
}

func htmlCodeSpan(out *bytes.Buffer, text []byte, opaque interface{}) int {
	out.WriteString("<code>")
	attrEscape(out, text)
	out.WriteString("</code>")
	return 1
}

func htmlDoubleEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	out.WriteString("<strong>")
	out.Write(text)
	out.WriteString("</strong>")
	return 1
}

func htmlEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	out.WriteString("<em>")
	out.Write(text)
	out.WriteString("</em>")
	return 1
}

func htmlImage(out *bytes.Buffer, link []byte, title []byte, alt []byte, opaque interface{}) int {
	options := opaque.(*htmlOptions)
	if len(link) == 0 {
		return 0
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
	return 1
}

func htmlLineBreak(out *bytes.Buffer, opaque interface{}) int {
	options := opaque.(*htmlOptions)
	out.WriteString("<br")
	out.WriteString(options.closeTag)
	return 1
}

func htmlLink(out *bytes.Buffer, link []byte, title []byte, content []byte, opaque interface{}) int {
	options := opaque.(*htmlOptions)

	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) {
		return 0
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
	return 1
}

func htmlRawTag(out *bytes.Buffer, text []byte, opaque interface{}) int {
	options := opaque.(*htmlOptions)
	if options.flags&HTML_SKIP_HTML != 0 {
		return 1
	}
	if options.flags&HTML_SKIP_STYLE != 0 && isHtmlTag(text, "style") {
		return 1
	}
	if options.flags&HTML_SKIP_LINKS != 0 && isHtmlTag(text, "a") {
		return 1
	}
	if options.flags&HTML_SKIP_IMAGES != 0 && isHtmlTag(text, "img") {
		return 1
	}
	out.Write(text)
	return 1
}

func htmlTripleEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	out.WriteString("<strong><em>")
	out.Write(text)
	out.WriteString("</em></strong>")
	return 1
}

func htmlStrikeThrough(out *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	out.WriteString("<del>")
	out.Write(text)
	out.WriteString("</del>")
	return 1
}

func htmlNormalText(out *bytes.Buffer, text []byte, opaque interface{}) {
	attrEscape(out, text)
}

func htmlTocHeader(out *bytes.Buffer, text func() bool, level int, opaque interface{}) {
	options := opaque.(*htmlOptions)
	marker := out.Len()

	for level > options.tocData.currentLevel {
		if options.tocData.currentLevel > 0 {
			out.WriteString("<li>")
		}
		out.WriteString("<ul>\n")
		options.tocData.currentLevel++
	}

	for level < options.tocData.currentLevel {
		out.WriteString("</ul>")
		if options.tocData.currentLevel > 1 {
			out.WriteString("</li>\n")
		}
		options.tocData.currentLevel--
	}

	out.WriteString("<li><a href=\"#toc_")
	out.WriteString(strconv.Itoa(options.tocData.headerCount))
	out.WriteString("\">")
	options.tocData.headerCount++

	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</a></li>\n")
}

func htmlTocFinalize(out *bytes.Buffer, opaque interface{}) {
	options := opaque.(*htmlOptions)
	for options.tocData.currentLevel > 1 {
		out.WriteString("</ul></li>\n")
		options.tocData.currentLevel--
	}

	if options.tocData.currentLevel > 0 {
		out.WriteString("</ul>\n")
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

	tag_i := i
	for ; i < len(tag); i, tag_i = i+1, tag_i+1 {
		if tag_i >= len(tagname) {
			break
		}

		if tag[i] != tagname[tag_i] {
			return false
		}
	}

	if i == len(tag) {
		return false
	}

	return isspace(tag[i]) || tag[i] == '>'
}
