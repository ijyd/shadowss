#!/bin/sh
timestamp=$(date)
echo "$timestamp restart ssservice" >> /var/ssservice.log
service ssservice restart