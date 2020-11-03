#!/bin/bash
go test -v -timeout 5s \
  ./... \
	| if [ "$1" != 'c' ]; then cat; else sed 's#.*PASS.*#\x1b[32;1m&\x1b[0m#;s#.*FAIL.*#\x1b[31;1m&\x1b[0m#'; fi
