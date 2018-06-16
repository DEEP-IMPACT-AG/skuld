#!/bin/bash

set -e

env

go vet

go get github.com/mitchellh/gox
gox -os "linux darwin windows" -arch "386 amd64"

