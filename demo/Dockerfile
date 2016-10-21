FROM alpine
MAINTAINER Tecsisa

USER root
RUN apk update && apk add ca-certificates

# API
COPY bin/demoapi /go/bin/demoapi

# WEB
COPY bin/demoweb /go/bin/demoweb

ENV PATH=$PATH:/go/bin

EXPOSE 8000 8001