#!/bin/bash

# Copyright 2019 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euxo pipefail

loudecho() {
  echo "###"
  echo "## ${1}"
  echo "#"
}

loudecho "Register gcloud as a Docker credential helper."
# Required for "docker buildx build --push".
# See https://github.com/kubernetes-csi/csi-release-tools/blob/master/prow.sh#L1243
gcloud auth configure-docker

loudecho "Set up Docker Buildx"
# See https://github.com/docker/setup-buildx-action
# and https://github.com/kubernetes-csi/csi-release-tools/blob/master/build.make#L132
export DOCKER_CLI_EXPERIMENTAL=enabled
trap "docker buildx rm multiarchimage-buildertest" EXIT
docker buildx create --use --name multiarchimage-buildertest

loudecho "Set up QEMU"
# See https://github.com/docker/setup-qemu-action
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

loudecho "Push manifest list containing debian and windows based images to GCR"
export REGISTRY=$REGISTRY_NAME
export TAG=$GIT_TAG
export ALL_OSVERSION_linux=debian
make all-push

loudecho "Push manifest containing amazon based image to GCR"
export REGISTRY=$REGISTRY_NAME
export TAG=$GIT_TAG-amazonlinux
export ALL_OSVERSION_linux=amazon
export ALL_OS=linux
make all-push
