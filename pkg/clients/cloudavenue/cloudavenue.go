/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package clientcloudavenue

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/mod/semver"
	"golang.org/x/sync/errgroup"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/model"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

var (
	_ model.ClientOpts = (*Opts)(nil)
	c                  = &internalClient{}
)

// Opts - Is a struct that contains the options for the vmware client.
type Opts struct {
	URL        string `env:"URL,overwrite"`      // Computed from Org if not provided
	Username   string `env:"USERNAME,overwrite"` // Required
	Password   string `env:"PASSWORD,overwrite"` // Required
	Org        string `env:"ORG,overwrite"`      // Required
	VDC        string `env:"VDC,overwrite"`
	Debug      bool   `env:"DEBUG,overwrite"`
	VCDVersion string `env:"VCD_VERSION,overwrite,default=37.2"`
	Dev        bool   `env:"DEV,overwrite"` // Only for development
}

func (o *Opts) Validate() error {
	l := envconfig.PrefixLookuper("CLOUDAVENUE_", envconfig.OsLookuper())
	config := &envconfig.Config{
		Target:   o,
		Lookuper: l,
	}
	if err := envconfig.ProcessWith(context.Background(), config); err != nil {
		return err
	}

	// Check if username is not empty
	if o.Username == "" {
		return fmt.Errorf("the username is %w", errors.ErrEmpty)
	}

	// Check if password is not empty
	if o.Password == "" {
		return fmt.Errorf("the password is %w", errors.ErrEmpty)
	}

	// Check if organization is not empty
	if o.Org == "" {
		return fmt.Errorf("the organization is %w", errors.ErrEmpty)
	}

	if !o.Dev {
		// Check if Organization has a valid format
		if ok := consoles.CheckOrganizationName(o.Org); !ok {
			return fmt.Errorf("the organization has an %w", errors.ErrInvalidFormat)
		}

		if o.URL == "" {
			console, err := consoles.FingByOrganizationName(o.Org)
			if err != nil {
				return err
			}
			if o.Debug {
				log.Default().Printf("Found console %s with URL %s", console.GetSiteID(), console.GetURL())
			}

			o.URL = console.GetURL()
		}
	}

	// Check if VDCVersion is not empty and semver format
	if o.VCDVersion == "" {
		return fmt.Errorf("the vcd version is %w", errors.ErrEmpty)
	}

	if semver.IsValid(o.VCDVersion) {
		return fmt.Errorf("the vcd version is %w", errors.ErrInvalidFormat)
	}

	return nil
}

type internalClient struct {
	token token
}

// Init - Initializes the client.
func Init(opts *Opts) (err error) {
	if err := opts.Validate(); err != nil {
		return err
	}

	c.token.username = opts.Username
	c.token.password = opts.Password
	c.token.org = opts.Org
	c.token.vdc = opts.VDC
	c.token.endpoint = opts.URL
	c.token.debug = opts.Debug
	c.token.vcdVersion = opts.VCDVersion
	c.token.endpoint = opts.URL

	return
}

type Client struct {
	*resty.Client
	Vmware   *govcd.VCDClient
	Org      *govcd.Org
	AdminOrg *govcd.AdminOrg
}

var cache *Client

// Refresh - Refreshes the client.
func (v *Client) Refresh() error {
	x, err := New()
	if err != nil {
		return err
	}

	*v = *x
	return nil
}

// New creates a new cloudavenue client.
func New() (*Client, error) {
	if cache != nil && !c.token.IsExpired() {
		return cache, nil
	}

	if err := c.token.RefreshToken(); err != nil {
		return nil, err
	}

	// wait group to wait for all goroutines to finish
	var wg errgroup.Group

	cache = &Client{}

	// Setup vmware client
	vmwareURL, err := url.Parse(fmt.Sprintf("%s/api", c.token.GetEndpoint()))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "Failed to parse vmware url", err)
	}

	cache.Vmware = govcd.NewVCDClient(
		*vmwareURL,
		false,
		govcd.WithAPIVersion(c.token.GetVCDVersion()),
	)

	if err := cache.Vmware.SetToken(c.token.GetOrganization(), govcd.AuthorizationHeader, c.token.GetToken()); err != nil {
		return nil, fmt.Errorf("%w : %w", errors.ErrConfigureVmwareClient, err)
	}

	// goroutine to get the org from client
	wg.Go(func() error {
		return cache.getOrg()
	})

	// goroutine to get the admin org from client
	wg.Go(func() error {
		return cache.getAdminOrg()
	})

	// Setup InfrAPI client
	wg.Go(func() error {
		cache.Client = resty.New().
			SetDebug(c.token.debug).
			SetHeader("Accept", "application/json;version="+c.token.vcdVersion).
			SetBaseURL(c.token.GetEndpoint()).
			SetAuthToken(c.token.GetToken())
		return nil
	})

	return cache, wg.Wait()
}

// GetUsername - Returns the username.
func (v *Client) GetUsername() string {
	return c.token.username
}

// GetOrganization - Returns the organization.
func (v *Client) GetOrganization() string {
	return c.token.GetOrganization()
}

// GetOrganizationID - Returns the organization ID.
func (v *Client) GetOrganizationID() string {
	return c.token.GetOrgID()
}

// GetEndpoint - Returns the API endpoint.
func (v *Client) GetEndpoint() string {
	return c.token.GetEndpoint()
}

// GetDebug - Returns the debug.
func (v *Client) GetDebug() bool {
	return c.token.debug
}

// GetURL - Returns the API endpoint.
func (v *Client) GetURL() string {
	return c.token.GetEndpoint()
}

// GetBearerToken - Returns the bearer token.
func GetBearerToken() string {
	return c.token.GetToken()
}

// MockClient - Returns the mock client.
func MockClient() *Client {
	if cache == nil {
		cache = &Client{
			Client: resty.New().
				SetHeader("Accept", "application/json;version="+c.token.vcdVersion).
				SetBaseURL("http://local.test").
				SetAuthToken(""),
		}
	}

	return cache
}
