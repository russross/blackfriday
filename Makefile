include $(GOROOT)/src/Make.inc

TARG=github.com/russross/blackfriday

GOFILES=markdown.go html.go smartypants.go

include $(GOROOT)/src/Make.pkg
