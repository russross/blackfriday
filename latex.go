//
// Black Friday Markdown Processor
// Originally based on http://github.com/tanoku/upskirt
// by Russ Ross <russ@russross.com>
//

//
//
// LaTeX rendering backend
//
//

package blackfriday

import (
	"bytes"
)

func LatexRenderer(flags int) *Renderer {
	// block-level rendering
	r := new(Renderer)
	r.BlockCode = latexBlockCode
	r.BlockQuote = latexBlockQuote
	//r.BlockHtml = ?
	r.Header = latexHeader
	r.HRule = latexHRule
	r.List = latexList
	r.ListItem = latexListItem
	r.Paragraph = latexParagraph
	r.Table = latexTable
	r.TableRow = latexTableRow
	r.TableCell = latexTableCell

	// inline rendering
	r.AutoLink = latexAutoLink
	r.CodeSpan = latexCodeSpan
	r.DoubleEmphasis = latexDoubleEmphasis
	r.Emphasis = latexEmphasis
	r.Image = latexImage
	r.LineBreak = latexLineBreak
	r.Link = latexLink
	//r.rawHtmlTag = ?
	r.StrikeThrough = latexStrikeThrough

	r.NormalText = latexNormalText

	r.DocumentHeader = latexDocumentHeader
	r.DocumentFooter = latexDocumentFooter

	r.Opaque = nil
	return r
}

// render code chunks using verbatim, or listings if we have a language
func latexBlockCode(out *bytes.Buffer, text []byte, lang string, opaque interface{}) {
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

func latexBlockQuote(out *bytes.Buffer, text []byte, opaque interface{}) {
	out.WriteString("\n\\begin{quotation}\n")
	out.Write(text)
	out.WriteString("\n\\end{quotation}\n")
}

//BlockHtml  func(out *bytes.Buffer, text []byte, opaque interface{})

func latexHeader(out *bytes.Buffer, text func() bool, level int, opaque interface{}) {
	marker := out.Len()

	switch level {
	case 1:
		out.WriteString("\n\\section{")
	case 2:
		out.WriteString("\n\\subsection{")
	case 3:
		out.WriteString("\n\\subsubsection{")
	case 4:
		out.WriteString("\n\\paragraph{")
	case 5:
		out.WriteString("\n\\subparagraph{")
	case 6:
		out.WriteString("\n\\textbf{")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("}\n")
}

func latexHRule(out *bytes.Buffer, opaque interface{}) {
	out.WriteString("\n\\HRule\n")
}

func latexList(out *bytes.Buffer, text func() bool, flags int, opaque interface{}) {
	marker := out.Len()
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString("\n\\begin{enumerate}\n")
	} else {
		out.WriteString("\n\\begin{itemize}\n")
	}
	if !text() {
		out.Truncate(marker)
		return
	}
	if flags&LIST_TYPE_ORDERED != 0 {
		out.WriteString("\n\\end{enumerate}\n")
	} else {
		out.WriteString("\n\\end{itemize}\n")
	}
}

func latexListItem(out *bytes.Buffer, text []byte, flags int, opaque interface{}) {
	out.WriteString("\n\\item ")
	out.Write(text)
}

func latexParagraph(out *bytes.Buffer, text func() bool, opaque interface{}) {
	marker := out.Len()
	out.WriteString("\n")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("\n")
}

func latexTable(out *bytes.Buffer, header []byte, body []byte, columnData []int, opaque interface{}) {
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

func latexTableRow(out *bytes.Buffer, text []byte, opaque interface{}) {
	if out.Len() > 0 {
		out.WriteString(" \\\\\n")
	}
	out.Write(text)
}

func latexTableCell(out *bytes.Buffer, text []byte, align int, opaque interface{}) {
	if out.Len() > 0 {
		out.WriteString(" & ")
	}
	out.Write(text)
}

func latexAutoLink(out *bytes.Buffer, link []byte, kind int, opaque interface{}) int {
	out.WriteString("\\href{")
	if kind == LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	out.Write(link)
	out.WriteString("}{")
	out.Write(link)
	out.WriteString("}")
	return 1
}

func latexCodeSpan(out *bytes.Buffer, text []byte, opaque interface{}) int {
	out.WriteString("\\texttt{")
	escapeSpecialChars(out, text)
	out.WriteString("}")
	return 1
}

func latexDoubleEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) int {
	out.WriteString("\\textbf{")
	out.Write(text)
	out.WriteString("}")
	return 1
}

func latexEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) int {
	out.WriteString("\\textit{")
	out.Write(text)
	out.WriteString("}")
	return 1
}

