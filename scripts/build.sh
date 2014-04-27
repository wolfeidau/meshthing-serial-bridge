#!/usr/bin/env bash
set -e

OWNER=wolfeidau
BIN_NAME=sniffer-bridge
PROJECT_NAME=meshthing-serial-bridge

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

GIT_COMMIT="$(git rev-parse HEAD)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"
VERSION="$(grep "const Version " version.go | sed -E 's/.*"(.+)"$/\1/' )"

# remove working build
rm -rf .gopath
mkdir -p .gopath/src/github.com/${OWNER}
ln -sf ../../../.. .gopath/src/github.com/${OWNER}/${PROJECT_NAME}
export GOPATH="$(pwd)/.gopath"

# move the working path and build
cd .gopath/src/github.com/${OWNER}/${PROJECT_NAME}
go build -ldflags "-X main.GitCommit ${GIT_COMMIT}${GIT_DIRTY}" -o ${BIN_NAME}
mv ${BIN_NAME} ./bin