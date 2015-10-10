GITHASH:=$(shell git rev-parse --short HEAD)

all: mmark

mmark: ../*.go
	go build -ldflags "-X main.githash=$(GITHASH)" -a -tags netgo -installsuffix netgo

.PHONY: clean
clean:
	rm -f mmark
