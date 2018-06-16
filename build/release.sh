#!/bin/bash
set -e

if [ -z ${CIRCLE_TAG} ]; then
	echo "No Tag = No Release"
else
	echo "Release: ${CIRCLE_TAG}"
	curl -sL https://git.io/goreleaser | bash
fi
