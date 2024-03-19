#!/bin/sh

set -xe

# This script is running inside a docker build container. It compiles the skuld
# binary and copies it to the /dist folder. Then the binaries are copied to a
# release image without the golang compile environment.
#
# Note that if this script is invoked from the buildx.sh script, the VERSION
# number is automatically updated from the buildx.sh script.

OS=$(uname -s)
ARCH=$(uname -m)
VERSION=v0.7.5

# Clone and compile with release flags.
cd /workdir
git clone --depth 1 --branch ${VERSION} https://github.com/DEEP-IMPACT-AG/skuld.git /workdir/skuld_src
cd /workdir/skuld_src
go build -ldflags "-s -w"

# Copy the binary to a separate folder for easy access by the "release" docker image.
mkdir /dist
cp /workdir/skuld_src/skuld /dist/skuld_${OS}_${ARCH}
cp /workdir/skuld_src/skuld /dist/
cd ..
