#!/bin/sh

set -e

#/go/bin/shadowss -c /conf/config.json &

/go/bin/shadowss-mu -config_path /conf/ &
