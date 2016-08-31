# tunnel


## ss-tunel dns请求流程

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
