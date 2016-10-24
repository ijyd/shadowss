#!/bin/sh

set -e


if [ "$1" = '/go/bin/vpslicense' ]; then
  echo "$@"
  exec "$@"
  exit
fi

if [  "X$ETCD_URL" == "X" ] || [  "X$MYSQL_URL" == "X" ]; then
  echo "not found etcd and mysql env, shutdown"
  exit
fi


DefaultARGS='--tls-cert-file=/keys/server.crt --tls-private-key-file=/keys/server.key  --swagger-path=/swagger-ui --storage-backend=etcd3 --etcd-servers='$ETCD_URL'  --etcd-certfile=/keys/client.pem --etcd-keyfile=/keys/client-key.pem --etcd-cafile=/keys/ca.pem  --storage-type=mysql  --server-list='$MYSQL_URL''

# first args is '-' attch with user args
# like as  ./entrypoint.sh -v=6 we will appead default args
if [ "${1:0:1}" = '-' ]; then
  set -- $DefaultARGS "$@"
	set -- /go/bin/vpskeeper  "$@"
fi

# support cmd vpskeeper
# we appand default args
if [ "$1" = 'vpskeeper' ]; then
	numa='numactl --interleave=all'
	if $numa true &> /dev/null; then
		set -- $numa "$@"
	fi

  set -- $DefaultARGS "$@"

  exec  /go/bin/vpskeeper  "$@"
fi

echo "$@"
exec "$@"
