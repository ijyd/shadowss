# API Server

This code here is to examplify what it takes to write your own API server.

To start this example api server, run:

```
$ go run examples/apiserver/server/main.go
```

```
sudo ./apiserver --alsologtostderr=true -v=9  --secure-port=18090 --insecure-port=18091 --etcd-servers="192.168.60.100:2379" --etcd-certfile="/home/seanchann/bin/etcd/files/client.pem" --etcd-keyfile="/home/seanchann/bin/etcd/files/client-key.pem" --etcd-cafile="/home/seanchann/bin/etcd/files/ca.pem" --bind-address=192.168.60.128 --advertise-address=192.168.60.128 --enable-swagger-ui --service-cluster-ip-range=192.168.60.0/24 --mysql-servers="sspanel:sspanel@tcp(192.168.60.100:23306)/sspanel" --anonymous-auth=false
```





*if enable swagger ui, you must place the swaager folder in /swagger-ui*
