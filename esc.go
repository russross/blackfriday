package blackfriday

import (
	"html"
	"io"
)

type escMap struct {
	char byte
	seq  []byte
}

var htmlEscaper = []escMap{
	{'&', []byte("&amp;")},
	{'<', []byte("&lt;")},
	{'>', []byte("&gt;")},
	{'"', []byte("&quot;")},
}

func escapeHTML(w io.Writer, s []byte) {
	var start, end int
	var sEnd byte
	for end < len(s) {
		sEnd = s[end]
		if sEnd == '&' || sEnd == '<' || sEnd == '>' || sEnd == '"' {
			for i := 0; i < len(htmlEscaper); i++ {
				if sEnd == htmlEscaper[i].char {
					w.Write(s[start:end])
					w.Write(htmlEscaper[i].seq)
					start = end + 1
					break
				}
			}
		}
		end++
	}
	if start < len(s) && end <= len(s) {
		w.Write(s[start:end])
	}
}

func escLink(w io.Writer, text []byte) {
	unesc := html.UnescapeString(string(text))
	escapeHTML(w, []byte(unesc))
}
