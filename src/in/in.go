package in

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/itsdalmo/ami-resource/src/manager"
	"github.com/itsdalmo/ami-resource/src/models"
	"io/ioutil"
	"path"
)

// Run (business logic)
func Run(request models.GetRequest, outputDir string) (models.GetResponse, error) {
	var response models.GetResponse

	if err := request.Source.Validate(); err != nil {
		return response, fmt.Errorf("invalid configuration: %s", err)
	}

	// Get image information
	manager, err := manager.New(request.Source)
	if err != nil {
		return response, fmt.Errorf("failed to create manager: %s", err)
	}
	images, err := manager.DescribeImages([]*ec2.Filter{
		{
			Name:   aws.String("image-id"),
			Values: []*string{aws.String(request.Version.ImageID)},
		},
	})
	if err != nil {
		return response, err
	}
	if len(images) == 0 {
		return response, fmt.Errorf("image not found: %s", request.Version.ImageID)
	}
	image := images[0]

	// Write image id
	imageID := aws.StringValue(image.ImageId)
	if err := ioutil.WriteFile(path.Join(outputDir, "id"), []byte(imageID), 0644); err != nil {
		return response, fmt.Errorf("failed to write image id: %s", err)
	}

	// Write packer json
	packerJSON, err := json.Marshal(models.GetPackerJSON{SourceAMI: imageID})
	if err != nil {
		return response, fmt.Errorf("failed to marshal packer json: %s", err)
	}
	if err := ioutil.WriteFile(path.Join(outputDir, "packer.json"), packerJSON, 0644); err != nil {
		return response, fmt.Errorf("failed to write packer json: %s", err)
	}

	// Return the response
	response.Version = models.Version{ImageID: imageID}
	response.Metadata = imageMetadata(image)

	return response, nil
}

func imageMetadata(image *ec2.Image) []models.Metadata {
	var m []models.Metadata

	m = append(m, models.Metadata{
		Name:  "name",
		Value: aws.StringValue(image.Name),
	})

	m = append(m, models.Metadata{
		Name:  "owner_id",
		Value: aws.StringValue(image.OwnerId),
	})

	m = append(m, models.Metadata{
		Name:  "creation_date",
		Value: aws.StringValue(image.CreationDate),
	})

	m = append(m, models.Metadata{
		Name:  "virtualization_type",
		Value: aws.StringValue(image.VirtualizationType),
	})

	m = append(m, models.Metadata{
		Name:  "root_device_type",
		Value: aws.StringValue(image.RootDeviceType),
	})

	return m
}
