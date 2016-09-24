# shadowsocks

## 协议规范


### tcp relay


```
/*
 * Shadowsocks TCP Relay Header:
 *
 *    +------+----------+----------+----------------+
 *    | ATYP | DST.ADDR | DST.PORT |    HMAC-SHA1   |
 *    +------+----------+----------+----------------+
 *    |  1   | Variable |    2     |      10        |
 *    +------+----------+----------+----------------+
 *
 *    If ATYP & ONETIMEAUTH_FLAG(0x10) == 1, Authentication (HMAC-SHA1) is enabled.
 *
 *    The key of HMAC-SHA1 is (IV + KEY) and the input is the whole header.
 *    The output of HMAC-SHA is truncated to 10 bytes (leftmost bits).
 */

/*
 * Shadowsocks Request's Chunk Authentication for TCP Relay's payload
 * (No chunk authentication for response's payload):
 *
 *    +------+-----------+-------------+------+
 *    | LEN  | HMAC-SHA1 |    DATA     |      ...
 *    +------+-----------+-------------+------+
 *    |  2   |    10     |  Variable   |      ...
 *    +------+-----------+-------------+------+
 *
 *    The key of HMAC-SHA1 is (IV + CHUNK ID)
 *    The output of HMAC-SHA is truncated to 10 bytes (leftmost bits).
 */
```

#### tcp server response

- first resp package

```
 /*
  *
  *    +-------------+-----------+
  *    |     IV      |    DATA   |
  *    +------+------+-----------+
  *    |  Variable   |  Variable |
  *    +-------------+-----------+
  *
  *    
  *   
  */

```


- follow resp package

*复用第一个应答的iv进行解包，server部分采用第一个应答生成的IV对后续数据进行加密*

```
 /*
  *
  *    +-----------+
  *    |    DATA   |
  *    +-----------+
  *    |  Variable |
  *    +-----------+
  *
  *    
  *   
  */

```


### udp relay

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


- shadowsocks UDP Request (before encrypted)
  - ATYP: ATYP & 0x10标识为支持hmac
  - DATA: 直接转发给远端的dst addr以及端口

```
+------+----------+----------+----------+-------------+
| ATYP | DST.ADDR | DST.PORT |   DATA   |  HMAC-SHA1  |
+------+----------+----------+----------+-------------+
|  1   | Variable |    2     | Variable |     10      |
+------+----------+----------+----------+-------------+

If ATYP & ONETIMEAUTH_FLAG(0x10) == 1, Authentication (HMAC-SHA1) is enabled.
The key of HMAC-SHA1 is (IV + KEY) (KEY is a password bytes to key like call evpBytesToKey) and the input is the whole packet(contains ATYP DST.ADDR DST.PORT DATA fields).
The output of HMAC-SHA is truncated to 10 bytes (leftmost bits).
```

- shadowsocks UDP Response (before encrypted)
  - ATYP为请求的应答类型
  - 应答头部分应保持与请求头的一致，除了ATYP的HMAC标识部分

```
+------+----------+----------+----------+
| ATYP | DST.ADDR | DST.PORT |   DATA   |
+------+----------+----------+----------+
|  1   | Variable |    2     | Variable |
+------+----------+----------+----------+
```

- shadowsocks UDP Request and Response (after encrypted)

```
+-------+--------------+
|   IV  |    PAYLOAD   |
+-------+--------------+
| Fixed |   Variable   |
+-------+--------------+
```

## run server

```
./shadowss --alsologtostderr=true --config-file="./server-multi-port.json" -v=6 --enable-udp-relay --storage-backend="etcd3" --etcd-servers="http://192.168.60.100:2379"
```

## run client


```
{
    "server": "47.89.189.237",
    "server_port": 18387,
    "local_address": "0.0.0.0",
    "local_port": 1080,
    "password": "713cfeb9dd0faf61609bf1c54bb3e766",
    "timeout": 60,
    "method": "aes-256-cfb"
}


/usr/bin/ss-tunnel -c /etc/shadowsocks.json -A -u -l 53000 -L 8.8.4.4:53
/usr/bin/ss-redir -c /etc/shadowsocks.json -A -u -v

./ss-local -c ./shadowsocks.json -A -u -v
curl --socks5-hostname 127.0.0.1:1080 www.google.com

dig @192.168.1.1 www.google.com

./shadowss --alsologtostderr=true --config-file="./server-multi-port.json" -v=6 --enable-udp-relay

"host":"0.0.0.0",
"port":38387,
"encrypt":"aes-256-cfb",
"password":"9785e42da78df215dcfdbc0bcc294a9f",
"enableOTA":true,
"timeout":60

```



## TODO

- cli

```
https://github.com/urfave/cli
```


- fail2ban

```
http://atjason.com/IT/ss_fail2ban.html
```
