#!/bin/bash

set -e

build_lib() {
    if [ -e lib ]; then
        echo removing lib
        rm -r lib
    fi
    babel --presets es2015,react --out-dir lib src
    node-sass --output lib/styles src/styles
}


build_mac_app() {
    if [ -e dist ]; then
        echo removing dist
        rm -r dist
    fi
    sudo spctl --master-disable
    build --macos
    sudo spctl --master-enable
}

main() {
    echo -e "\033[1m\033[34mBuilding lib\033[0m"
    build_lib
    echo

    echo -e "\033[1m\033[34mBuilding app\033[0m"
    build_mac_app
    echo
}

main
