// +build windows

/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mounter

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	disk "github.com/kubernetes-csi/csi-proxy/client/api/disk/v1"
	diskclient "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1"

	fs "github.com/kubernetes-csi/csi-proxy/client/api/filesystem/v1"
	fsclient "github.com/kubernetes-csi/csi-proxy/client/groups/filesystem/v1"

	volume "github.com/kubernetes-csi/csi-proxy/client/api/volume/v1"
	volumeclient "github.com/kubernetes-csi/csi-proxy/client/groups/volume/v1"

	"k8s.io/klog/v2"
	"k8s.io/mount-utils"
	utilexec "k8s.io/utils/exec"
)

var _ mount.Interface = &CSIProxyMounter{}

type CSIProxyMounter struct {
	FsClient     *fsclient.Client
	DiskClient   *diskclient.Client
	VolumeClient *volumeclient.Client
}

func normalizeWindowsPath(path string) string {
	normalizedPath := strings.Replace(path, "/", "\\", -1)
	if strings.HasPrefix(normalizedPath, "\\") {
		normalizedPath = "c:" + normalizedPath
	}
	return normalizedPath
}

// Mount just creates a soft link at target pointing to source.
func (mounter *CSIProxyMounter) Mount(source string, target string, fstype string, options []string) error {
	// Mount is called after the format is done.
	// TODO: Confirm that fstype is empty.
	linkRequest := &fs.CreateSymlinkRequest{
		SourcePath: normalizeWindowsPath(source),
		TargetPath: normalizeWindowsPath(target),
	}
	_, err := mounter.FsClient.CreateSymlink(context.Background(), linkRequest)
	if err != nil {
		return err
	}
	return nil
}

// Rmdir - delete the given directory
// TODO: Call separate rmdir for pod context and plugin context. v1alpha1 for CSI
//       proxy does a relaxed check for prefix as c:\var\lib\kubelet, so we can do
//       rmdir with either pod or plugin context.
func (mounter *CSIProxyMounter) Rmdir(path string) error {
	rmdirRequest := &fs.RmdirRequest{
		Path:  normalizeWindowsPath(path),
		Force: true,
	}
	_, err := mounter.FsClient.Rmdir(context.Background(), rmdirRequest)
	if err != nil {
		return err
	}
	return nil
}

// Unmount - Removes the directory - equivalent to unmount on Linux.
func (mounter *CSIProxyMounter) Unmount(target string) error {
	// WriteVolumeCache before unmount
	response, err := mounter.VolumeClient.GetVolumeIDFromTargetPath(context.Background(), &volume.GetVolumeIDFromTargetPathRequest{TargetPath: target})
	if err != nil || response == nil {
		klog.Warningf("GetVolumeIDFromTargetPath(%s) failed with error: %v, response: %v", target, err, response)
	} else {
		request := &volume.WriteVolumeCacheRequest{
			VolumeId: response.VolumeId,
		}
		if res, err := mounter.VolumeClient.WriteVolumeCache(context.Background(), request); err != nil {
			klog.Warningf("WriteVolumeCache(%s) failed with error: %v, response: %v", response.VolumeId, err, res)
		}
	}
	return mounter.Rmdir(target)
}

