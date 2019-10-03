#!/usr/bin/env bash
#
# This script update homebrew tap

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"

go run cmd/letitgo/main.go $(git describe --tags --abbrev=0)
