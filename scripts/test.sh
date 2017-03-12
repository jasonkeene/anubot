#!/bin/bash

set -e

passed[0]="ᕕ( ᐛ )ᕗ"
passed[1]="(╯✧∇✧)╯"
passed[2]="(づ｡◕‿‿◕｡)づ"
passed[3]="(๑˃̵ᴗ˂̵)و"

failed[0]="щ(ಥДಥщ)"
failed[1]="ヽ(｀⌒´メ)ノ"
failed[2]="(屮ಠ益ಠ)屮"
failed[3]="ヽ(#ﾟДﾟ)ﾉ┌┛"
failed[4]="(ノಠ益ಠ)ノ彡┻━┻"

function passed {
    local rand=$[ $RANDOM % 4 ]
    echo
    echo -e "  \033[32;48;5;2m  ${passed[$rand]}                    \033[0m"
    echo -e "  \033[30;48;5;2m  ${passed[$rand]} ALL TESTS PASSED!  \033[0m"
    echo -e "  \033[32;48;5;2m  ${passed[$rand]}                    \033[0m"
    echo
}

function failed {
    local rand=$[ $RANDOM % 5 ]
    echo
    echo -e "  \033[31;48;5;1m  ${failed[$rand]}               \033[0m"
    echo -e "  \033[30;48;5;1m  ${failed[$rand]} TEST FAILED!  \033[0m"
    echo -e "  \033[31;48;5;1m  ${failed[$rand]}               \033[0m"
    echo
}

function check_status {
    if [ $1 -ne 0 ]; then
        failed
    else
        passed
    fi
}
function handle_exit {
    status=$?
    for exit_func in "${exit_funcs[@]}"; do
        $exit_func
    done
    check_status $status
}
trap handle_exit EXIT

function checkpoint {
    echo
    echo -e "  \033[38;5;104m$@\033[0m"
    echo
}

function header {
    echo
    echo -e "  \033[30;48;5;104m                             \033[0m"
    echo -e "  \033[30;48;5;104m   ___ ___ _ _| |_ ___| |_   \033[0m"
    echo -e "  \033[30;48;5;104m  | .'|   | | | . | . |  _|  \033[0m"
    echo -e "  \033[30;48;5;104m  |__,|_|_|___|___|___|_|    \033[0m"
    echo -e "  \033[30;48;5;104m              test suite     \033[0m"
    echo -e "  \033[30;48;5;104m                             \033[0m"
}

function install_deps {
    if [ "$1" = "ci" ]; then
        checkpoint installing node dependencies
        npm install --silent --no-progress --depth=0
    fi
}

function build {
    checkpoint building lib
    babel --presets es2015,react --out-dir lib src
}

function test {
    checkpoint running tests
    jasmine
}

function lint {
    checkpoint running linters
    echo jshint lib
    jshint lib
    echo jshint spec
    jshint spec
}

function main {
    header
    install_deps $1
    build
    test
    lint
}
main $@
