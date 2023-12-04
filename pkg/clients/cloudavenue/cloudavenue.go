package clientcloudavenue

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/sethvargo/go-envconfig"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
)

var c = &internalClient{}

// Opts - Is a struct that contains the options for the vmware client
type Opts struct {
	Endpoint   string `env:"ENDPOINT"`
	Username   string `env:"USERNAME"`
	Password   string `env:"PASSWORD"`
	Org        string `env:"ORG"`
	VDC        string `env:"VDC"`
	Debug      bool   `env:"DEBUG,default=false"`
	VCDVersion string `env:"VCD_VERSION,default=37.2"`
}

type internalClient struct {
	token token
}

// Init - Initializes the client
func Init(opts Opts) (err error) {
	l := envconfig.PrefixLookuper("CLOUDAVENUE_", envconfig.OsLookuper())
	if err := envconfig.ProcessWith(context.Background(), &opts, l); err != nil {
		return err
	}

	c.token.username = opts.Username
	c.token.password = opts.Password
	c.token.org = opts.Org
	c.token.vdc = opts.VDC
	c.token.endpoint = opts.Endpoint
	c.token.debug = opts.Debug
	c.token.vcdVersion = opts.VCDVersion

	if c.token.endpoint == "" {
		console, err := consoles.FingByOrganizationName(opts.Org)
		if err != nil {
			return err
		}
		if opts.Debug {
			log.Default().Printf("Found console %s with URL %s", console.GetSiteID(), console.GetURL())
		}

		c.token.endpoint = console.GetURL()
	}

	return
}

type Client struct {
	*resty.Client
	Vmware *govcd.VCDClient
}

// New creates a new cloudavenue client.
func New() (*Client, error) {
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

	vmware.Client.UsingBearerToken = true
	vmware.Client.VCDAuthHeader = govcd.BearerTokenHeader
	vmware.Client.VCDToken = c.token.GetToken()
	vmware.Client.APIVersion = c.token.GetVCDVersion()
	vmware.QueryHREF = vmware.Client.VCDHREF
	vmware.QueryHREF.Path += "/query"

	return &Client{
		Client: x,
		Vmware: vmware,
	}, nil
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

// GetBearerToken - Returns the bearer token
func GetBearerToken() string {
	return c.token.GetToken()
}
