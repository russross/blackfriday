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
	"strconv"
	"time"
)

// XML renderer configuration options.
const (
	XML_STANDALONE = 1 << iota // create standalone document
)

// Xml is a type that implements the Renderer interface for XML2RFV3 output.
//
// Do not create this directly, instead use the XmlRenderer function.
type Xml struct {
	flags        int // XML_* options
	sectionLevel int // current section level
	docLevel     int // frontmatter/mainmatter or backmatter

	// Store the IAL we see for this block element
	ial []*IAL

	// TitleBlock in TOML
	titleBlock *title
}

func (options *Xml) SetIAL(i []*IAL)        { options.ial = append(options.ial, i...) }
func (options *Xml) GetAndResetIAL() []*IAL { i := options.ial; options.ial = nil; return i }

// XmlRenderer creates and configures a Xml object, which
// satisfies the Renderer interface.
//
// flags is a set of XML_* options ORed together (currently no such options
// are defined).
func XmlRenderer(flags int) Renderer {
	return &Xml{flags: flags}
}

func (options *Xml) GetFlags() int {
	return options.flags
}

func (options *Xml) GetState() int {
	return 0
}

// render code chunks using verbatim, or listings if we have a language
func (options *Xml) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	s := renderIAL(options.GetAndResetIAL())
	if lang == "" {
		out.WriteString("<sourcecode" + s + ">\n")
	} else {
		out.WriteString("\n<sourcecode" + s + "type=\"" + lang + "\">\n")
	}
	out.Write(text)
	if lang == "" {
		out.WriteString("</sourcecode>\n")
	} else {
		out.WriteString("</sourcecode>\n")
	}
}

func (options *Xml) TitleBlock(out *bytes.Buffer, text []byte) {}

func (options *Xml) TitleBlockTOML(out *bytes.Buffer, block *title) {
	if options.flags&XML_STANDALONE == 0 {
		return
	}
	options.titleBlock = block
	out.WriteString("<rfc xmlns:xi=\"http://www.w3.org/2001/XInclude\" ipr=\"" +
		options.titleBlock.Ipr + "\" category=\"" +
		options.titleBlock.Category + "\" docName=\"" + options.titleBlock.DocName + "\">\n")
	out.WriteString("<front>\n")
	out.WriteString("<title abbrev=\"" + options.titleBlock.Abbrev + "\">")
	out.WriteString(options.titleBlock.Title + "</title>\n\n")

	year := ""
	if options.titleBlock.Date.Year() > 0 {
		year = " year=\"" + strconv.Itoa(options.titleBlock.Date.Year()) + "\""
	}
	month := ""
	if options.titleBlock.Date.Month() > 0 {
		month = " month=\"" + time.Month(options.titleBlock.Date.Month()).String() + "\""
	}
	day := ""
	if options.titleBlock.Date.Day() > 0 {
		day = " day=\"" + strconv.Itoa(options.titleBlock.Date.Day()) + "\""
	}
	out.WriteString("<date" + year + month + day + "/>\n\n")

	out.WriteString("<area>" + options.titleBlock.Area + "</area>\n")
	out.WriteString("<workgroup>" + options.titleBlock.Workgroup + "</workgroup>\n")
	for _, k := range options.titleBlock.Keyword {
		out.WriteString("<keyword>" + k + "</keyword>\n")
	}
	for _, a := range options.titleBlock.Author {
		out.WriteString("<author>\n")
		out.WriteString("<initials>" + a.Initials + "</initials>\n")
		out.WriteString("<surname>" + a.Surname + "</surname>\n")
		out.WriteString("<fullname>" + a.Fullname + "</fullname>\n")
		out.WriteString("<role>" + a.Role + "</role>\n")
		out.WriteString("<ascii>" + a.Ascii + "</ascii>\n")
		out.WriteString("</author>\n")
	}
	// Author information
	out.WriteString("\n")
}

func (options *Xml) BlockQuote(out *bytes.Buffer, text []byte) {
	s := renderIAL(options.GetAndResetIAL())
	out.WriteString("<blockquote" + s + ">\n")
	out.Write(text)
	out.WriteString("</blockquote>\n")
}

func (options *Xml) Abstract(out *bytes.Buffer, text []byte) {
	out.WriteString("<abstract>\n")
	out.Write(text)
	out.WriteString("</abstract>\n")
}

func (options *Xml) Aside(out *bytes.Buffer, text []byte) {
	out.WriteString("<aside>\n")
	out.Write(text)
	out.WriteString("</aside>\n")
}

func (options *Xml) Note(out *bytes.Buffer, text []byte) {
	out.WriteString("<note>\n")
	out.Write(text)
	out.WriteString("</note>\n")
}

func (options *Xml) BlockHtml(out *bytes.Buffer, text []byte) {
	// a pretty lame thing to do...
	out.WriteString("\n\\begin{verbatim}\n")
	out.Write(text)
	out.WriteString("\n\\end{verbatim}\n")
}

func (options *Xml) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	// set amount of open in options, so we know what to close after we finish
	// parsing the doc.
	//marker := out.Len()
	//out.Truncate(marker)

	id = "a"
	if level <= options.sectionLevel {
		// close previous ones
		for i := options.sectionLevel - level + 1; i > 0; i-- {
			out.WriteString("</section>\n")
		}
	}
	// new section
	out.WriteString("\n<section anchor=\"" + id + "\">\n")
	out.WriteString("<name>")
	text() // check bool here
	out.WriteString("</name>\n")
	options.sectionLevel = level
	return
}

func (options *Xml) HRule(out *bytes.Buffer) {
	// not used
}

