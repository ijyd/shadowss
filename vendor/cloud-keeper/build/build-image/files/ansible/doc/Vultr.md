# Ansible Vultr

## 1. Download Vultr Library

```
  # cd /usr/share/ansible
  # git clone https://github.com/tundrax/ansible-vultr.git vultr
```

## 2. Config ansible.cfg

```
# vi /etc/ansible/ansible.cfg
```

### modify library path

```
library  = /usr/share/ansible/vultr
```

## 3. Disable host_key_checking 

```
# vi /etc/ansible/ansible.cfg

change 
host_key_checking = False
```

## 4. Create node

```
# ansible-playbook deploy/create-node.yml 
```

## 5. Delete node

```
# ansible-playbook deploy/delete-node.yml 
```

## 6. Depoly and Run shadowss 

```
# ansible-playbook deploy/deploy_ss.yml 
```
