// Functions to parse block-level elements.

package mmark

import "log"

func Printf(p *parser, format string, v ...interface{}) {
	if p != nil {
		// We don't track newlines seen, so we don't error on a specific line.
		//log.Printf("%s: line %d: "+format, "mmark", p.lineNumber, v)
		log.Printf("%s: "+format, "mmark", v)
		return
	}
	log.Printf("%s: "+format, "mmark", v)
}