func latexImage(out *bytes.Buffer, link []byte, title []byte, alt []byte, opaque interface{}) int {
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
	return 1
}

func latexLineBreak(out *bytes.Buffer, opaque interface{}) int {
	out.WriteString(" \\\\\n")
	return 1
}

func latexLink(out *bytes.Buffer, link []byte, title []byte, content []byte, opaque interface{}) int {
	out.WriteString("\\href{")
	out.Write(link)
	out.WriteString("}{")
	out.Write(content)
	out.WriteString("}")
	return 1
}

func latexRawHtmlTag(out *bytes.Buffer, tag []byte, opaque interface{}) int {
	return 0
}

func latexTripleEmphasis(out *bytes.Buffer, text []byte, opaque interface{}) int {
	out.WriteString("\\textbf{\\textit{")
	out.Write(text)
	out.WriteString("}}")
	return 1
}

func latexStrikeThrough(out *bytes.Buffer, text []byte, opaque interface{}) int {
	out.WriteString("\\sout{")
	out.Write(text)
	out.WriteString("}")
	return 1
}

func needsBackslash(c byte) bool {
	for _, r := range []byte("_{}%$&\\~") {
		if c == r {
			return true
		}
	}
	return false
}

func escapeSpecialChars(out *bytes.Buffer, text []byte) {
	for i := 0; i < len(text); i++ {
		// directly copy normal characters
		org := i

		for i < len(text) && !needsBackslash(text[i]) {
			i++
		}
		if i > org {
			out.Write(text[org:i])
		}

		// escape a character
		if i >= len(text) {
			break
		}
		out.WriteByte('\\')
		out.WriteByte(text[i])
	}
}

func latexNormalText(out *bytes.Buffer, text []byte, opaque interface{}) {
	escapeSpecialChars(out, text)
}

// header and footer
func latexDocumentHeader(out *bytes.Buffer, opaque interface{}) {
	out.WriteString("\\documentclass{article}\n")
	out.WriteString("\n")
	out.WriteString("\\usepackage{graphicx}\n")
	out.WriteString("\\usepackage{listings}\n")
	out.WriteString("\\usepackage[margin=1in]{geometry}\n")
	out.WriteString("\\usepackage[utf8]{inputenc}\n")
	out.WriteString("\\usepackage{verbatim}\n")
	out.WriteString("\\usepackage[normalem]{ulem}\n")
	out.WriteString("\\usepackage{hyperref}\n")
	out.WriteString("\n")
	out.WriteString("\\hypersetup{colorlinks,%\n")
	out.WriteString("  citecolor=black,%\n")
	out.WriteString("  filecolor=black,%\n")
	out.WriteString("  linkcolor=black,%\n")
	out.WriteString("  urlcolor=black,%\n")
	out.WriteString("  pdfstartview=FitH,%\n")
	out.WriteString("  breaklinks=true,%\n")
	out.WriteString("  pdfauthor={Black Friday Markdown Processor}}\n")
	out.WriteString("\n")
	out.WriteString("\\newcommand{\\HRule}{\\rule{\\linewidth}{0.5mm}}\n")
	out.WriteString("\\addtolength{\\parskip}{0.5\\baselineskip}\n")
	out.WriteString("\\parindent=0pt\n")
	out.WriteString("\n")
	out.WriteString("\\begin{document}\n")
}

func latexDocumentFooter(out *bytes.Buffer, opaque interface{}) {
	out.WriteString("\n\\end{document}\n")
}
