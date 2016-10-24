#!/bin/sh
timestamp=$(date)
echo "$timestamp restart ssservice" >> /var/ssservice.log
/usr/sbin/service ssservice restart