// XML2RFC v3 rendering backend

package mmark

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

// XML renderer configuration options.
const (
	XML_STANDALONE = 1 << iota // create standalone document
)

var words2119 = map[string]bool{
	"MUST":        true,
	"MUST NOT":    true,
	"REQUIRED":    true,
	"SHALL":       true,
	"SHALL NOT":   true,
	"SHOULD":      true,
	"SHOULD NOT":  true,
	"RECOMMENDED": true,
	"MAY":         true,
	"OPTIONAL":    true,
}

// Xml is a type that implements the Renderer interface for XML2RFV3 output.
//
// Do not create this directly, instead use the XmlRenderer function.
type Xml struct {
	flags        int // XML_* options
	sectionLevel int // current section level
	docLevel     int // frontmatter/mainmatter or backmatter

	// Store the IAL we see for this block element
	ial *IAL

	// TitleBlock in TOML
	titleBlock *title
}

// XmlRenderer creates and configures a Xml object, which
// satisfies the Renderer interface.
//
// flags is a set of XML_* options ORed together
func XmlRenderer(flags int) Renderer { return &Xml{flags: flags} }
func (options *Xml) GetFlags() int   { return options.flags }
func (options *Xml) GetState() int   { return 0 }
func (options *Xml) SetIAL(i *IAL)   { options.ial = i }
func (options *Xml) IAL() *IAL       { i := options.ial; options.ial = nil; return i }

// render code chunks using verbatim, or listings if we have a language
func (options *Xml) BlockCode(out *bytes.Buffer, text []byte, lang string, caption []byte) {
	// Tick of language for sourcecode...
	s := options.IAL().String()
	if len(caption) > 0 {
		out.WriteString("<figure" + s + ">\n")
		s = ""
		out.WriteString("<name>")
		out.Write(caption)
		out.WriteString("</name>\n")
	}

	if lang == "" {
		out.WriteString("<artwork" + s + ">\n")
	} else {
		out.WriteString("\n<sourcecode" + s + "type=\"" + lang + "\">\n")
	}
	out.Write(text)
	if lang == "" {
		out.WriteString("</artwork>\n")
	} else {
		out.WriteString("</sourcode>\n")
	}
	if len(caption) > 0 {
		out.WriteString("</figure>\n")
	}
}

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
		out.WriteString("<author")
		out.WriteString(" initials=\"" + a.Initials + "\"")
		out.WriteString(" surname=\"" + a.Surname + "\"")
		out.WriteString(" fullname=\"" + a.Fullname + "\">")

		out.WriteString("<organization>" + a.Organization + "</organization>\n")
		out.WriteString("<address>\n")
		out.WriteString("<email>" + a.Address.Email + "</email>\n")
		out.WriteString("</address>\n")
		out.WriteString("<role>" + a.Role + "</role>\n")
		out.WriteString("<ascii>" + a.Ascii + "</ascii>\n")
		out.WriteString("</author>\n")
	}
	out.WriteString("\n")
}

func (options *Xml) BlockQuote(out *bytes.Buffer, text []byte) {
	s := options.IAL().String()
	out.WriteString("<blockquote" + s + ">\n")
	out.Write(text)
	out.WriteString("</blockquote>\n")
}

func (options *Xml) Abstract(out *bytes.Buffer, text []byte) {
	s := options.IAL().String()
	out.WriteString("<abstract" + s + ">\n")
	out.Write(text)
	out.WriteString("</abstract>\n")
}

func (options *Xml) Aside(out *bytes.Buffer, text []byte) {
	s := options.IAL().String()
	out.WriteString("<aside" + s + ">\n")
	out.Write(text)
	out.WriteString("</aside>\n")
}

func (options *Xml) Note(out *bytes.Buffer, text []byte) {
	s := options.IAL().String()
	out.WriteString("<note" + s + ">\n")
	out.Write(text)
	out.WriteString("</note>\n")
}

func (options *Xml) CommentHtml(out *bytes.Buffer, text []byte) {
	return
}

func (options *Xml) BlockHtml(out *bytes.Buffer, text []byte) {
	// not supported, don't know yet if this is useful
	return
}

func (options *Xml) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	// set amount of open in options, so we know what to close after we finish
	// parsing the doc.
	//marker := out.Len()
	//out.Truncate(marker)
	if level <= options.sectionLevel {
		// close previous ones
		for i := options.sectionLevel - level + 1; i > 0; i-- {
			out.WriteString("</section>\n")
		}
	}

	ial := options.ial
	if ial != nil {
		id = ial.GetOrDefaultId(id)
	}
	if id != "" {
		id = " anchor=\"" + id + "\""
	}

	// new section
	out.WriteString("\n<section" + id + "\"" + ial.String() + ">")
	out.WriteString("<name>")
	text() // check bool here
	out.WriteString("</name>\n")
	options.sectionLevel = level
	return
}

func (options *Xml) HRule(out *bytes.Buffer) {
	// not used
}

