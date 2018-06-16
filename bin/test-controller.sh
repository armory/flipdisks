#!/bin/bash -xe

cd "$(dirname "$0")/.."
ROOT=$(pwd)

cd controller/
dep ensure
go test -v ./...
