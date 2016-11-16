#!/bin/bash

CURRENT_DIR=`dirname $(readlink -f $0)`
ROOT_DIR=${CURRENT_DIR%%hack}

InsideRootDir=/go/src/cloud-keeper

AutoGenCmd=update_autogen.sh


docker run -v $ROOT_DIR:$InsideRootDir --rm  -i gcr.io/google_containers/kube-cross:v1.7.1-2 bash -c "cd /go/src/cloud-keeper/hack && ./update_autogen.sh"
