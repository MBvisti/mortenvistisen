#!/bin/bash

latestCommitSha=$(git rev-parse --short HEAD)

echo "starting to build docker image based on commit: ${latestCommitSha}"

docker build --target app --build-arg="COMMIT_SHA=${latestCommitSha}" -t mbvofdocker/thelab:blog-"${latestCommitSha}" -t mbvofdocker/thelab:blog-latest .

docker build --target worker --build-arg="COMMIT_SHA=${latestCommitSha}" -t mbvofdocker/thelab:blog-worker-"${latestCommitSha}" -t mbvofdocker/thelab:blog-worker-latest .

echo "build done, starting to push"

docker push mbvofdocker/thelab:blog-${latestCommitSha}
docker push mbvofdocker/thelab:blog-latest
docker push mbvofdocker/thelab:blog-worker-${latestCommitSha} 
docker push mbvofdocker/thelab:blog-worker-latest

echo "pushed to mbvofdocker/blog, deleting images"

docker image rm mbvofdocker/thelab:blog-${latestCommitSha} mbvofdocker/thelab:blog-latest mbvofdocker/thelab:blog-worker-${latestCommitSha} mbvofdocker/thelab:blog-worker-latest

echo "cleanup done"

echo "releasing"

ssh admin@188.245.71.73 /bin/bash << EOF
	cd blog;
	docker compose pull;
	docker rollout blog;
	docker rollout blog-worker;
	docker image prune -a -f
EOF
