#!/bin/bash

set -e

if [ "$1" = "ci" ]; then
    echo installing node dependencies
    pushd app
        npm install
    popd
else
    echo running go imports
    pushd src/anubot
        if [ ! "$(which goimports)" ]; then
            go install golang.org/x/tools/cmd/goimports
        fi
        goimports -w .
    popd
fi

echo building binaries

mkdir -p bin
ls src/anubot/cmd | while read line; do
    echo building $line
    go build -o bin/$line anubot/cmd/$line
done

echo running go tests

pushd src/anubot
    go test -race ./...
popd

echo running jasmine tests

pushd app
    babel --presets es2015,react --out-dir lib src
    jasmine
    jshint lib
    jshint spec
popd

echo 'All Tests Passed!'
