//
// Black Friday Markdown Processor
// Originally based on http://github.com/tanoku/upskirt
// by Russ Ross <russ@russross.com>
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
	HTML_EXPAND_TABS
	HTML_SAFELINK
	HTML_TOC
	HTML_HARD_WRAP
	HTML_GITHUB_BLOCKCODE
	HTML_USE_XHTML
	HTML_USE_SMARTYPANTS
	HTML_SMARTYPANTS_FRACTIONS
	HTML_SMARTYPANTS_LATEX_DASHES
)

type htmlOptions struct {
	flags     int
	close_tag string // how to end singleton tags: usually " />\n", possibly ">\n"
	toc_data  struct {
		header_count  int
		current_level int
	}
	smartypants *SmartypantsRenderer
}

var xhtml_close = " />\n"
var html_close = ">\n"

func HtmlRenderer(flags int) *Renderer {
	// configure the rendering engine
	r := new(Renderer)
	if flags&HTML_GITHUB_BLOCKCODE == 0 {
		r.blockcode = htmlBlockcode
	} else {
		r.blockcode = htmlBlockcodeGithub
	}
	r.blockquote = htmlBlockquote
	if flags&HTML_SKIP_HTML == 0 {
		r.blockhtml = htmlRawBlock
	}
	r.header = htmlHeader
	r.hrule = htmlHrule
	r.list = htmlList
	r.listitem = htmlListitem
	r.paragraph = htmlParagraph
	r.table = htmlTable
	r.tableRow = htmlTablerow
	r.tableCell = htmlTablecell

	r.autolink = htmlAutolink
	r.codespan = htmlCodespan
	r.doubleEmphasis = htmlDoubleEmphasis
	r.emphasis = htmlEmphasis
	if flags&HTML_SKIP_IMAGES == 0 {
		r.image = htmlImage
	}
	r.linebreak = htmlLinebreak
	if flags&HTML_SKIP_LINKS == 0 {
		r.link = htmlLink
	}
	r.rawHtmlTag = htmlRawTag
	r.tripleEmphasis = htmlTripleEmphasis
	r.strikethrough = htmlStrikethrough

	var cb *SmartypantsRenderer
	if flags&HTML_USE_SMARTYPANTS == 0 {
		r.normalText = htmlNormalText
	} else {
		cb = Smartypants(flags)
		r.normalText = htmlSmartypants
	}

	close_tag := html_close
	if flags&HTML_USE_XHTML != 0 {
		close_tag = xhtml_close
	}
	r.opaque = &htmlOptions{flags: flags, close_tag: close_tag, smartypants: cb}
	return r
}

func HtmlTocRenderer(flags int) *Renderer {
	// configure the rendering engine
	r := new(Renderer)
	r.header = htmlTocHeader

	r.codespan = htmlCodespan
	r.doubleEmphasis = htmlDoubleEmphasis
	r.emphasis = htmlEmphasis
	r.tripleEmphasis = htmlTripleEmphasis
	r.strikethrough = htmlStrikethrough

	r.documentFooter = htmlTocFinalize

	close_tag := ">\n"
	if flags&HTML_USE_XHTML != 0 {
		close_tag = " />\n"
	}
	r.opaque = &htmlOptions{flags: flags | HTML_TOC, close_tag: close_tag}
	return r
}

func attrEscape(ob *bytes.Buffer, src []byte) {
	for i := 0; i < len(src); i++ {
		// directly copy normal characters
		org := i
		for i < len(src) && src[i] != '<' && src[i] != '>' && src[i] != '&' && src[i] != '"' {
			i++
		}
		if i > org {
			ob.Write(src[org:i])
		}

		// escape a character
		if i >= len(src) {
			break
		}
		switch src[i] {
		case '<':
			ob.WriteString("&lt;")
		case '>':
			ob.WriteString("&gt;")
		case '&':
			ob.WriteString("&amp;")
		case '"':
			ob.WriteString("&quot;")
		}
	}
}

func htmlHeader(ob *bytes.Buffer, text []byte, level int, opaque interface{}) {
	options := opaque.(*htmlOptions)

	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}

	if options.flags&HTML_TOC != 0 {
		ob.WriteString(fmt.Sprintf("<h%d id=\"toc_%d\">", level, options.toc_data.header_count))
		options.toc_data.header_count++
	} else {
		ob.WriteString(fmt.Sprintf("<h%d>", level))
	}

	ob.Write(text)
	ob.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func htmlRawBlock(ob *bytes.Buffer, text []byte, opaque interface{}) {
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
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.Write(text[org:sz])
	ob.WriteByte('\n')
}

