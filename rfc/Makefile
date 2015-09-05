MMARK:=../mmark/mmark -xml2 -page
MMARK3:=../mmark/mmark -xml -page

objects := $(patsubst %.md,%.md.txt,$(wildcard *.md))
objectsv3xml := $(patsubst %.md,%.md.3.xml,$(wildcard *.md))

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

.PHONY: validate
validate: $(objectsv3xml)
	for i in $^; do echo xmllint --xinclude $$i | jing -c ../xml2rfcv3.rnc /dev/stdin; done
