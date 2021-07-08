module github.com/kubernetes-sigs/aws-ebs-csi-driver

require (
	github.com/aws/aws-sdk-go v1.35.37
	github.com/container-storage-interface/spec v1.3.0
	github.com/elazarl/goproxy v0.0.0-20181111060418-2ce16c963a8a // indirect
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/kubernetes-csi/csi-proxy/client v1.0.0-rc.1
	github.com/kubernetes-csi/csi-test v2.0.0+incompatible
	github.com/kubernetes-csi/external-snapshotter/client/v4 v4.0.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.7.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073
	google.golang.org/grpc v1.34.0
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.0
	k8s.io/client-go v1.21.0
	k8s.io/component-base v0.21.0
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.8.0
	k8s.io/kubernetes v1.21.0
	k8s.io/mount-utils v0.21.0
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920
)

replace (
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.0
	k8s.io/apiserver => k8s.io/apiserver v0.21.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.0
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.0
	k8s.io/code-generator => k8s.io/code-generator v0.21.0
	k8s.io/component-base => k8s.io/component-base v0.21.0
	k8s.io/component-helpers => k8s.io/component-helpers v0.21.0
	k8s.io/controller-manager => k8s.io/controller-manager v0.21.0
	k8s.io/cri-api => k8s.io/cri-api v0.17.4-beta.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.0
	k8s.io/kubectl => k8s.io/kubectl v0.21.0
	k8s.io/kubelet => k8s.io/kubelet v0.21.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.0
	k8s.io/metrics => k8s.io/metrics v0.21.0
	k8s.io/mount-utils => k8s.io/mount-utils v0.21.0-beta.0
	k8s.io/node-api => k8s.io/node-api v0.21.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.0
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.21.0
	k8s.io/sample-controller => k8s.io/sample-controller v0.21.0
	vbom.ml/util => github.com/fvbommel/util v0.0.2 // Mitigate https://github.com/fvbommel/util/issues/6
)

go 1.16
