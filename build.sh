#!/bin/sh -eu

rm -rf dist
mkdir dist
go build -o dist/kafka-data-import kafka-data-import.go
