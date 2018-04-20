#!/bin/bash -xe

cd "$(dirname "$0")/.."
ROOT=$(pwd)

cd webclient
yarn install
yarn start
