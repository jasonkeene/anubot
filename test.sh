#!/bin/bash

set -ex

echo building binaries

mkdir -p bin
ls src/anubot/cmd | while read line; do
    go build -o bin/$line anubot/cmd/$line
done

echo running ginkgo tests

pushd src/anubot
    ginkgo unfocus
    goimports -w .
    ginkgo -r -race -randomizeAllSpecs
popd

echo running jasmine tests

pushd app
    jasmine
popd

echo 'All Tests Passed!'
