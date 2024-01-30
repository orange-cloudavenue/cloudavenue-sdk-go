package clientcloudavenue

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/sethvargo/go-envconfig"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"golang.org/x/mod/semver"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/model"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

var (
	_ model.ClientOpts = (*Opts)(nil)
	c                  = &internalClient{}
)

// Opts - Is a struct that contains the options for the vmware client
type Opts struct {
	Endpoint   string `env:"ENDPOINT,overwrite"` // Deprecated - use URL instead
	URL        string `env:"URL,overwrite"`      // Computed from Org if not provided
	Username   string `env:"USERNAME,overwrite"` // Required
	Password   string `env:"PASSWORD,overwrite"` // Required
	Org        string `env:"ORG,overwrite"`      // Required
	VDC        string `env:"VDC,overwrite"`
	Debug      bool   `env:"DEBUG,overwrite"`
	VCDVersion string `env:"VCD_VERSION,overwrite,default=37.2"`
}

func (o *Opts) Validate() error {
	l := envconfig.PrefixLookuper("CLOUDAVENUE_", envconfig.OsLookuper())
	if err := envconfig.ProcessWith(context.Background(), o, l); err != nil {
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

	// Check if Organization has a valid format
	if ok := consoles.CheckOrganizationName(o.Org); !ok {
		return fmt.Errorf("the organization has an %w", errors.ErrInvalidFormat)
	}

	if o.Endpoint == "" && o.URL == "" {
		console, err := consoles.FingByOrganizationName(o.Org)
		if err != nil {
			return err
		}
		if o.Debug {
			log.Default().Printf("Found console %s with URL %s", console.GetSiteID(), console.GetURL())
		}

		o.URL = console.GetURL()
		o.Endpoint = o.URL
	}

	if o.URL == "" && o.Endpoint != "" {
		o.URL = o.Endpoint
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

// Init - Initializes the client
func Init(opts *Opts) (err error) {
	if err := opts.Validate(); err != nil {
		return err
	}

	c.token.username = opts.Username
	c.token.password = opts.Password
	c.token.org = opts.Org
	c.token.vdc = opts.VDC
	c.token.endpoint = opts.Endpoint
	c.token.debug = opts.Debug
	c.token.vcdVersion = opts.VCDVersion
	c.token.endpoint = opts.URL

	return
}

type Client struct {
	*resty.Client
	Vmware *govcd.VCDClient
}

var cache *Client

// New creates a new cloudavenue client.
func New() (*Client, error) {
	if cache != nil {
		return cache, nil
	}

	if err := c.token.RefreshToken(); err != nil {
		return nil, err
	}

	// Setup InfrAPI client
	x := resty.New().
		SetDebug(c.token.debug).
		SetHeader("Accept", "application/json;version="+c.token.vcdVersion).
		SetBaseURL(c.token.GetEndpoint()).
		SetAuthToken(c.token.GetToken())

	// Setup vmware client
	vmwareURL, err := url.Parse(fmt.Sprintf("%s/api", c.token.GetEndpoint()))
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "Failed to parse vmware url", err)
	}

	vmware := govcd.NewVCDClient(
		*vmwareURL,
		false,
		govcd.WithAPIVersion(c.token.GetVCDVersion()),
	)

	if err := vmware.SetToken(c.token.GetOrganization(), govcd.AuthorizationHeader, c.token.GetToken()); err != nil {
		return nil, fmt.Errorf("%w : %w", errors.ErrConfigureVmwareClient, err)
	}

	cache = &Client{
		Client: x,
		Vmware: vmware,
	}

	return cache, nil
}

// GetUsername - Returns the username
func (v *Client) GetUsername() string {
	return c.token.username
}

// GetOrganization - Returns the organization
func (v *Client) GetOrganization() string {
	return c.token.GetOrganization()
}

// GetOrganizationID - Returns the organization ID
func (v *Client) GetOrganizationID() string {
	return c.token.GetOrgID()
}

// GetEndpoint - Returns the API endpoint
func (v *Client) GetEndpoint() string {
	return c.token.GetEndpoint()
}

// GetDebug - Returns the debug
func (v *Client) GetDebug() bool {
	return c.token.debug
}

// GetURL - Returns the API endpoint
func (v *Client) GetURL() string {
	return c.token.GetEndpoint()
}

// GetBearerToken - Returns the bearer token
func GetBearerToken() string {
	return c.token.GetToken()
}
