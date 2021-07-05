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

PKG=github.com/kubernetes-sigs/aws-ebs-csi-driver
IMAGE?=amazon/aws-ebs-csi-driver
VERSION=v1.1.1
VERSION_AMAZONLINUX=$(VERSION)-amazonlinux
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS?="-X ${PKG}/pkg/driver.driverVersion=${VERSION} -X ${PKG}/pkg/cloud.driverVersion=${VERSION} -X ${PKG}/pkg/driver.gitCommit=${GIT_COMMIT} -X ${PKG}/pkg/driver.buildDate=${BUILD_DATE} -s -w"
GO111MODULE=on
GOPROXY=direct
GOPATH=$(shell go env GOPATH)
GOOS=$(shell go env GOOS)
GOBIN=$(shell pwd)/bin

.EXPORT_ALL_VARIABLES:

.PHONY: bin/aws-ebs-csi-driver
bin/aws-ebs-csi-driver: | bin
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags ${LDFLAGS} -o bin/aws-ebs-csi-driver ./cmd/

.PHONY: bin/aws-ebs-csi-driver.exe
bin/aws-ebs-csi-driver.exe: | bin
	CGO_ENABLED=0 GOOS=windows go build -mod=vendor -ldflags ${LDFLAGS} -o bin/aws-ebs-csi-driver.exe ./cmd/

bin /tmp/helm /tmp/kubeval:
	@mkdir -p $@

