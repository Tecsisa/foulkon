FROM alpine
MAINTAINER Tecsisa

RUN apk update && apk add ca-certificates
# Worker
COPY bin/worker /go/bin/worker
COPY config_env_vars.toml /config_env_vars.toml
# Proxy
COPY bin/proxy /go/bin/proxy
COPY proxy_env_vars.toml /proxy_env_vars.toml

ADD entrypoint.sh /go/bin/entrypoint.sh

EXPOSE 8000 8001

ENTRYPOINT ["/go/bin/entrypoint.sh"]
