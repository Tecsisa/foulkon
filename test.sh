echo -e '\n################################\n# Starting authorizr unit test #\n################################\n'

go list ./... | grep -v '/vendor/' | grep -v '/database/' | PATH=$TEMPDIR:$PATH xargs -n1 go test ${GOTEST_FLAGS:--cover -timeout=900s}

echo -e '\n############################\n# Starting connectors test #\n############################\n'

# Postgres
echo -e '########################\n# PostgreSQL connector #\n########################'
echo $(echo -e 'Starting PostgreSQL container postgrestest with id ') \
$(docker run --name postgrestest -p 54320:5432 -e POSTGRES_PASSWORD=password -d postgres) \
$(echo -e '\n')
echo -e 'Starting authorizr test for PostgreSQL connector \n'
go test ./database/postgresql ${GOTEST_FLAGS:--cover -timeout=900s}
echo -e '\nRemoving PostgreSQL container' $(docker rm -f postgrestest)

echo -e '\n############################\n# Connectors test finished #\n############################\n'
