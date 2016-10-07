#vps keeper

## simple run

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --insecure-port=18088  --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(localhost:13306)/sspanel" --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="http://192.168.60.100:2379"
```

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --secure-port=18088 --tls-cert-file="../../keys/server.crt" --tls-private-key-file="../../keys/server.key"  --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="http://192.168.60.128:2379" --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(192.168.60.100:23306)/sspanel"
```

## run on aliyun

```
sudo ./vpskeeper --alsologtostderr=true -v=6  --secure-port=18088 --tls-cert-file="./keys/server.crt" --tls-private-key-file="./keys/server.key"  --swagger-path="../../third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="http://172.22.0.2:2379" --storage-type="mysql"  --server-list="d15e047257e6a6:ee7021982796bb4cd3d750f655db85fe@tcp(vps.cjwhwe9kqqzt.us-west-2.rds.amazonaws.com:3306)/vps"

sudo ./vpskeeper --alsologtostderr=true -v=6  --insecure-port=18088  --swagger-path="../../third_party/swagger-ui" --etcd-servers="http://172.22.0.2:2379" --storage-type="mysql"  --server-list="d15e047257e6a6:ee7021982796bb4cd3d750f655db85fe@tcp(vps.cjwhwe9kqqzt.us-west-2.rds.amazonaws.com:3306)/vps"
```


## pagination
格式:  `{URL}?pagination=page=1,perPage=2`

`paga` 指定输出第几页。

`perPage` 指定每一页包含的资源条数。dock

当访问的页面超出范围，则返回全部资源

*Http Header*

在分页请求的HTTP应答包中， 系统会在头域填加`Link`，作为前后一页和最后一页做操作提示：
rel=prev提示前一页页码，rel=next提示后一页页码，rel=last提示最后一页的页码。

例如:
> Link: /api/v1beta1/namespace/default/users?pagination=page=1,perPage=1; rel= **prev** ,/api/v1beta1/users?pagination=page=3,perPage=1; rel= **next** ,/api/v1beta1/users?pagination=page=5,perPage=1; rel= **last**


```
etcdctl --endpoints="192.168.60.100:2379" get --prefix=true "/registry"
./etcdctl --endpoints="172.22.0.2:2379" get --prefix=true "/registry"
```
