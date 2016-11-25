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
	&& apk add --no-cache bash  git  dmidecode openssh ansible py-setuptools

RUN easy_install-2.7 pip\
  && pip install 'dopy>=0.3.5,<=0.3.5' \
  && mkdir /etc/ansible/ \
  && touch /etc/ansible/ansible.cfg \
  && echo "[defaults]" > /etc/ansible/ansible.cfg \
  && echo "library        = /go/bin/ansible/vultr/libraray/vultr" >>  /etc/ansible/ansible.cfg \
  && echo "host_key_checking = False" >>  /etc/ansible/ansible.cfg


COPY . /go/src/cloud-keeper/
#COPY /build/build-image/files/keys/ /keys/
#COPY /build/build-image/files/entrypoint.sh /entrypoint.sh
#COPY /third_party/swagger-ui /swagger-ui


# need ansible to deploy shadowss
#COPY /build/build-image/files/ansible /go/bin/ansible

#COPY /build/build-image/files/vpslicense /go/bin/vpslicense
#RUN chmod +x /go/bin/vpslicense

RUN set -ex \
  && cd /go/src/cloud-keeper/\
  && cp -rf contrib/keys/ /keys/\
  && cp -rf build/build-image/files/entrypoint.sh /entrypoint.sh\
  && cp -rf third_party/swagger-ui /swagger-ui\
  && cp -rf build/build-image/files/ansible /go/bin/ansible\
  && cp -rf build/build-image/files/vpslicense /go/bin/vpslicense\
  && chmod +x /go/bin/vpslicense


RUN cd /go/src/cloud-keeper/ \
	&& cp vendor/* /go/src/ -rf \
	&& cd cmd/vpskeeper/ \
	&& go build -v  -o  /go/bin/vpskeeper \
	&& chmod +x /entrypoint.sh \
	&& rm -rf /go/src/


ENTRYPOINT ["/entrypoint.sh"]
