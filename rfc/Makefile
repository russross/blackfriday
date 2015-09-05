MMARK:=../mmark/mmark -xml2 -page
MMARK3:=../mmark/mmark -xml -page

objects := $(patsubst %.md,%.md.txt,$(wildcard *.md))

all: $(objects)

%.md.txt: %.md
	$(MMARK) $< > $<.xml
	xml2rfc --text $<.xml && rm $<.xml

%.md.2.xml: %.md
	$(MMARK) $< > $<.2.xml

%.md.3.xml: %.md
	$(MMARK3) $< > $<.3.xml

.PHONY: clean
clean:
	rm -f *.md.txt *md.[23].xml
