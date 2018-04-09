#!/usr/bin/env bash

echo "----> Building Demo..."

echo "----> Compiling demo example apps..."
# Make sure $GOPATH is set
CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/demo/api || exit 1
CGO_ENABLED=0 go install github.com/Tecsisa/foulkon/demo/web || exit 1

echo "----> Moving compiled files to GOROOT path..."
mkdir bin/ 2>/dev/null
cp $GOPATH/bin/api ./bin
cp $GOPATH/bin/web ./bin

echo "----> Building Docker demo images..."
docker build -t tecsisa/foulkondemo -f demo/docker/Dockerfile . >/dev/null || exit 1

echo "----> Starting Docker Compose..."
docker-compose -f demo/docker/docker-compose.yml up --force-recreate --abort-on-container-exit