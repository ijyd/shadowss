#!/bin/bash

CURRENT_DIR=`dirname $(readlink -f $0)`
ROOT_DIR=${CURRENT_DIR%%examples/apiserver/hack}

InsideRootDir=/go/src/apistack

AutoGenCmd=update_autogen.sh


docker run -v $ROOT_DIR:$InsideRootDir --rm  -i gcr.io/google_containers/kube-cross:v1.7.1-2 bash -c "cd $InsideRootDir/examples/apiserver/hack && ./update_autogen.sh"
