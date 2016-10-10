#!/usr/bin/env bash

export PATH=$TEMPDIR:$PATH
export GOPATH=$GOPATH

echo "" > coverage.txt

echo "--> Running tests"
echo -e '----> Running unit tests'
for d in $(go list ./... | grep -v '/vendor/' | egrep -v '/database/|cmd/|auth/oidc|foulkon/foulkon'); do
    go test -race -coverprofile=profile.out -covermode=atomic $d || exit 1
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

echo -e '\n----> Running connector tests'
# Postgres
echo -e '--------> Running PostgreSQL connector'
echo $(echo -e 'Starting PostgreSQL (Docker container) postgrestest with id ') \
$(docker run --name postgrestest -p 54320:5432 -e POSTGRES_PASSWORD=password -d postgres) \
$(echo -e '\n\n')
go test ./database/postgresql ${GOTEST_FLAGS:--race -coverprofile=profile.out -covermode=atomic -timeout=900s} || exit 1
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
echo -e 'Removing PostgreSQL container' $(docker rm -f postgrestest) '\n'