func (options *Xml) List(out *bytes.Buffer, text func() bool, flags, start int) {
	marker := out.Len()
	s := options.IAL().String()
	switch {
	case flags&LIST_TYPE_ORDERED != 0:
		if start <= 1 {
			out.WriteString("<ol" + s + ">\n")
		} else {
			out.WriteString(fmt.Sprintf("<ol"+s+" start=\"%d\">\n", start))
		}
	case flags&LIST_TYPE_DEFINITION != 0:
		out.WriteString("<dl" + s + ">\n")
	default:
		out.WriteString("<ul" + s + ">\n")
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

func (options *Xml) Paragraph(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	out.WriteString("<t>")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</t>\n")
}

func (options *Xml) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int, caption []byte) {
	s := options.IAL().String()
	out.WriteString("<table" + s + ">\n")
	if caption != nil {
		out.WriteString("<name>")
		out.Write(caption)
		out.WriteString("</name>\n")
	}
	out.WriteString("<thead>\n")
	out.Write(header)
	out.WriteString("</thead>\n")
	out.Write(body)
	out.WriteString("</table>\n")
}

func (options *Xml) TableRow(out *bytes.Buffer, text []byte) {
	out.WriteString("<tr>")
	out.Write(text)
	out.WriteString("</tr>\n")
}

func (options *Xml) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	a := ""
	switch align {
	case TABLE_ALIGNMENT_LEFT:
		a = " align=\"left\""
	case TABLE_ALIGNMENT_RIGHT:
		a = " align=\"right\""
	default:
		a = " align=\"center\""
	}
	out.WriteString("<th" + a + ">")
	out.Write(text)
	out.WriteString("</th>")

}

func (options *Xml) TableCell(out *bytes.Buffer, text []byte, align int) {
	out.WriteString("<td>")
	out.Write(text)
	out.WriteString("</td>")
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
	if len(title) == 0 {
		out.WriteString("<xref target=\"" + string(link) + "\"/>")
		return
	}
	out.WriteString("<xref target=\"" + string(link) + "\" section=\"" + string(title) + "\"/>")
}

func (options *Xml) References(out *bytes.Buffer, citations map[string]*citation) {
	if options.flags&XML_STANDALONE == 0 {
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
						f = referenceFile(c)
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
						f = referenceFile(c)
					}
					out.WriteString("\t<xi:include href=\"" + f + "\"/>\n")
				}
			}
			out.WriteString("</references>\n")
		}
	}
}

// create reference file
func referenceFile(c *citation) string {
	if len(c.link) < 4 {
		return ""
	}
	switch string(c.link[:3]) {
	case "RFC":
		return "reference.RFC." + string(c.link[3:]) + ".xml"
	case "I-D":
		return "reference.I-D.draft-" + string(c.link[4:]) + ".xml"
	}
	return ""
}

func (options *Xml) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.WriteString("<eref target=\"")
	if kind == LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	out.Write(link)
	out.WriteString("\"/>")
}

func (options *Xml) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString("<tt>")
	convertEntity(out, text)
	out.WriteString("</tt>")
}

func (options *Xml) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	// Check for 2119 Keywords
	s := string(text)
	if _, ok := words2119[s]; ok {
		out.WriteString("<bcp14>")
		out.Write(text)
		out.WriteString("</bcp14>")
		return
	}
	out.WriteString("<strong>")
	out.Write(text)
	out.WriteString("</strong>")
}

func (options *Xml) Emphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<em>")
	out.Write(text)
	out.WriteString("</em>")
}

func (options *Xml) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	// use title as caption is we have it
	// check the extension of the local include to set the type of the thing.
	s := options.IAL().String()
	if bytes.HasPrefix(link, []byte("http://")) || bytes.HasPrefix(link, []byte("https://")) {
		// link to external entity
		out.WriteString("<artwork" + s)
		out.WriteString(" alt=\"")
		out.Write(alt)
		out.WriteString("\"")
		out.WriteString(" src=\"")
		out.Write(link)
		out.WriteString("\"/>")
	} else {
		// local file, xi:include it
		out.WriteString("<artwork" + s)
		out.WriteString(" alt=\"")
		out.Write(alt)
		out.WriteString("\">")
		out.WriteString("<xi:include href=\"")
		out.Write(link)
		out.WriteString("\"/>\n")
		out.WriteString("</artwork>\n")
	}
}

func (options *Xml) LineBreak(out *bytes.Buffer) {
	out.WriteString("\n<br/>\n")
}

func (options *Xml) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	// kill # in the link name
	if link[0] == '#' {
		link = link[1:]
	}
	out.WriteString("<xref target=\"")
	out.Write(link)
	out.WriteString("\"/>")
	//	out.Write(content)
}

func (options *Xml) RawHtmlTag(out *bytes.Buffer, tag []byte) {
}

func (options *Xml) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<strong><em>")
	out.Write(text)
	out.WriteString("</em></strong>")
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
	if !first || options.flags&XML_STANDALONE == 0 {
		return
	}
	out.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
}

func (options *Xml) DocumentFooter(out *bytes.Buffer, first bool) {
	if !first || options.flags&XML_STANDALONE == 0 {
		return
	}
	// close any option section tags
	for i := options.sectionLevel; i > 0; i-- {
		out.WriteString("</section>\n")
		options.sectionLevel--
	}
	switch options.docLevel {
	case DOC_FRONT_MATTER:
		out.WriteString("\n</front>\n")
	case DOC_MAIN_MATTER:
		out.WriteString("\n</middle>\n")
	case DOC_BACK_MATTER:
		out.WriteString("\n</back>\n")
	}
	out.WriteString("</rfc>\n")
}

func (options *Xml) DocumentMatter(out *bytes.Buffer, matter int) {
	// we default to frontmatter already openened in the documentHeader
	for i := options.sectionLevel; i > 0; i-- {
		out.WriteString("</section>\n")
		options.sectionLevel--
	}
	switch matter {
	case DOC_FRONT_MATTER:
		// already open
	case DOC_MAIN_MATTER:
		out.WriteString("</front>\n")
		out.WriteString("\n<middle>\n")
	case DOC_BACK_MATTER:
		out.WriteString("\n</middle>\n")
		out.WriteString("<back>\n")
	}
	options.docLevel = matter
}

// quotes &quot;
var entityConvert = map[byte]string{
	'<': "&lt;",
	'>': "&gt;",
	'&': "&amp;",
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
