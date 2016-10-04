#!/bin/sh
newip=$1

if [ -z $1 ] 
then
   echo "-------------------------------------------"
   echo "--------./change_ssh.sh newip--------"
   exit
fi

echo "$newip" >> /etc/ansible/hosts

echo "" >> ~/.ssh/config
echo "Host $newip" >> ~/.ssh/config
echo " HostName $newip" >> ~/.ssh/config
echo " User jyd" >> ~/.ssh/config
echo " IdentityFile /home/id_rsa" >> ~/.ssh/config
