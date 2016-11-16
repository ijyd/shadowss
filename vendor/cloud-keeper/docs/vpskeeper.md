#vps keeper

## simple run

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --insecure-port=18088  --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(localhost:13306)/sspanel" --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="http://192.168.60.100:2379"
```

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --secure-port=18088 --tls-cert-file="../../contrib/keys/server.crt" --tls-private-key-file="../../contrib/keys/server.key" --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="192.168.60.128:2379" --etcd-certfile="../../contrib/keys/client.pem" --etcd-keyfile="../../contrib/keys/client-key.pem" --etcd-cafile="../../contrib/keys/ca.pem" --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(192.168.60.100:23306)/sspanel"
```

## run on aliyun

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --secure-port=18088 --tls-cert-file="./keys/server.crt" --tls-private-key-file="./keys/server.key"  --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="http://172.22.0.2:2379" --storage-type="mysql"  --server-list="d15e047257e6a6:ee7021982796bb4cd3d750f655db85fe@tcp(vps.cjwhwe9kqqzt.us-west-2.rds.amazonaws.com:3306)/vps"

sudo ./vpskeeper --alsologtostderr=true -v=6  --insecure-port=18088  --swagger-path="../../third_party/swagger-ui" --etcd-servers="http://172.22.0.2:2379" --storage-type="mysql"  --server-list="d15e047257e6a6:ee7021982796bb4cd3d750f655db85fe@tcp(vps.cjwhwe9kqqzt.us-west-2.rds.amazonaws.com:3306)/vps"
```


3.0

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --secure-port=18090 --insecure-port=18091 --etcd-servers="192.168.60.100:2379" --etcd-certfile="/home/seanchann/bin/etcd/files/client.pem" --etcd-keyfile="/home/seanchann/bin/etcd/files/client-key.pem" --etcd-cafile="/home/seanchann/bin/etcd/files/ca.pem"  --bind-address=192.168.60.128 --advertise-address=192.168.60.128 --enable-swagger-ui --service-cluster-ip-range=192.168.60.0/24 --mysql-servers="sspanel:sspanel@tcp(192.168.60.100:23306)/sspanel"  --anonymous-auth=false
```
