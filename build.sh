#!/bin/sh

GOOS=$1
GOARCH=$2
EXE=$3
GOOS=$GOOS GOARCH=$GOARCH go build -o ./ ./...
zip -r -j ../$VERSION/reliza-cli-$VERSION-$GOOS-$GOARCH.zip ./reliza-cli$EXE