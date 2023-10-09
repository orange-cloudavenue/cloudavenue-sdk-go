package cloudavenue

import (
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
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
	}

	// * Client Netbackup
	if opts.Netbackup != (clientnetbackup.Opts{}) {
		if err := clientnetbackup.Init(opts.Netbackup); err != nil {
			return nil, err
		}
	}

	return &Client{}, nil
}
