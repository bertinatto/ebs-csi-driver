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

package cloud

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	dm "github.com/kubernetes-sigs/aws-ebs-csi-driver/pkg/cloud/devicemanager"
	"github.com/kubernetes-sigs/aws-ebs-csi-driver/pkg/cloud/mocks"
	"github.com/kubernetes-sigs/aws-ebs-csi-driver/pkg/util"
)

const (
	defaultZone = "test-az"
	expZone     = "us-west-2b"
)

func TestCreateDisk(t *testing.T) {
	testCases := []struct {
		name                 string
		volumeName           string
		volState             string
		diskOptions          *DiskOptions
		expDisk              *Disk
		cleanUpFailedVolume  bool
		expErr               error
		expCreateVolumeErr   error
		expDescVolumeErr     error
		expCreateVolumeInput *ec2.CreateVolumeInput
	}{
		{
			name:       "success: normal",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(1),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: defaultZone,
			},
			expErr: nil,
		},
		{
			name:       "success: normal with gp2 options",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(1),
				VolumeType:    VolumeTypeGP2,
				Tags:          map[string]string{VolumeNameTagKey: "vol-test"},
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: defaultZone,
			},
			expErr: nil,
		},
		{
			name:       "success: normal with io2 options",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(1),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeIO2,
				IOPSPerGB:     100,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(100),
				MultiAttachEnabled: aws.Bool(false),
			},
			expErr: nil,
		},
		{
			name:       "success: normal with gp3 options",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(1),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeGP3,
				IOPS:          3000,
				Throughput:    125,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops: aws.Int64(3000),
			},
			expErr: nil,
		},
		{
			name:       "success: normal with provided zone",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: expZone,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: expZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               nil,
		},
		{
			name:       "success: normal with encrypted volume",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: expZone,
				Encrypted:        true,
				KmsKeyID:         "arn:aws:kms:us-east-1:012345678910:key/abcd1234-a123-456a-a12b-a123b4cd56ef",
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: expZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               nil,
		},
		{
			name:       "success: outpost volume",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: expZone,
				OutpostArn:       "arn:aws:outposts:us-west-2:111111111111:outpost/op-0aaa000a0aaaa00a0",
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: expZone,
				OutpostArn:       "arn:aws:outposts:us-west-2:111111111111:outpost/op-0aaa000a0aaaa00a0",
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               nil,
		},
		{
			name:       "success: empty outpost arn",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: expZone,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: expZone,
				OutpostArn:       "",
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               nil,
		},
		{
			name:       "fail: ec2.CreateVolume returned CreateVolume error",
			volumeName: "vol-test-name-error",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: expZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               fmt.Errorf("could not create volume in EC2: CreateVolume generic error"),
			expCreateVolumeErr:   fmt.Errorf("CreateVolume generic error"),
		},
		{
			name:       "fail: ec2.DescribeVolumes error after volume created",
			volumeName: "vol-test-name-error",
			volState:   "creating",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: "",
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               fmt.Errorf("failed to get an available volume in EC2: DescribeVolumes generic error"),
			expDescVolumeErr:     fmt.Errorf("DescribeVolumes generic error"),
			cleanUpFailedVolume:  true,
		},
		{
			name:       "fail: Volume is not ready to use, volume stuck in creating status and controller timed out waiting for the condition",
			volumeName: "vol-test-name-error",
			volState:   "creating",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: "",
			},
			cleanUpFailedVolume:  true,
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               fmt.Errorf("failed to get an available volume in EC2: timed out waiting for the condition"),
		},
		{
			name:       "success: normal from snapshot",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:    util.GiBToBytes(1),
				Tags:             map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				AvailabilityZone: expZone,
				SnapshotID:       "snapshot-test",
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: expZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{},
			expErr:               nil,
		},
		{
			name:       "success: io1 with too low iopsPerGB and AllowIOPSPerGBIncrease",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:          util.GiBToBytes(4),
				Tags:                   map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:             VolumeTypeIO1,
				IOPSPerGB:              1,
				AllowIOPSPerGBIncrease: true,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(100),
				MultiAttachEnabled: aws.Bool(false),
			},
			expErr: nil,
		},
		{
			name:       "fail: io1 with too low iopsPerGB",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(4),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeIO1,
				IOPSPerGB:     1,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: nil,
			expErr:               fmt.Errorf("invalid combination of volume size 4 GB and iopsPerGB 1: the resulting IOPS 4 is too low for AWS, it must be at least 100"),
		},
		{
			name:       "success: small io1 with too high iopsPerGB",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(4),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeIO1,
				IOPSPerGB:     10000,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(200),
				MultiAttachEnabled: aws.Bool(false),
			},
			expErr: nil,
		},
		{
			name:       "success: large io1 with too high iopsPerGB",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(4000),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeIO1,
				IOPSPerGB:     10000,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4000,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(64000),
				MultiAttachEnabled: aws.Bool(false),
			},
			expErr: nil,
		},
		{
			name:       "success: io2 with too low iopsPerGB and AllowIOPSPerGBIncrease",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:          util.GiBToBytes(4),
				Tags:                   map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:             VolumeTypeIO2,
				IOPSPerGB:              1,
				AllowIOPSPerGBIncrease: true,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(100),
				MultiAttachEnabled: aws.Bool(false),
			},
			expErr: nil,
		},
		{
			name:       "fail: io2 with too low iopsPerGB",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(4),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeIO2,
				IOPSPerGB:     1,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: nil,
			expErr:               fmt.Errorf("invalid combination of volume size 4 GB and iopsPerGB 1: the resulting IOPS 4 is too low for AWS, it must be at least 100"),
		},
		{
			name:       "success: small io2 with too high iopsPerGB",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(4),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeIO2,
				IOPSPerGB:     10000,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(2000),
				MultiAttachEnabled: aws.Bool(false),
			},
			expErr: nil,
		},
		{
			name:       "success: large io2 with too high iopsPerGB",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes: util.GiBToBytes(4000),
				Tags:          map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:    VolumeTypeIO2,
				IOPSPerGB:     100000,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      4000,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(64000),
				MultiAttachEnabled: aws.Bool(false),
			},
			expErr: nil,
		},
		{
			name:       "success: io2 with multi-attach enabled",
			volumeName: "vol-test-name",
			diskOptions: &DiskOptions{
				CapacityBytes:      util.GiBToBytes(1),
				Tags:               map[string]string{VolumeNameTagKey: "vol-test", AwsEbsDriverTagKey: "true"},
				VolumeType:         VolumeTypeIO2,
				IOPSPerGB:          100,
				MultiAttachEnabled: true,
			},
			expDisk: &Disk{
				VolumeID:         "vol-test",
				CapacityGiB:      1,
				AvailabilityZone: defaultZone,
			},
			expCreateVolumeInput: &ec2.CreateVolumeInput{
				Iops:               aws.Int64(100),
				MultiAttachEnabled: aws.Bool(true),
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			volState := tc.volState
			if volState == "" {
				volState = "available"
			}

			vol := &ec2.Volume{
				VolumeId:           aws.String(tc.diskOptions.Tags[VolumeNameTagKey]),
				Size:               aws.Int64(util.BytesToGiB(tc.diskOptions.CapacityBytes)),
				State:              aws.String(volState),
				AvailabilityZone:   aws.String(tc.diskOptions.AvailabilityZone),
				OutpostArn:         aws.String(tc.diskOptions.OutpostArn),
				MultiAttachEnabled: aws.Bool(tc.diskOptions.MultiAttachEnabled),
			}
			snapshot := &ec2.Snapshot{
				SnapshotId: aws.String(tc.diskOptions.SnapshotID),
				VolumeId:   aws.String("snap-test-volume"),
				State:      aws.String("completed"),
			}
			ctx := context.Background()

			if tc.expCreateVolumeInput != nil {
				matcher := eqCreateVolume(tc.expCreateVolumeInput)
				mockEC2.EXPECT().CreateVolumeWithContext(gomock.Eq(ctx), matcher).Return(vol, tc.expCreateVolumeErr)
				mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{vol}}, tc.expDescVolumeErr).AnyTimes()
				if len(tc.diskOptions.SnapshotID) > 0 {
					mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{Snapshots: []*ec2.Snapshot{snapshot}}, nil).AnyTimes()
				}
				if tc.cleanUpFailedVolume == true {
					mockEC2.EXPECT().DeleteVolumeWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DeleteVolumeOutput{}, nil)
				}
				if len(tc.diskOptions.AvailabilityZone) == 0 {
					mockEC2.EXPECT().DescribeAvailabilityZonesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeAvailabilityZonesOutput{
						AvailabilityZones: []*ec2.AvailabilityZone{
							{ZoneName: aws.String(defaultZone)},
						},
					}, nil)
				}
			}

			disk, err := c.CreateDisk(ctx, tc.volumeName, tc.diskOptions)
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("CreateDisk() failed: expected no error, got: %v", err)
				} else if tc.expErr.Error() != err.Error() {
					t.Fatalf("CreateDisk() failed: expected error %q, got: %q", tc.expErr, err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("CreateDisk() failed: expected error, got nothing")
				} else {
					if tc.expDisk.CapacityGiB != disk.CapacityGiB {
						t.Fatalf("CreateDisk() failed: expected capacity %d, got %d", tc.expDisk.CapacityGiB, disk.CapacityGiB)
					}
					if tc.expDisk.VolumeID != disk.VolumeID {
						t.Fatalf("CreateDisk() failed: expected capacity %q, got %q", tc.expDisk.VolumeID, disk.VolumeID)
					}
					if tc.expDisk.AvailabilityZone != disk.AvailabilityZone {
						t.Fatalf("CreateDisk() failed: expected availabilityZone %q, got %q", tc.expDisk.AvailabilityZone, disk.AvailabilityZone)
					}
					if tc.expDisk.OutpostArn != disk.OutpostArn {
						t.Fatalf("CreateDisk() failed: expected outpoustArn %q, got %q", tc.expDisk.OutpostArn, disk.OutpostArn)
					}
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestDeleteDisk(t *testing.T) {
	testCases := []struct {
		name     string
		volumeID string
		expResp  bool
		expErr   error
	}{
		{
			name:     "success: normal",
			volumeID: "vol-test-1234",
			expResp:  true,
			expErr:   nil,
		},
		{
			name:     "fail: DeleteVolume returned generic error",
			volumeID: "vol-test-1234",
			expResp:  false,
			expErr:   fmt.Errorf("DeleteVolume generic error"),
		},
		{
			name:     "fail: DeleteVolume returned not found error",
			volumeID: "vol-test-1234",
			expResp:  false,
			expErr:   awserr.New("InvalidVolume.NotFound", "", nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			ctx := context.Background()
			mockEC2.EXPECT().DeleteVolumeWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DeleteVolumeOutput{}, tc.expErr)

			ok, err := c.DeleteDisk(ctx, tc.volumeID)
			if err != nil && tc.expErr == nil {
				t.Fatalf("DeleteDisk() failed: expected no error, got: %v", err)
			}

			if err == nil && tc.expErr != nil {
				t.Fatal("DeleteDisk() failed: expected error, got nothing")
			}

			if tc.expResp != ok {
				t.Fatalf("DeleteDisk() failed: expected return %v, got %v", tc.expResp, ok)
			}

			mockCtrl.Finish()
		})
	}
}

func TestAttachDisk(t *testing.T) {
	testCases := []struct {
		name     string
		volumeID string
		nodeID   string
		expErr   error
	}{
		{
			name:     "success: normal",
			volumeID: "vol-test-1234",
			nodeID:   "node-1234",
			expErr:   nil,
		},
		{
			name:     "fail: AttachVolume returned generic error",
			volumeID: "vol-test-1234",
			nodeID:   "node-1234",
			expErr:   fmt.Errorf(""),
		},
		{
			name:     "fail: AttachVolume returned error volumeInUse",
			volumeID: "vol-test-1234",
			nodeID:   "node-1234",
			expErr:   awserr.New("VolumeInUse", "Volume is in use", nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			vol := &ec2.Volume{
				VolumeId:    aws.String(tc.volumeID),
				Attachments: []*ec2.VolumeAttachment{{Device: aws.String("/dev/xvdba"), InstanceId: aws.String("node-1234"), State: aws.String("attached")}},
			}

			ctx := context.Background()
			mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{vol}}, nil).AnyTimes()
			mockEC2.EXPECT().DescribeInstancesWithContext(gomock.Eq(ctx), gomock.Any()).Return(newDescribeInstancesOutput(tc.nodeID), nil)
			mockEC2.EXPECT().AttachVolumeWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.VolumeAttachment{}, tc.expErr)

			devicePath, err := c.AttachDisk(ctx, tc.volumeID, tc.nodeID)
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("AttachDisk() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("AttachDisk() failed: expected error, got nothing")
				}
				if !strings.HasPrefix(devicePath, "/dev/") {
					t.Fatal("AttachDisk() failed: expected valid device path, got empty string")
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestDetachDisk(t *testing.T) {
	testCases := []struct {
		name     string
		volumeID string
		nodeID   string
		expErr   error
	}{
		{
			name:     "success: normal",
			volumeID: "vol-test-1234",
			nodeID:   "node-1234",
			expErr:   nil,
		},
		{
			name:     "fail: DetachVolume returned generic error",
			volumeID: "vol-test-1234",
			nodeID:   "node-1234",
			expErr:   fmt.Errorf("DetachVolume generic error"),
		},
		{
			name:     "fail: DetachVolume returned not found error",
			volumeID: "vol-test-1234",
			nodeID:   "node-1234",
			expErr:   fmt.Errorf("DetachVolume not found error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			vol := &ec2.Volume{
				VolumeId:    aws.String(tc.volumeID),
				Attachments: nil,
			}

			ctx := context.Background()
			mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{vol}}, nil).AnyTimes()
			mockEC2.EXPECT().DescribeInstancesWithContext(gomock.Eq(ctx), gomock.Any()).Return(newDescribeInstancesOutput(tc.nodeID), nil)
			switch tc.name {
			case "fail: DetachVolume returned not found error":
				mockEC2.EXPECT().DetachVolumeWithContext(gomock.Eq(ctx), gomock.Any()).Return(nil, awserr.New("InvalidVolume.NotFound", "foo", fmt.Errorf("")))
			default:
				mockEC2.EXPECT().DetachVolumeWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.VolumeAttachment{}, tc.expErr)
			}

			err := c.DetachDisk(ctx, tc.volumeID, tc.nodeID)
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("DetachDisk() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("DetachDisk() failed: expected error, got nothing")
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestGetDiskByName(t *testing.T) {
	testCases := []struct {
		name             string
		volumeName       string
		volumeCapacity   int64
		availabilityZone string
		outpostArn       string
		expErr           error
	}{
		{
			name:             "success: normal",
			volumeName:       "vol-test-1234",
			volumeCapacity:   util.GiBToBytes(1),
			availabilityZone: expZone,
			expErr:           nil,
		},
		{
			name:             "success: outpost volume",
			volumeName:       "vol-test-1234",
			volumeCapacity:   util.GiBToBytes(1),
			availabilityZone: expZone,
			outpostArn:       "arn:aws:outposts:us-west-2:111111111111:outpost/op-0aaa000a0aaaa00a0",
			expErr:           nil,
		},
		{
			name:           "fail: DescribeVolumes returned generic error",
			volumeName:     "vol-test-1234",
			volumeCapacity: util.GiBToBytes(1),
			expErr:         fmt.Errorf("DescribeVolumes generic error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			vol := &ec2.Volume{
				VolumeId:         aws.String(tc.volumeName),
				Size:             aws.Int64(util.BytesToGiB(tc.volumeCapacity)),
				AvailabilityZone: aws.String(tc.availabilityZone),
				OutpostArn:       aws.String(tc.outpostArn),
			}

			ctx := context.Background()
			mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{vol}}, tc.expErr)

			disk, err := c.GetDiskByName(ctx, tc.volumeName, tc.volumeCapacity)
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("GetDiskByName() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("GetDiskByName() failed: expected error, got nothing")
				}
				if disk.CapacityGiB != util.BytesToGiB(tc.volumeCapacity) {
					t.Fatalf("GetDiskByName() failed: expected capacity %d, got %d", util.BytesToGiB(tc.volumeCapacity), disk.CapacityGiB)
				}
				if tc.availabilityZone != disk.AvailabilityZone {
					t.Fatalf("GetDiskByName() failed: expected availabilityZone %q, got %q", tc.availabilityZone, disk.AvailabilityZone)
				}
				if tc.outpostArn != disk.OutpostArn {
					t.Fatalf("GetDiskByName() failed: expected outpostArn %q, got %q", tc.outpostArn, disk.OutpostArn)
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestGetDiskByID(t *testing.T) {
	testCases := []struct {
		name             string
		volumeID         string
		availabilityZone string
		outpostArn       string
		attachments      *ec2.VolumeAttachment
		expErr           error
	}{
		{
			name:             "success: normal",
			volumeID:         "vol-test-1234",
			availabilityZone: expZone,
			attachments:      &ec2.VolumeAttachment{},
			expErr:           nil,
		},
		{
			name:             "success: outpost volume",
			volumeID:         "vol-test-1234",
			availabilityZone: expZone,
			outpostArn:       "arn:aws:outposts:us-west-2:111111111111:outpost/op-0aaa000a0aaaa00a0",
			attachments:      &ec2.VolumeAttachment{},
			expErr:           nil,
		},
		{
			name:             "success: attached instance list",
			volumeID:         "vol-test-1234",
			availabilityZone: expZone,
			outpostArn:       "arn:aws:outposts:us-west-2:111111111111:outpost/op-0aaa000a0aaaa00a0",
			attachments: &ec2.VolumeAttachment{
				InstanceId: aws.String("test-instance"),
				State:      aws.String("attached")},
			expErr: nil,
		},
		{
			name:     "fail: DescribeVolumes returned generic error",
			volumeID: "vol-test-1234",
			expErr:   fmt.Errorf("DescribeVolumes generic error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			ctx := context.Background()
			mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(
				&ec2.DescribeVolumesOutput{
					Volumes: []*ec2.Volume{
						{
							VolumeId:         aws.String(tc.volumeID),
							AvailabilityZone: aws.String(tc.availabilityZone),
							OutpostArn:       aws.String(tc.outpostArn),
							Attachments:      []*ec2.VolumeAttachment{tc.attachments},
						},
					},
				},
				tc.expErr,
			)

			disk, err := c.GetDiskByID(ctx, tc.volumeID)
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("GetDisk() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("GetDisk() failed: expected error, got nothing")
				}
				if disk.VolumeID != tc.volumeID {
					t.Fatalf("GetDisk() failed: expected ID %q, got %q", tc.volumeID, disk.VolumeID)
				}
				if tc.availabilityZone != disk.AvailabilityZone {
					t.Fatalf("GetDiskByName() failed: expected availabilityZone %q, got %q", tc.availabilityZone, disk.AvailabilityZone)
				}
				if disk.OutpostArn != tc.outpostArn {
					t.Fatalf("GetDisk() failed: expected outpostArn %q, got %q", tc.outpostArn, disk.OutpostArn)
				}
				if len(disk.Attachments) > 0 && disk.Attachments[0] != aws.StringValue(tc.attachments.InstanceId) {
					t.Fatalf("GetDisk() failed: expected attachment instance %q, got %q", aws.StringValue(tc.attachments.InstanceId), disk.Attachments[0])
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestCreateSnapshot(t *testing.T) {
	testCases := []struct {
		name            string
		snapshotName    string
		snapshotOptions *SnapshotOptions
		expInput        *ec2.CreateSnapshotInput
		expSnapshot     *Snapshot
		expErr          error
	}{
		{
			name:         "success: normal",
			snapshotName: "snap-test-name",
			snapshotOptions: &SnapshotOptions{
				Tags: map[string]string{
					SnapshotNameTagKey: "snap-test-name",
					AwsEbsDriverTagKey: "true",
					"extra-tag-key":    "extra-tag-value",
				},
			},
			expInput: &ec2.CreateSnapshotInput{
				VolumeId: aws.String("snap-test-volume"),
				DryRun:   aws.Bool(false),
				TagSpecifications: []*ec2.TagSpecification{
					{
						ResourceType: aws.String("snapshot"),
						Tags: []*ec2.Tag{
							{
								Key:   aws.String(SnapshotNameTagKey),
								Value: aws.String("snap-test-name"),
							},
							{
								Key:   aws.String(AwsEbsDriverTagKey),
								Value: aws.String("true"),
							},
							{
								Key:   aws.String("extra-tag-key"),
								Value: aws.String("extra-tag-value"),
							},
						},
					},
				},
				Description: aws.String("Created by AWS EBS CSI driver for volume snap-test-volume"),
			},
			expSnapshot: &Snapshot{
				SourceVolumeID: "snap-test-volume",
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			ec2snapshot := &ec2.Snapshot{
				SnapshotId: aws.String(tc.snapshotOptions.Tags[SnapshotNameTagKey]),
				VolumeId:   aws.String("snap-test-volume"),
				State:      aws.String("completed"),
			}

			ctx := context.Background()
			mockEC2.EXPECT().CreateSnapshotWithContext(gomock.Eq(ctx), eqCreateSnapshotInput(tc.expInput)).Return(ec2snapshot, tc.expErr)
			mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{Snapshots: []*ec2.Snapshot{ec2snapshot}}, nil).AnyTimes()

			snapshot, err := c.CreateSnapshot(ctx, tc.expSnapshot.SourceVolumeID, tc.snapshotOptions)
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("CreateSnapshot() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("CreateSnapshot() failed: expected error, got nothing")
				} else {
					if snapshot.SourceVolumeID != tc.expSnapshot.SourceVolumeID {
						t.Fatalf("CreateSnapshot() failed: expected source volume ID %s, got %v", tc.expSnapshot.SourceVolumeID, snapshot.SourceVolumeID)
					}
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestDeleteSnapshot(t *testing.T) {
	testCases := []struct {
		name         string
		snapshotName string
		expErr       error
	}{
		{
			name:         "success: normal",
			snapshotName: "snap-test-name",
			expErr:       nil,
		},
		{
			name:         "fail: delete snapshot return generic error",
			snapshotName: "snap-test-name",
			expErr:       fmt.Errorf("DeleteSnapshot generic error"),
		},
		{
			name:         "fail: delete snapshot return not found error",
			snapshotName: "snap-test-name",
			expErr:       awserr.New("InvalidSnapshot.NotFound", "", nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			ctx := context.Background()
			mockEC2.EXPECT().DeleteSnapshotWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DeleteSnapshotOutput{}, tc.expErr)

			_, err := c.DeleteSnapshot(ctx, tc.snapshotName)
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("DeleteSnapshot() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("DeleteSnapshot() failed: expected error, got nothing")
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestResizeDisk(t *testing.T) {
	testCases := []struct {
		name                string
		volumeID            string
		existingVolume      *ec2.Volume
		existingVolumeError awserr.Error
		modifiedVolume      *ec2.ModifyVolumeOutput
		modifiedVolumeError awserr.Error
		descModVolume       *ec2.DescribeVolumesModificationsOutput
		reqSizeGiB          int64
		expErr              error
	}{
		{
			name:     "success: normal",
			volumeID: "vol-test",
			existingVolume: &ec2.Volume{
				VolumeId:         aws.String("vol-test"),
				Size:             aws.Int64(1),
				AvailabilityZone: aws.String(defaultZone),
			},
			modifiedVolume: &ec2.ModifyVolumeOutput{
				VolumeModification: &ec2.VolumeModification{
					VolumeId:          aws.String("vol-test"),
					TargetSize:        aws.Int64(2),
					ModificationState: aws.String(ec2.VolumeModificationStateCompleted),
				},
			},
			reqSizeGiB: 2,
			expErr:     nil,
		},
		{
			name:     "success: normal modifying state",
			volumeID: "vol-test",
			existingVolume: &ec2.Volume{
				VolumeId:         aws.String("vol-test"),
				Size:             aws.Int64(1),
				AvailabilityZone: aws.String(defaultZone),
			},
			modifiedVolume: &ec2.ModifyVolumeOutput{
				VolumeModification: &ec2.VolumeModification{
					VolumeId:          aws.String("vol-test"),
					TargetSize:        aws.Int64(2),
					ModificationState: aws.String(ec2.VolumeModificationStateModifying),
				},
			},
			descModVolume: &ec2.DescribeVolumesModificationsOutput{
				VolumesModifications: []*ec2.VolumeModification{
					{
						VolumeId:          aws.String("vol-test"),
						TargetSize:        aws.Int64(2),
						ModificationState: aws.String(ec2.VolumeModificationStateCompleted),
					},
				},
			},
			reqSizeGiB: 2,
			expErr:     nil,
		},
		{
			name:     "success: with previous expansion",
			volumeID: "vol-test",
			existingVolume: &ec2.Volume{
				VolumeId:         aws.String("vol-test"),
				Size:             aws.Int64(2),
				AvailabilityZone: aws.String(defaultZone),
			},
			descModVolume: &ec2.DescribeVolumesModificationsOutput{
				VolumesModifications: []*ec2.VolumeModification{
					{
						VolumeId:          aws.String("vol-test"),
						TargetSize:        aws.Int64(2),
						ModificationState: aws.String(ec2.VolumeModificationStateCompleted),
					},
				},
			},
			reqSizeGiB: 2,
			expErr:     nil,
		},
		{
			name:                "fail: volume doesn't exist",
			volumeID:            "vol-test",
			existingVolumeError: awserr.New("InvalidVolume.NotFound", "", nil),
			reqSizeGiB:          2,
			expErr:              fmt.Errorf("ResizeDisk generic error"),
		},
		{
			name:     "failure: volume in modifying state",
			volumeID: "vol-test",
			existingVolume: &ec2.Volume{
				VolumeId:         aws.String("vol-test"),
				Size:             aws.Int64(1),
				AvailabilityZone: aws.String(defaultZone),
			},
			descModVolume: &ec2.DescribeVolumesModificationsOutput{
				VolumesModifications: []*ec2.VolumeModification{
					{
						VolumeId:          aws.String("vol-test"),
						TargetSize:        aws.Int64(2),
						ModificationState: aws.String(ec2.VolumeModificationStateModifying),
					},
				},
			},
			reqSizeGiB: 2,
			expErr:     fmt.Errorf("ResizeDisk generic error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			// reduce number of steps to reduce test time
			volumeModificationWaitSteps = 3
			c := newCloud(mockEC2)

			ctx := context.Background()
			if tc.existingVolume != nil || tc.existingVolumeError != nil {
				mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(
					&ec2.DescribeVolumesOutput{
						Volumes: []*ec2.Volume{
							tc.existingVolume,
						},
					}, tc.existingVolumeError)

				if tc.expErr == nil && aws.Int64Value(tc.existingVolume.Size) != tc.reqSizeGiB {
					resizedVolume := &ec2.Volume{
						VolumeId:         aws.String("vol-test"),
						Size:             aws.Int64(tc.reqSizeGiB),
						AvailabilityZone: aws.String(defaultZone),
					}
					mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(
						&ec2.DescribeVolumesOutput{
							Volumes: []*ec2.Volume{
								resizedVolume,
							},
						}, tc.existingVolumeError)
				}
			}
			if tc.modifiedVolume != nil || tc.modifiedVolumeError != nil {
				mockEC2.EXPECT().ModifyVolumeWithContext(gomock.Eq(ctx), gomock.Any()).Return(tc.modifiedVolume, tc.modifiedVolumeError).AnyTimes()
			}
			if tc.descModVolume != nil {
				mockEC2.EXPECT().DescribeVolumesModificationsWithContext(gomock.Eq(ctx), gomock.Any()).Return(tc.descModVolume, nil).AnyTimes()
			} else {
				emptyOutput := &ec2.DescribeVolumesModificationsOutput{}
				mockEC2.EXPECT().DescribeVolumesModificationsWithContext(gomock.Eq(ctx), gomock.Any()).Return(emptyOutput, nil).AnyTimes()
			}

			newSize, err := c.ResizeDisk(ctx, tc.volumeID, util.GiBToBytes(tc.reqSizeGiB))
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("ResizeDisk() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("ResizeDisk() failed: expected error, got nothing")
				} else {
					if tc.reqSizeGiB != newSize {
						t.Fatalf("ResizeDisk() failed: expected capacity %d, got %d", tc.reqSizeGiB, newSize)
					}
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestGetSnapshotByName(t *testing.T) {
	testCases := []struct {
		name            string
		snapshotName    string
		snapshotOptions *SnapshotOptions
		expSnapshot     *Snapshot
		expErr          error
	}{
		{
			name:         "success: normal",
			snapshotName: "snap-test-name",
			snapshotOptions: &SnapshotOptions{
				Tags: map[string]string{
					SnapshotNameTagKey: "snap-test-name",
					AwsEbsDriverTagKey: "true",
					"extra-tag-key":    "extra-tag-value",
				},
			},
			expSnapshot: &Snapshot{
				SourceVolumeID: "snap-test-volume",
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			ec2snapshot := &ec2.Snapshot{
				SnapshotId: aws.String(tc.snapshotOptions.Tags[SnapshotNameTagKey]),
				VolumeId:   aws.String("snap-test-volume"),
				State:      aws.String("completed"),
			}

			ctx := context.Background()
			mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{Snapshots: []*ec2.Snapshot{ec2snapshot}}, nil)

			_, err := c.GetSnapshotByName(ctx, tc.snapshotOptions.Tags[SnapshotNameTagKey])
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("GetSnapshotByName() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("GetSnapshotByName() failed: expected error, got nothing")
				}
			}

			mockCtrl.Finish()
		})
	}
}

func TestGetSnapshotByID(t *testing.T) {
	testCases := []struct {
		name            string
		snapshotName    string
		snapshotOptions *SnapshotOptions
		expSnapshot     *Snapshot
		expErr          error
	}{
		{
			name:         "success: normal",
			snapshotName: "snap-test-name",
			snapshotOptions: &SnapshotOptions{
				Tags: map[string]string{
					SnapshotNameTagKey: "snap-test-name",
					AwsEbsDriverTagKey: "true",
					"extra-tag-key":    "extra-tag-value",
				},
			},
			expSnapshot: &Snapshot{
				SourceVolumeID: "snap-test-volume",
			},
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			ec2snapshot := &ec2.Snapshot{
				SnapshotId: aws.String(tc.snapshotOptions.Tags[SnapshotNameTagKey]),
				VolumeId:   aws.String("snap-test-volume"),
				State:      aws.String("completed"),
			}

			ctx := context.Background()
			mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{Snapshots: []*ec2.Snapshot{ec2snapshot}}, nil)

			_, err := c.GetSnapshotByID(ctx, tc.snapshotOptions.Tags[SnapshotNameTagKey])
			if err != nil {
				if tc.expErr == nil {
					t.Fatalf("GetSnapshotByName() failed: expected no error, got: %v", err)
				}
			} else {
				if tc.expErr != nil {
					t.Fatal("GetSnapshotByName() failed: expected error, got nothing")
				}
			}

			mockCtrl.Finish()
		})
	}
}
func TestListSnapshots(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "success: normal",
			testFunc: func(t *testing.T) {
				expSnapshots := []*Snapshot{
					{
						SourceVolumeID: "snap-test-volume1",
						SnapshotID:     "snap-test-name1",
					},
					{
						SourceVolumeID: "snap-test-volume2",
						SnapshotID:     "snap-test-name2",
					},
				}
				ec2Snapshots := []*ec2.Snapshot{
					{
						SnapshotId: aws.String(expSnapshots[0].SnapshotID),
						VolumeId:   aws.String("snap-test-volume1"),
						State:      aws.String("completed"),
					},
					{
						SnapshotId: aws.String(expSnapshots[1].SnapshotID),
						VolumeId:   aws.String("snap-test-volume2"),
						State:      aws.String("completed"),
					},
				}

				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockEC2 := mocks.NewMockEC2(mockCtl)
				c := newCloud(mockEC2)

				ctx := context.Background()

				mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{Snapshots: ec2Snapshots}, nil)

				_, err := c.ListSnapshots(ctx, "", 0, "")
				if err != nil {
					t.Fatalf("ListSnapshots() failed: expected no error, got: %v", err)
				}
			},
		},
		{
			name: "success: with volume ID",
			testFunc: func(t *testing.T) {
				sourceVolumeID := "snap-test-volume"
				expSnapshots := []*Snapshot{
					{
						SourceVolumeID: sourceVolumeID,
						SnapshotID:     "snap-test-name1",
					},
					{
						SourceVolumeID: sourceVolumeID,
						SnapshotID:     "snap-test-name2",
					},
				}
				ec2Snapshots := []*ec2.Snapshot{
					{
						SnapshotId: aws.String(expSnapshots[0].SnapshotID),
						VolumeId:   aws.String(sourceVolumeID),
						State:      aws.String("completed"),
					},
					{
						SnapshotId: aws.String(expSnapshots[1].SnapshotID),
						VolumeId:   aws.String(sourceVolumeID),
						State:      aws.String("completed"),
					},
				}

				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockEC2 := mocks.NewMockEC2(mockCtl)
				c := newCloud(mockEC2)

				ctx := context.Background()

				mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{Snapshots: ec2Snapshots}, nil)

				resp, err := c.ListSnapshots(ctx, sourceVolumeID, 0, "")
				if err != nil {
					t.Fatalf("ListSnapshots() failed: expected no error, got: %v", err)
				}

				if len(resp.Snapshots) != len(expSnapshots) {
					t.Fatalf("Expected %d snapshots, got %d", len(expSnapshots), len(resp.Snapshots))
				}

				for _, snap := range resp.Snapshots {
					if snap.SourceVolumeID != sourceVolumeID {
						t.Fatalf("Unexpected source volume.  Expected %s, got %s", sourceVolumeID, snap.SourceVolumeID)
					}
				}
			},
		},
		{
			name: "success: max results, next token",
			testFunc: func(t *testing.T) {
				maxResults := 5
				nextTokenValue := "nextTokenValue"
				var expSnapshots []*Snapshot
				for i := 0; i < maxResults*2; i++ {
					expSnapshots = append(expSnapshots, &Snapshot{
						SourceVolumeID: "snap-test-volume1",
						SnapshotID:     fmt.Sprintf("snap-test-name%d", i),
					})
				}

				var ec2Snapshots []*ec2.Snapshot
				for i := 0; i < maxResults*2; i++ {
					ec2Snapshots = append(ec2Snapshots, &ec2.Snapshot{
						SnapshotId: aws.String(expSnapshots[i].SnapshotID),
						VolumeId:   aws.String(fmt.Sprintf("snap-test-volume%d", i)),
						State:      aws.String("completed"),
					})
				}

				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockEC2 := mocks.NewMockEC2(mockCtl)
				c := newCloud(mockEC2)

				ctx := context.Background()

				firstCall := mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{
					Snapshots: ec2Snapshots[:maxResults],
					NextToken: aws.String(nextTokenValue),
				}, nil)
				secondCall := mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{
					Snapshots: ec2Snapshots[maxResults:],
				}, nil)
				gomock.InOrder(
					firstCall,
					secondCall,
				)

				firstSnapshotsResponse, err := c.ListSnapshots(ctx, "", 5, "")
				if err != nil {
					t.Fatalf("ListSnapshots() failed: expected no error, got: %v", err)
				}

				if len(firstSnapshotsResponse.Snapshots) != maxResults {
					t.Fatalf("Expected %d snapshots, got %d", maxResults, len(firstSnapshotsResponse.Snapshots))
				}

				if firstSnapshotsResponse.NextToken != nextTokenValue {
					t.Fatalf("Expected next token value '%s' got '%s'", nextTokenValue, firstSnapshotsResponse.NextToken)
				}

				secondSnapshotsResponse, err := c.ListSnapshots(ctx, "", 0, firstSnapshotsResponse.NextToken)
				if err != nil {
					t.Fatalf("CreateSnapshot() failed: expected no error, got: %v", err)
				}

				if len(secondSnapshotsResponse.Snapshots) != maxResults {
					t.Fatalf("Expected %d snapshots, got %d", maxResults, len(secondSnapshotsResponse.Snapshots))
				}

				if secondSnapshotsResponse.NextToken != "" {
					t.Fatalf("Expected next token value to be empty got %s", secondSnapshotsResponse.NextToken)
				}
			},
		},
		{
			name: "fail: AWS DescribeSnapshotsWithContext error",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockEC2 := mocks.NewMockEC2(mockCtl)
				c := newCloud(mockEC2)

				ctx := context.Background()

				mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(nil, errors.New("test error"))

				if _, err := c.ListSnapshots(ctx, "", 0, ""); err == nil {
					t.Fatalf("ListSnapshots() failed: expected an error, got none")
				}
			},
		},
		{
			name: "fail: no snapshots ErrNotFound",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()
				mockEC2 := mocks.NewMockEC2(mockCtl)
				c := newCloud(mockEC2)

				ctx := context.Background()

				mockEC2.EXPECT().DescribeSnapshotsWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeSnapshotsOutput{}, nil)

				if _, err := c.ListSnapshots(ctx, "", 0, ""); err != nil {
					if err != ErrNotFound {
						t.Fatalf("Expected error %v, got %v", ErrNotFound, err)
					}
				} else {
					t.Fatalf("Expected error, got none")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestWaitForAttachmentState(t *testing.T) {
	testCases := []struct {
		name             string
		volumeID         string
		expectedState    string
		expectedInstance string
		expectedDevice   string
		alreadyAssigned  bool
		expectError      bool
	}{
		{
			name:             "success: attached",
			volumeID:         "vol-test-1234",
			expectedState:    volumeAttachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  false,
			expectError:      false,
		},
		{
			name:             "success: detached",
			volumeID:         "vol-test-1234",
			expectedState:    volumeDetachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  false,
			expectError:      false,
		},
		{
			name:             "success: disk not found, assumed detached",
			volumeID:         "vol-test-1234",
			expectedState:    volumeDetachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  false,
			expectError:      false,
		},
		{
			name:             "failure: disk not found, expected attached",
			volumeID:         "vol-test-1234",
			expectedState:    volumeAttachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  false,
			expectError:      true,
		},
		{
			name:             "failure: unexpected device",
			volumeID:         "vol-test-1234",
			expectedState:    volumeAttachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdbb",
			alreadyAssigned:  false,
			expectError:      true,
		},
		{
			name:             "failure: unexpected instance",
			volumeID:         "vol-test-1234",
			expectedState:    volumeAttachedState,
			expectedInstance: "1235",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  false,
			expectError:      true,
		},
		{
			name:             "failure: already assigned but wrong state",
			volumeID:         "vol-test-1234",
			expectedState:    volumeAttachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  true,
			expectError:      true,
		},
		{
			name:             "success: multiple attachments",
			volumeID:         "vol-test-1234",
			expectedState:    volumeAttachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  false,
			expectError:      false,
		},
		{
			name:             "failure: disk still attaching",
			volumeID:         "vol-test-1234",
			expectedState:    volumeAttachedState,
			expectedInstance: "1234",
			expectedDevice:   "/dev/xvdba",
			alreadyAssigned:  false,
			expectError:      true,
		},
	}

	volumeAttachmentStatePollSteps = 1

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockEC2 := mocks.NewMockEC2(mockCtrl)
			c := newCloud(mockEC2)

			attachedVol := &ec2.Volume{
				VolumeId:    aws.String(tc.volumeID),
				Attachments: []*ec2.VolumeAttachment{{Device: aws.String("/dev/xvdba"), InstanceId: aws.String("1234"), State: aws.String("attached")}},
			}

			attachingVol := &ec2.Volume{
				VolumeId:    aws.String(tc.volumeID),
				Attachments: []*ec2.VolumeAttachment{{Device: aws.String("/dev/xvdba"), InstanceId: aws.String("1234"), State: aws.String("attaching")}},
			}

			detachedVol := &ec2.Volume{
				VolumeId:    aws.String(tc.volumeID),
				Attachments: []*ec2.VolumeAttachment{{Device: aws.String("/dev/xvdba"), InstanceId: aws.String("1234"), State: aws.String("detached")}},
			}

			multipleAttachmentsVol := &ec2.Volume{
				VolumeId:    aws.String(tc.volumeID),
				Attachments: []*ec2.VolumeAttachment{{Device: aws.String("/dev/xvdba"), InstanceId: aws.String("1235"), State: aws.String("attached")}, {Device: aws.String("/dev/xvdba"), InstanceId: aws.String("1234"), State: aws.String("attached")}},
			}

			ctx := context.Background()

			switch tc.name {
			case "success: detached", "failure: already assigned but wrong state":
				mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{detachedVol}}, nil).AnyTimes()
			case "success: disk not found, assumed detached", "failure: disk not found, expected attached":
				mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(nil, awserr.New("InvalidVolume.NotFound", "foo", fmt.Errorf(""))).AnyTimes()
			case "success: multiple attachments":
				mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{multipleAttachmentsVol}}, nil).AnyTimes()
			case "failure: disk still attaching":
				mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{attachingVol}}, nil).AnyTimes()
			default:
				mockEC2.EXPECT().DescribeVolumesWithContext(gomock.Eq(ctx), gomock.Any()).Return(&ec2.DescribeVolumesOutput{Volumes: []*ec2.Volume{attachedVol}}, nil).AnyTimes()
			}

			attachment, err := c.WaitForAttachmentState(ctx, tc.volumeID, tc.expectedState, tc.expectedInstance, tc.expectedDevice, tc.alreadyAssigned)
			if tc.expectError {
				if err == nil {
					t.Fatal("WaitForAttachmentState() failed: expected error, got nothing")
				}
			} else {
				if err != nil {
					t.Fatalf("WaitForAttachmentState() failed: expected no error, got %v", err)
				}

				if tc.expectedState == volumeAttachedState {
					if attachment == nil {
						t.Fatal("WaiForAttachmentState() failed: expected attachment, got nothing")
					}
				} else {
					if attachment != nil {
						t.Fatalf("WaiForAttachmentState() failed: expected no attachment, got %v", attachment)
					}
				}
			}
		})
	}
}

func newCloud(mockEC2 EC2) Cloud {
	return &cloud{
		region: "test-region",
		dm:     dm.NewDeviceManager(),
		ec2:    mockEC2,
	}
}

func newDescribeInstancesOutput(nodeID string) *ec2.DescribeInstancesOutput {
	return &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{{
			Instances: []*ec2.Instance{
				{InstanceId: aws.String(nodeID)},
			},
		}},
	}
}

type eqCreateSnapshotInputMatcher struct {
	expected *ec2.CreateSnapshotInput
}

func eqCreateSnapshotInput(expected *ec2.CreateSnapshotInput) gomock.Matcher {
	return &eqCreateSnapshotInputMatcher{expected}
}

func (m *eqCreateSnapshotInputMatcher) Matches(x interface{}) bool {
	input, ok := x.(*ec2.CreateSnapshotInput)
	if !ok {
		return false
	}

	if input != nil {
		for _, ts := range input.TagSpecifications {
			// Because these tags are generated from a map
			// which has a random order.
			sort.SliceStable(ts.Tags, func(i, j int) bool {
				return *ts.Tags[i].Key < *ts.Tags[j].Key
			})
		}
	}

	return reflect.DeepEqual(m.expected, input)
}

func (m *eqCreateSnapshotInputMatcher) String() string {
	return m.expected.String()
}

type eqCreateVolumeMatcher struct {
	expected *ec2.CreateVolumeInput
}

func eqCreateVolume(expected *ec2.CreateVolumeInput) gomock.Matcher {
	return &eqCreateVolumeMatcher{expected}
}

func (m *eqCreateVolumeMatcher) Matches(x interface{}) bool {
	input, ok := x.(*ec2.CreateVolumeInput)
	if !ok {
		return false
	}

	if input == nil {
		return false
	}

	iops := reflect.DeepEqual(m.expected.Iops, input.Iops)
	multiAttach := reflect.DeepEqual(m.expected.MultiAttachEnabled, input.MultiAttachEnabled)
	return iops && multiAttach
}

func (m *eqCreateVolumeMatcher) String() string {
	return m.expected.String()
}
