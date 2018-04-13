package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/itsdalmo/ami-resource/manager"
	"github.com/itsdalmo/ami-resource/models"
	"log"
	"os"
	"sort"
	"time"
)

func main() {
	var request models.CheckRequest
	if err := json.NewDecoder(os.Stdin).Decode(&request); err != nil {
		log.Fatalf("failed to unmarshal request: %s", err)
	}

	response, err := Run(request)
	if err != nil {
		log.Fatalf("check failed: %s", err)
	}

	if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
		log.Fatalf("failed to marshal response: %s", err)
	}
}

// Run (business logic)
func Run(request models.CheckRequest) (models.CheckResponse, error) {
	var response models.CheckResponse

	if err := request.Source.Validate(); err != nil {
		return response, fmt.Errorf("invalid configuration: %s", err)
	}

	// List images
	manager, err := manager.New(request.Source)
	if err != nil {
		return response, fmt.Errorf("failed to create manager: %s", err)
	}

	var filters []*ec2.Filter
	for key, value := range request.Source.Filters {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(key),
			Values: []*string{aws.String(value)},
		})
	}

	images, err := manager.DescribeImages(filters)
	if err != nil {
		return response, err
	}

	// Find new versions
	versions := NewVersions(Images(images), request.Version.ImageID)
	for _, version := range versions {
		response = append(response, models.Version{ImageID: version})
	}

	return response, nil
}

// NewVersions is exported for testing purposes.
func NewVersions(images Images, previous string) []string {
	var response []string

	// Make sure to sort the images
	sort.Sort(images)

	latest, foundLatest := images.latestVersion()
	if !foundLatest {
		// No version
		return response
	}

	last, foundPrevious := images.indexOf(previous)
	if !foundPrevious {
		// No (valid) previous version, only return latest.
		response = append(response, latest)
		return response
	}

	// Create an array with all versions since last.
	for _, image := range images[last:] {
		response = append(response, aws.StringValue(image.ImageId))
	}
	return response
}

// Images is a wrapper around an ec2.Image slice with sort implemented.
type Images []*ec2.Image

func (p Images) Len() int {
	return len(p)
}

func (p Images) Less(i, j int) bool {
	iDate := parseCreationDate(p[i])
	jDate := parseCreationDate(p[j])
	return jDate.After(iDate)
}

func (p Images) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Images) latestVersion() (string, bool) {
	if len(p) == 0 {
		return "", false
	}
	return aws.StringValue(p[len(p)-1].ImageId), true
}

func (p Images) indexOf(version string) (int, bool) {
	for i, image := range p {
		if aws.StringValue(image.ImageId) == version {
			return i, true
		}
	}
	return -1, false
}

func parseCreationDate(i *ec2.Image) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", aws.StringValue(i.CreationDate))
	if err != nil {
		panic(fmt.Errorf("failed to parse image creation date: %s", err.Error()))
	}
	return t
}