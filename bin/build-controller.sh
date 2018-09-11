#!/bin/bash -xe

cd "$(dirname "$0")"/..

cd controller
rm -rf build/ || true  # remove the old build

dep ensure
GOOS=linux GOARCH=arm GOARM=7 go build -o build/main cmd/main.go
