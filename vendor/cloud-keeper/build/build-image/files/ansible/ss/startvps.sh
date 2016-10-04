#!/bin/sh

pid=$(ps -ef | grep shadowss  | grep sspanel | awk '{print $2}')

if [ ! -z  $pid ]
then
 kill -9 $pid
fi

project_path=/home/shadowsocks-node

$project_path/shadowss \
  --config-file="$project_path/server-multi-port.json" \
  --enable-udp-relay \
  --log-dir="/var/log/shadowss" \
  --storage-type="mysql" \
  --sync-user-interval=20 \
  --server-list="sspanel:sspanel@tcp(47.89.189.237:13306)/sspanel" &

