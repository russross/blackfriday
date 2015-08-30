// Functions to parse block-level elements.

package mmark

import "log"

func (p *parser) Printf(format string, v ...interface{}) {
	log.Printf("%s: line %d: " + format, "mmark", p.lineNumber, v)
}
