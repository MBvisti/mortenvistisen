#!/bin/bash

latestCommitSha=$(git rev-parse --short HEAD)

echo "starting to build docker image based on commit: ${latestCommitSha}"

docker build --target app --build-arg="COMMIT_SHA=${latestCommitSha}" -t mbvofdocker/blog:"${latestCommitSha}" -t mbvofdocker/blog:latest .

docker build --target worker --build-arg="COMMIT_SHA=${latestCommitSha}" -t mbvofdocker/blog:"${latestCommitSha}"-worker -t mbvofdocker/blog:latest-worker .

echo "build done, starting to push"

docker push mbvofdocker/blog:${latestCommitSha}
docker push mbvofdocker/blog:latest
docker push mbvofdocker/blog:${latestCommitSha}-worker
docker push mbvofdocker/blog:latest-worker

echo "pushed to mbvofdocker/blog, deleting images"

docker image rm mbvofdocker/blog:${latestCommitSha} mbvofdocker/blog:latest mbvofdocker/blog:${latestCommitSha}-worker mbvofdocker/blog:latest-worker

echo "cleanup done"
