package mmark

// Pretty print XMl markup

import (
	"bytes"
	"encoding/xml"
	"io"
)

func prettyPass(p *parser, input []byte) ([]byte, error) {
	var b bytes.Buffer
	i := newIndentWriter(&b)
	decoder := xml.NewDecoder(bytes.NewReader(input))
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch x := token.(type) {
		case xml.StartElement:
			i.WriteString("<")
			i.WriteString(x.Name.Local)
			for _, a := range x.Attr {
				i.WriteString(" ")
				i.WriteString(a.Name.Local)
				i.WriteString("=\"")
				i.WriteString(a.Value)
				i.WriteString("\"")
			}
			i.WriteString(">")
		case xml.EndElement:
			i.WriteString("</")
			i.WriteString(x.Name.Local)
			i.WriteString(">")
		case xml.CharData:
			i.Write(x)
		case xml.Attr:
			i.WriteString(x.Name.Local)
			i.WriteString(x.Value)
		case xml.ProcInst:
			i.WriteString("<?")
			i.WriteString(x.Target)
			i.Write(x.Inst)
			i.WriteString("?>")
		case xml.Directive:
			i.WriteString("<!")
			i.Write(x)
			i.WriteString(">")
		case xml.Comment:
			i.WriteString("<!-- ")
			i.Write(x)
			i.WriteString(" -->")
		}
	}
	return b.Bytes(), nil
}

type indentWriter struct {
	b      *bytes.Buffer
	indent []byte // current indentation
	last   byte   // last written byte
}

func newIndentWriter(b *bytes.Buffer) *indentWriter {
	return &indentWriter{b, nil, 0}
}

func (i *indentWriter) WriteString(s string) {
	if len(s) == 0 {
		return
	}
	defer func() {
		i.last = s[len(s)-1]
	}()
	if s[0] == '<' {
		if len(s) > 1 && s[1] == '/' {
			i.indent = i.indent[:len(i.indent)-2]
			if i.last == '\n' {
				i.b.Write(i.indent)
			}
			i.b.WriteString(s)
			return
		}
		if i.last == '\n' {
			i.b.Write(i.indent)
		}
		i.b.WriteString(s)
		// only for real open tags we increment the indentation
		if !(len(s) > 1 && (s[1] == '!' || s[1] == '?')) {
			i.indent = append(i.indent, []byte("  ")...)

		}
		return
	}
	i.b.WriteString(s)
}

func (i *indentWriter) Write(b []byte) {
	if len(b) == 0 {
		return
	}
	i.b.Write(b)
	i.last = b[len(b)-1]
}
