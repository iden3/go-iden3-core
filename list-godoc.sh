#!/bin/sh

for path in `find . -type d | grep -v ".git" | grep -v "testVectors" | tail -n+2 | sed 's#./##'`; do
    printf -- "- [![GoDoc](https://godoc.org/github.com/iden3/go-iden3-core/%s?status.svg)](https://godoc.org/github.com/iden3/go-iden3-core/%s) %s\n" $path $path $path
done
