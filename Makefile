MMARK2:=./mmark/mmark -xml2 -page
MMARK3:=./mmark/mmark -xml -page

all:
	@echo "not defined"

mmark/mmark:
	( cd mmark; make )

mmark2rfc2.txt: mmark2rfc2.xml
	xml2rfc --text mmark2rfc2.xml
	@ls -l mmark2rfc2.txt

mmark2rfc2.xml: mmark2rfc.md mmark/mmark
	$(MMARK2) mmark2rfc.md > mmark2rfc2.xml

#mmark2rfc3.txt: mmark2rfc3.xml
#	xml2rfc --text mmark2rfc3.xml
#	@ls -l mmark2rfc3.txt

mmark2rfc3.xml: mmark2rfc.md mmark/mmark
	$(MMARK3) mmark2rfc.md > mmark2rfc3.xml

.PHONY: clean
clean:
	rm -f mmark2rfc2.xml mmark2rfc3.xml mmark2rfc2.txt

.PHONY: validate
validate: mmark2rfc3.xml
	xmllint --xinclude mmark2rfc3.xml | jing -c xml2rfcv3.rnc /dev/stdin

.PHONY: release
release: mmark/mmark
