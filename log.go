// Functions to parse block-level elements.

package mmark

import "log"

func Printf(p *parser, format string, v ...interface{}) {
	if p != nil {
		log.Printf("mmark: "+format, v...)
		return
	}
	log.Printf("mmark: "+format, v...)
}
