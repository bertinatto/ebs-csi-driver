package cloud

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/bertinatto/ebs-csi-driver/pkg/util"
)

const (
	// TODO: what should be the default size?
	// DefaultVolumeSize represents the default volume size.
	DefaultVolumeSize int64 = 1 * 1024 * 1024 * 1024

	// VolumeNameTagKey is the key value that refers to the volume's name.
	VolumeNameTagKey = "com.amazon.aws.csi.volume"

	// VolumeTypeIO1 represents a provisioned IOPS SSD type of volume.
	VolumeTypeIO1 = "io1"

	// VolumeTypeGP2 represents a general purpose SSD type of volume.
	VolumeTypeGP2 = "gp2"

	// VolumeTypeSC1 represents a cold HDD (sc1) type of volume.
	VolumeTypeSC1 = "sc1"

	// VolumeTypeST1 represents a throughput-optimized HDD type of volume.
	VolumeTypeST1 = "st1"

	// MinTotalIOPS represents the minimum Input Output per second.
	MinTotalIOPS int64 = 100

	// MaxTotalIOPS represents the maximum Input Output per second.
	MaxTotalIOPS int64 = 20000

	// DefaultVolumeType specifies which storage to use for newly created Volumes.
	DefaultVolumeType = VolumeTypeGP2
)

type Compute interface {
	GetMetadata() *Metadata
	CreateDisk(string, *DiskOptions) (*Disk, error)
	DeleteDisk(string) (bool, error)
	AttachDisk(string, string) error
	DetachDisk(string, string) error
	GetDiskByNameAndSize(string, int64) (*Disk, error)
}

type DiskOptions struct {
	CapacityBytes int64
	Tags          map[string]string
	VolumeType    string
	IOPSPerGB     int64
}

type Disk struct {
	VolumeID    string
	CapacityGiB int64
}

type Cloud struct {
	metadata *Metadata
	ec2      *ec2.EC2
}

var _ Compute = &Cloud{}

func NewCloud() (*Cloud, error) {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize AWS session: %v", err)
	}

	metadata, err := NewMetadata()
	if err != nil {
		return nil, fmt.Errorf("could not get metadata from AWS: %v", err)
	}

	provider := []credentials.Provider{
		&credentials.EnvProvider{},
		&ec2rolecreds.EC2RoleProvider{Client: ec2metadata.New(sess)},
		&credentials.SharedCredentialsProvider{},
	}

	awsConfig := &aws.Config{
		Region:      &metadata.Region,
		Credentials: credentials.NewChainCredentials(provider),
	}
	awsConfig = awsConfig.WithCredentialsChainVerboseErrors(true)

	return &Cloud{
		metadata: metadata,
		ec2:      ec2.New(session.New(awsConfig)),
	}, nil
}

func (c *Cloud) GetMetadata() *Metadata {
	return c.metadata
}

