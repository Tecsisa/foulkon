#!/usr/bin/env bash
echo "==> Building..."
#Make sure $GOPATH is set
CGO_ENABLED=0 go install github.com/tecsisa/authorizr/cmd/worker
CGO_ENABLED=0 go install github.com/tecsisa/authorizr/cmd/proxy

mkdir bin/ 2>/dev/null
cp $GOPATH/bin/worker ./bin
cp $GOPATH/bin/proxy ./bin

echo "==> Building Docker images..."
docker build -t tecsisa/authorizr-proxy -f Dockerfile_proxy .
docker build -t tecsisa/authorizr-worker -f Dockerfile_worker .
echo "Docker images built!"
docker images | grep "tecsisa/authorizr"
