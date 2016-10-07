## Build Ansible image

```
$ docker build -t ansible -f Dockerfile .
```

## Run 
```
$ docker run -it --name ansible --rm -v /home/scao/work/vps/ansible:/home/  ansible /bin/bash
```