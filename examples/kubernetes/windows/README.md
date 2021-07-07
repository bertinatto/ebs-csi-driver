## Windows

**This example requires pre-release versions of csi-proxy and the driver that do not exist yet. It is intended for developers only for now. Only basic read/write (mount/unmount and attach/detach) functionality has been tested, other features like resize don't work yet.**

This example shows how to create a EBS volume and consume it from a Windows container dynamically.


## Prerequisites

1. A 1.18+ Windows node. Windows support has only been tested on 1.18 EKS Windows nodes. https://docs.aws.amazon.com/eks/latest/userguide/windows-support.html
2. [csi-proxy](https://github.com/kubernetes-csi/csi-proxy) v1.0.0+ installed on the Windows node.
3. The driver v1.2.0+ Node plugin (DaemonSet) installed on the Windows node.
4. The driver v1.2.0+ Controller plugin (Deployment) installed on a Linux node (as it cannot run on a Windows node).

## Usage

1. Create a sample app along with the StorageClass and the PersistentVolumeClaim:
```
kubectl apply -f specs/
```

2. Validate the volume was created and `volumeHandle` contains an EBS volumeID:
```
kubectl describe pv
```

3. Validate the pod can write data to the volume:
```
kubectl exec -it windows-server-iis-7c5fc8f6c5-t5mk9 -- powershell

PS C:\> New-Item -Path data -Name "testfile1.txt" -ItemType "file" -Value "This 
is a text string."


    Directory: C:\data


Mode                LastWriteTime         Length Name
----                -------------         ------ ----
-a----         4/7/2021  12:31 AM             22 testfile1.txt
```

4. Validate a different pod can read data from the volume:
```
kubectl delete po windows-server-iis-7c5fc8f6c5-t5mk9

kubectl exec -it windows-server-iis-7c5fc8f6c5-j44qv -- powershell

PS C:\> ls data 


    Directory: C:\data 


Mode                LastWriteTime         Length Name
----                -------------         ------ ----
-a----         4/7/2021  12:31 AM             22 testfile1.txt
```

5. OPTIONAL: In case you want to run some e2e tests with Windows pods, make sure to cordon Linux nodes for the duration of the test and modify the vpc-admission-webhook so that the Pods created as part of the tests get scheduled to the Windows nodes.
```
kubectl cordon -l kubernetes.io/os=linux
# edit the webhook such that OSLabelSelectorOverride=all, otherwise the webhook
# won't mutate Pods created by the test and they won't run
kubectl edit deployment -n kube-system vpc-admission-webhook
deployment.apps/vpc-admission-webhook edited
ginkgo -nodes=1 -v --focus="External.Storage.*default.fs.*should.store.data" ./tests/e2e-kubernetes/ -- -kubeconfig=$KUBECONFIG -gce-zone=us-west-2a -node-os-distro=windows
```

6. Cleanup resources:
```
kubectl delete -f specs/
kubectl uncordon -l kubernetes.io/os=linux
```
