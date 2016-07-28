# build shadowsocket with docker

## build  image for shadowsocket server

```
cd /shadowsocket-go/build/build-image
docker build --rm -t bjjyd/shadowsocket-go -f Dockerfile ../../
```
## startup shadowsocket server

```
docker run --rm -it bjjyd/shadowsocket-go /bin/bash
docker run -d -p 18387-18388:8387-8388  bjjyd/shadowsocket-go
```
