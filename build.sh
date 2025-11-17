#!/bin/bash

PARAM_TEST=test

PACKAGE=github.com/marko-gacesa/gamatet/internal/values
VAL_TIME=`date -u +%FT%T%Z`
VAL_SHA=`git rev-parse HEAD`
VAL_TAG=`git describe --tags`
EXE_NAME=gmt

PARAMS_GO=()
PARAMS_LINKER=(-X "$PACKAGE.BuildTime=$VAL_TIME" -X "$PACKAGE.GitSHA=$VAL_SHA" -X "$PACKAGE.VersionTag=$VAL_TAG")
if [[ "$1" == "$PARAM_TEST" ]]; then
  PARAMS_GO+=(-race)
else
  PARAMS_LINKER+=(-s -w)
fi

PARAMS_LINKER_GO="${PARAMS_LINKER[@]}"

go build -o "$EXE_NAME" "${PARAMS_GO[@]}" -ldflags "$PARAMS_LINKER_GO" main.go
