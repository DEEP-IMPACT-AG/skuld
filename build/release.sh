#!/bin/bash

set +e

git describe --exact-match --tags HEAD
ON_TAG=$?

set -e

if [ "${ON_TAG}" -eq "0" ]; then
	curl -sL https://git.io/goreleaser | bash
else
	echo "No Tag = No Release"
fi
