#Make sure $GOPATH is set
CGO_ENABLED=0 go install github.com/tecsisa/authorizr/cmd/worker
CGO_ENABLED=0 go install github.com/tecsisa/authorizr/cmd/proxy

mkdir bin/ 2>/dev/null
cp $GOPATH/bin/worker ./bin
cp $GOPATH/bin/proxy ./bin

docker build -t tecsisa/authorizr .

echo "Docker image built!"
docker images | grep "tecsisa/authorizr"
