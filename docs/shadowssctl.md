# shadowssctl

管理用户以及节点服务器

## 调试

```
./shadowssctl --alsologtostderr=true -v=6 --host=192.168.60.143 --port=18088 --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(localhost:13306)/sspanel" --swagger-path="../../third_party/swagger-ui"
```


### Logins


- token

添加token到http的头中：

```
Authorization: Bearer d0f67631b4e426ae
```

- Post

```
{
  "kind": "Login",
  "apiVersion": "v1",
  "metadata":{
    "name": "1"
  },
  "spec":{
      "authname": "78a3511574bc@gmail.com",
      "auth": "0a21be745687cc463407680ceac64564"
  }
}
```

###　apiserver

- Get

```
{
  "kind": "APIServer",
  "apiVersion": "v1",
  "metadata":{
    "name": "1"
  },
  "spec":{
      "server": [
      {
        "host":"1.1.1.1",
        "port":10006,
      },
      ]
  }
}
```

###　Node

- Get

```
{
  "kind": "Node",
  "apiVersion": "v1",
  "metadata":{
    "name": "1"
  },
  "spec":{
      "server": [
      {
        "host":"1.1.1.1",
        "status":true
      },
      ],
      "account":{
        "id": 1,
        "port":8154,
        "method":"aes-256-cfb"
      }
  }
}
```
