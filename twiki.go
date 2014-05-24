//
// Blackfriday Markdown Processor
// Available at http://github.com/hanguofeng/blackfriday
//
// Copyright © 2014 HanGuofeng <hanguofeng@gmail.com>.
// Distributed under the Simplified BSD License.
// See README.md for details.
//

//
//
// TWiki rendering backend
//
//

package blackfriday

import (
	"bytes"
	"strconv"
	"strings"
)

// TWiki is a type that implements the Renderer interface for TWiki output.
//
// Do not create this directly, instead use the TWikiRenderer function.
type TWiki struct {
	flags    int    // TWiki_* options
	closeTag string // how to end singleton tags: either " />\n" or ">\n"
	title    string // document title
	css      string // optional css file url (used with HTML_COMPLETE_PAGE)

	// table of contents data
	tocMarker    int
	headerCount  int
	currentLevel int
	toc          *bytes.Buffer

	smartypants *smartypantsRenderer
}

// TWikiRenderer creates and configures an TWiki object, which
// satisfies the Renderer interface.
//
// flags is a set of TWiki_* options ORed together.
func TWikiRenderer(flags int) Renderer {
	return &TWiki{}
}

func (options *TWiki) GetFlags() int {
	return options.flags
}

func (options *TWiki) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()

	out.WriteString("\n---")
	out.WriteString(strings.Repeat("+", level))
	out.WriteString(" ")

	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n")
}

func (options *TWiki) BlockHtml(out *bytes.Buffer, text []byte) {

	doubleSpace(out)
	out.Write(text)
	out.WriteByte('\n')
}

func (options *TWiki) HRule(out *bytes.Buffer) {
	doubleSpace(out)
	out.WriteString("-------")
}

func (options *TWiki) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	doubleSpace(out)
	out.WriteString("<verbatim>\n")
	out.Write(text)
	out.WriteString("</verbatim>\n")
}

func (options *TWiki) BlockQuote(out *bytes.Buffer, text []byte) {
	doubleSpace(out)
	out.WriteString("<blockquote>\n")
	out.Write(text)
	out.WriteString("</blockquote>\n")
}

func (options *TWiki) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	doubleSpace(out)
	out.WriteString("<table>\n<thead>\n")
	out.Write(header)
	out.WriteString("</thead>\n\n<tbody>\n")
	out.Write(body)
	out.WriteString("</tbody>\n</table>\n")
}

func (options *TWiki) TableRow(out *bytes.Buffer, text []byte) {
	doubleSpace(out)
	out.WriteString("<tr>\n")
	out.Write(text)
	out.WriteString("\n</tr>\n")
}

func (options *TWiki) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	doubleSpace(out)

	out.WriteString("<th>")
	out.Write(text)
	out.WriteString("</th>")
}

func (options *TWiki) TableCell(out *bytes.Buffer, text []byte, align int) {
	doubleSpace(out)

	out.WriteString("<td>")
	out.Write(text)
	out.WriteString("</td>")
}

func (options *TWiki) Footnotes(out *bytes.Buffer, text func() bool) {
	options.HRule(out)
	options.List(out, text, LIST_TYPE_ORDERED)
	out.WriteString("\n")
}

func (options *TWiki) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	if flags&LIST_ITEM_CONTAINS_BLOCK != 0 || flags&LIST_ITEM_BEGINNING_OF_LIST != 0 {
		doubleSpace(out)
	}

	out.WriteString("\n    * ")
	out.Write(text)
	out.WriteString("\n")
}

func (options *TWiki) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	//doubleSpace(out)

	if !text() {
		out.Truncate(marker)
		return
	}
}

func (options *TWiki) ListItem(out *bytes.Buffer, text []byte, flags int) {
	if flags&LIST_ITEM_CONTAINS_BLOCK != 0 || flags&LIST_ITEM_BEGINNING_OF_LIST != 0 {
		doubleSpace(out)
	}
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString("\n1.")
	} else {
		out.WriteString("\n   * ")
	}
	out.Write(text)
	out.WriteString("\n")
}

func (options *TWiki) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	doubleSpace(out)

	out.WriteString("\n")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n")
}

