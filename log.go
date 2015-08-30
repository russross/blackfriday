// Functions to parse block-level elements.

package mmark

import "log"

func Printf(p *parser, format string, v ...interface{}) {
	if p != nil {
		log.Printf("%s: line %d: "+format, "mmark", p.lineNumber, v)
		return
	}
	log.Printf("%s: "+format, "mmark", v)
}
