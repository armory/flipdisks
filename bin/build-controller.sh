#!/bin/bash -xe

cd "$(dirname "$0")"/..

rm -rf controller/build/ || true  # remove the old build
GOOS=linux GOARCH=arm GOARM=7 go build -o controller/build/main controller/cmd/main.go
