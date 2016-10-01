#!/bin/sh

set -e

sudo /go/bin/vpskeeper  --secure-port=18088 \
  --tls-cert-file="/keys/server.crt" --tls-private-key-file="/keys/server.key"\
  --storage-backend="etcd3" --etcd-servers="$ETCD_URL"\
  --storage-type="mysql"  --server-list="MYSQL_URL"
