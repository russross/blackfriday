// Functions and definition to be backwards compatible with RFC7328 markdown.

package mmark

import "bytes"

func (p *parser) rfc7328Index(out *bytes.Buffer, text []byte) int {
	//if p.flags&EXTENSION_RFC7328 == 0 {
	//return false
	//}
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
		p.r.Index(out, text[1:itemEnd], text[subItemStart:len(text)], false)
		return len(text)
	}
	p.r.Index(out, text[1:itemEnd], nil, false)
	return len(text)
}

func (p *parser) rfc7328Caption(out *bytes.Buffer, text []byte) int {
	//if p.flags&EXTENSION_RFC7328 == 0 {
	//return
	//}
	text = bytes.TrimSpace(text)
	println(string(text))
	return 0
}
