#!/usr/bin/env bash
echo "--> Building..."
#Make sure $GOPATH is set
CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/cmd/worker
CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/cmd/proxy

mkdir bin/ 2>/dev/null
cp $GOPATH/bin/worker ./bin
cp $GOPATH/bin/proxy ./bin

echo "----> Building Docker images..."
docker build -t tecsisa/foulkon:latest -f scripts/docker/Dockerfile .

echo "----> Pushing images to Docker hub"
docker login -u ${DOCKER_HUB_USER} -p ${DOCKER_HUB_KEY}
docker push tecsisa/foulkon:latest
