# DigitalOcean Ansible

## 1. update dopy
```
# pip install 'dopy>=0.3.5,<=0.3.5'
```

### 2. Create droplet
```
# ansible-playbook deploy/create_droplet_v2.yml  
```

### 3. Delete droplet
```
# ansible-playbook deploy/delete_droplet_v2.yml  
```

### 4. Depoly and start shadowss
```
# ansible-playbook deploy/deploy_ss.yml  
```

