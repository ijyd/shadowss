# shadowssctl

管理用户以及节点服务器

## 调试

```
./shadowssctl --alsologtostderr=true -v=6 --host=192.168.60.143 --port=18088 --storage-type="mysql"  --server-list="sspanel:sspanel@tcp(localhost:13306)/sspanel" --swagger-path="../../third_party/swagger-ui"
```


### API访问权限


- token

添加token到http的头中：

```
Authorization: Bearer d0f67631b4e426ae
```


3f5b7306bd64f36dd935fd890808409c
