#!/bin/bash

set -e

if [ "$1" = "ci" ]; then
    echo installing node dependencies
    npm install --silent --no-progress --depth=0
fi

echo building lib
babel --presets es2015,react --out-dir lib src

echo running jasmine tests
jasmine

echo running linters
jshint lib
jshint spec

echo 'All Tests Passed!'
