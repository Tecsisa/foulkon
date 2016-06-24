go list ./... | grep -v '/vendor/' | sudo -E PATH=$TEMPDIR:$PATH xargs -n1 go test ${GOTEST_FLAGS:--cover -timeout=900s}
