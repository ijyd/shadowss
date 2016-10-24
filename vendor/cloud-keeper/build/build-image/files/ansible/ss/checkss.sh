#!/bin/sh

shadowsspid=$(pgrep shadowss)

if [ -z $shadowsspid ]
then
	date >> /var/checkss.log
	echo "shadowss dead, restart now" >> /var/checkss.log
	/usr/sbin/service ssservice start
fi