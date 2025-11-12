#!/bin/bash

PACKAGE=github.com/marko-gacesa/gamatet/internal/values
VAL_TIME=`date -u +%FT%T%Z`
VAL_SHA=`git rev-parse HEAD`
VAL_TAG=`git describe --tags`
EXE_NAME=gmt

go build -ldflags="-s -w -X $PACKAGE.BuildTime=$VAL_TIME -X $PACKAGE.GitSHA=$VAL_SHA -X $PACKAGE.VersionTag=$VAL_TAG" -o "$EXE_NAME" main.go
