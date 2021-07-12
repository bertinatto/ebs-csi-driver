# v1.1.1

### Bug fixes
- update inFlight cache to avoid race condition on volume operation ([#924](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/924), [@AndyXiangLi](https://github.com/AndyXiangLi))

# v1.1.0

## Notable changes
- Helm chart cleaned up ([#856](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/856), [@krmichel](https://github.com/krmichel))

### New features
* Add podAnnotations to snapshotController StatefulSet ([#884](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/884), [@snstanton](https://github.com/snstanton))
* Support custom pod labels in Helm chart ([#905](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/905), [@eytanhanig](https://github.com/eytanhanig))

### Bug fixes
* fix naming mistake in clusterrolebinding, expose env var to controller via downward api ([#874](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/874), [@vdhanan](https://github.com/vdhanan))
* Fix kustomize RBAC bindings to have namespace kube-system ([#878](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/878), [@wongma7](https://github.com/wongma7))
* rename node clusterrolebinding to make auto upgrade work ([#894](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/894), [@vdhanan](https://github.com/vdhanan))
* remove hardcoded namespace for pod disruption budget ([#895](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/895), [@vdhanan](https://github.com/vdhanan))
* Only initialize the in-cluster kube client when metadata service is actually unavailable ([#897](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/897), [@chrisayoub](https://github.com/chrisayoub))
* Reduce default log level to 2 ([#903](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/903), [@wongma7](https://github.com/wongma7))
* Add pod disruption budgets that got missed in a rebase ([#906](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/906), [@krmichel](https://github.com/krmichel))
* remove WellKnownTopologyKey from PV Topology ([#912](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/912), [@Elbehery](https://github.com/Elbehery))
* Skip volume expansion if block node ([#916](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/916), [@gnufied](https://github.com/gnufied))

### Misc.
* Add eksctl support to e2e scripts ([#852](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/852), [@wongma7](https://github.com/wongma7))
* release v1.0.0 ([#865](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/865), [@vdhanan](https://github.com/vdhanan))
* add self as owner ([#866](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/866), [@vdhanan](https://github.com/vdhanan))
* bump helm chart version ([#881](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/881), [@vdhanan](https://github.com/vdhanan))
* add custom useragent suffix ([#910](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/910), [@vdhanan](https://github.com/vdhanan))
* Bump chart-releaser-action to v1.2.1 ([#914](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/914), [@gliptak](https://github.com/gliptak))

# v1.0.0
## Notable changes
- With this release, the EBS CSI Driver is now Generally Available!

### New features
* add options to enable aws sdk debug log and add more logs when driver… ([#830](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/830), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Emit AWS API operation duration/error/throttle metrics ([#842](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/842), [@wongma7](https://github.com/wongma7))
* add pod disruption budget for csi controller ([#857](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/857), [@vdhanan](https://github.com/vdhanan))

### Bug fixes
* Resize filesystem when restore a snapshot to larger size volume ([#753](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/753), [@AndyXiangLi](https://github.com/AndyXiangLi))
* handling describe instances consistency issue ([#801](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/801), [@vdhanan](https://github.com/vdhanan))
* Cap IOPS when calculating from iopsPerGB ([#809](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/809), [@jsafrane](https://github.com/jsafrane))
* Fix broken gomocks ([#843](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/843), [@wongma7](https://github.com/wongma7))
* Fix missing import ([#849](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/849), [@wongma7](https://github.com/wongma7))
* instance metadata issue fix ([#855](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/855), [@vdhanan](https://github.com/vdhanan))

### Misc.
* release v0.10.0 ([#820](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/820), [@vdhanan](https://github.com/vdhanan))
* release v0.10.1 ([#827](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/827), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Rebase 1.21 ([#828](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/828), [@jsafrane](https://github.com/jsafrane))
* update installation command to use latest stable version ([#832](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/832), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Bump/reconcile sidecar versions in helm/kustomize ([#834](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/834), [@wongma7](https://github.com/wongma7))
* update IAM policy sample and add new driver level tag ([#835](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/835), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Switch to non-deprecated apiVersion ([#836](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/836), [@dntosas](https://github.com/dntosas))
* Update readme file to provide more info on driver options and tagging ([#844](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/844), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Add empty StorageClasses to static example ([#850](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/850), [@wongma7](https://github.com/wongma7))
* Add additional logging for outpost arn handling ([#851](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/851), [@ayberk](https://github.com/ayberk))

# v0.10.1
## Notable changes
* support volume partition, users can specify partition in the pv and driver will mount the device on the specified partition ([#824](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/824), [@AndyXiangLi](https://github.com/AndyXiangLi))

### Misc.
* Warn users of migrating without draining ([#822](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/822), [@partcyborg](https://github.com/partcyborg))

# v0.10.0

## Notable changes
- Prep for Windows support: Copy pkg/mounter and refactor to use k8s.io/mount-utils ([#786](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/786), [@wongma7](https://github.com/wongma7))
- Add well-known topology label ([#773](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/773), [@ayberk](https://github.com/ayberk))
- Update livenessprobe image version from 2.1.0 to 2.2.0 ([#756](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/756), [@mowangdk](https://github.com/mowangdk))
- Remove arm overlay ([#719](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/719), [@ayberk](https://github.com/ayberk))
- Add readiness probe so controller does not report "Ready" prematurely ([#751](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/751), [@vdhanan](https://github.com/vdhanan))
- Add toleration time to NoExecute effect ([#776](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/776), [@AndyXiangLi](https://github.com/AndyXiangLi))

### New features
* Add ability to specify topologySpreadConstraints ([#770](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/770), [@arcivanov](https://github.com/arcivanov))

### Bug fixes
* delete leaked volume if driver don't know the volume status ([#771](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/771), [@AndyXiangLi](https://github.com/AndyXiangLi))
* modify error message when request volume is in use with other node ([#698](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/698), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Make CreateVolume idempotent ([#744](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/744), [@chrishenzie](https://github.com/chrishenzie))

### Misc.
* Add documentation for release process ([#610](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/610), [@ayberk](https://github.com/ayberk))
* feat: Add option to provision StorageClasses ([#697](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/697), [@gazal-k](https://github.com/gazal-k))
* Refactor inFlight key to add lock per volumeId ([#702](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/702), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Add support for node existing service accounts ([#704](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/704), [@mper0003](https://github.com/mper0003))
* More controll over snapshot-controller scheduling ([#708](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/708), [@alex-berger](https://github.com/alex-berger))
* Remove hardcoded snapshot controller image references ([#711](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/711), [@ig0rsky](https://github.com/ig0rsky))
* release 0.9.0 ([#718](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/718), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Move cr.yaml out of github workflows ([#720](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/720), [@ayberk](https://github.com/ayberk))
* Bump chart version ([#724](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/724), [@ayberk](https://github.com/ayberk))
* Integrate external e2e test in the testsuits ([#726](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/726), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Allow all fields to be set on StorageClasses ([#730](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/730), [@haines](https://github.com/haines))
* [chart] Allow resources override for node DaemonSet + priorityClassName ([#732](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/732), [@dntosas](https://github.com/dntosas))
* [chart]  Add storage class annotation and label handling ([#734](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/734), [@nicholasmhughes](https://github.com/nicholasmhughes))
* Updated installation to use latest 0.9 release ([#735](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/735), [@PhilThurston](https://github.com/PhilThurston))
* patch stable release to use gcr image ([#740](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/740), [@AndyXiangLi](https://github.com/AndyXiangLi))
* correct kustomization gcr image repo ([#742](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/742), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Update ECR overlay ([#745](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/745), [@ayberk](https://github.com/ayberk))
* Set enableVolumeScheduling to true by default in the helm chart ([#752](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/752), [@mtougeron](https://github.com/mtougeron))
* Sets the imagePullSecrets if the value is set in the chart ([#755](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/755), [@mtougeron](https://github.com/mtougeron))
* Update test k8s version to 1.18.16 ([#759](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/759), [@ayberk](https://github.com/ayberk))
* add a document separator for storageclass template file ([#762](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/762), [@nvnmandadhi](https://github.com/nvnmandadhi))
* Allow setting http proxy and no proxy environment values ([#765](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/765), [@rubroboletus](https://github.com/rubroboletus))
* Fix error message when IOPSPerGB is missing in io1 volumes ([#767](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/767), [@jsafrane](https://github.com/jsafrane))
* removed harcoded NAMESPACE from helm chart ([#768](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/768), [@alexandrst88](https://github.com/alexandrst88))
* Aws client config: increase MaxRetries ([#769](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/769), [@josselin-c](https://github.com/josselin-c))
* Update chart version ([#772](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/772), [@ayberk](https://github.com/ayberk))
* Add self as reviewer ([#774](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/774), [@AndyXiangLi](https://github.com/AndyXiangLi))
* go mod tidy ([#777](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/777), [@vdhanan](https://github.com/vdhanan))
* Removing prestop hook for node-driver-registrar ([#778](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/778), [@tsunny](https://github.com/tsunny))
* hack/e2e: Support passing helm values as values.yaml and make other similar files optional ([#787](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/787), [@wongma7](https://github.com/wongma7))
* Print csi plugin logs at end of e2e test ([#789](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/789), [@wongma7](https://github.com/wongma7))
* Update snapshot controller resources ([#791](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/791), [@tirumerla](https://github.com/tirumerla))
* Remove storageclass from static example ([#794](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/794), [@wongma7](https://github.com/wongma7))
* Don't exit script prematurely if test fails ([#802](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/802), [@wongma7](https://github.com/wongma7))
* csi.storage.k8s.io/fstype is case sensitive ([#807](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/807), [@jsafrane](https://github.com/jsafrane))
* fix deploy stable ecr error kustomization file ([#808](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/808), [@ABNER-1](https://github.com/ABNER-1))
* release v0.9.1 ([#813](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/813), [@vdhanan](https://github.com/vdhanan))
* Use the old topology key for e2e tests ([#814](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/814), [@ayberk](https://github.com/ayberk))
* Track driver deploy time in e2e test pipeline ([#815](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/815), [@AndyXiangLi](https://github.com/AndyXiangLi))
* AWS EBS CSI Driver Helm chart to inject environment variables ([#817](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/817), [@tomdymond](https://github.com/tomdymond))

# v0.9.1

## Notable changes
- Change helm deploy settings: default tolerationAllTaints to false, NoExecute toleration time is 300s and will tolerate `CriticalAddonsOnly`

### New features
* Integrate external e2e test in the testsuits ([#726](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/726), [@AndyXiangLi](https://github.com/AndyXiangLi))

### Bug fixes
* delete leaked volume if driver don't know the volume status ([#771](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/771), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Update test k8s version to 1.18.16 ([#759](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/759), [@ayberk](https://github.com/ayberk))

# v0.9.0

## Notable changes
- All images (including sidecars) are Multiarch
- Enable volume stats metrics on Node service

### New features
* Feature: Add ability to customize node daemonset nodeselector ([#647](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/647), [@pliu](https://github.com/pliu))
* add volume stats metrics - ([#677](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/677), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Add support for existing service accounts ([#688](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/688), [@ayberk](https://github.com/ayberk))
* NodeExpandVolume no-op for raw block ([#695](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/695), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Allow specifying --volume-attach-limit in the helm chart ([#700](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/700), [@keznikl](https://github.com/keznikl))

### Bug fixes
* Fix outdated ecr login command ([#680](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/680), [@wongma7](https://github.com/wongma7))
* Update sidecars to newer version ([#707](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/707), [@AndyXiangLi](https://github.com/AndyXiangLi))

### Misc.
* Update README.md ([#607](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/607), [@robisoh88](https://github.com/robisoh88))
* Add self to OWNERS ([#638](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/638), [@ayberk](https://github.com/ayberk))
* Bring Go to 1.15.6 in Travis ([#648](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/648), [@gliptak](https://github.com/gliptak))
* Fix overlays not being updated for gcr migration ([#649](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/649), [@wongma7](https://github.com/wongma7))
* Arm overlay ([#653](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/653), [@ayberk](https://github.com/ayberk))
* Use buildx in cloudbuild ([#658](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/658), [@wongma7](https://github.com/wongma7))
* (Try to) fix cloudbuild ([#659](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/659), [@wongma7](https://github.com/wongma7))
* Fix stray argument in cloudbuild.yaml ([#661](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/661), [@wongma7](https://github.com/wongma7))
* Add note for gp3 on outposts ([#665](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/665), [@ayberk](https://github.com/ayberk))
* Call hack/prow.sh from cloudbuild ([#666](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/666), [@wongma7](https://github.com/wongma7))
* cloudbuild: Set _STAGING_PROJECT ([#668](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/668), [@wongma7](https://github.com/wongma7))
* add import snapshot e2e test ([#678](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/678), [@AndyXiangLi](https://github.com/AndyXiangLi))
* Prefix helm chart releases with "helm-chart-" ([#682](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/682), [@wongma7](https://github.com/wongma7))
* Release 0.8.1 ([#683](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/683), [@wongma7](https://github.com/wongma7))
* Push debian target to Docker Hub ([#686](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/686), [@wongma7](https://github.com/wongma7))
* Adds patch for ebs-csi-controller-sa to volumeattachments/status ([#690](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/690), [@cuppett](https://github.com/cuppett))
* Add a prerequisite to dynamic provisioning ([#691](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/691), [@ronenl1](https://github.com/ronenl1))
* Refactor e2e testing scripts to be more reusable and use them instead… ([#694](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/694), [@wongma7](https://github.com/wongma7))
* Update to golang@1.15.6 ([#699](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/699), [@ialidzhikov](https://github.com/ialidzhikov))
* add e2e test for volume resizing ([#705](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/705), [@AndyXiangLi](https://github.com/AndyXiangLi))

# v0.8.1

## Notable changes
- Images in k8s.gcr.io are multiarch.

### Bug fixes
* release-0.8: Use buildx in cloudbuild ([#670](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/670), [@wongma7](https://github.com/wongma7))

# v0.8.0

## Notable changes
- gp3 is now the default volume type. gp3 is **not** supported on outposts. Outpost customers need to use a different type for their volumes.
- Images will be built on a Debian base by default. Images built on Amazon Linux will still be available but with the tag suffix `-amazonlinux`.
- Images will be published to k8s.gcr.io in addition to ECR and Docker Hub.

### New features
* Chart option to disable default toleration of all taints ([#526](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/526), [@risinger](https://github.com/risinger))
* Apply extra volume tags to EBS snapshots ([#568](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/568), [@chrishenzie](https://github.com/chrishenzie))
* [helm] add tag options and update csi-provisioner ([#577](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/577), [@kcking](https://github.com/kcking))
* vendor: bump aws sdk for AssumeRoleWithWebIdentity support ([#614](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/614), [@abhinavdahiya](https://github.com/abhinavdahiya))
* Add EBS gp3 support ([#633](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/633), [@samjo-nyang](https://github.com/samjo-nyang))
* Apply resource constraints to all sidecar containers ([#640](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/640), [@tirumerla](https://github.com/tirumerla))

### Bug fixes
* Fix the name of the snapshot controller leader election RoleBinding ([#601](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/601), [@robbie-demuth](https://github.com/robbie-demuth))

### Misc.
* Post-release v0.7.0 ([#576](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/576), [@ayberk](https://github.com/ayberk))
* Fixing Helm install command ([#578](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/578), [@danil-smirnov](https://github.com/danil-smirnov))
* Fix markdown issue in README.md ([#579](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/579), [@ialidzhikov](https://github.com/ialidzhikov))
* Document behavior wrt minimum and maximum iops ([#582](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/582), [@wongma7](https://github.com/wongma7))
* Set CSIMigrationAWSComplete for migration tests ([#593](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/593), [@wongma7](https://github.com/wongma7))
* Bump migration kops and k8s version ([#602](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/602), [@wongma7](https://github.com/wongma7))
* Update hack/run-e2e-test to be more idempotent and pleasant to use ([#616](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/616), [@wongma7](https://github.com/wongma7))
* Post-release v0.7.1 ([#619](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/619), [@ayberk](https://github.com/ayberk))
* Move chart to charts directory and add workflow to publish new chart versions ([#624](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/624), [@krmichel](https://github.com/krmichel))
* docs(readme): update link to developer documentation ([#629](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/629), [@BondAnthony](https://github.com/BondAnthony))
* Update ecr overlay image tag ([#630](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/630), [@ayberk](https://github.com/ayberk))
* Add cloudbuild.yaml for image pushing to gcr ([#632](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/632), [@wongma7](https://github.com/wongma7))
* Add latest tags to cloudbuild ([#634](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/634), [@wongma7](https://github.com/wongma7))
* Fix target name in cloudbuild.yaml from amazon to amazonlinux ([#636](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/636), [@wongma7](https://github.com/wongma7))
* Suffix amazonlinux image with -amazonlinux and push debian image to GitHub ([#639](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/639), [@wongma7](https://github.com/wongma7))
* Set up QEMU to build for arm64 ([#641](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/641), [@wongma7](https://github.com/wongma7))

# v0.7.1
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.7.1/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.7.1.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.7.1.zip) | `0c8b1e539f5852e54b5f4ab48cb3054ac52145db3d692cdc6b3ac683c39ebf11951c5ff3823a83666605a56a30b38953d20f392397c16bf39a5727c66ddf0827`
[v0.7.1.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.7.1.tar.gz) | `157ed2c7aa00635a61438a1574bd7e124676bcabd9e27cfe865c7bbb3194609894536b1eb38a12a8e5bfa71b540e0f1cde12000b02d90b390d17987fc913042e`

## Notable changes
This release includes a fix for the helm chart to point to the correct image.

# v0.7.0
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.7.0/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.7.0.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.7.0.zip) | `6e1117ce046d0030c3008b3eec8ba3196c516adf0ecef8909fcfd3d68e63624a73a992033356e208bf0d5563f7dec2e40675f0fee7f322bd4f69d7b03750961a`
[v0.7.0.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.7.0.tar.gz) | `4dc3402ffa3dcc59c9af1f7d776a3f53a288f62a31c05cde00aeceeef6000be16ca6cdae08712b4f7f64c9e89ceeaa13df7f1ca4bf3d62ba62845b52cc13eadf`

## Notable changes
### New features
* Add arm support ([#527](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/527), [@leakingtapan](https://github.com/leakingtapan))
* Add EBS IO2 support ([#558](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/558), [@ayberk](https://github.com/ayberk))
* Create volumes in outpost for outpost instances ([#561](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/561), [@ayberk](https://github.com/ayberk))

### Improvements
* Make EBS controllerexpansion idempotent ([#552](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/552), [@gnufied](https://github.com/gnufied))
* Add overlay for ECR images ([#570](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/570), [@ayberk](https://github.com/ayberk))

# v0.6.0
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.6.0/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.6.0.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.6.0.zip) | `67dc79703c2d022cbc53a370e8ac7279bf4345030a3ecc5b2bdff2b722ec807b712f2cd6eae79598edb87e15d92e683e98dde7c25e52f705233bc3ece649c693`
[v0.6.0.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.6.0.tar.gz) | `a3b5e95ec05ce6b4e6eb22ae00c7898cb876f21719354636dae5d323934c7a0bb32a7a8e89abdfcc6b0a0827c7169a349cba9dce32b7bf25e7287a2ec0387f21`

## Notable changes
### New features
* Allow volume attach limit overwrite via command line parameter ([#522](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/522), [@rfranzke](https://github.com/rfranzke))
* Add tags that the in-tree volume plugin uses ([#530](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/530), [@jsafrane](https://github.com/jsafrane))

### Bug fixes
* Adding amd64 as nodeSelector to avoid arm64 archtectures (#471) ([#472](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/472), [@hugoprudente](https://github.com/hugoprudente))
* Update stable overlay to 0.5.0 ([#495](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/495), [@wongma7](https://github.com/wongma7))

### Improvements
* Update aws-sdk to v1.29.11 to get IMDSv2 support ([#463](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/463), [@msau42](https://github.com/msau42))
* Fix e2e test ([#468](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/468), [@leakingtapan](https://github.com/leakingtapan))
* Generate deployment manifests from helm chart ([#475](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/475), [@krmichel](https://github.com/krmichel))
* Correct golint warning ([#478](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/478), [@gliptak](https://github.com/gliptak))
* Bump Go to 1.14.1 ([#479](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/479), [@gliptak](https://github.com/gliptak))
* Add mount unittest ([#481](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/481), [@gliptak](https://github.com/gliptak))
* Remove volume IOPS limit ([#483](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/483), [@jacobmarble](https://github.com/jacobmarble))
* Additional mount unittest ([#484](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/484), [@gliptak](https://github.com/gliptak))
* docs/README: add missing "--namespace" flag to "helm" command ([#486](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/486), [@gyuho](https://github.com/gyuho))
* Add nodeAffinity to avoid Fargate worker nodes ([#488](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/488), [@bgsilvait](https://github.com/bgsilvait))
* remove deprecated "beta.kubernetes.io/os" nodeSelector ([#489](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/489), [@gyuho](https://github.com/gyuho))
* Update kubernetes-csi/external-snapshotter components to v2.1.1 ([#490](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/490), [@ialidzhikov](https://github.com/ialidzhikov))
* Improve csi-snapshotter ClusterRole ([#491](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/491), [@ialidzhikov](https://github.com/ialidzhikov))
* Fix migration test ([#500](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/500), [@leakingtapan](https://github.com/leakingtapan))
* Add missing IAM permissions ([#501](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/501), [@robbiet480](https://github.com/robbiet480))
* Fixed resizing docs to refer the right path to example spec ([#504](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/504), [@amuraru](https://github.com/amuraru))
* optimization: cache go mod during docker build ([#513](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/513), [@leakingtapan](https://github.com/leakingtapan))

# v0.5.0
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.5.0/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.5.0.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.5.0.zip) | `c53327e090352a7f79ee642dbf8c211733f4a2cb78968ec688a1eade55151e65f1f97cd228d22168317439f1db9f3d2f07dcaa2873f44732ad23aaf632cbef3a`
[v0.5.0.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.5.0.tar.gz) | `ec4963d34c601cdf718838d90b8aa6f36b16c9ac127743e73fbe76118a606d41aced116aaaab73370c17bcc536945d5ccd735bc5a4a00f523025c8e41ddedcb8`

## Notable changes
### New features
* Add a cmdline option to add extra volume tags ([#353](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/353), [@jieyu](https://github.com/jieyu))
* Switch to use kustomize for manifest ([#360](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/360), [@leakingtapan](https://github.com/leakingtapan))
* enable users to set ec2-endpoint for nonstandard regions ([#369](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/369), [@amdonov](https://github.com/amdonov))
* Add standard volume type ([#379](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/379), [@leakingtapan](https://github.com/leakingtapan))
* Update aws sdk version to enable EKS IAM for SA ([#386](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/386), [@leakingtapan](https://github.com/leakingtapan))
* Implement different driver modes and AWS Region override for controller service ([#438](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/438), [@rfranzke](https://github.com/rfranzke))
* Add manifest files for snapshotter 2.0 ([#452](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/452), [@leakingtapan](https://github.com/leakingtapan))

### Bug fixes
* Return success if instance or volume are not found ([#375](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/375), [@bertinatto](https://github.com/bertinatto))
* Patch k8scsi sidecars CVE-2019-11255 ([#413](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/413), [@jnaulty](https://github.com/jnaulty))
* Handle mount flags in NodeStageVolume ([#430](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/430), [@bertinatto](https://github.com/bertinatto))

### Improvements
* Run upstream e2e test suites with migration  ([#341](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/341), [@wongma7](https://github.com/wongma7))
* Use new test framework for test orchestration ([#359](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/359), [@leakingtapan](https://github.com/leakingtapan))
* Update to use 1.16 cluster with inline test enabled ([#362](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/362), [@leakingtapan](https://github.com/leakingtapan))
* Enable leader election ([#380](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/380), [@leakingtapan](https://github.com/leakingtapan))
* Update go mod and mount library ([#388](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/388), [@leakingtapan](https://github.com/leakingtapan))
* Refactor NewCloud by pass in region ([#394](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/394), [@leakingtapan](https://github.com/leakingtapan))
* helm: provide an option to set extra volume tags ([#396](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/396), [@jieyu](https://github.com/jieyu))
* Allow override for csi-provisioner image ([#401](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/401), [@gliptak](https://github.com/gliptak))
* Enable volume expansion e2e test for CSI migration ([#407](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/407), [@leakingtapan](https://github.com/leakingtapan))
* Swith to use kops 1.16 ([#409](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/409), [@leakingtapan](https://github.com/leakingtapan))
* Added tolerations for node support ([#420](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/420), [@zerkms](https://github.com/zerkms))
* Update helm chart to better match available values and add the ability to add annotations ([#423](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/423), [@krmichel](https://github.com/krmichel))
* [helm] Also add toleration support to controller ([#433](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/433), [@jyaworski](https://github.com/jyaworski))
* Add ec2:ModifyVolume action ([#434](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/434), [@zodiac12k](https://github.com/zodiac12k))
* Schedule the EBS CSI DaemonSet on all nodes by default ([#441](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/441), [@pcfens](https://github.com/pcfens))

# v0.4.0
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.4.0/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.4.0.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.4.0.zip) | `2f46b54211178ad1e55926284b9f6218be874038a1a62ef364809a5d2c37b7bbbe58a2cc4991b9cf44cbfe4966c61dd6c16df0790627dffac4f7df9ffc084a0c`
[v0.4.0.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.4.0.tar.gz) | `0199df52ac1e19ee6b04efb80439024dde11de3d8fc292ce10527f2e658b393d8bfd4e37a6ec321cb415c9bdbee83ff5dbdf58e2336d03fe5d1b2717ccb11169`

## Action Required
* Update Kubernetes cluster to 1.14+ before installing the driver, since the released driver manifest assumes 1.14+ cluster.
* storageclass parameter's `fstype` key is deprecated in favor of `csi.storage.k8s.io/fstype` key. Please update the key in you stroage parameters.

## Changes since v0.3.0
See [details](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/compare/v0.3.0...v0.4.0) for all the changes.

### Notable changes
* Make secret optional ([#247](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/247), [@leakingtapan](https://github.com/leakingtapan/))
* Add support for XFS filesystem ([#253](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/253), [@leakingtapan](https://github.com/leakingtapan/))
* Upgrade CSI spec to 1.1.0 ([#263](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/263), [@leakingtapan](https://github.com/leakingtapan/))
* Refactor controller unit test with proper mock ([#269](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/269), [@zacharya](https://github.com/zacharya/))
* Refactor device path allocator ([#274](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/274), [@leakingtapan](https://github.com/leakingtapan/))
* Implementing ListSnapshots ([#286](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/286), [@zacharya](https://github.com/zacharya/))
* Add max number of volumes that can be attached to an instance ([#289](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/289), [@bertinatto](https://github.com/bertinatto/))
* Add helm chart ([#303](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/303), [@leakingtapan](https://github.com/leakingtapan/))
* Add volume expansion ([#271](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/271), [@bertinatto](https://github.com/bertinatto/))
* Remove cluster-driver-registrar ([#322](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/322), [@jsafrane](https://github.com/jsafrane/))
* Upgrade to golang 1.12 ([#329](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/329), [@leakingtapan](https://github.com/leakingtapan/))
* Fix bugs by passing fstype correctly ([#335](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/335), [@leakingtapan](https://github.com/leakingtapan/))
* Output junit to ARTIFACTS for testgrid ([#340](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/340), [@wongma7](https://github.com/wongma7/))

# v0.3.0
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.3.0/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.3.0.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.3.0.zip) | `27a7a1cd4fc7a8afa1c0dd8fb3ce4cb1d9fc7439ebdbeba7ac0bfb0df723acb654a92f88270bc68ab4dd6c8943febf779efa8cbebdf3ea2ada145ff7ce426870`
[v0.3.0.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.3.0.tar.gz) | `9126a3493f958aaa4727bc62b1a5c545ac8795f08844a605541aac3d38dea8769cee12c7db94f44179a91af7e8702174bba2533b4e30eb3f32f9b8338101a5db`

## Action Required
* None

## Upgrade Driver
Driver upgrade should be performed one version at a time by using following steps:
1. Delete the old driver controller service and node service along with other resources including cluster roles, cluster role bindings and service accounts.
1. Deploy the new driver controller service and node service along with other resources including cluster roles, cluster role bindings and service accounts.

## Changes since v0.2.0
See [details](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/compare/v0.2.0...v0.3.0) for all the changes.

### Notable changes
* Strip symbol for production build ([#201](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/201), [@leakingtapan](https://github.com/leakingtapan/))
* Remove vendor directory ([#198](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/198), [@leakingtapan](https://github.com/leakingtapan/))
* Use same mount to place in the csi.sock, remove obsolete volumes ([#212](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/212), [@frittentheke](https://github.com/frittentheke/))
* Add snapshot support ([#131](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/131), [@tsmetana](https://github.com/tsmetana/))
* Add snapshot examples ([#210](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/210), [@tsmetana](https://github.com/tsmetana/))
* Implement raw block volume support ([#215](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/215), [@leakingtapan](https://github.com/leakingtapan/))
* Add unit tests for ControllerPublish and ControllerUnpublish requests ([#219](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/219), [@sreis](https://github.com/sreis/))
* New block volume e2e tests ([#226](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/226), [@dkoshkin](https://github.com/dkoshkin/))
* Implement device path discovery for NVMe support ([#231](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/231), [@leakingtapan](https://github.com/leakingtapan/))
* Cleanup README and examples ([@232](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/232), [@dkoshkin](https://github.com/dkoshkin/))
* New volume snapshot e2e tests ([#235](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/235), [@dkoshkin](https://github.com/dkoshkin/))

# v0.2.0
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.2.0/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.2.0.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.2.0.zip) | `a9733881c43dfb788f6c657320b6b4acdd8ee9726649c850282f8a7f15f816a6aa5db187a5d415781a76918a30ac227c03a81b662027c5b192ab57a050bf28ee`
[v0.2.0.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.2.0.tar.gz) | `0d7a3efd0c1b0c6bf01b08c3cbd48d867aeab1cf1f7f12274f42d561f64526c0345f23d5947ddada7a333046f101679eea620c9ab8985f9d4d1c8c3f28de49ce`

## Action Required
* Upgrade the Kubernetes cluster to 1.13+ before deploying the driver. Since CSI 1.0 is only supported starting from Kubernetes 1.13.

## Upgrade Driver
Driver upgrade should be performed one version at a time by using following steps:
1. Delete the old driver controller service and node service along with other resources including cluster roles, cluster role bindings and service accounts.
1. Deploy the new driver controller service and node service along with other resources including cluster roles, cluster role bindings and service accounts.

## Changes since v0.1.0
See [details](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/compare/v0.1.0...v0.2.0) for all the changes.

### Notable changes
* Update to CSI 1.0 ([#122](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/122), [@bertinatto](https://github.com/bertinatto/))
* Add mountOptions support ([#130](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/130), [@bertinatto](https://github.com/bertinatto/))
* Resolve memory addresses in log messages ([#132](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/132), [@bertinatto](https://github.com/bertinatto/))
* Add version flag ([#136](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/136), [@dkoshkin](https://github.com/dkoshkin/))
* Wait for volume to become available ([#126](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/126), [@bertinatto](https://github.com/bertinatto/))
* Add first few e2e test cases #151 ([#151](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/151/commits), [@dkoshkin](https://github.com/dkoshkin/))
* Make test-integration uses aws-k8s-tester ([#153](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/153), [@kschumy](https://github.com/kschumy))
* Rename VolumeNameTagKey ([#161](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/issues/161), [@leakingtapan](https://github.com/leakingtapan/))
* CSI image version and deployment manifests updates  ([#171](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/171), [@dkoshkin](https://github.com/dkoshkin/))
* Update driver manifest files ([#181](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/181), [@leakingtapan](https://github.com/leakingtapan/))
* More e2e tests ([#173](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/173), [@dkoshkin](https://github.com/dkoshkin/))
* Update run-e2e-test script to setup cluster ([#186](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/186), [@leakingtapan](https://github.com/leakingtapan/))
* Check if target path is mounted before unmounting ([#183](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/183), [@sreis](https://github.com/sreis/))

# v0.1.0
[Documentation](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/v0.1.0/docs/README.md)

## Downloads for v0.1.0

filename  | sha512 hash
--------- | ------------
[v0.1.0.zip](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.1.0.zip) | `03841418496e292c3f91cee7942b545395bce049e9c4d2305532545fb82ad2e5189866afec2ed937924e144142b0b915a9467bac42e9f2b881181aba6aa80a68`
[v0.1.0.tar.gz](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/archive/v0.1.0.tar.gz) | `106b6c2011acd42b0f10117b7f104ab188dde798711e98119137cf3d8265e381df09595b8e861c0c9fdcf8772f4a711e338e822602e98bfd68f54f9e1c7f8f16`

## Changelog since initial commit

### Notable changes
* Update driver name and topology key ([#105](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/105), [@leakingtapan](https://github.com/leakingtapan/))
* Add support for creating encrypted volume and unit test ([#80](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/80), [@leakingtapan](https://github.com/leakingtapan/))
* Implement support for storage class parameter - volume type ([#73](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/73), [@leakingtapan](https://github.com/leakingtapan/))
* Implement support for storage class parameter - fsType ([#67](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/67), [@leakingtapan](https://github.com/leakingtapan/))
* Add missing capability and clusterrole permission to enable tology awareness scheduling ([#61](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/commit/2873e0b), [@leakingtapan](https://github.com/leakingtapan/))
* Wait for correct attachment state ([#58](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/58), [@bertinatto](https://github.com/bertinatto/))
* Implement topology awareness support for dynamic provisioning ([#42](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/42), [@leakingtapan](https://github.com/leakingtapan/))
* Wait for volume status in e2e test ([#34](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/34), [@bertinatto](https://github.com/bertinatto/))
* Update cloud provider interface to take in context ([#45](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/pull/45), [@leakingtapan](https://github.com/leakingtapan/))
* Initial driver implementation ([9ba4c5d](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/commit/9ba4c5d), [@bertinatto](https://github.com/bertinatto/))
