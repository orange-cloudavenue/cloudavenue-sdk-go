package cloudavenue

import (
	"os"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	clientS3 "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/s3"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

type Client struct {
	V1 v1.V1
	// TODO add config
}

// Opts - Is a struct that contains the options for the SDK
type ClientOpts struct {
	CloudAvenue *clientcloudavenue.Opts
	Netbackup   *clientnetbackup.Opts
}

// New creates a new instance of the Client struct.
// It initializes the CloudAvenue and Netbackup options if they are nil.
// It then initializes the CloudAvenue client and retrieves a new refresh token.
// It fetches the console information for the organization.
// If the S3 service is enabled, it initializes the S3 client.
// If the Netbackup service is enabled, it initializes the Netbackup client.
// Finally, it returns a pointer to the Client struct and nil error if successful.
// Otherwise, it returns nil and the error encountered.
func New(opts *ClientOpts) (*Client, error) {
	if opts.CloudAvenue == nil {
		opts.CloudAvenue = new(clientcloudavenue.Opts)
	}

	if opts.Netbackup == nil {
		opts.Netbackup = new(clientnetbackup.Opts)
	}

	// * Client CloudAvenue
	if err := clientcloudavenue.Init(opts.CloudAvenue); err != nil {
		return nil, err
	}

	// New refresh token
	cavClient, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	if os.Getenv("CLOUDAVENUE_DEV") == "true" {
		return &Client{}, nil
	}

	console, err := consoles.FingByOrganizationName(cavClient.GetOrganization())
	if err != nil {
		return nil, err
	}

	// * Client S3
	if console.Services().S3.IsEnabled() {
		if err := clientS3.Init(clientS3.Opts{
			Username:         cavClient.GetUsername(),
			OrganizationName: cavClient.GetOrganization(),
			Debug:            cavClient.GetDebug(),
			CAVToken:         clientcloudavenue.GetBearerToken(),
		}); err != nil {
			return nil, err
		}
	}

	// * Client Netbackup
	if console.Services().Netbackup.IsEnabled() {
		if err := clientnetbackup.Init(opts.Netbackup, cavClient.GetOrganization()); err != nil {
			return nil, err
		}
	}

	return &Client{}, nil
}

// * Expose particular functions

type ClientConfig struct{}

func (c *Client) Config() ClientConfig {
	return ClientConfig{}
}

func (cc ClientConfig) GetOrganization() (string, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return "", err
	}

	return c.GetOrganization(), nil
}

func (cc ClientConfig) GetUsername() (string, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return "", err
	}

	return c.GetUsername(), nil
}

func (cc ClientConfig) GetURL() (string, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return "", err
	}

	return c.GetURL(), nil
}
