#!/bin/bash

export GOOS="linux"
export GOARCH="amd64"
export CGO_ENABLED="0"

build_lambda() {
  go build -o bootstrap -tags lambda.norpc
  zip main.zip bootstrap
  rm -rf bootstrap
}

build_lambda