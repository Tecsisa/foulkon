#!/usr/bin/env bash
echo "--> Building..."
#Make sure $GOPATH is set
CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/cmd/worker
CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/cmd/proxy

mkdir bin/ 2>/dev/null
cp $GOPATH/bin/worker ./bin
cp $GOPATH/bin/proxy ./bin

echo "----> Building Docker images..."
docker build -t tecsisa/foulkon-proxy:latest -f scripts/docker/Dockerfile_proxy .
docker build -t tecsisa/foulkon-worker:latest -f scripts/docker/Dockerfile_worker .
echo "----> Pushing images to Docker hub"
curl -u ${DOCKER_HUB_USER}:${DOCKER_HUB_KEY} https://cloud.docker.com/api/app/v1/service/
docker push tecsisa/foulkon-worker:latest
docker push tecsisa/foulkon-proxy:latest
