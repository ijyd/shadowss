# build shadowsocket with docker

## build  image for shadowsocket server

```
cd /shadowss/build/build-image
docker build --rm -t ijyd/shadowss -f Dockerfile ../../
```
## startup shadowsocket server

```
docker run --rm -it ijyd/shadowss /bin/bash
docker run -d -p 20000-20100:20000-20100  ijyd/shadowss
```
