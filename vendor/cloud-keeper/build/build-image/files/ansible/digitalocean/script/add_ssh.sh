#!/bin/sh
newip=$1

if [ -z $1 ] 
then
   echo "-------------------------------------------"
   echo "--------./change_ssh.sh newip--------"
   exit
fi

echo "$newip" >> /etc/ansible/hosts

echo "" >> /root/.ssh/config
echo "Host $newip" >> /root/.ssh/config
echo " HostName $newip" >> /root/.ssh/config
echo " User jyd" >> /root/.ssh/config
echo " IdentityFile /home/id_rsa" >> /root/.ssh/config
