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
// LaTeX rendering backend
//
//

package blackfriday

import "bytes"

// Latex is a type that implements the Renderer interface for LaTeX output.
//
// Do not create this directly, instead use the LatexRenderer function.
type Latex struct {
	w HtmlWriter
}

// LatexRenderer creates and configures a Latex object, which
// satisfies the Renderer interface.
//
// flags is a set of LATEX_* options ORed together (currently no such options
// are defined).
func LatexRenderer(flags int) Renderer {
	var writer HtmlWriter
	return &Latex{
		w: writer,
	}
}

func (r *Latex) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

// render code chunks using verbatim, or listings if we have a language
func (r *Latex) BlockCode(text []byte, lang string) {
	if lang == "" {
		r.w.WriteString("\n\\begin{verbatim}\n")
	} else {
		r.w.WriteString("\n\\begin{lstlisting}[language=")
		r.w.WriteString(lang)
		r.w.WriteString("]\n")
	}
	r.w.Write(text)
	if lang == "" {
		r.w.WriteString("\n\\end{verbatim}\n")
	} else {
		r.w.WriteString("\n\\end{lstlisting}\n")
	}
}

func (r *Latex) TitleBlock(text []byte) {

}

func (r *Latex) BlockQuote(text []byte) {
	r.w.WriteString("\n\\begin{quotation}\n")
	r.w.Write(text)
	r.w.WriteString("\n\\end{quotation}\n")
}

func (r *Latex) BlockHtml(text []byte) {
	// a pretty lame thing to do...
	r.w.WriteString("\n\\begin{verbatim}\n")
	r.w.Write(text)
	r.w.WriteString("\n\\end{verbatim}\n")
}

func (r *Latex) BeginHeader(level int, id string) {
	switch level {
	case 1:
		r.w.WriteString("\n\\section{")
	case 2:
		r.w.WriteString("\n\\subsection{")
	case 3:
		r.w.WriteString("\n\\subsubsection{")
	case 4:
		r.w.WriteString("\n\\paragraph{")
	case 5:
		r.w.WriteString("\n\\subparagraph{")
	case 6:
		r.w.WriteString("\n\\textbf{")
	}
}

func (r *Latex) EndHeader(level int, id string, header []byte) {
	r.w.WriteString("}\n")
}

func (r *Latex) HRule() {
	r.w.WriteString("\n\\HRule\n")
}

func (r *Latex) BeginList(flags ListType) {
	if flags&ListTypeOrdered != 0 {
		r.w.WriteString("\n\\begin{enumerate}\n")
	} else {
		r.w.WriteString("\n\\begin{itemize}\n")
	}
}

func (r *Latex) EndList(flags ListType) {
	if flags&ListTypeOrdered != 0 {
		r.w.WriteString("\n\\end{enumerate}\n")
	} else {
		r.w.WriteString("\n\\end{itemize}\n")
	}
}

func (r *Latex) ListItem(text []byte, flags ListType) {
	r.w.WriteString("\n\\item ")
	r.w.Write(text)
}

func (r *Latex) BeginParagraph() {
	r.w.WriteString("\n")
}

func (r *Latex) EndParagraph() {
	r.w.WriteString("\n")
}

func (r *Latex) Table(header []byte, body []byte, columnData []int) {
	r.w.WriteString("\n\\begin{tabular}{")
	for _, elt := range columnData {
		switch elt {
		case TableAlignmentLeft:
			r.w.WriteByte('l')
		case TableAlignmentRight:
			r.w.WriteByte('r')
		default:
			r.w.WriteByte('c')
		}
	}
	r.w.WriteString("}\n")
	r.w.Write(header)
	r.w.WriteString(" \\\\\n\\hline\n")
	r.w.Write(body)
	r.w.WriteString("\n\\end{tabular}\n")
}

func (r *Latex) TableRow(text []byte) {
	r.w.WriteString(" \\\\\n")
	r.w.Write(text)
}

func (r *Latex) TableHeaderCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString(" & ")
	}
	out.Write(text)
}

func (r *Latex) TableCell(out *bytes.Buffer, text []byte, align int) {
	if out.Len() > 0 {
		out.WriteString(" & ")
	}
	out.Write(text)
}

// TODO: this
func (r *Latex) BeginFootnotes() {
}

// TODO: this
func (r *Latex) EndFootnotes() {
}

func (r *Latex) FootnoteItem(name, text []byte, flags ListType) {

}

func (r *Latex) AutoLink(link []byte, kind LinkType) {
	r.w.WriteString("\\href{")
	if kind == LinkTypeEmail {
		r.w.WriteString("mailto:")
	}
	r.w.Write(link)
	r.w.WriteString("}{")
	r.w.Write(link)
	r.w.WriteString("}")
}

