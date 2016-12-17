#!/bin/bash

build_binary(){
  LDFLAGS=$1
  DIR=$2
  if [[ -n $LDFLAGS ]]; then
    cd $2
    echo $2
    echo $LDFLAGS
    local -ra build_cmd=(go build -ldflags "$LDFLAGS"  -v -o  /go/bin/vpskeeper)
    go_output=$("${build_cmd[@]}")
    echo $go_output
    #go build -ldflags '-X cloud-keeper/vendor/apistack/pkg/version.buildDate=2016-12-10T04:01:19Z -X cloud-keeper/vendor/apistack/pkg/version.gitCommit=75b725aa2fdca48b6367785c7d1976f608c5163f -X cloud-keeper/vendor/apistack/pkg/version.gitTreeState=dirty -X cloud-keeper/vendor/apistack/pkg/version.gitVersion=v3.0.1-alpha.1-dirty-dirty -X cloud-keeper/vendor/apistack/pkg/version.gitMajor=3 -X cloud-keeper/vendor/apistack/pkg/version.gitMinor=0+' -v -o  /go/bin/vpskeeper
  else
    echo "Not build, need specify ldflags for version package"
  fi

}

build_binary "$1" $2
#./build.sh '-X cloud-keeper/vendor/apistack/pkg/version.buildDate=2016-12-10T04:01:19Z -X cloud-keeper/vendor/apistack/pkg/version.gitCommit=75b725aa2fdca48b6367785c7d1976f608c5163f -X cloud-keeper/vendor/apistack/pkg/version.gitTreeState=dirty -X cloud-keeper/vendor/apistack/pkg/version.gitVersion=v3.0.1-alpha.1-dirty-dirty -X cloud-keeper/vendor/apistack/pkg/version.gitMajor=3 -X cloud-keeper/vendor/apistack/pkg/version.gitMinor=0+' /go/src/cloud-keeper/cmd/vpskeeper
