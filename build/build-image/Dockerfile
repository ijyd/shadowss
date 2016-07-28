FROM golang:alpine

MAINTAINER seanchann <seanchann@foxmail.com>

RUN set -ex \
	&& apk add --no-cache --virtual bash

COPY . /go/src/github.com/shadowsocks/shadowsocks-go/

RUN cd /go/src/github.com/shadowsocks/shadowsocks-go/ \
	&& cp vendor/* /go/src/ -rf \
	&& cp sample-config/server-multi-port.json /config.json \
	&& cd cmd/shadowsocks-server/ \
	&& go build -v -o  /go/bin/shadowss


CMD ["/go/bin/shadowss", "-c", "/config.json"]
