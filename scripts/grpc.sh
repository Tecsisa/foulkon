#!/bin/bash
set -x
GOPATH=$(go env GOPATH)

# requires:
# - protoc (brew install grpc),
# - protoc-gen-go (go install github.com/golang/protobuf/protoc-gen-go)
protoc -I. \
  -I$GOPATH/src \
  --go_out=plugins=grpc:$GOPATH/src \
  grpc/*.proto

# requires:
# - github.com/ckaznocha/protoc-gen-lint
protoc -I. \
  -I$GOPATH/src \
  --lint_out=. \
  grpc/*.proto