func htmlHrule(ob *bytes.Buffer, opaque interface{}) {
	options := opaque.(*htmlOptions)

	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.WriteString("<hr")
	ob.WriteString(options.close_tag)
}

func htmlBlockcode(ob *bytes.Buffer, text []byte, lang string, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}

	if lang != "" {
		ob.WriteString("<pre><code class=\"")

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
					ob.WriteByte(' ')
				}
				attrEscape(ob, []byte(lang[org:]))
			}
		}

		ob.WriteString("\">")
	} else {
		ob.WriteString("<pre><code>")
	}

	if len(text) > 0 {
		attrEscape(ob, text)
	}

	ob.WriteString("</code></pre>\n")
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
func htmlBlockcodeGithub(ob *bytes.Buffer, text []byte, lang string, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}

	if len(lang) > 0 {
		ob.WriteString("<pre lang=\"")

		i := 0
		for i < len(lang) && !isspace(lang[i]) {
			i++
		}

		if lang[0] == '.' {
			attrEscape(ob, []byte(lang[1:i]))
		} else {
			attrEscape(ob, []byte(lang[:i]))
		}

		ob.WriteString("\"><code>")
	} else {
		ob.WriteString("<pre><code>")
	}

	if len(text) > 0 {
		attrEscape(ob, text)
	}

	ob.WriteString("</code></pre>\n")
}


func htmlBlockquote(ob *bytes.Buffer, text []byte, opaque interface{}) {
	ob.WriteString("<blockquote>\n")
	ob.Write(text)
	ob.WriteString("</blockquote>")
}

func htmlTable(ob *bytes.Buffer, header []byte, body []byte, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.WriteString("<table><thead>\n")
	ob.Write(header)
	ob.WriteString("\n</thead><tbody>\n")
	ob.Write(body)
	ob.WriteString("\n</tbody></table>")
}

func htmlTablerow(ob *bytes.Buffer, text []byte, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	ob.WriteString("<tr>\n")
	ob.Write(text)
	ob.WriteString("\n</tr>")
}

func htmlTablecell(ob *bytes.Buffer, text []byte, align int, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	switch align {
	case TABLE_ALIGNMENT_LEFT:
		ob.WriteString("<td align=\"left\">")
	case TABLE_ALIGNMENT_RIGHT:
		ob.WriteString("<td align=\"right\">")
	case TABLE_ALIGNMENT_CENTER:
		ob.WriteString("<td align=\"center\">")
	default:
		ob.WriteString("<td>")
	}

	ob.Write(text)
	ob.WriteString("</td>")
}

func htmlList(ob *bytes.Buffer, text []byte, flags int, opaque interface{}) {
	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}
	if flags&LIST_TYPE_ORDERED != 0 {
		ob.WriteString("<ol>\n")
	} else {
		ob.WriteString("<ul>\n")
	}
	ob.Write(text)
	if flags&LIST_TYPE_ORDERED != 0 {
		ob.WriteString("</ol>\n")
	} else {
		ob.WriteString("</ul>\n")
	}
}

func htmlListitem(ob *bytes.Buffer, text []byte, flags int, opaque interface{}) {
	ob.WriteString("<li>")
	size := len(text)
	for size > 0 && text[size-1] == '\n' {
		size--
	}
	ob.Write(text[:size])
	ob.WriteString("</li>\n")
}

func htmlParagraph(ob *bytes.Buffer, text []byte, opaque interface{}) {
	options := opaque.(*htmlOptions)
	i := 0

	if ob.Len() > 0 {
		ob.WriteByte('\n')
	}

	if len(text) == 0 {
		return
	}

	for i < len(text) && isspace(text[i]) {
		i++
	}

	if i == len(text) {
		return
	}

	ob.WriteString("<p>")
	if options.flags&HTML_HARD_WRAP != 0 {
		for i < len(text) {
			org := i
			for i < len(text) && text[i] != '\n' {
				i++
			}

			if i > org {
				ob.Write(text[org:i])
			}

			if i >= len(text) {
				break
			}

			ob.WriteString("<br>")
			ob.WriteString(options.close_tag)
			i++
		}
	} else {
		ob.Write(text[i:])
	}
	ob.WriteString("</p>\n")
}

