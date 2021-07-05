[![Build Status](https://travis-ci.org/kubernetes-sigs/aws-ebs-csi-driver.svg?branch=master)](https://travis-ci.org/kubernetes-sigs/aws-ebs-csi-driver)
[![Coverage Status](https://coveralls.io/repos/github/kubernetes-sigs/aws-ebs-csi-driver/badge.svg?branch=master)](https://coveralls.io/github/kubernetes-sigs/aws-ebs-csi-driver?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubernetes-sigs/aws-ebs-csi-driver)](https://goreportcard.com/report/github.com/kubernetes-sigs/aws-ebs-csi-driver)

# Amazon Elastic Block Store (EBS) CSI driver

## Overview

The [Amazon Elastic Block Store](https://aws.amazon.com/ebs/) Container Storage Interface (CSI) Driver provides a [CSI](https://github.com/container-storage-interface/spec/blob/master/spec.md) interface used by Container Orchestrators to manage the lifecycle of Amazon EBS volumes.

## CSI Specification Compatibility Matrix
| AWS EBS CSI Driver \ CSI Version       | v0.3.0| v1.0.0 | v1.1.0 |
|----------------------------------------|-------|--------|--------|
| master branch                          | no    | no     | yes    |
| v1.1.x                                 | no    | no     | yes    |
| v1.0.0                                 | no    | no     | yes    |
| v0.10.x                                | no    | no     | yes    |
| v0.9.x                                 | no    | no     | yes    |
| v0.8.x                                 | no    | no     | yes    |
| v0.7.1                                 | no    | no     | yes    |
| v0.6.0                                 | no    | no     | yes    |
| v0.5.0                                 | no    | no     | yes    |
| v0.4.0                                 | no    | no     | yes    |
| v0.3.0                                 | no    | yes    | no     |
| v0.2.0                                 | no    | yes    | no     |
| v0.1.0                                 | yes   | no     | no     |

## Features
The following CSI gRPC calls are implemented:
* **Controller Service**: CreateVolume, DeleteVolume, ControllerPublishVolume, ControllerUnpublishVolume, ControllerGetCapabilities, ValidateVolumeCapabilities, CreateSnapshot, DeleteSnapshot, ListSnapshots
* **Node Service**: NodeStageVolume, NodeUnstageVolume, NodePublishVolume, NodeUnpublishVolume, NodeGetCapabilities, NodeGetInfo
* **Identity Service**: GetPluginInfo, GetPluginCapabilities, Probe

## CreateVolume Parameters
There are several optional parameters that could be passed into `CreateVolumeRequest.parameters` map, these parameters can be configured in StorageClass, see [example](../examples/kubernetes/storageclass):

| Parameters                  | Values                                 | Default  | Description         |
|-----------------------------|----------------------------------------|----------|---------------------|
| "csi.storage.k8s.io/fstype" | xfs, ext2, ext3, ext4                  | ext4     | File system type that will be formatted during volume creation. This parameter is case sensitive! |
| "type"                      | io1, io2, gp2, gp3, sc1, st1,standard  | gp3*     | EBS volume type     |
| "iopsPerGB"                 |                                        |          | I/O operations per second per GiB. Required when io1 or io2 volume type is specified. If this value multiplied by the size of a requested volume produces a value above the maximum IOPs allowed for the volume type, as documented [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html), AWS will cap the IOPS to maximum supported value. If the value is lower than minimal supported IOPS value per volume, either error is returned (the default behavior) or the value is increased to fit into the supported range when `allowautoiopspergbincrease` is `"true"`.|
| "allowAutoIOPSPerGBIncrease"| true, false                            | false    | When `"true"`, the CSI driver increases IOPS for a volume when `iopsPerGB * <volume size>` is too low to fit into IOPS range supported by AWS. This allows dynamic provisioning to always succeed, even when user specifies too small PVC capacity or `iopsPerGB` value. On the other hand, it may introduce additional costs, as such volumes have higher IOPS than requested in `iopsPerGB`.|
| "iops"                      |                                        | 3000     | I/O operations per second. Only effetive when gp3 volume type is specified. If empty, it will set to 3000 as documented [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html). |
| "throughput"                |                                        | 125      | Throughput in MiB/s. Only effective when gp3 volume type is specified. If empty, it will set to 125MiB/s as documented [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html). |
| "encrypted"                 |                                        |          | Whether the volume should be encrypted or not. Valid values are "true" or "false" |
| "kmsKeyId"                  |                                        |          | The full ARN of the key to use when encrypting the volume. When not specified, the default KMS key is used |

**Notes**:
* `gp3` is currently not supported on outposts. Outpost customers need to use a different type for their volumes.
* Unless explicitly noted, all parameters are case insensitive (e.g. "kmsKeyId", "kmskeyid" and any other combination of upper/lowercase characters can be used).

## Tagging
To help manage volumes in the aws account, CSI driver will automatically add tags to the volumes it manages.

| TagKey                 | TagValue                  | sample                                                              | Description         |
|------------------------|---------------------------|---------------------------------------------------------------------|---------------------|
| CSIVolumeName          | pvcName                   | CSIVolumeName = pvc-a3ab0567-3a48-4608-8cb6-4e3b1485c808            | add to all volumes, for recording associated pvc id and checking if a given volume was already created so that ControllerPublish/CreateVolume is idempotent. |
| CSISnapshotName        | volumeSnapshotContentName | CSISnapshotName = snapcontent-69477690-803b-4d3e-a61a-03c7b2592a76  | add to all snapshots, for recording associated VolumeSnapshot id and checking if a given snapshot was already created                                    |
| ebs.csi.aws.com/cluster| true                      | ebs.csi.aws.com/cluster = true                                      | add to all volumes and snapshots, for allowing users to use a policy to limit csi driver's permission to just the resources it manages.                      |
| kubernetes.io/cluster/X| owned                     | kubernetes.io/cluster/aws-cluster-id-1 = owned                      | add to all volumes and snapshots if k8s-tag-cluster-id argument is set to X.|
| extra-key              | extra-value               | extra-key = extra-value                                             | add to all volumes and snapshots if extraTags argument is set|

## Driver Options
There are couple driver options that can be passed as arguments when starting driver container.

| Option argument             | value sample                                      | default                                             | Description         |
|-----------------------------|---------------------------------------------------|-----------------------------------------------------|---------------------|
| endpoint                    | tcp://127.0.0.1:10000/                            | unix:///var/lib/csi/sockets/pluginproxy/csi.sock    | added to all volumes, for checking if a given volume was already created so that ControllerPublish/CreateVolume is idempotent. |
| volume-attach-limit         | 1,2,3 ...                                         | -1                                                  | Value for the maximum number of volumes attachable per node. If specified, the limit applies to all nodes. If not specified, the value is approximated from the instance type.    |
| extra-tags                  | key1=value1,key2=value2                           |                                                     | Extra tags to attach to each dynamically provisioned resource.|
| k8s-tag-cluster-id          | aws-cluster-id-1                                  |                                                     | ID of the Kubernetes cluster used for tagging provisioned EBS volumes.|
| aws-sdk-debug-log           | true                                              | false                                               | if true, driver will enable the aws sdk debug log level|

# EBS CSI Driver on Kubernetes
Following sections are Kubernetes specific. If you are Kubernetes user, use followings for driver features, installation steps and examples.

## Kubernetes Version Compatibility Matrix
| AWS EBS CSI Driver \ Kubernetes Version| v1.12 | v1.13 | v1.14 | v1.15 | v1.16 | v1.17 | v1.18+ |
|----------------------------------------|-------|-------|-------|-------|-------|-------|-------|
| master branch                          | no    | no+   | no    | no    | no    | yes   | yes   |
| v1.1.0                                 | no    | no+   | no    | no    | no    | yes   | yes   |
| v1.0.0                                 | no    | no+   | no    | no    | no    | yes   | yes   |
| v0.10.x                                | no    | no+   | no    | no    | no    | yes   | yes   |
| v0.9.x                                 | no    | no+   | no    | no    | no    | yes   | yes   |
| v0.8.x                                 | no    | no+   | yes   | yes   | yes   | yes   | yes   |
| v0.7.1                                 | no    | no+   | yes   | yes   | yes   | yes   | yes   |
| v0.6.0                                 | no    | no+   | yes   | yes   | yes   | yes   | yes   |
| v0.5.0                                 | no    | no+   | yes   | yes   | yes   | yes   | yes   |
| v0.4.0                                 | no    | no+   | yes   | yes   | no    | no    | no    |
| v0.3.0                                 | no    | no+   | yes   | no    | no    | no    | no    |
| v0.2.0                                 | no    | yes   | yes   | no    | no    | no    | no    |
| v0.1.0                                 | yes   | yes   | yes   | no    | no    | no    | no    |

**Note**: for the entry with `+` sign, it means the driver's default released manifest doesn't work with corresponding Kubernetes version, but the driver container image is compatiable with the Kubernetes version if an older version's manifest is used.

## Container Images:
|AWS EBS CSI Driver Version | Image                                            |
|---------------------------|--------------------------------------------------|
|v1.1.1                     |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v1.1.1 |
|v1.1.0                     |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v1.1.0 |
|v1.0.0                     |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v1.0.0 |
|v0.10.1                    |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v0.10.1|
|v0.10.0                    |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v0.10.0|
|v0.9.1                     |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v0.9.1 |
|v0.9.0                     |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v0.9.0 |
|v0.8.1                     |k8s.gcr.io/provider-aws/aws-ebs-csi-driver:v0.8.1 |
|v0.7.1                     |amazon/aws-ebs-csi-driver:v0.7.1                  |
|v0.6.0                     |amazon/aws-ebs-csi-driver:v0.6.0                  |
|v0.5.0                     |amazon/aws-ebs-csi-driver:v0.5.0                  |
|v0.4.0                     |amazon/aws-ebs-csi-driver:v0.4.0                  |
|v0.3.0                     |amazon/aws-ebs-csi-driver:v0.3.0                  |
|v0.2.0                     |amazon/aws-ebs-csi-driver:0.2.0                   |
|v0.1.0                     |amazon/aws-ebs-csi-driver:0.1.0-alpha             |

## Features
* **Static Provisioning** - create a new or migrating existing EBS volumes, then create persistence volume (PV) from the EBS volume and consume the PV from container using persistence volume claim (PVC).
* **Dynamic Provisioning** - uses persistence volume claim (PVC) to request the Kuberenetes to create the EBS volume on behalf of user and consumes the volume from inside container. Storage class's **allowedTopologies** could be used to restrict which AZ the volume should be provisioned in. The topology key should be **topology.ebs.csi.aws.com/zone**.
* **Mount Option** - mount options could be specified in persistence volume (PV) to define how the volume should be mounted.
* **NVMe** - consume NVMe EBS volume from EC2 [Nitro instance](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instance-types.html#ec2-nitro-instances).
* **[Block Volume](https://kubernetes-csi.github.io/docs/raw-block.html)** - consumes the EBS volume as a raw block device for latency sensitive application eg. MySql. The corresponding CSI feature (`CSIBlockVolume`) is GA since Kubernetes 1.18.
* **[Volume Snapshot](https://kubernetes-csi.github.io/docs/snapshot-restore-feature.html)** - creating volume snapshots and restore volume from snapshot. The corresponding CSI feature (`VolumeSnapshotDataSource`) is beta since Kubernetes 1.17.
* **[Volume Resizing](https://kubernetes-csi.github.io/docs/volume-expansion.html)** - expand the volume size. The corresponding CSI feature (`ExpandCSIVolumes`) is beta since Kubernetes 1.16.

## Prerequisites
* If you are managing EBS volumes using static provisioning, get yourself familiar with [EBS volume](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AmazonEBS.html).
* Get yourself familiar with how to setup Kubernetes on AWS and have a working Kubernetes cluster:
  * Enable flag `--allow-privileged=true` for `kubelet` and `kube-apiserver`
  * Enable `kube-apiserver` feature gates `--feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,CSIBlockVolume=true,VolumeSnapshotDataSource=true`
  * Enable `kubelet` feature gates `--feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,CSIBlockVolume=true`
* If you intend to use the csi-snapshotter functionality you will need to first install the [CSI Snapshotter](https://github.com/kubernetes-csi/external-snapshotter)

## Installation
#### Set up driver permission
The driver requires IAM permission to talk to Amazon EBS to manage the volume on user's behalf. [The example policy here](./example-iam-policy.json) defines these permissions. There are several methods to grant the driver IAM permission:
* Using IAM [instance profile](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_switch-role-ec2_instance-profiles.html) - attach the policy to the instance profile IAM role and turn on access to [instance metadata](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) for the instance(s) on which the driver Deployment will run
* EKS only: Using [IAM roles for ServiceAccounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) - create an IAM role, attach the policy to it, then follow the IRSA documentation to associate the IAM role with the driver Deployment service account, which if you are installing via helm is determined by value `serviceAccount.controller.name`, `ebs-csi-controller-sa` by default
* Using secret object - create an IAM user, attach the policy to it, put that user's credentials in [secret manifest](../deploy/kubernetes/secret.yaml), then deploy the secret
```sh
curl https://raw.githubusercontent.com/kubernetes-sigs/aws-ebs-csi-driver/master/deploy/kubernetes/secret.yaml > secret.yaml
# Edit the secret with user credentials
kubectl apply -f secret.yaml
```

#### Config node toleration settings
By default, driver tolerates taint `CriticalAddonsOnly` and has `tolerationSeconds` configured as `300`, to deploy the driver on all nodes, please set helm `Value.node.tolerateAllTaints` to true before deployment

#### Deploy driver
Please see the compatibility matrix above before you deploy the driver

If you want to deploy the stable driver without alpha features:
```sh
kubectl apply -k "github.com/kubernetes-sigs/aws-ebs-csi-driver/deploy/kubernetes/overlays/stable/?ref=release-0.10"
```

If you want to deploy the driver with alpha features:
```sh
kubectl apply -k "github.com/kubernetes-sigs/aws-ebs-csi-driver/deploy/kubernetes/overlays/alpha/?ref=master"
```

Verify driver is running:
```sh
kubectl get pods -n kube-system
```

Alternatively, you could also install the driver using helm:

Add the aws-ebs-csi-driver Helm repository:
```sh
helm repo add aws-ebs-csi-driver https://kubernetes-sigs.github.io/aws-ebs-csi-driver
helm repo update
```

Then install a release of the driver using the chart
```sh
helm upgrade --install aws-ebs-csi-driver \
    --namespace kube-system \
    aws-ebs-csi-driver/aws-ebs-csi-driver
```

#### Deploy driver with debug mode
To view driver debug logs, run the CSI driver with `-v=5` command line option

To enable aws sdk debug logs, run the CSI driver with `--aws-sdk-debug-log=true` command line option.

## Examples
Make sure you follow the [Prerequisites](README.md#Prerequisites) before the examples:
* [Dynamic Provisioning](../examples/kubernetes/dynamic-provisioning)
* [Block Volume](../examples/kubernetes/block-volume)
* [Volume Snapshot](../examples/kubernetes/snapshot)
* [Configure StorageClass](../examples/kubernetes/storageclass)
* [Volume Resizing](../examples/kubernetes/resizing)


## Migrating from in-tree EBS plugin
Starting from Kubernetes 1.17, CSI migration is supported as beta feature (alpha since 1.14). If you have persistent volumes that are created with in-tree `kubernetes.io/aws-ebs` plugin, you can migrate to use EBS CSI driver. To turn on the migration, set `CSIMigration` and `CSIMigrationAWS` feature gates to `true` for `kube-controller-manager`. Then drain Nodes and set the same feature gates to `true` for `kubelet`.

To make sure dynamically provisioned EBS volumes have all tags that the in-tree volume plugin used:
* Run the external-provisioner sidecar with `--extra-create-metadata=true` cmdline option. The helm chart sets this option true by default.
* Run the CSI driver with `--k8s-tag-cluster-id=<ID of the Kubernetes cluster>` command line option.

To make sure that the CSI driver has permission to Attach, Detach, and Delete volumes that were dynamically provisioned and tagged by the in-tree plugin prior to migration being turned on, the IAM policy has to grant permission to operate on volumes with tag `kubernetes.io/cluster/<ID of the Kubernetes cluster>": "owned"` like in [the example policy](./example-iam-policy.json#L85).

**Warning**:
* kubelet *must* be drained of all pods with mounted EBS volumes ***before*** changing its CSI migration feature flags.  Failure to do this will cause deleted pods to get stuck in `Terminating`, requiring a forced delete which can cause filesystem corruption. See [#679](../../../issues/679) for more details.

## Development
Please go through [CSI Spec](https://github.com/container-storage-interface/spec/blob/master/spec.md) and [General CSI driver development guideline](https://kubernetes-csi.github.io/docs/developing.html) to get some basic understanding of CSI driver before you start.

### Requirements
* Golang 1.15.+
* [Ginkgo](https://github.com/onsi/ginkgo) in your PATH for integration testing and end-to-end testing
* Docker 17.05+ for releasing

### Dependency
Dependencies are managed through go module. To build the project, first turn on go mod using `export GO111MODULE=on`, then build the project using: `make`

### Testing
* To execute all unit tests, run: `make test`
* To execute sanity test run: `make test-sanity`
* To execute integration tests, run: `make test-integration`
* To execute e2e tests, run: `make test-e2e-single-az` and `make test-e2e-multi-az`

### Release Process
Please see [Release Process](./RELEASE.md).

**Notes**:
* Sanity tests make sure the driver complies with the CSI specification
* EC2 instance is required to run integration test, since it is exercising the actual flow of creating EBS volume, attaching it and read/write on the disk. See [Integration Testing](../tests/integration/README.md) for more details.
* E2E tests exercises various driver functionalities in Kubernetes cluster. See [E2E Testing](../tests/e2e/README.md) for more details.

### Helm and manifests
The helm chart for this project is in the `charts/aws-ebs-csi-driver` directory.  The manifests for this project are in the `deploy/kubernetes` directory.  All of the manifests except kustomize patches are generated by running `helm template`.  This keeps the helm chart and the manifests in sync.

When updating the helm chart:
* Generate manifests: `make generate-kustomize`
* There are values files in `deploy/kubernetes/values` used for generating some of the manifests
* When adding a new resource template to the helm chart please update the `generate-kustomize` make target, the `deploy/kubernetes/values` files, and the appropriate kustomization.yaml file(s).

## Milestone
[Milestones page](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/milestones)