func (options *TWiki) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	skipRanges := htmlEntity.FindAllIndex(link, -1)
	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) && kind != LINK_TYPE_EMAIL {
		// mark it but don't link it if it is not a safe link: no smartypants
		out.WriteString("<tt>")
		entityEscapeWithSkip(out, link, skipRanges)
		out.WriteString("</tt>")
		return
	}

	out.WriteString("<a href=\"")
	if kind == LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	entityEscapeWithSkip(out, link, skipRanges)

	if options.flags&HTML_NOFOLLOW_LINKS != 0 && !isRelativeLink(link) {
		out.WriteString("\" rel=\"nofollow")
	}
	// blank target only add to external link
	if options.flags&HTML_HREF_TARGET_BLANK != 0 && !isRelativeLink(link) {
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

func (options *TWiki) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString("<verbatim>")
	attrEscape(out, text)
	out.WriteString("</verbatim>")
}

func (options *TWiki) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("*")
	out.Write(text)
	out.WriteString("* ")
}

func (options *TWiki) Emphasis(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	out.WriteString("_")
	out.Write(text)
	out.WriteString("_")
}

func (options *TWiki) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	if options.flags&HTML_SKIP_IMAGES != 0 {
		return
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
	return
}

func (options *TWiki) LineBreak(out *bytes.Buffer) {
	out.WriteString("\n")
}

func (options *TWiki) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if options.flags&HTML_SKIP_LINKS != 0 {
		// write the link text out but don't link it, just mark it with typewriter font
		out.WriteString("<tt>")
		attrEscape(out, content)
		out.WriteString("</tt>")
		return
	}

	if options.flags&HTML_SAFELINK != 0 && !isSafeLink(link) {
		// write the link text out but don't link it, just mark it with typewriter font
		out.WriteString("<tt>")
		attrEscape(out, content)
		out.WriteString("</tt>")
		return
	}

	out.WriteString("<a href=\"")
	attrEscape(out, link)
	if len(title) > 0 {
		out.WriteString("\" title=\"")
		attrEscape(out, title)
	}
	if options.flags&HTML_NOFOLLOW_LINKS != 0 && !isRelativeLink(link) {
		out.WriteString("\" rel=\"nofollow")
	}
	// blank target only add to external link
	if options.flags&HTML_HREF_TARGET_BLANK != 0 && !isRelativeLink(link) {
		out.WriteString("\" target=\"_blank")
	}

	out.WriteString("\">")
	out.Write(content)
	out.WriteString("</a>")
	return
}

func (options *TWiki) RawHtmlTag(out *bytes.Buffer, text []byte) {
	if options.flags&HTML_SKIP_HTML != 0 {
		return
	}
	if options.flags&HTML_SKIP_STYLE != 0 && isHtmlTag(text, "style") {
		return
	}
	if options.flags&HTML_SKIP_LINKS != 0 && isHtmlTag(text, "a") {
		return
	}
	if options.flags&HTML_SKIP_IMAGES != 0 && isHtmlTag(text, "img") {
		return
	}
	out.Write(text)
}

func (options *TWiki) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("__")
	out.Write(text)
	out.WriteString("__")
}

func (options *TWiki) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.WriteString("<del>")
	out.Write(text)
	out.WriteString("</del>")
}

func (options *TWiki) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	slug := slugify(ref)
	out.WriteString(`<sup class="footnote-ref" id="fnref:`)
	out.Write(slug)
	out.WriteString(`"><a rel="footnote" href="#fn:`)
	out.Write(slug)
	out.WriteString(`">`)
	out.WriteString(strconv.Itoa(id))
	out.WriteString(`</a></sup>`)
}

func (options *TWiki) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (options *TWiki) NormalText(out *bytes.Buffer, text []byte) {
	if options.flags&HTML_USE_SMARTYPANTS != 0 {
		options.Smartypants(out, text)
	} else {
		attrEscape(out, text)
	}
}

func (options *TWiki) Smartypants(out *bytes.Buffer, text []byte) {
	smrt := smartypantsData{false, false}

	// first do normal entity escaping
	var escaped bytes.Buffer
	attrEscape(&escaped, text)
	text = escaped.Bytes()

	mark := 0
	for i := 0; i < len(text); i++ {
		if action := options.smartypants[text[i]]; action != nil {
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

func (options *TWiki) DocumentHeader(out *bytes.Buffer) {
}

func (options *TWiki) DocumentFooter(out *bytes.Buffer) {
}
