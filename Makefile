include $(GOROOT)/src/Make.inc

TARG=github.com/russross/blackfriday

GOFILES=markdown.go html.go smartypants.go

package:

include $(GOROOT)/src/Make.pkg

markdown: package
	make -C example
