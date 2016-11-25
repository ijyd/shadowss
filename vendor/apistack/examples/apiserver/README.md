# API Server

This code here is to examplify what it takes to write your own API server.

To start this example api server, run:

```
$ go run examples/apiserver/server/main.go
```

```
./apiserver --alsologtostderr=true -v=9 --tls-cert-file=../../contrib/keys/server.crt --tls-private-key-file=../../contrib/keys/server.key  --secure-port=18090 --insecure-port=18091 --etcd-servers="192.168.60.100:2379" --etcd-certfile="/home/seanchann/bin/etcd/files/client.pem" --etcd-keyfile="/home/seanchann/bin/etcd/files/client-key.pem" --etcd-cafile="/home/seanchann/bin/etcd/files/ca.pem" --enable-swagger-ui --mysql-servers="sspanel:sspanel@tcp(192.168.60.100:23306)/sspanel" --anonymous-auth=false --aws-region="us-west-2" --aws-table="vpstest" --aws-cred-accessid="AKIAI2AE5UHGZBW2WV7Q" --aws-cred-accesskey="RXt4mxyX7wlXGgHg01QUMoK+ybb0GePbrGQxTM7p"
```





*if enable swagger ui, you must place the swaager folder in /swagger-ui*
