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

watch_app_source() {
    babel --presets es2015,react --watch --out-dir lib src &
    BABEL_PID=$!
    echo babel started as pid $BABEL_PID
}

watch_app_styles() {
    node-sass --watch --output lib/styles src/styles &
    NODE_PID=$!
    echo node-stats started as pid $NODE_PID
}

kill_watchers() {
    if [ -n "$BABEL_PID" ]; then
        echo tearing down babel pid $BABEL_PID
        kill $BABEL_PID
    fi
    if [ -n "$NODE_PID" ]; then
        echo tearing down node-scss pid $NODE_PID
        kill $NODE_PID
    fi
}

main() {
    echo -e "\033[1m\033[34mBuilding lib\033[0m"
    build_lib
    echo

    trap kill_watchers EXIT
    echo -e "\033[1m\033[34mWatching Application Source\033[0m"
    watch_app_source
    echo

    echo -e "\033[1m\033[34mWatching Application Styles\033[0m"
    watch_app_styles
    echo

    echo -e "\033[1m\033[34mStarting Application\033[0m"
    electron .
}

main
