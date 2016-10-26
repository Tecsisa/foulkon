FROM tecsisa/foulkon:v0.2.0
MAINTAINER Tecsisa

# Copy resources of this API
COPY demo/docker/demo-proxy.toml /proxy.toml

# API
COPY bin/api /go/bin/api

# WEB
COPY bin/web /go/bin/web
COPY demo/web/tmpl /tmpl

# Entrypoint
ADD demo/docker/entrypoint.sh /go/bin/entrypoint.sh
RUN chmod 750 /go/bin/*

ENV PATH=$PATH:/go/bin

EXPOSE 8000 8001 8100 8101

ENTRYPOINT ["/go/bin/entrypoint.sh"]