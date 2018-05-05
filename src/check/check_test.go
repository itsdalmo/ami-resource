package check_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/itsdalmo/ami-resource/src/check"
	"testing"
)

func TestNewVersions(t *testing.T) {
	images := check.Images([]*ec2.Image{
		{
			ImageId:      aws.String("ami-00000001"),
			CreationDate: aws.String("2018-04-18T11:34:03.000Z"),
		},
		{
			ImageId:      aws.String("ami-00000003"),
			CreationDate: aws.String("2018-04-17T12:34:03.000Z"),
		},
		{
			ImageId:      aws.String("ami-00000002"),
			CreationDate: aws.String("2018-04-18T11:33:03.000Z"),
		},
	})

	t.Run("returns latest version on first check", func(t *testing.T) {
		previous := ""
		expected := []string{"ami-00000001"}

		actual := check.NewVersions(images, previous)

		if len(actual) != len(expected) {
			t.Errorf("length of expected and actual does not match")
		}

		for i := range expected {
			if actual[i] != expected[i] {
				t.Errorf("expected '%s' and got '%s'", expected[i], actual[i])
			}
		}
	})

	t.Run("returns last version if there is nothing newer", func(t *testing.T) {
		previous := "ami-00000001"
		expected := []string{"ami-00000001"}

		actual := check.NewVersions(images, previous)

		if len(actual) != len(expected) {
			t.Errorf("length of expected and actual does not match")
		}
		for i := range expected {
			if actual[i] != expected[i] {
				t.Errorf("expected '%s' and got '%s'", expected[i], actual[i])
			}
		}
	})

	t.Run("returns all images since last version", func(t *testing.T) {
		previous := "ami-00000002"
		expected := []string{
			"ami-00000002",
			"ami-00000001",
		}

		actual := check.NewVersions(images, previous)

		if len(actual) != len(expected) {
			t.Errorf("length of expected and actual does not match")
		}
		for i := range expected {
			if actual[i] != expected[i] {
				t.Errorf("expected '%s' and got '%s'", expected[i], actual[i])
			}
		}
	})
}
