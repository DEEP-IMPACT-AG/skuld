#!/usr/bin/env bash
#
# This script builds skuld for Intel and multiple Arm processors, using
# docker buildx cross compilers. To build the latest skuld release, do:
#
# 1. Update the VERSION variable, to build the tagged version from git.
# 2. Run: docker login deepimpact
# 3. Run: ./buildx.sh
#
# That's it!

USER=deepimpact
REPO=skuld
VERSION=v0.7.5
IMAGE=${REPO}
BUILDER_CONTAINER_NAME=${REPO}-docker-builder
PLATFORM="linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6"

# Update the version in the compile script.
sed -i "s/^VERSION=.*/VERSION=${VERSION}/g" build/compile.sh

# Create the builder container if it doesn't exist
if [[ $(docker container ls -a | grep ${BUILDER_CONTAINER_NAME} | wc -l) -eq 1 ]]; then
	echo "Builder container exists, skipping creation."
else
	docker buildx create --name ${BUILDER_CONTAINER_NAME}
	docker buildx use ${BUILDER_CONTAINER_NAME}
	docker buildx inspect --bootstrap
fi

# Build and push the images
docker buildx build --platform=${PLATFORM} \
	--tag ${USER}/${IMAGE}:${VERSION} \
	--tag ${USER}/${IMAGE}:latest \
	-f Dockerfile \
	--push \
	.
