#!/bin/bash

latestCommitSha=$(git rev-parse --short HEAD)

echo "starting to build docker image based on commit: ${latestCommitSha}"

docker build --build-arg="${latestCommitSha}" -t mbvofdocker/blog:"${latestCommitSha}" -t mbvofdocker/blog:latest .

echo "build done, starting to push"

docker push mbvofdocker/blog:${latestCommitSha}
docker push mbvofdocker/blog:latest

echo "pushed to mbvofdocker/blog, deleting images"

docker image rm mbvofdocker/blog:${latestCommitSha} mbvofdocker/blog:latest

echo "cleanup done"
