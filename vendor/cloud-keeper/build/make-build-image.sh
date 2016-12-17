#!/bin/bash

source common.sh


build_image() {
  dockerfile=$ROOT/build/build-image/Dockerfile
  build_context=$ROOT
  image_tag=dockerhub.bj-jyd.cn/library/vpskeeper

  apistack::version::get_version_vars
  APISTACK_GO_PACKAGE="cloud-keeper/vendor/apistack"
  LDFLAGS=$(apistack::version::ldflags)


  local -r image=$1
  local -r context_dir=$2
  local -r pull="${3:-true}"
  local -r tag="${image_tag}:latest"
  # tag image with gitversion
  local -r gitversiontag="${image_tag}:${APISTACK_GIT_VERSION}"
  local -ra build_cmd=(docker build --build-arg "LDFLAGS=$LDFLAGS"  --rm -t "${gitversiontag}" -f "$dockerfile" "${build_context}")

  local docker_output
  docker_output=$("${build_cmd[@]}")


  # return git version image name
  echo $gitversiontag
}


build_image
