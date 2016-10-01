# docker image build

```
cd /cloud-keeper/deploy/build-image
docker build --rm -t bjjyd/keeper -f Dockerfile ../../
```


# simple start

```
docker run -e "ETCD_URL=http://192.168.60.100:2379" -e "MYSQL_URL=sspanel:sspanel@tcp(localhost:13306)/sspanel" bjjyd/keeper
```
