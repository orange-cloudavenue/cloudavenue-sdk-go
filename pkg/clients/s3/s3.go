package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/go-resty/resty/v2"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	"github.com/sethvargo/go-envconfig"
)

var c = internalClient{}

// Opts - Is a struct that contains the options for the S3 client
type Opts struct {
	OSEEndpoint      string `env:"ENDPOINT"`
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
	config := &envconfig.Config{
		Target:   &opts,
		Lookuper: l,
	}
	if err := envconfig.ProcessWith(context.Background(), config); err != nil {
		return err
	}

	c.token.cavToken = opts.CAVToken
	c.token.organizationName = opts.OrganizationName
	c.token.oseEndpoint = opts.OSEEndpoint
	c.token.s3Endpoint = opts.S3Endpoint
	c.token.debug = opts.Debug
	c.token.userName = opts.Username

	if c.token.oseEndpoint == "" {
		console, err := consoles.FingByOrganizationName(opts.OrganizationName)
		if err != nil {
			return err
		}
		if opts.Debug {
			log.Default().Printf("Found console %s with URL %s", console.GetSiteID(), console.GetURL())
		}

		if !console.Services().S3.IsEnabled() {
			return fmt.Errorf("S3 service is not available in location %s", console.GetSiteID())
		}
		c.token.oseEndpoint = console.Services().S3.GetEndpoint()
	}

	return err
}

type Client struct {
	*s3.Client
}

// New creates a new S3 client.
func New() (*Client, error) {
	if err := c.token.RefreshAccessKey(); err != nil {
		return nil, err
	}

	//config := &aws.Config{}
	//config.WithRegion("region01")
	//config.WithCredentials(credentials.NewStaticCredentialsProvider(c.token.GetAccessKey(), c.token.GetSecretKey(), ""))
	//config.WithEndpoint(c.token.GetEndpointS3())
	//if c.token.debug {
	//	config.WithLogLevel(aws.LogDebugWithHTTPBody)
	//}

	log := aws.LogRequest
	if c.token.debug {
		log = aws.LogRequestWithBody
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("region01"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.token.GetAccessKey(), c.token.GetSecretKey(), "")),
		config.WithBaseEndpoint(c.token.GetEndpointS3()),
		config.WithClientLogMode(log),
	)
	if err != nil {
		return nil, err
	}

	return &Client{s3.NewFromConfig(cfg)}, nil
}

// NewOSE - Return a new OSE client
func NewOSE() *resty.Client {
	return resty.New().
		SetDebug(GetDebug()).
		SetBaseURL(GetOSEEndpoint()).
		SetAuthToken(GetOSEToken())
}

// GetDebug - Returns the debug flag
func GetDebug() bool {
	return c.token.debug
}

// GetOrganizationName - Returns the organization name
func GetOrganizationName() string {
	return c.token.organizationName
}

// GetOrganizationID - Returns the organization ID
func GetOrganizationID() string {
	return c.token.organizationID
}

// GetOSEEndpoint - Returns the OSE endpoint
func GetOSEEndpoint() string {
	return c.token.oseEndpoint
}

// GetOSEToken - Returns the OSE token
func GetOSEToken() string {
	return c.token.cavToken
}
