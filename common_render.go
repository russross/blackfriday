// Common functions used in the render backends.

package mmark

import "bytes"

// blockCodePrefix adds the prefix to each line of text and returns it as a byte slice.
// If prefix is empty, text is returned as-is.
func blockCodePrefix(prefix string, text []byte) []byte {
	if prefix == "" {
		return text
	}

	nl := bytes.Count(text, []byte{'\n'})
	prefixText := bytes.Replace(text, []byte{'\n'}, []byte("\n"+prefix), nl-1)
	prefixText = append([]byte(prefix), prefixText...)
	return prefixText
}
