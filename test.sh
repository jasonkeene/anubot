#!/bin/bash

set -ex

pushd src/anubot
    ginkgo unfocus
    goimports -w .
    ginkgo -r -race -randomizeAllSpecs
popd

pushd app
    jasmine
popd

echo 'All Tests Passed!'
