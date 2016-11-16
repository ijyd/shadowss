#!/bin/bash

CURRENT_DIR=`dirname $(readlink -f $0)`
ROOT_DIR=${CURRENT_DIR%%hack}

ConversionGen=conversion-gen
DeepcopyGen=deepcopy-gen
DefaulterGen=defaulter-gen
GoHeaderFile=$ROOT_DIR/contrib/gengo/boilerplate/boilerplate.go.txt

cp -rf /go/src/apistack/examples/apiserver/hack/vendor /go/src/apistack/examples/apiserver/

hash $ConversionGen &> /dev/null
if [ $? -eq 1 ]; then
  go install  apistack/examples/apiserver/cmd/libs/go2idl/conversion-gen/
fi

hash $DeepcopyGen &> /dev/null
if [ $? -eq 1 ]; then
  go install  apistack/examples/apiserver/cmd/libs/go2idl/deepcopy-gen/
fi

hash $DefaulterGen &> /dev/null
if [ $? -eq 1 ]; then
  go install  apistack/examples/apiserver/cmd/libs/go2idl/defaulter-gen/
fi

function deepcopy {
  #for api

  $DeepcopyGen --bounding-dirs="gofreezer/pkg/api" --input-dirs="apistack/examples/apiserver/pkg/api" --go-header-file=$GoHeaderFile
  $DeepcopyGen --input-dirs="apistack/examples/apiserver/pkg/api/v1" --go-header-file=$GoHeaderFile

  #for testtype
  $DeepcopyGen --bounding-dirs="gofreezer/pkg/api,cloud-keeper/pkg/api" --input-dirs="apistack/examples/apiserver/pkg/apis/testgroup" --go-header-file=$GoHeaderFile
  $DeepcopyGen --bounding-dirs="apistack/examples/apiserver/pkg/api/v1" --input-dirs="apistack/examples/apiserver/pkg/apis/testgroup/v1" --go-header-file=$GoHeaderFile
}

function conversion {
  #for conversion v1 to api
  $ConversionGen --extra-peer-dirs="gofreezer/pkg/api" --input-dirs="apistack/examples/apiserver/pkg/api/v1" --go-header-file=$GoHeaderFile

  #for apis/batch
  $ConversionGen --extra-peer-dirs="gofreezer/pkg/api" --input-dirs="apistack/examples/apiserver/pkg/apis/testgroup/v1" --go-header-file=$GoHeaderFile
}


deepcopy
conversion

rm -rf  /go/src/apistack/examples/apiserver/vendor