func (c *Cloud) CreateDisk(volumeName string, diskOptions *DiskOptions) (*Disk, error) {
	var createType string
	var iops int64
	capacityGiB := util.BytesToGiB(diskOptions.CapacityBytes)

	switch diskOptions.VolumeType {
	case VolumeTypeGP2, VolumeTypeSC1, VolumeTypeST1:
		createType = diskOptions.VolumeType
	case VolumeTypeIO1:
		createType = diskOptions.VolumeType
		iops = capacityGiB * diskOptions.IOPSPerGB
		if iops < MinTotalIOPS {
			iops = MinTotalIOPS
		}
		if iops > MaxTotalIOPS {
			iops = MaxTotalIOPS
		}
	case "":
		createType = DefaultVolumeType
	default:
		return nil, fmt.Errorf("invalid AWS VolumeType %q", diskOptions.VolumeType)
	}

	var tags []*ec2.Tag
	for key, value := range diskOptions.Tags {
		tags = append(tags, &ec2.Tag{Key: &key, Value: &value})
	}
	tagSpec := ec2.TagSpecification{
		ResourceType: aws.String("volume"),
		Tags:         tags,
	}

	m := c.GetMetadata()
	request := &ec2.CreateVolumeInput{
		AvailabilityZone:  aws.String(m.AvailabilityZone),
		Size:              aws.Int64(capacityGiB),
		VolumeType:        aws.String(createType),
		TagSpecifications: []*ec2.TagSpecification{&tagSpec},
	}
	if iops > 0 {
		request.Iops = aws.Int64(iops)
	}

	response, err := c.ec2.CreateVolume(request)
	if err != nil {
		return nil, fmt.Errorf("could not create volume in EC2: %v", err)
	}

	volumeID := aws.StringValue(response.VolumeId)
	if len(volumeID) == 0 {
		return nil, fmt.Errorf("volume ID was not returned by CreateVolume")
	}

	size := aws.Int64Value(response.Size)
	if size == 0 {
		return nil, fmt.Errorf("disk size was not returned by CreateVolume")
	}

	return &Disk{CapacityGiB: size, VolumeID: volumeID}, nil
}

var ErrVolumeNotFound = errors.New("Volume was not found")

func (c *Cloud) DeleteDisk(volumeID string) (bool, error) {
	request := &ec2.DeleteVolumeInput{VolumeId: &volumeID}
	if _, err := c.ec2.DeleteVolume(request); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "InvalidVolume.NotFound" {
				return false, ErrVolumeNotFound
			}
		}
		return false, fmt.Errorf("DeleteDisk could not delete volume: %v", err)
	}
	return true, nil
}

func (c *Cloud) AttachDisk(volumeID, nodeID string) error {
	// TODO: choose a valid and non-duplicate device name
	device := "/dev/xvdbc"
	request := &ec2.AttachVolumeInput{
		Device:     aws.String(device),
		InstanceId: aws.String(nodeID),
		VolumeId:   aws.String(volumeID),
	}

	_, err := c.ec2.AttachVolume(request)
	if err != nil {
		return fmt.Errorf("could not attach volume %q to node %q: %v", volumeID, nodeID, err)
	}

	return nil
}

func (c *Cloud) DetachDisk(volumeID, nodeID string) error {
	request := &ec2.DetachVolumeInput{
		InstanceId: aws.String(nodeID),
		VolumeId:   aws.String(volumeID),
	}

	_, err := c.ec2.DetachVolume(request)
	if err != nil {
		return fmt.Errorf("could not detach volume %q from node %q: %v", volumeID, nodeID, err)
	}

	return nil
}

var ErrMultiDisks = errors.New("Multiple disks with same name")
var ErrDiskExistsDiffSize = errors.New("There is already a disk with same name and different size")

func (c *Cloud) GetDiskByNameAndSize(name string, capacityBytes int64) (*Disk, error) {
	var volumes []*ec2.Volume
	var nextToken *string

	request := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:" + VolumeNameTagKey),
				Values: []*string{aws.String(name)},
			},
		},
	}
	for {
		response, err := c.ec2.DescribeVolumes(request)
		if err != nil {
			return nil, err
		}
		for _, volume := range response.Volumes {
			volumes = append(volumes, volume)
		}
		nextToken = response.NextToken
		if aws.StringValue(nextToken) == "" {
			break
		}
		request.NextToken = nextToken
	}

	if len(volumes) > 1 {
		return nil, ErrMultiDisks
	}

	if len(volumes) == 0 {
		return nil, nil
	}

	volSizeBytes := aws.Int64Value(volumes[0].Size)
	if volSizeBytes != util.BytesToGiB(capacityBytes) {
		return nil, ErrDiskExistsDiffSize
	}

	return &Disk{
		VolumeID:    aws.StringValue(volumes[0].VolumeId),
		CapacityGiB: volSizeBytes,
	}, nil
}
