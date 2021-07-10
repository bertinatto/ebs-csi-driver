# Helm chart

# v2.0.0
* Remove support for Helm 2
* Remove deprecated values
* No longer install snapshot controller or its CRDs
* Reorganize additional values

[Upgrade instructions](/docs/README.md#upgrading-from-version-1x-to-2x-of-the-helm-chart)

# v1.2.4
* Bump app/driver version to `v1.1.1`
* Install VolumeSnapshotClass, VolumeSnapshotContent, VolumeSnapshot CRDs if enableVolumeSnapshot is true
* Only run csi-snapshotter sidecar if enableVolumeSnapshot is true or if CRDs are already installed
