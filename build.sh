#!/bin/sh -eu

rm -rf dist
mkdir dist
go build -a -v -tags musl -o dist/kafka-data-import kafka-data-import.go
