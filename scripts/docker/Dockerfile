FROM alpine
MAINTAINER Tecsisa

USER root
RUN apk update && apk add ca-certificates

# Worker
COPY bin/worker /go/bin/worker
COPY dist/worker_env_vars.toml /worker.toml

# Proxy
COPY bin/proxy /go/bin/proxy
COPY dist/proxy_env_vars.toml /proxy.toml

# Entrypoint
ADD scripts/docker/entrypoint.sh /go/bin/entrypoint.sh
RUN chmod 750 /go/bin/*

ENV PATH=$PATH:/go/bin

EXPOSE 8000 8001

ENTRYPOINT ["/go/bin/entrypoint.sh"]
