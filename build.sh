#!/bin/sh -eu

rm -rf dist
mkdir dist
go build -o dist/lift-data-import lift-data-import.go
