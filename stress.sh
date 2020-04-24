#!/bin/sh

set -e

for i in `seq 0 7`
do
    echo -e "\n=== $i ===\n"
    go test -count=1 ./...
done
