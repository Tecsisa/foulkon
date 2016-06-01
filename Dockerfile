FROM alpine
MAINTAINER Tecsisa

COPY bin/authorizr /go/bin/authorizr
COPY config_env_vars.toml /config_env_vars.toml
ENTRYPOINT ["/go/bin/authorizr", "-config-file=config_env_vars.toml"]

EXPOSE 8000
