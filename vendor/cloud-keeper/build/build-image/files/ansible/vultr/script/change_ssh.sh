#!/bin/sh
oldip=$1
newip=$2

if [ -z $1 ] || [ -z $2 ]
then
   echo "-------------------------------------------"
   echo "--------./change_ssh.sh oldip newip--------"
   exit
fi

sed -i "/$oldip/d" /etc/ansible/hosts
echo "$newip" >> /etc/ansible/hosts

sed -i "s/$oldip/$newip/g" /root/.ssh/config
