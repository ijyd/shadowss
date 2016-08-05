# shadowsocks libev client test

shadowss --alsologtostderr=true --config-file="/home/seanchann/pgitlab/cloud/src/shadowsocks-go/sample-config/server-multi-port.json" -v=6


go run main.go --alsologtostderr=true --config-file="/home/seanchann/pgitlab/cloud/src/shadowsocks-go/sample-config/server-multi-port.json" -v=6

https://github.com/EasyPi/docker-shadowsocks-libev/blob/master/docker-compose.yml

shadowsocks-server -p 18387 -k barfoo -m aes-128-cfb -c config.json -t 60
```
./ss-local -s 47.89.189.237 -p 18087 -k 44c096c8e1e5fd7a85ea0bbe4623778d -m aes-256-cfb -l 11080 -b 0.0.0.0
./ss-local -s 192.168.137.191 -p 18387 -k barfoo -m aes-128-cfb -l 11080 -b 0.0.0.0
```
"host":"0.0.0.0",
"port":18387,
"encrypt":"aes-128-cfb",
"password":"barfoo",
"enableOTA":false,
"timeout":60
