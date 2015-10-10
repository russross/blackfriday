#!/bin/bash

FILES="mmark/mmark CONVERSION_RFC7328.md mmark2rfc.md README.md skel.md misc rfc"
VERSION=$(mmark/mmark -version)

dir=$(mktemp -d)
mkdir ${dir}/mmark
trap "rm -rf $dir" EXIT

cp -r $FILES $dir/mmark
source <(go tool dist env)
( cd $dir; \
    tar --verbose --create --bzip2 --file /tmp/mmark-v$VERSION-$GOOS-$GOARCH.tar.bz2 mmark ) && \
ls /tmp/mmark-v$VERSION-$GOOS-$GOARCH.tar.bz2