func (mounter *CSIProxyMounter) List() ([]mount.MountPoint, error) {
	return []mount.MountPoint{}, fmt.Errorf("List not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) IsMountPointMatch(mp mount.MountPoint, dir string) bool {
	return mp.Path == dir
}

// IsLikelyMountPoint - If the directory does not exists, the function will return os.ErrNotExist error.
//   If the path exists, call to CSI proxy will check if its a link, if its a link then existence of target
//   path is checked.
func (mounter *CSIProxyMounter) IsLikelyNotMountPoint(path string) (bool, error) {
	isExists, err := mounter.ExistsPath(path)
	if err != nil {
		return false, err
	}

	if !isExists {
		return true, os.ErrNotExist
	}

	response, err := mounter.FsClient.IsSymlink(context.Background(),
		&fs.IsSymlinkRequest{
			Path: normalizeWindowsPath(path),
		})
	if err != nil {
		return false, err
	}
	return !response.IsSymlink, nil
}

func (mounter *CSIProxyMounter) PathIsDevice(pathname string) (bool, error) {
	return false, fmt.Errorf("PathIsDevice not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) DeviceOpened(pathname string) (bool, error) {
	return false, fmt.Errorf("DeviceOpened not implemented for CSIProxyMounter")
}

// GetDeviceNameFromMount returns the volume ID for a mount path.
func (mounter *CSIProxyMounter) GetDeviceNameFromMount(mountPath, pluginMountDir string) (string, error) {
	req := &volume.GetVolumeIDFromTargetPathRequest{TargetPath: normalizeWindowsPath(mountPath)}
	resp, err := mounter.VolumeClient.GetVolumeIDFromTargetPath(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.VolumeId, nil
}

func (mounter *CSIProxyMounter) MakeRShared(path string) error {
	return fmt.Errorf("MakeRShared not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) MakeFile(pathname string) error {
	return fmt.Errorf("MakeFile not implemented for CSIProxyMounter")
}

// MakeDir - Creates a directory. The CSI proxy takes in context information.
// Currently the make dir is only used from the staging code path, hence we call it
// with Plugin context..
func (mounter *CSIProxyMounter) MakeDir(pathname string) error {
	mkdirReq := &fs.MkdirRequest{
		Path: normalizeWindowsPath(pathname),
	}
	_, err := mounter.FsClient.Mkdir(context.Background(), mkdirReq)
	if err != nil {
		klog.V(4).Infof("Error: %v", err)
		return err
	}

	return nil
}

// ExistsPath - Checks if a path exists. Unlike util ExistsPath, this call does not perform follow link.
func (mounter *CSIProxyMounter) ExistsPath(path string) (bool, error) {
	isExistsResponse, err := mounter.FsClient.PathExists(context.Background(),
		&fs.PathExistsRequest{
			Path: normalizeWindowsPath(path),
		})
	if err != nil {
		return false, err
	}
	return isExistsResponse.Exists, err
}

func (mounter *CSIProxyMounter) EvalHostSymlinks(pathname string) (string, error) {
	return "", fmt.Errorf("EvalHostSymlinks is not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) GetMountRefs(pathname string) ([]string, error) {
	return []string{}, fmt.Errorf("GetMountRefs is not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) GetFSGroup(pathname string) (int64, error) {
	return -1, fmt.Errorf("GetFSGroup is not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) GetSELinuxSupport(pathname string) (bool, error) {
	return false, fmt.Errorf("GetSELinuxSupport is not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) GetMode(pathname string) (os.FileMode, error) {
	return 0, fmt.Errorf("GetMode is not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) MountSensitive(source string, target string, fstype string, options []string, sensitiveOptions []string) error {
	return fmt.Errorf("MountSensitive is not implemented for CSIProxyMounter")
}

func (mounter *CSIProxyMounter) MountSensitiveWithoutSystemd(source string, target string, fstype string, options []string, sensitiveOptions []string) error {
	return fmt.Errorf("MountSensitiveWithoutSystemd is not implemented for CSIProxyMounter")
}

// Rescan would trigger an update storage cache via the CSI proxy.
func (mounter *CSIProxyMounter) Rescan() error {
	// Call Rescan from disk APIs of CSI Proxy.
	if _, err := mounter.DiskClient.Rescan(context.Background(), &disk.RescanRequest{}); err != nil {
		return err
	}
	return nil
}

// FindDiskByLun - given a lun number, find out the corresponding disk
func (mounter *CSIProxyMounter) FindDiskByLun(lun string) (diskNum string, err error) {
	findDiskByLunResponse, err := mounter.DiskClient.ListDiskLocations(context.Background(), &disk.ListDiskLocationsRequest{})
	if err != nil {
		return "", err
	}

	// List all disk locations and match the lun id being requested for.
	// If match is found then return back the disk number.
	for diskID, location := range findDiskByLunResponse.DiskLocations {
		if strings.EqualFold(location.LUNID, lun) {
			return strconv.Itoa(int(diskID)), nil
		}
	}
	return "", fmt.Errorf("could not find disk id for lun: %s", lun)
}

// FormatAndMount - accepts the source disk number, target path to mount, the fstype to format with and options to be used.
func (mounter *CSIProxyMounter) FormatAndMount(source string, target string, fstype string, options []string) error {
	diskNumber, err := strconv.Atoi(source)
	if err != nil {
		return err
	}

	// Call PartitionDisk CSI proxy call to partition the disk and return the volume id
	partionDiskRequest := &disk.PartitionDiskRequest{
		DiskNumber: uint32(diskNumber),
	}
	_, err = mounter.DiskClient.PartitionDisk(context.Background(), partionDiskRequest)
	if err != nil {
		return err
	}

	// List the volumes on the given disk.
	volumeIDsRequest := &volume.ListVolumesOnDiskRequest{
		DiskNumber: uint32(diskNumber),
	}
	volumeIdResponse, err := mounter.VolumeClient.ListVolumesOnDisk(context.Background(), volumeIDsRequest)
	if err != nil {
		return err
	}

	// TODO: consider partitions and choose the right partition.
	// For now just choose the first volume.
	volumeID := volumeIdResponse.VolumeIds[0]

	// Check if the volume is formatted.
	isVolumeFormattedRequest := &volume.IsVolumeFormattedRequest{
		VolumeId: volumeID,
	}
	isVolumeFormattedResponse, err := mounter.VolumeClient.IsVolumeFormatted(context.Background(), isVolumeFormattedRequest)
	if err != nil {
		return err
	}

	// If the volume is not formatted, then format it, else proceed to mount.
	if !isVolumeFormattedResponse.Formatted {
		formatVolumeRequest := &volume.FormatVolumeRequest{
			VolumeId: volumeID,
			// TODO: Accept the filesystem and other options
		}
		_, err = mounter.VolumeClient.FormatVolume(context.Background(), formatVolumeRequest)
		if err != nil {
			return err
		}
	}

	// Mount the volume by calling the CSI proxy call.
	mountVolumeRequest := &volume.MountVolumeRequest{
		VolumeId:   volumeID,
		TargetPath: normalizeWindowsPath(target),
	}
	_, err = mounter.VolumeClient.MountVolume(context.Background(), mountVolumeRequest)
	if err != nil {
		return err
	}
	return nil
}

// ResizeVolume resizes the volume to the maximum available size.
func (mounter *CSIProxyMounter) ResizeVolume(devicePath string) error {
	req := &volume.ResizeVolumeRequest{VolumeId: devicePath, SizeBytes: 0}

	_, err := mounter.VolumeClient.ResizeVolume(context.Background(), req)
	if err != nil {
		return err
	}

	return nil
}

// GetVolumeSizeInBytes returns the size of the volume in bytes.
func (mounter *CSIProxyMounter) GetVolumeSizeInBytes(devicePath string) (int64, error) {
	req := &volume.GetVolumeStatsRequest{VolumeId: devicePath}

	resp, err := mounter.VolumeClient.GetVolumeStats(context.Background(), req)
	if err != nil {
		return -1, err
	}

	return resp.TotalBytes, nil
}

// NewCSIProxyMounter - creates a new CSI Proxy mounter struct which encompassed all the
// clients to the CSI proxy - filesystem, disk and volume clients.
func NewCSIProxyMounter() (*CSIProxyMounter, error) {
	fsClient, err := fsclient.NewClient()
	if err != nil {
		return nil, err
	}
	diskClient, err := diskclient.NewClient()
	if err != nil {
		return nil, err
	}
	volumeClient, err := volumeclient.NewClient()
	if err != nil {
		return nil, err
	}
	return &CSIProxyMounter{
		FsClient:     fsClient,
		DiskClient:   diskClient,
		VolumeClient: volumeClient,
	}, nil
}

func NewSafeMounter() (*mount.SafeFormatAndMount, error) {
	csiProxyMounter, err := NewCSIProxyMounter()
	if err != nil {
		return nil, err
	}
	return &mount.SafeFormatAndMount{
		Interface: csiProxyMounter,
		Exec:      utilexec.New(),
	}, nil
}
