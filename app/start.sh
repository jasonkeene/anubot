#!/bin/bash

set -e

build_server() {
    if [ -e anubot-server ]; then
        echo removing anubot-server
        rm anubot-server
    fi
    go build -o anubot-server anubot/cmd/api-server
}

build_app_views() {
    babel --presets es2015,react --out-dir lib/views src/views
}

build_app() {
    cp src/app.js lib/app.js
}

build_app_styles() {
    node-sass --output lib/styles src/styles
}

main() {
    echo -e "\033[1m\033[34mBuilding API Server\033[0m"
    echo
    build_server
    echo

    echo -e "\033[1m\033[34mBuilding Application Views\033[0m"
    echo
    build_app_views
    echo

    echo -e "\033[1m\033[34mBuilding Application Styles\033[0m"
    echo
    build_app_styles
    echo

    echo -e "\033[1m\033[34mBuilding Application\033[0m"
    echo
    build_app
    echo

    echo -e "\033[1m\033[34mStarting Application\033[0m"
    echo
    electron .
}

main
