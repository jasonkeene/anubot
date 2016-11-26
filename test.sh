#!/bin/bash

set -e

if [ "$1" = "ci" ]; then
    echo installing node dependencies
    pushd app
        npm install --quiet
    popd

    echo installing go tools
    go get -u golang.org/x/tools/cmd/goimports
    go get -u github.com/tsenart/deadcode
    go get -u github.com/golang/lint/golint
    go get -u github.com/opennota/check/cmd/aligncheck
    go get -u github.com/opennota/check/cmd/structcheck
    go get -u github.com/opennota/check/cmd/varcheck
    go get -u github.com/kisielk/errcheck
    go get -u github.com/gordonklaus/ineffassign
    go get -u github.com/mvdan/interfacer/cmd/interfacer
    go get -u github.com/mdempsky/unconvert
    go get -u honnef.co/go/simple/cmd/gosimple
    go get -u honnef.co/go/staticcheck/cmd/staticcheck
    go get -u honnef.co/go/unused/cmd/unused
    go get -u github.com/client9/misspell/cmd/misspell
else
    echo building binaries

    mkdir -p bin
    ls src/anubot/cmd | while read line; do
        echo building $line
        go build -o bin/$line anubot/cmd/$line
    done
fi

echo running go tests

pushd src/anubot
    go test -race ./...

    goimports -w .
    gofmt -s -w .
    misspell -w .

    go vet ./...
    deadcode .
    golint ./...
    aligncheck ./...
    structcheck ./...
    varcheck ./...
    errcheck ./...
    ineffassign .
    interfacer ./...
    unconvert -v -apply ./...
    gosimple ./...
    staticcheck ./...
    unused ./...
popd

echo running jasmine tests

pushd app
    babel --presets es2015,react --out-dir lib src

    jasmine

    jshint lib
    jshint spec

    misspell -w src
popd

echo 'All Tests Passed!'
