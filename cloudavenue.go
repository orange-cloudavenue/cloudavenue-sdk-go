package cloudavenue

import (
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	clientS3 "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/s3"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

type Client struct {
	V1 v1.V1
}

// Opts - Is a struct that contains the options for the SDK
type ClientOpts struct {
	CloudAvenue clientcloudavenue.Opts
	Netbackup   clientnetbackup.Opts
}

func New(opts ClientOpts) (*Client, error) {
	// * Client CloudAvenue
	if opts.CloudAvenue != (clientcloudavenue.Opts{}) {
		if err := clientcloudavenue.Init(opts.CloudAvenue); err != nil {
			return nil, err
		}

		// New refresh token
		_, err := clientcloudavenue.New()
		if err != nil {
			return nil, err
		}

		if err := clientS3.Init(clientS3.Opts{
			Username:         opts.CloudAvenue.Username,
			OrganizationName: opts.CloudAvenue.Org,
			Debug:            opts.CloudAvenue.Debug,
			CAVToken:         clientcloudavenue.GetBearerToken(),
		}); err != nil {
			return nil, err
		}
	}

	// * Client Netbackup
	if opts.Netbackup != (clientnetbackup.Opts{}) {
		if err := clientnetbackup.Init(opts.Netbackup); err != nil {
			return nil, err
		}
	}

	return &Client{}, nil
}
