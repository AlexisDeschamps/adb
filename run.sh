#!/bin/sh

flags=""
if [[ -d adb-wayne-config ]]; then
  . adb-wayne-config/env
  flags="-prod"
fi

./adb $flags
