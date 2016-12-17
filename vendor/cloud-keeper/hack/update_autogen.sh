#!/bin/bash

CURRENT_DIR=`dirname $(readlink -f $0)`
ROOT_DIR=${CURRENT_DIR%%hack}

ConversionGen=conversion-gen
DeepcopyGen=deepcopy-gen
DefaulterGen=defaulter-gen
GoHeaderFile=$ROOT_DIR/contrib/gengo/boilerplate/boilerplate.go.txt


hash $ConversionGen &> /dev/null
if [ $? -eq 1 ]; then
  go install  cloud-keeper/cmd/libs/go2idl/conversion-gen/
fi

hash $DeepcopyGen &> /dev/null
if [ $? -eq 1 ]; then
  go install  cloud-keeper/cmd/libs/go2idl/deepcopy-gen/
fi

hash $DefaulterGen &> /dev/null
if [ $? -eq 1 ]; then
  go install  cloud-keeper/cmd/libs/go2idl/defaulter-gen/
fi

function deepcopy {
  #for api
  $DeepcopyGen --bounding-dirs="gofreezer/pkg/api"  --input-dirs="cloud-keeper/pkg/api" --go-header-file=$GoHeaderFile
  $DeepcopyGen --input-dirs="cloud-keeper/pkg/api/v1" --go-header-file=$GoHeaderFile

  #for batch
  $DeepcopyGen --bounding-dirs="gofreezer/pkg/api,cloud-keeper/pkg/api" --input-dirs="cloud-keeper/pkg/apis/batch" --go-header-file=$GoHeaderFile
  $DeepcopyGen --bounding-dirs="cloud-keeper/pkg/api/v1" --input-dirs="cloud-keeper/pkg/apis/batch/v1alpha1" --go-header-file=$GoHeaderFile


  #for abacpolicys
  $DeepcopyGen --bounding-dirs="gofreezer/pkg/api,apistack/pkg/apis/abac" --input-dirs="cloud-keeper/pkg/apis/abacpolicys" --go-header-file=$GoHeaderFile
  $DeepcopyGen --bounding-dirs="apistack/pkg/apis/abac/v1beta1" --input-dirs="cloud-keeper/pkg/apis/abacpolicys/v1beta1" --go-header-file=$GoHeaderFile
}

function conversion {
  #for conversion v1 to api
  $ConversionGen --extra-peer-dirs="gofreezer/pkg/api" --input-dirs="cloud-keeper/pkg/api/v1" --go-header-file=$GoHeaderFile

  #for apis/batch
  $ConversionGen --extra-peer-dirs="gofreezer/pkg/api" --input-dirs="cloud-keeper/pkg/apis/batch/v1alpha1" --go-header-file=$GoHeaderFile


  #for apis/abacpolicys
  $ConversionGen --extra-peer-dirs="gofreezer/pkg/api,apistack/pkg/apis/abac" --input-dirs="cloud-keeper/pkg/apis/abacpolicys/v1beta1" --go-header-file=$GoHeaderFile
}


deepcopy
conversion
