FROM golang:alpine

MAINTAINER seanchann <zhouxiaoqiang.mstech@foxmail.com>

RUN set -ex \
	&& apk add --no-cache --virtual bash

COPY . /go/src/shadowss/
COPY /build/build-image/files/entrypoint.sh /entrypoint.sh
COPY /build/build-image/files/config.json //conf/config.json

RUN cd /go/src/shadowss/ \
	&& cp vendor/* /go/src/ -rf \
	&& cd cmd/shadowss/ \
	&& go build -v  -o  /go/bin/shadowss \
	&& chmod +x /entrypoint.sh


CMD ["entrypoint.sh"]
