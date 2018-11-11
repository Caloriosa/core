#!/bin/sh

mkdir -p ./bin

for i in `ls cmd`; do
  go build -v -o ./bin/caloriosa-${i} ./cmd/$i
done
