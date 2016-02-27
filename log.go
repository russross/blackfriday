package mmark

import "log"

func printf(p *parser, format string, v ...interface{}) {
	if test {
		return
	}
	if p != nil {
		log.Printf("mmark: "+format, v...)
		return
	}
	log.Printf("mmark: "+format, v...)
}
