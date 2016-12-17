#!/bin/bash

CURRENT_DIR=`dirname $(readlink -f $0)`
ROOT=${CURRENT_DIR%%build}

cd $ROOT
API_ROOT=$ROOT
source vendor/apistack/hack/libs/version.sh

#    import from apistack
#    APISTACK_GIT_COMMIT - The git commit id corresponding to this
#          source code.
#    APISTACK_GIT_TREE_STATE - "clean" indicates no changes since the git commit id
#        "dirty" indicates source code changes after the git commit id
#    APISTACK_GIT_VERSION - "vX.Y" used to indicate the last release version.
#    APISTACK_GIT_MAJOR - The major part of the version
#    APISTACK_GIT_MINOR - The minor component of the version

function get_version_build_ldflags() {
  apistack::version::get_version_vars
  echo ${APISTACK_GIT_VERSION}
  #apistack::version::save_version_vars .version.txt

  APISTACK_GO_PACKAGE="cloud-keeper/vendor/apistack"
  ldflags=$(apistack::version::ldflags)
  echo $ldflags
}
