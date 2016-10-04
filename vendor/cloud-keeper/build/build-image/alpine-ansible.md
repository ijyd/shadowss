## downgrade dopy
```
# apk add --no-cache py-setuptools
# easy_install-2.7 pip
# pip install 'dopy>=0.3.5,<=0.3.5'
```

## config vultr module
```
# mkdir /etc/ansible/
# touch /etc/ansible/ansible.cfg
# echo "[defaults]" > /etc/ansible/ansible.cfg
# echo "library        = /go/bin/ansible/vultr/libraray/vultr" >>  /etc/ansible/ansible.cfg
```
