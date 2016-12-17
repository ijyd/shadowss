#vps keeper

## simple run

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --insecure-port=18088  --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(localhost:13306)/sspanel" --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="http://192.168.60.100:2379"
```

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --secure-port=18088 --tls-cert-file="../../contrib/keys/server.crt" --tls-private-key-file="../../contrib/keys/server.key" --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="192.168.60.128:2379" --etcd-certfile="../../contrib/keys/client.pem" --etcd-keyfile="../../contrib/keys/client-key.pem" --etcd-cafile="../../contrib/keys/ca.pem" --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(192.168.60.100:23306)/sspanel"
```


# 3.0 start up

```
./vpskeeper --alsologtostderr=true -v=9  --tls-cert-file=../../contrib/keys/server.crt --tls-private-key-file=../../contrib/keys/server.key --secure-port=18090 --insecure-port=18091 --etcd-servers="192.168.60.128:2379" --etcd-certfile="/home/seanchann/bin/etcd/auth/etcdclient.pem" --etcd-keyfile="/home/seanchann/bin/etcd/auth/etcdclient-key.pem" --etcd-cafile="/home/seanchann/bin/etcd/auth/ca.pem" --enable-swagger-ui  --mysql-servers="sspanel:sspanel@tcp(192.168.60.100:23306)/sspanel"  --aws-region="us-west-2" --aws-table="vpsdevelopment" --authorization-mode="ABAC" --authorization-policy-file=../../contrib/authorizer/policy_file.json --anonymous-auth=false
```
