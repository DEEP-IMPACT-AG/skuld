#!/bin/bash
set -e

if [ -z ${CIRCLE_TAG} ]; then
	echo "No Tag = No Release"
else
	echo "Release: ${CIRCLE_TAG}"
	openssl aes-256-cbc -d \
              -in .circleci/credentials.enc \
              -out credentials \
              -k ${SNAP_CRED_PASSWORD}
    snapcraft login --with credentials
    rm credentials
	goreleaser
	find dist -name "skuld*.snap" | xargs snapcraft push --release edge
fi
