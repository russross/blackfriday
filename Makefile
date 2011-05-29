include $(GOROOT)/src/Make.inc

TARG=github.com/russross/blackfriday

GOFILES=markdown.go block.go inline.go html.go smartypants.go

include $(GOROOT)/src/Make.pkg

markdown: package
	make -C example
