#!/bin/bash

set -e

echo building binaries

mkdir -p bin
ls src/anubot/cmd | while read line; do
    echo building $line
    go build -o bin/$line anubot/cmd/$line
done

echo running go tests

pushd src/anubot
    if [ ! "$(which goimports)" ]; then
        go install golang.org/x/tools/cmd/goimports
    fi
    goimports -w .
    go test -race ./...
popd

echo running jasmine tests

pushd app
    babel --presets es2015,react --out-dir lib src;
    jasmine
    jshint lib
    jshint spec
popd

echo 'All Tests Passed!'
