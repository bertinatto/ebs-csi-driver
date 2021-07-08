// +build windows

/*
Copyright 2019 The Kubernetes Authors.

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

package driver

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	diskapi "github.com/kubernetes-csi/csi-proxy/client/api/disk/v1"
	diskclient "github.com/kubernetes-csi/csi-proxy/client/groups/disk/v1"
	"github.com/kubernetes-sigs/aws-ebs-csi-driver/pkg/mounter"
	"k8s.io/klog"
)

// findDevicePath finds disk number of device
// https://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/ec2-windows-volumes.html#list-nvme-powershell
func (d *nodeService) findDevicePath(devicePath, volumeID, _ string) (string, error) {

	diskClient, err := diskclient.NewClient()
	if err != nil {
		return "", fmt.Errorf("error creating csi-proxy disk client: %q", err)
	}
	defer diskClient.Close()

	response, err := diskClient.ListDiskIDs(context.TODO(), &diskapi.ListDiskIDsRequest{})
	if err != nil {
		return "", fmt.Errorf("error listing disk ids: %q", err)
	}

	diskIDs := response.GetDiskIDs()

	foundDiskNumber := ""
	for diskNumber, diskID := range diskIDs {
		serialNumber := diskID.GetSerialNumber()
		cleanVolumeID := strings.ReplaceAll(volumeID, "-", "")
		if strings.Contains(serialNumber, cleanVolumeID) {
			foundDiskNumber = strconv.Itoa(int(diskNumber))
			break
		}
	}

	if foundDiskNumber == "" {
		return "", fmt.Errorf("disk number for device path %q volume id %q not found", devicePath, volumeID)
	}

	return foundDiskNumber, nil
}

func (d *nodeService) preparePublishTarget(target string) error {
	// On Windows, Mount will create the parent of target and mklink (create a symbolic link) at target later, so don't create a
	// directory at target now. Otherwise mklink will error: "Cannot create a file when that file already exists".
	// Instead, delete the target if it already exists (like if it was created by kubelet <1.20)
	// https://github.com/kubernetes/kubernetes/pull/88759
	klog.V(4).Infof("NodePublishVolume: removing dir %s", target)
	exists, err := d.mounter.PathExists(target)
	if err != nil {
		return fmt.Errorf("error checking path %q exists: %v", target, err)
	}

	proxyMounter, ok := (d.mounter.(*NodeMounter)).SafeFormatAndMount.Interface.(*mounter.CSIProxyMounter)
	if !ok {
		return fmt.Errorf("failed to cast mounter to csi proxy mounter")
	}

	if exists {
		if err := proxyMounter.Rmdir(target); err != nil {
			return fmt.Errorf("error Rmdir target %q: %v", target, err)
		}
	}
	return nil
}

// IsBlock checks if the given path is a block device
func (d *nodeService) IsBlockDevice(fullPath string) (bool, error) {
	return false, errors.New("unsupported")
}

func (d *nodeService) getBlockSizeBytes(devicePath string) (int64, error) {
	return 0, errors.New("unsupported")
}
