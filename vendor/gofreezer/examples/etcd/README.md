# run

```
./etcd --alsologtostderr=true -v=6 --host=192.168.60.107 --port=18088  --swagger-path="../../../cloud-keeper/third_party/swagger-ui" --storage-backend="etcd3" --etcd-servers="http://192.168.60.100:2379"
```

## etcdctl

```
etcdctl --endpoints="192.168.60.100:2379" get --prefix=true "/registry"
```
