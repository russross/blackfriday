#!/bin/bash
FILES="mmark/mmark CONVERSION_RFC7328.md mmark2rfc.md README.md rfc"
( cd mmark ; make clean; make )
VERSION=$(mmark/mmark -version)

# Linux
export GOOS=linux GOARCH=amd64
( cd mmark ; make clean; make )

dir=$(mktemp -d)
mkdir ${dir}/mmark
trap "rm -rf $dir" EXIT

cp -r $FILES $dir/mmark
( cd $dir; \
    tar --verbose --create --bzip2 --file /tmp/mmark-v$VERSION-$GOOS-$GOARCH.tar.bz2 mmark )

# Darwin
export GOOS=darwin GOARCH=amd64
( cd mmark ; make clean; make )

dir=$(mktemp -d)
mkdir ${dir}/mmark
trap "rm -rf $dir" EXIT

cp -r $FILES $dir/mmark
( cd $dir; \
    zip -r /tmp/mmark-v$VERSION-$GOOS-$GOARCH.zip mmark )

ls /tmp/mmark-v$VERSION-linux-$GOARCH.tar.bz2
ls /tmp/mmark-v$VERSION-darwin-$GOARCH.zip
