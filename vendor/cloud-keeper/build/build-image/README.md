# docker image build

```
cd /cloud-keeper/deploy/build-image
docker build --rm -t bjjyd/vpskeeper -f Dockerfile ../../
```


# simple start

first get hwcode for license

```
docker run  -it --privileged -v /dev/mem:/dev/mem bjjyd/vpskeeper /go/bin/vpslicense gencode
```


```
docker run  --privileged -v /home/seanchann/pgitlab/cloud/src/cloud-keeper/.license:/go/bin/.license -v /dev/mem:/dev/mem  -e "ETCD_URL=http://192.168.60.100:2379" -e "MYSQL_URL=sspanel:sspanel@tcp(192.168.60.128:13306)/sspanel" bjjyd/vpskeeper
```


run witd devel
```
docker run  -it --privileged -v /home/seanchann/pgitlab/cloud/src/cloud-keeper/.license:/go/bin/.license -v /dev/mem:/dev/mem -v /home/seanchann/pgitlab/cloud/src/cloud-keeper/:/go/src/cloud-keeper/ -e "ETCD_URL=http://192.168.60.100:2379" -e "MYSQL_URL=sspanel:sspanel@tcp(localhost:13306)/sspanel" bjjyd/vpskeeper /bin/bash
```