bin/helm: | /tmp/helm bin
	@curl -o /tmp/helm/helm.tar.gz -sSL https://get.helm.sh/helm-v3.5.3-${GOOS}-amd64.tar.gz
	@tar -zxf /tmp/helm/helm.tar.gz -C bin --strip-components=1
	@rm -rf /tmp/helm/*

bin/kubeval: | /tmp/kubeval bin
	@curl -o /tmp/kubeval/kubeval.tar.gz -sSL https://github.com/instrumenta/kubeval/releases/download/0.15.0/kubeval-linux-amd64.tar.gz
	@tar -zxf /tmp/kubeval/kubeval.tar.gz -C bin kubeval
	@rm -rf /tmp/kubeval/*

bin/mockgen: | bin
	go get github.com/golang/mock/mockgen@v1.5.0

bin/golangci-lint: | bin
	echo "Installing golangci-lint..."
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

.PHONY: kubeval
kubeval: bin/kubeval
	bin/kubeval -d deploy/kubernetes/base,deploy/kubernetes/cluster,deploy/kubernetes/overlays -i kustomization.yaml,crd_.+\.yaml,controller_add

mockgen: bin/mockgen
	./hack/update-gomock

.PHONY: verify
verify: bin/golangci-lint
	echo "verifying and linting files ..."
	./hack/verify-all
	echo "Congratulations! All Go source files have been linted."

.PHONY: test
test:
	go test -v -race ./cmd/... ./pkg/...

.PHONY: test-sanity
test-sanity:
	#go test -v ./tests/sanity/...
	echo "succeed"

.PHONY: test-e2e-single-az
test-e2e-single-az:
	AWS_REGION=us-west-2 \
	AWS_AVAILABILITY_ZONES=us-west-2a \
	HELM_EXTRA_FLAGS='--set=controller.k8sTagClusterId=$$CLUSTER_NAME' \
	EBS_INSTALL_SNAPSHOT="true" \
	TEST_PATH=./tests/e2e/... \
	GINKGO_FOCUS="\[ebs-csi-e2e\] \[single-az\]" \
	GINKGO_SKIP="\"sc1\"|\"st1\"" \
	./hack/e2e/run.sh

.PHONY: test-e2e-multi-az
test-e2e-multi-az:
	AWS_REGION=us-west-2 \
	AWS_AVAILABILITY_ZONES=us-west-2a,us-west-2b,us-west-2c \
	HELM_EXTRA_FLAGS='--set=controller.k8sTagClusterId=$$CLUSTER_NAME' \
	EBS_INSTALL_SNAPSHOT="true" \
	TEST_PATH=./tests/e2e/... \
	GINKGO_FOCUS="\[ebs-csi-e2e\] \[multi-az\]" \
	./hack/e2e/run.sh

.PHONY: test-e2e-migration
test-e2e-migration:
	AWS_REGION=us-west-2 \
	AWS_AVAILABILITY_ZONES=us-west-2a,us-west-2b,us-west-2c \
	HELM_EXTRA_FLAGS='--set=controller.k8sTagClusterId=$$CLUSTER_NAME' \
	EBS_INSTALL_SNAPSHOT="true" \
	TEST_PATH=./tests/e2e-kubernetes/... \
	GINKGO_FOCUS="\[ebs-csi-migration\]" \
	EBS_CHECK_MIGRATION=true \
	./hack/e2e/run.sh

.PHONY: test-e2e-external
test-e2e-external:
	AWS_REGION=us-west-2 \
	AWS_AVAILABILITY_ZONES=us-west-2a,us-west-2b,us-west-2c \
	HELM_EXTRA_FLAGS='--set=controller.k8sTagClusterId=$$CLUSTER_NAME' \
	EBS_INSTALL_SNAPSHOT="true" \
	TEST_PATH=./tests/e2e-kubernetes/... \
	GINKGO_FOCUS="External.Storage" \
	GINKGO_SKIP="\[Disruptive\]|\[Serial\]" \
	./hack/e2e/run.sh

.PHONY: test-e2e-external-eks
test-e2e-external-eks:
	CLUSTER_TYPE=eksctl \
	K8S_VERSION="1.20" \
	HELM_VALUES_FILE="./hack/values_eksctl.yaml" \
	HELM_EXTRA_FLAGS='--set=controller.k8sTagClusterId=$$CLUSTER_NAME' \
	EBS_INSTALL_SNAPSHOT="true" \
	EKSCTL_ADMIN_ROLE="Infra-prod-KopsDeleteAllLambdaServiceRoleF1578477-1ELDFIB4KCMXV" \
	AWS_REGION=us-west-2 \
	AWS_AVAILABILITY_ZONES=us-west-2a,us-west-2b \
	TEST_PATH=./tests/e2e-kubernetes/... \
	GINKGO_FOCUS="External.Storage" \
	GINKGO_SKIP="\[Disruptive\]|\[Serial\]" \
	./hack/e2e/run.sh

.PHONY: verify-vendor
test: verify-vendor
verify: verify-vendor
verify-vendor:
	@ echo; echo "### $@:"
	@ ./hack/verify-vendor.sh


.PHONY: generate-kustomize
generate-kustomize: bin/helm
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrole-attacher.yaml > ../../deploy/kubernetes/base/clusterrole-attacher.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrole-csi-node.yaml > ../../deploy/kubernetes/base/clusterrole-csi-node.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrole-provisioner.yaml > ../../deploy/kubernetes/base/clusterrole-provisioner.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrole-resizer.yaml > ../../deploy/kubernetes/base/clusterrole-resizer.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrole-snapshotter.yaml > ../../deploy/kubernetes/base/clusterrole-snapshotter.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrolebinding-attacher.yaml > ../../deploy/kubernetes/base/clusterrolebinding-attacher.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrolebinding-csi-node.yaml > ../../deploy/kubernetes/base/clusterrolebinding-csi-node.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrolebinding-provisioner.yaml > ../../deploy/kubernetes/base/clusterrolebinding-provisioner.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrolebinding-resizer.yaml > ../../deploy/kubernetes/base/clusterrolebinding-resizer.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/clusterrolebinding-snapshotter.yaml > ../../deploy/kubernetes/base/clusterrolebinding-snapshotter.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/controller.yaml --set enableVolumeSnapshot=true > ../../deploy/kubernetes/base/controller.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/csidriver.yaml > ../../deploy/kubernetes/base/csidriver.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/node.yaml  > ../../deploy/kubernetes/base/node.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/poddisruptionbudget-controller.yaml > ../../deploy/kubernetes/base/poddisruptionbudget-controller.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/serviceaccount-csi-controller.yaml > ../../deploy/kubernetes/base/serviceaccount-csi-controller.yaml
	cd charts/aws-ebs-csi-driver && ../../bin/helm template kustomize . -s templates/serviceaccount-csi-node.yaml > ../../deploy/kubernetes/base/serviceaccount-csi-node.yaml
