//
// Blackfriday Markdown Processor
// Available at http://github.com/russross/blackfriday
//
// Copyright © 2014 Miek Gieben <miek@miek.nl>.

//
//
// XML2RFC v3 rendering backend
//
//

package blackfriday

import (
	"bytes"
)

// Xml is a type that implements the Renderer interface for XML2RFV3 output.
//
// Do not create this directly, instead use the XmlRenderer function.
type Xml struct {
}

// XmlRenderer creates and configures a Xml object, which
// satisfies the Renderer interface.
//
// flags is a set of XML_* options ORed together (currently no such options
// are defined).
func XmlRenderer(flags int) Renderer {
	return &Xml{}
}

func (options *Xml) GetFlags() int {
	return 0
}

func (options *Xml) GetState() int {
	return 0
}

// render code chunks using verbatim, or listings if we have a language
func (options *Xml) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	if lang == "" {
		out.WriteString("\n\\begin{verbatim}\n")
	} else {
		out.WriteString("\n\\begin{lstlisting}[language=")
		out.WriteString(lang)
		out.WriteString("]\n")
	}
	out.Write(text)
	if lang == "" {
		out.WriteString("\n\\end{verbatim}\n")
	} else {
		out.WriteString("\n\\end{lstlisting}\n")
	}
}

func (options *Xml) TitleBlock(out *bytes.Buffer, text []byte) {

}

func (options *Xml) BlockQuote(out *bytes.Buffer, text []byte) {
	out.WriteString("\n<blockquote>\n")
	out.Write(text)
	out.WriteString("\n</blockquote>\n")
}

func (options *Xml) BlockHtml(out *bytes.Buffer, text []byte) {
	// a pretty lame thing to do...
	out.WriteString("\n\\begin{verbatim}\n")
	out.Write(text)
	out.WriteString("\n\\end{verbatim}\n")
}

func (options *Xml) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()

	switch level {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		fallthrough
	case 4:
		fallthrough
	case 5:
		fallthrough
	case 6:
		// Don't know  if we need to close one
		out.WriteString("\n</section>\n")
		out.WriteString("\n<section>\n")
	}
	if !text() {
		//		out.WriteString("<name>")
		out.Truncate(marker)
		//		out.WriteString("</name>")
		return
	}
}

func (options *Xml) HRule(out *bytes.Buffer) {
	out.WriteString("\n\\HRule\n")
}

func (options *Xml) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString("\n<ol>\n")
	} else {
		out.WriteString("\n</ol>\n")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString("\n</ol>\n")
	} else {
		out.WriteString("\n</ol>\n")
	}
}

func (options *Xml) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("\n<li> ")
	out.Write(text)
	out.WriteString("\n</li> ")
}

func (options *Xml) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	out.WriteString("<t>\n")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n</t>\n")
}

func (options *Xml) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	out.WriteString("\n\\begin{tabular}{")
	for _, elt := range columnData {
		switch elt {
		case TABLE_ALIGNMENT_LEFT:
			out.WriteByte('l')
		case TABLE_ALIGNMENT_RIGHT:
			out.WriteByte('r')
		default:
			out.WriteByte('c')
		}
	}
	out.WriteString("}\n")
	out.Write(header)
	out.WriteString(" \\\\\n\\hline\n")
	out.Write(body)
	out.WriteString("\n\\end{tabular}\n")
}

func (options *Xml) TableRow(out *bytes.Buffer, text []byte) {
	if out.Len() > 0 {
		out.WriteString(" \\\\\n")
	}
	out.Write(text)
}

func (options *Xml) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString(" & ")
	}
	out.Write(text)
}

func (options *Xml) TableCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString(" & ")
	}
	out.Write(text)
}

// TODO: this
func (options *Xml) Footnotes(out *bytes.Buffer, text func() bool) {

}

func (options *Xml) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {

}

func (options *Xml) Index(out *bytes.Buffer, primary, secondary []byte) {
	out.WriteString("<iref item=\"" + string(primary) + "\"")
	out.WriteString(" subitem=\"" + string(secondary) + "\"" + "/>")
}

func (options *Xml) Citation(out *bytes.Buffer, link, title []byte) {
	out.WriteString("<xref target=\"" + string(link) + "\"/>")
}

func (options *Xml) References(out *bytes.Buffer, citations map[string]*citation, first bool) {
	if !first {
		return
	}
	// count the references
	refi, refn := 0, 0
	for _, c := range citations {
		if c.typ == 'i' {
			refi++
		}
		if c.typ == 'n' {
			refn++
		}
	}
	if refi+refn > 0 {
		println("References")
		if refi > 0 {
			println("Informative References")
			for k, c := range citations {
				if c.typ == 'i' {
					println(k)
				}
			}
		}
		if refn > 0 {
			println("Normative References")
			for k, c := range citations {
				if c.typ == 'n' {
					println(k)
				}
			}
		}
	}
}

func (options *Xml) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.WriteString("\\href{")
	if kind == LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	out.Write(link)
	out.WriteString("}{")
	out.Write(link)
	out.WriteString("}")
}

func (options *Xml) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString("<tt>")
	out.Write(text)
	out.WriteString("</tt>")
}

func (options *Xml) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<b>")
	out.Write(text)
	out.WriteString("</b>")
}

func (options *Xml) Emphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<i>")
	out.Write(text)
	out.WriteString("</i>")
}

func (options *Xml) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	if bytes.HasPrefix(link, []byte("http://")) || bytes.HasPrefix(link, []byte("https://")) {
		// treat it like a link
		out.WriteString("\\href{")
		out.Write(link)
		out.WriteString("}{")
		out.Write(alt)
		out.WriteString("}")
	} else {
		out.WriteString("\\includegraphics{")
		out.Write(link)
		out.WriteString("}")
	}
}

func (options *Xml) LineBreak(out *bytes.Buffer) {
	out.WriteString(" \\\\\n")
}

func (options *Xml) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.WriteString("\\href{")
	out.Write(link)
	out.WriteString("}{")
	out.Write(content)
	out.WriteString("}")
}

func (options *Xml) RawHtmlTag(out *bytes.Buffer, tag []byte) {
}

func (options *Xml) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("\\textbf{\\textit{")
	out.Write(text)
	out.WriteString("}}")
}

func (options *Xml) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.WriteString("\\sout{")
	out.Write(text)
	out.WriteString("}")
}

// TODO: this
func (options *Xml) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {

}

func (options *Xml) Entity(out *bytes.Buffer, entity []byte) {
	// TODO: convert this into a unicode character or something
	out.Write(entity)
}

func (options *Xml) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

// header and footer
func (options *Xml) DocumentHeader(out *bytes.Buffer, first bool) {
	if !first {
		return
	}
	out.WriteString("\n<rfc>\n")
	out.WriteString("\n")
}

func (options *Xml) DocumentFooter(out *bytes.Buffer, first bool) {
	if !first {
		return
	}
	out.WriteString("\n</rfc>\n")
}
