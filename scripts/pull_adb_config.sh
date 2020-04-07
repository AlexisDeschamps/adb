#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $DIR/..

if [[ -d adb-wayne-config ]] ; then
    cd adb-wayne-config
    git pull
else
    git clone 'git@github.com:dxe/adb-wayne-config'
fi
