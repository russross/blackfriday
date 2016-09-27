all: mmark/mmark

mmark/mmark:
	( cd mmark; make)

.PHONY: clean
clean:
	( cd ./mmark; make clean )
