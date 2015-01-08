all:
	( cd mmark; go build )

mmark2rfc.txt: mmark2rfc.md
	./mmark/mmark -xml2 -page mmark2rfc.md > x.xml && xml2rfc --text x.xml && rm x.xml && mv x.txt mmark2rfc.txt

mmark2rfc.html: mmark2rfc.md
	./mmark/mmark -page mmark2rfc.md > mmark2rfc.html
