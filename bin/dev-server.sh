#!/bin/bash -xe

cd "$(dirname "$0")/.."
ROOT=$(pwd)

cd server
go get github.com/cespare/reflex
dep ensure
reflex -s go run cmd/server.go