func (options *Xml) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	switch {
	case flags&LIST_TYPE_ORDERED != 0:
		out.WriteString("<ol>\n")
	case flags&LIST_TYPE_DEFINITION != 0:
		out.WriteString("<dl>\n")
	default:
		out.WriteString("<ul>\n")
	}

	if !text() {
		out.Truncate(marker)
		return
	}
	switch {
	case flags&LIST_TYPE_ORDERED != 0:
		out.WriteString("</ol>\n")
	case flags&LIST_TYPE_DEFINITION != 0:
		out.WriteString("</dl>\n")
	default:
		out.WriteString("</ul>\n")
	}
}

func (options *Xml) ListItem(out *bytes.Buffer, text []byte, flags int) {
	if flags&LIST_TYPE_DEFINITION != 0 && flags&LIST_TYPE_TERM == 0 {
		out.WriteString("<dd>")
		out.Write(text)
		out.WriteString("</dd>\n")
		return
	}
	if flags&LIST_TYPE_TERM != 0 {
		out.WriteString("<dt>")
		out.Write(text)
		out.WriteString("</dt>\n")
		return
	}
	out.WriteString("<li>")
	out.Write(text)
	out.WriteString("</li>\n")
}

func (options *Xml) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	out.WriteString("<t>")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</t>\n")
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

func (options *Xml) Footnotes(out *bytes.Buffer, text func() bool) {
	// not used
}

func (options *Xml) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	// not used
}

func (options *Xml) Index(out *bytes.Buffer, primary, secondary []byte) {
	out.WriteString("<iref item=\"" + string(primary) + "\"")
	out.WriteString(" subitem=\"" + string(secondary) + "\"" + "/>")
}

func (options *Xml) Citation(out *bytes.Buffer, link, title []byte) {
	out.WriteString("<xref target=\"" + string(link) + "\"/>")
}

func (options *Xml) References(out *bytes.Buffer, citations map[string]*citation, first bool) {
	if !first || options.flags&XML_STANDALONE == 0  {
		return
	}
	// close any option section tags
	for i := options.sectionLevel; i > 0; i-- {
		out.WriteString("</section>\n")
		options.sectionLevel--
	}
	switch options.docLevel {
	case DOC_FRONT_MATTER:
		out.WriteString("</front>\n")
		out.WriteString("<back>\n")
	case DOC_MAIN_MATTER:
		out.WriteString("</middle>\n")
		out.WriteString("<back>\n")
	case DOC_BACK_MATTER:
		// nothing to do
	}
	options.docLevel = DOC_BACK_MATTER
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
	// output <xi:include href="<references file>.xml"/>, we use file it its not empty, otherwise
	// we construct one for RFCNNNN and I-D.something something.
	if refi+refn > 0 {
		if refi > 0 {
			out.WriteString("<references title=\"Informative References\">\n")
			for _, c := range citations {
				if c.typ == 'i' {
					f := string(c.filename)
					if f == "" {
						f = referenceFile(c.link)
					}
					out.WriteString("\t<xi:include href=\"" + f + "\"/>\n")
				}
			}
			out.WriteString("</references>\n")
		}
		if refn > 0 {
			out.WriteString("<references title=\"Normative References\">\n")
			for _, c := range citations {
				if c.typ == 'n' {
					f := string(c.filename)
					if f == "" {
						f = referenceFile(c.link)
					}
					out.WriteString("\t<xi:include href=\"" + f + "\"/>\n")
				}
			}
			out.WriteString("</references>\n")
		}
	}
}

// create reference file
func referenceFile(id []byte) string {
	if len(id) < 4 {
		return ""
	}
	switch string(id[:3]) {
	case "RFC":
		return "reference.RFC." + string(id[:3]) + ".xml"
	case "I-D":
		return "reference.I-D.draft-" + string(id[4:]) + ".xml"
	}
	return ""
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
	convertEntity(out, text)
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
	out.WriteString("\n<vspace/>\n")
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
	out.Write(text)
}

func (options *Xml) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	// not used
}

func (options *Xml) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (options *Xml) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

// header and footer
func (options *Xml) DocumentHeader(out *bytes.Buffer, first bool) {
	if !first || options.flags&XML_STANDALONE == 0  {
		return
	}
	out.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
}

func (options *Xml) DocumentFooter(out *bytes.Buffer, first bool) {
	if !first || options.flags&XML_STANDALONE == 0  {
		return
	}
	// close any option section tags
	for i := options.sectionLevel; i > 0; i-- {
		out.WriteString("</section>\n")
		options.sectionLevel--
	}
	switch options.docLevel {
	case DOC_FRONT_MATTER:
		out.WriteString("</front>\n")
	case DOC_MAIN_MATTER:
		out.WriteString("</middle>\n")
	case DOC_BACK_MATTER:
		out.WriteString("</back>\n")
	}
	out.WriteString("</rfc>\n")
}

func (options *Xml) DocumentMatter(out *bytes.Buffer, matter int) {
	// we default to frontmatter already openened in the documentHeader
	switch matter {
	case DOC_FRONT_MATTER:
		// already open
	case DOC_MAIN_MATTER:
		out.WriteString("</front>\n")
		out.WriteString("<middle>\n")
	case DOC_BACK_MATTER:
		out.WriteString("</middle>\n")
		out.WriteString("<back>\n")
	}
	options.docLevel = matter
}

// TODO(miek): extend
var entityConvert = map[byte]string{
	'<': "&lt;",
	'>': "&gt;",
}

func convertEntity(out *bytes.Buffer, text []byte) {
	for i := 0; i < len(text); i++ {
		if s, ok := entityConvert[text[i]]; ok {
			out.WriteString(s)
			continue
		}
		out.WriteByte(text[i])
	}
}
