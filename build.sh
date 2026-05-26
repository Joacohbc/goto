#!/bin/bash

VERSION=${1:-"unknown"}

go build -ldflags="-X 'goto/src/cmd.VersionGoto=$VERSION'" -o goto.bin src/*.go
mv goto.bin ~/.config/goto/goto.bin
