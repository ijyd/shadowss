FROM golang:alpine

MAINTAINER seanchann <seanchann@foxmail.com>

RUN echo "http://dl-1.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories; \
    echo "http://dl-2.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories; \
    echo "http://dl-3.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories; \
    echo "http://dl-4.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories; \
    echo "http://dl-5.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories; \
		echo "http://dl-1.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
		echo "http://dl-2.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
		echo "http://dl-3.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
		echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
		echo "http://dl-5.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories

RUN set -ex \
	&& apk add --no-cache bash git dmidecode

COPY . /go/src/cloud-keeper/
COPY /deploy/build-image/files/keys/ /keys/
COPY /deploy/build-image/files/entrypoint.sh /entrypoint.sh


#RUN go get -u github.com/Masterminds/glide/...

RUN cd /go/src/cloud-keeper/ \
	&& cp vendor/* /go/src/ -rf \
	&& cd cmd/vpskeeper/ \
	&& go build -v  -o  /go/bin/vpskeeper \
	&& chmod +x /entrypoint.sh \
	&& rm -rf /go/src/


CMD ["/entrypoint.sh"]
