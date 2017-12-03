#!/bin/sh

for i in `ls cmd`; do
  go build -v ./cmd/$i
done