func htmlAutolink(ob *bytes.Buffer, link []byte, kind int, opaque interface{}) int {
	options := opaque.(*htmlOptions)

	if len(link) == 0 {
		return 0
	}
	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) && kind != LINK_TYPE_EMAIL {
		return 0
	}

	ob.WriteString("<a href=\"")
	if kind == LINK_TYPE_EMAIL {
		ob.WriteString("mailto:")
	}
	ob.Write(link)
	ob.WriteString("\">")

	/*
	 * Pretty print: if we get an email address as
	 * an actual URI, e.g. `mailto:foo@bar.com`, we don't
	 * want to print the `mailto:` prefix
	 */
	if bytes.HasPrefix(link, []byte("mailto:")) {
		attrEscape(ob, link[7:])
	} else {
		attrEscape(ob, link)
	}

	ob.WriteString("</a>")

	return 1
}

func htmlCodespan(ob *bytes.Buffer, text []byte, opaque interface{}) int {
	ob.WriteString("<code>")
	attrEscape(ob, text)
	ob.WriteString("</code>")
	return 1
}

func htmlDoubleEmphasis(ob *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	ob.WriteString("<strong>")
	ob.Write(text)
	ob.WriteString("</strong>")
	return 1
}

func htmlEmphasis(ob *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	ob.WriteString("<em>")
	ob.Write(text)
	ob.WriteString("</em>")
	return 1
}

func htmlImage(ob *bytes.Buffer, link []byte, title []byte, alt []byte, opaque interface{}) int {
	options := opaque.(*htmlOptions)
	if len(link) == 0 {
		return 0
	}
	ob.WriteString("<img src=\"")
	attrEscape(ob, link)
	ob.WriteString("\" alt=\"")
	if len(alt) > 0 {
		attrEscape(ob, alt)
	}
	if len(title) > 0 {
		ob.WriteString("\" title=\"")
		attrEscape(ob, title)
	}

	ob.WriteByte('"')
	ob.WriteString(options.close_tag)
	return 1
}

func htmlLinebreak(ob *bytes.Buffer, opaque interface{}) int {
	options := opaque.(*htmlOptions)
	ob.WriteString("<br")
	ob.WriteString(options.close_tag)
	return 1
}

func htmlLink(ob *bytes.Buffer, link []byte, title []byte, content []byte, opaque interface{}) int {
	options := opaque.(*htmlOptions)

	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) {
		return 0
	}

	ob.WriteString("<a href=\"")
	if len(link) > 0 {
		ob.Write(link)
	}
	if len(title) > 0 {
		ob.WriteString("\" title=\"")
		attrEscape(ob, title)
	}
	ob.WriteString("\">")
	if len(content) > 0 {
		ob.Write(content)
	}
	ob.WriteString("</a>")
	return 1
}

func htmlRawTag(ob *bytes.Buffer, text []byte, opaque interface{}) int {
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
	ob.Write(text)
	return 1
}

func htmlTripleEmphasis(ob *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	ob.WriteString("<strong><em>")
	ob.Write(text)
	ob.WriteString("</em></strong>")
	return 1
}

func htmlStrikethrough(ob *bytes.Buffer, text []byte, opaque interface{}) int {
	if len(text) == 0 {
		return 0
	}
	ob.WriteString("<del>")
	ob.Write(text)
	ob.WriteString("</del>")
	return 1
}

func htmlNormalText(ob *bytes.Buffer, text []byte, opaque interface{}) {
	attrEscape(ob, text)
}

func htmlTocHeader(ob *bytes.Buffer, text []byte, level int, opaque interface{}) {
	options := opaque.(*htmlOptions)
	for level > options.toc_data.current_level {
		if options.toc_data.current_level > 0 {
			ob.WriteString("<li>")
		}
		ob.WriteString("<ul>\n")
		options.toc_data.current_level++
	}

	for level < options.toc_data.current_level {
		ob.WriteString("</ul>")
		if options.toc_data.current_level > 1 {
			ob.WriteString("</li>\n")
		}
		options.toc_data.current_level--
	}

	ob.WriteString("<li><a href=\"#toc_")
	ob.WriteString(strconv.Itoa(options.toc_data.header_count))
	ob.WriteString("\">")
	options.toc_data.header_count++

	if len(text) > 0 {
		ob.Write(text)
	}
	ob.WriteString("</a></li>\n")
}

func htmlTocFinalize(ob *bytes.Buffer, opaque interface{}) {
	options := opaque.(*htmlOptions)
	for options.toc_data.current_level > 1 {
		ob.WriteString("</ul></li>\n")
		options.toc_data.current_level--
	}

	if options.toc_data.current_level > 0 {
		ob.WriteString("</ul>\n")
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
