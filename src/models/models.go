package models

import (
	"errors"
)

// Source represents the configuration for the resource.
type Source struct {
	Filters            map[string]string `json:"filters"`
	AWSAccessKeyID     string            `json:"aws_access_key_id"`
	AWSSecretAccessKey string            `json:"aws_secret_access_key"`
	AWSSessionToken    string            `json:"aws_session_token"`
	AWSRegion          string            `json:"aws_region"`
}

// Validate the source configuration.
func (s *Source) Validate() error {
	if s.AWSAccessKeyID == "" {
		return errors.New("aws_access_key_id must be set")
	}
	if s.AWSSecretAccessKey == "" {
		return errors.New("aws_secret_access_key must be set")
	}
	if s.AWSRegion == "" {
		return errors.New("aws_region must be set")
	}
	return nil
}

// Metadata for the resource.
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Version for the resource.
type Version struct {
	ImageID string `json:"ami"`
}

// CheckRequest ...
type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

// CheckResponse ...
type CheckResponse []Version

// GetParameters for the resource.
type GetParameters struct{}

// GetRequest ...
type GetRequest struct {
	Source  Source        `json:"source"`
	Version Version       `json:"version"`
	Params  GetParameters `json:"params"`
}

// GetResponse ...
type GetResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// GetPackerJSON ...
type GetPackerJSON struct {
	SourceAMI string `json:"source_ami"`
}

// PutParameters for the resource.
type PutParameters struct{}

// PutRequest ...
type PutRequest struct {
	Source Source        `json:"source"`
	Params PutParameters `json:"params"`
}

// PutResponse ...
type PutResponse struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}
