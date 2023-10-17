package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sethvargo/go-envconfig"
)

var c = internalClient{}

// Opts - Is a struct that contains the options for the S3 client
type Opts struct {
	OSEEndpoint      string `env:"ENDPOINT,default=https://s3console1.cloudavenue.orange-business.com"`
	S3Endpoint       string `env:"S3_ENDPOINT,default=https://s3-region01.cloudavenue.orange-business.com"`
	CAVToken         string `env:"CAV_TOKEN"`
	Debug            bool   `env:"DEBUG,default=false"`
	OrganizationName string `env:"ORGANIZATION_NAME"`
	Username         string `env:"USERNAME"`
}

type internalClient struct {
	token token
}

// Init - Initializes the client
func Init(opts Opts) (err error) {
	l := envconfig.PrefixLookuper("S3_", envconfig.OsLookuper())
	if err := envconfig.ProcessWith(context.Background(), &opts, l); err != nil {
		return err
	}

	c.token.cavToken = opts.CAVToken
	c.token.organizationName = opts.OrganizationName
	c.token.oseEndpoint = opts.OSEEndpoint
	c.token.s3Endpoint = opts.S3Endpoint
	c.token.debug = opts.Debug
	c.token.userName = opts.Username

	return
}

type Client struct {
	*s3.S3
}

// New creates a new S3 client.
func New() (*Client, error) {
	if err := c.token.RefreshAccessKey(); err != nil {
		return nil, err
	}

	config := &aws.Config{}
	config.WithRegion("region01")
	config.WithCredentials(credentials.NewStaticCredentials(c.token.GetAccessKey(), c.token.GetSecretKey(), ""))
	config.WithEndpoint(c.token.GetEndpointS3())
	if c.token.debug {
		config.WithLogLevel(aws.LogDebugWithHTTPBody)
	}

	s, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return &Client{s3.New(s)}, nil
}
