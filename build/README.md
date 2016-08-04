# build shadowsocket with docker

## build  image for shadowsocket server

```
cd /shadowsocket-go/build/build-image
docker build --rm -t bjjyd/shadowsocket-go -f Dockerfile ../../
```
## startup shadowsocket server

```
docker run --rm -it bjjyd/shadowsocket-go /bin/bash
docker run -d -p 20000-20100:20000-20100  bjjyd/shadowsocket-go
```
