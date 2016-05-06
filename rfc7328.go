// Functions and definition to be backwards compatible with RFC7328 markdown.

package mmark

import "bytes"

func (p *parser) rfc7328Index(out *bytes.Buffer, text []byte) int {
	if p.flags&EXTENSION_RFC7328 == 0 {
		return 0
	}
	text = bytes.TrimSpace(text)
	// look for ^item1^ subitem
	if text[0] != '^' || len(text) < 3 {
		return 0
	}

	itemEnd := 0
	for i := 1; i < len(text); i++ {
		if text[i] == '^' {
			itemEnd = i
			break
		}
	}
	if itemEnd == 0 {
		return 0
	}

	// Check the sub item, if there.
	// skip whitespace
	i := itemEnd + 1
	for i < len(text) && isspace(text[i]) {
		i++
	}

	// rewind
	outSize := out.Len()
	outBytes := out.Bytes()
	if outSize > 0 && outBytes[outSize-1] == '^' {
		out.Truncate(outSize - 1)
	}

	subItemStart := i
	if subItemStart != len(text) {
		printf(p, "rfc 7328 style index parsed to: ((%s, %s))", string(text[1:itemEnd]), text[subItemStart:])
		p.r.Index(out, text[1:itemEnd], text[subItemStart:], false)
		return len(text)
	}
	printf(p, "rfc 7328 style index parsed to: ((%s))", string(text[1:itemEnd]))
	p.r.Index(out, text[1:itemEnd], nil, false)
	return len(text)
}

func (p *parser) rfc7328Caption(out *bytes.Buffer, text []byte) int {
	if p.flags&EXTENSION_RFC7328 == 0 {
		return 0
	}
	// Parse stuff like:
	// ^[fig:minimal::A minimal template.xml.]
	// If we don't find double colon it is not a inline note masking as a caption
	text = bytes.TrimSpace(text)
	colons := bytes.Index(text, []byte("::"))
	if colons == -1 {
		return 0
	}
	caption := []byte{}
	anchor := text[:colons]
	if colons+2 < len(text) {
		caption = text[colons+2:]
	}
	if len(anchor) == 0 && len(caption) == 0 {
		return 0
	}

	// rewind
	outSize := out.Len()
	outBytes := out.Bytes()
	if outSize > 0 && outBytes[outSize-1] == '^' {
		out.Truncate(outSize - 1)
	}
	// It is somewhat hard to now go back to the original start of the figure
	// and marge this new content in (there already may be a #id, etc. etc.).
	// For now just log that we have seen this line and return a positive integer
	// indicating this wasn't a footnote.
	if len(anchor) > 0 {
		printf(p, "rfc 7328 style anchor seen: consider adding '{#%s}' IAL before the figure/table", string(anchor))
	}
	if len(caption) > 0 {
		printf(p, "rfc 7328 style caption seen: consider adding 'Figure: %s' or 'Table: %s' after the figure/table", string(caption), string(caption))
	}
	return len(text)
}
