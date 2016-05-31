FROM alpine
MAINTAINER Tecsisa

COPY bin/authorizr /go/bin/authorizr
ENTRYPOINT ["/go/bin/authorizr"]

EXPOSE 8000
