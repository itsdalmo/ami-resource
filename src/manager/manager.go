package manager

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/itsdalmo/ami-resource/src/models"
	"os"
)

func newSession() (*session.Session, error) {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// Manager for AWS API calls.
type Manager struct {
	ec2    *ec2.EC2
	source models.Source
}

// New creates/initializes a new manager.
func New(source models.Source) (*Manager, error) {
	if source.AWSAccessKeyID != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", source.AWSAccessKeyID)
	}
	if source.AWSSecretAccessKey != "" {
		os.Setenv("AWS_SECRET_ACCESS_KEY", source.AWSSecretAccessKey)
	}
	if source.AWSSessionToken != "" {
		os.Setenv("AWS_SESSION_TOKEN", source.AWSSessionToken)
	}
	if source.AWSRegion != "" {
		os.Setenv("AWS_DEFAULT_REGION", source.AWSRegion)
	}

	sess, err := newSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err.Error())
	}
	return &Manager{
		ec2: ec2.New(sess, &aws.Config{
			Region: aws.String(source.AWSRegion),
		}),
		source: source,
	}, nil
}

// DescribeImages is a thin wrapper around the aws sdk function.
func (m *Manager) DescribeImages(filters []*ec2.Filter) ([]*ec2.Image, error) {
	result, err := m.ec2.DescribeImages(
		&ec2.DescribeImagesInput{Filters: filters},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to describe images: %s", err.Error())
	}
	return result.Images, nil
}
