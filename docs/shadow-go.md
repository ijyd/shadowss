# shadowsocks

## 协议规范

```
SOCKS5 UDP Request
+----+------+------+----------+----------+----------+
|RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
+----+------+------+----------+----------+----------+
| 2  |  1   |  1   | Variable |    2     | Variable |
+----+------+------+----------+----------+----------+
SOCKS5 UDP Response
|RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
+----+------+------+----------+----------+----------+
+----+------+------+----------+----------+----------+
| 2  |  1   |  1   | Variable |    2     | Variable |
+----+------+------+----------+----------+----------+
```


```
shadowsocks UDP Request (before encrypted)
+------+----------+----------+----------+-------------+
| ATYP | DST.ADDR | DST.PORT |   DATA   |  HMAC-SHA1  |
+------+----------+----------+----------+-------------+
|  1   | Variable |    2     | Variable |     10      |
+------+----------+----------+----------+-------------+

If ATYP & ONETIMEAUTH_FLAG(0x10) == 1, Authentication (HMAC-SHA1) is enabled.
The key of HMAC-SHA1 is (IV + KEY) (KEY is a password bytes to key like call evpBytesToKey) and the input is the whole packet(contains ATYP DST.ADDR DST.PORT DATA fields).
The output of HMAC-SHA is truncated to 10 bytes (leftmost bits).

shadowsocks UDP Response (before encrypted)
+------+----------+----------+----------+
| ATYP | DST.ADDR | DST.PORT |   DATA   |
+------+----------+----------+----------+
|  1   | Variable |    2     | Variable |
+------+----------+----------+----------+

shadowsocks UDP Request and Response (after encrypted)
+-------+--------------+
|   IV  |    PAYLOAD   |
+-------+--------------+
| Fixed |   Variable   |
+-------+--------------+
```

## run server

```
./shadowss --alsologtostderr=true --config-file="./server-multi-port.json" -v=6 --enable-udp-relay --storage-type="mysql" --sync-user-interval=20 --server-list="sspanel:sspanel@tcp(192.168.60.132:13306)/sspanel"
```

test with dig

```
dig @192.168.1.1 www.google.com

```


# 调试问题汇总

- server返回的数据，在客户端做解密的生成的文本数据总是在字节流的末尾多一个字节
- 从端口53000收到数据后发送给60249端口，后面应该还有一个步骤从chinadns的53530到63651端口，然后就会有包从53返回给客户端


# ss-tunel dns请求流程

- client dns 发起请求到op的dnsmasq(53)
- dnsmasq将请求转发到`server=127.0.0.1#53530`
- 53530端口是本机的chinadns的监听端口(`/usr/bin/chinadns -p 53530 -s 114.114.114.114,127.0.0.1:53000 -c /etc/ignore.list -d -m`)
- chinadns配置了两个dns的请求服务器，`114.114.114.114`和`127.0.0.1:53000`
  - chinadns的策略是同时向这个地址发起dns的请求报文，返回的dns地址列表选择为国外网址的地址返回给dnsmasq
- `127.0.0.1:53000`端口的进程就是`ss-tunnel`:`/usr/bin/ss-tunnel -c /etc/oc-go.json -u -l 53000 -L 8.8.4.4:53`
- `ss-tunnel`获取到这个请求后，进行shadowsockets的封包，直接转发给远端的server
- 远端的server，去除shadowsocket的请求头，发起到请求头中的host以及port的请求，请求数据为请求中的data数据块
- 远端获取到relay请求的结果，携带原请求头以及结果，再返回给`ss-tunnel`。返回的数据，其实就是dns的resp的包体
- `ss-tunnel`收到返回请求后，检查如果是一个有效的应答包，就会去除shadows的头字段，把数据部分发给chinadns
- chinadns获取到的就是一个dns的resp，如果是一个有效的dns resp，那么它就会将此响应透传给dnsmasq
- 最终dnsmasq将dns应答发给client

## ss-tunenl

```
/usr/bin/ss-tunnel -c /var/etc/shadowsocks.json -A -u -l 53000 -L 8.8.4.4:53 -f /var/run/ss-tunnel.pid
```
