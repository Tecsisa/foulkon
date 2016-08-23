#!/usr/bin/env bash
echo "--> Running tests"
echo -e '----> Running unit tests'
go list ./... | grep -v '/vendor/' | egrep -v '/database/|auth|cmd/|foulkon/foulkon' | PATH=$TEMPDIR:$PATH xargs -n1 go test ${GOTEST_FLAGS:--cover -timeout=900s}

echo -e '\n----> Running connector tests'
# Postgres
echo -e '--------> Running PostgreSQL connector'
echo $(echo -e 'Starting PostgreSQL (Docker container) postgrestest with id ') \
$(docker run --name postgrestest -p 54320:5432 -e POSTGRES_PASSWORD=password -d postgres) \
$(echo -e '\n\n')
go test ./database/postgresql ${GOTEST_FLAGS:--cover -timeout=900s}
echo -e 'Removing PostgreSQL container' $(docker rm -f postgrestest) '\n'