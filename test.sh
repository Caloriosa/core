#!/bin/sh

for i in `ls cmd`; do
  go test -v ./cmd/$i
done
