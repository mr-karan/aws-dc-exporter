package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/directconnect"
)

// FetchAWSData represents the set of methods used to interact with AWS API.
type FetchAWSData interface {
	GetConnections() (*directconnect.Connections, error)
	GetVirtualInterfaces() (*directconnect.DescribeVirtualInterfacesOutput, error)
}

// NewDCClient creates an AWS Session and returns an initialized client with the session object embedded.
func (hub *Hub) NewDCClient(awsCreds *AWSCreds) (*DCClient, error) {
	// Initialize default config.
	config := &aws.Config{
		Region: aws.String(awsCreds.Region),
	}
	// Override Access Key and Secret Key env vars if specified in config.
	if awsCreds.AccessKey != "" && awsCreds.SecretKey != "" {
		config.Credentials = credentials.NewStaticCredentials(awsCreds.AccessKey, awsCreds.SecretKey, "")
	}
	// Initialize session with custom config embedded.
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: *config,
	})
	if err != nil {
		hub.logger.Errorf("Error creating AWS Session %s", err)
		return nil, fmt.Errorf("could not create aws session")
	}
	// Initialize EC2 Client.
	var client *directconnect.DirectConnect
	if awsCreds.RoleARN != "" {
		// Assume Role if specified
		hub.logger.Debugf("Assuming Role: %v", awsCreds.RoleARN)
		creds := stscreds.NewCredentials(sess, awsCreds.RoleARN)
		client = directconnect.New(sess, &aws.Config{Credentials: creds})
	} else {
		client = directconnect.New(sess)
	}
	return &DCClient{
		client: client,
	}, nil
}

// GetConnections returns the API response of `DescribeConnections` API Call.
func (e *DCClient) GetConnections() (*directconnect.Connections, error) {
	// Construct request params for the API Request.
	params := &directconnect.DescribeConnectionsInput{}
	resp, err := e.client.DescribeConnections(params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetVirtualInterfaces returns the API response of `DescribeVirtualInterfaces` API Call.
func (e *DCClient) GetVirtualInterfaces() (*directconnect.DescribeVirtualInterfacesOutput, error) {
	// Construct request params for the API Request.
	params := &directconnect.DescribeVirtualInterfacesInput{}
	resp, err := e.client.DescribeVirtualInterfaces(params)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
