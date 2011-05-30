include $(GOROOT)/src/Make.inc

TARG=markdown

GOFILES=main.go

markdown: ../_obj/github.com/russross/blackfriday.a

GC += -I../_obj/
LD += -L../_obj/

include $(GOROOT)/src/Make.cmd
