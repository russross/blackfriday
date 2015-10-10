MMARK:=./mmark/mmark -xml2 -page
MMARK3:=./mmark/mmark -xml -page

objects := README.md.txt mmark2rfc.md.txt

%.md.txt: %.md
	$(MMARK) $< > $<.xml
	xml2rfc --text $<.xml && rm $<.xml

all: mmark/mmark $(objects)

mmark/mmark:
	( cd mmark; make )

mmark2rfc.md.3.xml: mmark2rfc.md mmark/mmark
	$(MMARK3) $< > $<.3.xml

.PHONY: clean
clean:
	rm -f *.md.txt *md.[23].xml

.PHONY: validate
validate: mmark2rfc.md.3.xml
	xmllint --xinclude $< | jing -c xml2rfcv3.rnc /dev/stdin
