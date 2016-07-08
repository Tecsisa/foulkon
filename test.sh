echo -e '- Starting authorizr unit test'
go list ./... | grep -v '/vendor/' | grep -v '/database/' | PATH=$TEMPDIR:$PATH xargs -n1 go test ${GOTEST_FLAGS:--cover -timeout=900s}
echo -e '- Starting connectors test'
echo -e '-- Starting test for Postgres connector'
echo $(echo -e '--- Starting Postgres container postgrestest with id ') $(docker run --name postgrestest -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres)
echo -e '--- Starting authorizr test for postgres connector'
go test ./database/postgresql ${GOTEST_FLAGS:--cover -timeout=900s}
echo $(echo -e '--- Removing Postgres container ') $(docker rm -f postgrestest)
