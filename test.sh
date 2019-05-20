#!/bin/bash
set -e
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $DIR
for d in */ ; do
	echo "$d"
	cd $d
	go fmt ./...
	go mod tidy -v
	go vet ./...
	go test -cover -race -v ./...
	staticcheck ./... # download from https://github.com/dominikh/go-tools/releases
	cd ..
done
