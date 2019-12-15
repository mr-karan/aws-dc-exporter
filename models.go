package main

import (
	"time"

	"github.com/aws/aws-sdk-go/service/directconnect/directconnectiface"
	"github.com/sirupsen/logrus"
)

// Hub represents the structure for all app wide functions and structs
type Hub struct {
	logger    *logrus.Logger
	config    config
	version   string
	exporters []Exporter
}

// cfgApp represents the structure to hold App specific configuration.
type cfgApp struct {
	LogLevel string `koanf:"log_level"`
	Jobs     []Job  `koanf:"jobs"`
}

// cfgServer represents the structure to hold Server specific configuration
type cfgServer struct {
	Name         string        `koanf:"name"`
	Address      string        `koanf:"address"`
	ReadTimeout  time.Duration `koanf:"read_timeout"`
	WriteTimeout time.Duration `koanf:"write_timeout"`
	MaxBodySize  int           `koanf:"max_body_size"`
}

// config represents the structure to hold configuration loaded from an external data source.
type config struct {
	App    cfgApp    `koanf:"app"`
	Server cfgServer `koanf:"server"`
}

// AWSCreds represents the structure to hold AWS Credentials required to create AWS session.
type AWSCreds struct {
	Region    string `koanf:"region"`
	RoleARN   string `koanf:"role_arn"`
	AccessKey string `koanf:"access_key"`
	SecretKey string `koanf:"secret_key"`
}

// Job represents a job name with AWS configs ideally meant to scrape from multiple AWS accounts/roles.
type Job struct {
	Name     string   `koanf:"name"`
	AWSCreds AWSCreds `koanf:"aws_creds"`
}

// Exporter represents the structure to hold Prometheus Descriptors. It implements prometheus.Collector
type Exporter struct {
	client FetchAWSData // Implements FetchAWSData interface which is a set of methods to interact with AWS.
	job    *Job         // Holds the Job metadata.
}

// DCClient represents the structure to hold DC Client object required to create AWS session and
// interact with AWS API SDK.
type DCClient struct {
	client directconnectiface.DirectConnectAPI
}
