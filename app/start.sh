#!/bin/bash

set -e

build_server() {
    if [ -e anubot-server ]; then
        echo removing anubot-server
        rm anubot-server
    fi
    go build -o anubot-server anubot/cmd/api-server
}

main() {
    build_server
    electron .
}

main
