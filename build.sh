#Make sure $GOPATH is set
CGO_ENABLED=0 go install github.com/tecsisa/authorizr

mkdir bin/ 2>/dev/null
cp $GOPATH/bin/authorizr ./bin

docker build -t tecsisa/authorizr .

echo "Docker image built!"
docker images | grep "tecsisa/authorizr"
