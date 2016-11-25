#!/bin/sh
timestamp=$(date)
echo "$timestamp restart ssservice" >> /var/ssservice.log
/usr/bin/service ssservice restart