func (r *Latex) CodeSpan(text []byte) {
	r.w.WriteString("\\texttt{")
	r.escapeSpecialChars(text)
	r.w.WriteString("}")
}

func (r *Latex) DoubleEmphasis(text []byte) {
	r.w.WriteString("\\textbf{")
	r.w.Write(text)
	r.w.WriteString("}")
}

func (r *Latex) Emphasis(text []byte) {
	r.w.WriteString("\\textit{")
	r.w.Write(text)
	r.w.WriteString("}")
}

func (r *Latex) Image(link []byte, title []byte, alt []byte) {
	if bytes.HasPrefix(link, []byte("http://")) || bytes.HasPrefix(link, []byte("https://")) {
		// treat it like a link
		r.w.WriteString("\\href{")
		r.w.Write(link)
		r.w.WriteString("}{")
		r.w.Write(alt)
		r.w.WriteString("}")
	} else {
		r.w.WriteString("\\includegraphics{")
		r.w.Write(link)
		r.w.WriteString("}")
	}
}

func (r *Latex) LineBreak() {
	r.w.WriteString(" \\\\\n")
}

func (r *Latex) Link(link []byte, title []byte, content []byte) {
	r.w.WriteString("\\href{")
	r.w.Write(link)
	r.w.WriteString("}{")
	r.w.Write(content)
	r.w.WriteString("}")
}

func (r *Latex) RawHtmlTag(tag []byte) {
}

func (r *Latex) TripleEmphasis(text []byte) {
	r.w.WriteString("\\textbf{\\textit{")
	r.w.Write(text)
	r.w.WriteString("}}")
}

func (r *Latex) StrikeThrough(text []byte) {
	r.w.WriteString("\\sout{")
	r.w.Write(text)
	r.w.WriteString("}")
}

// TODO: this
func (r *Latex) FootnoteRef(ref []byte, id int) {
}

func needsBackslash(c byte) bool {
	for _, r := range []byte("_{}%$&\\~#") {
		if c == r {
			return true
		}
	}
	return false
}

func (r *Latex) escapeSpecialChars(text []byte) {
	for i := 0; i < len(text); i++ {
		// directly copy normal characters
		org := i

		for i < len(text) && !needsBackslash(text[i]) {
			i++
		}
		if i > org {
			r.w.Write(text[org:i])
		}

		// escape a character
		if i >= len(text) {
			break
		}
		r.w.WriteByte('\\')
		r.w.WriteByte(text[i])
	}
}

func (r *Latex) Entity(entity []byte) {
	// TODO: convert this into a unicode character or something
	r.w.Write(entity)
}

func (r *Latex) NormalText(text []byte) {
	r.escapeSpecialChars(text)
}

// header and footer
func (r *Latex) DocumentHeader() {
	r.w.WriteString("\\documentclass{article}\n")
	r.w.WriteString("\n")
	r.w.WriteString("\\usepackage{graphicx}\n")
	r.w.WriteString("\\usepackage{listings}\n")
	r.w.WriteString("\\usepackage[margin=1in]{geometry}\n")
	r.w.WriteString("\\usepackage[utf8]{inputenc}\n")
	r.w.WriteString("\\usepackage{verbatim}\n")
	r.w.WriteString("\\usepackage[normalem]{ulem}\n")
	r.w.WriteString("\\usepackage{hyperref}\n")
	r.w.WriteString("\n")
	r.w.WriteString("\\hypersetup{colorlinks,%\n")
	r.w.WriteString("  citecolor=black,%\n")
	r.w.WriteString("  filecolor=black,%\n")
	r.w.WriteString("  linkcolor=black,%\n")
	r.w.WriteString("  urlcolor=black,%\n")
	r.w.WriteString("  pdfstartview=FitH,%\n")
	r.w.WriteString("  breaklinks=true,%\n")
	r.w.WriteString("  pdfauthor={Blackfriday Markdown Processor v")
	r.w.WriteString(VERSION)
	r.w.WriteString("}}\n")
	r.w.WriteString("\n")
	r.w.WriteString("\\newcommand{\\HRule}{\\rule{\\linewidth}{0.5mm}}\n")
	r.w.WriteString("\\addtolength{\\parskip}{0.5\\baselineskip}\n")
	r.w.WriteString("\\parindent=0pt\n")
	r.w.WriteString("\n")
	r.w.WriteString("\\begin{document}\n")
}

func (r *Latex) DocumentFooter() {
	r.w.WriteString("\n\\end{document}\n")
}

func (r *Latex) Render(ast *Node) []byte {
	return nil
}
