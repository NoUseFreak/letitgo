#!/usr/bin/env bash
#
# This script tests the application from source.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"

echo "===> Running tests..."
go test ./...

echo "===> Running errcheck..."
(cd /tmp && go get -u github.com/kisielk/errcheck)
errcheck  ./...

echo "===> Running staticcheck..."
(cd /tmp && go get -u honnef.co/go/tools/cmd/staticcheck)
staticcheck  ./...