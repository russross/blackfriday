package mmark

import "bytes"

// linespan implements a minimal line iterator over '\n' delimited content
type linespan struct{ begin, end int }

// next updates begin and end to point to the next line
func (sc *linespan) next(content []byte) bool {
	sc.begin = sc.end
	if sc.begin >= len(content) {
		return false
	}

	off := bytes.IndexByte(content[sc.begin:], '\n')
	if off >= 0 {
		sc.end = sc.begin + off + 1
		return true
	}

	sc.end = len(content)
	return true
}
