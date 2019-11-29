#!/bin/bash -xe

cd "$(dirname "$0")/.."

cd controller/
go test -v ./...
