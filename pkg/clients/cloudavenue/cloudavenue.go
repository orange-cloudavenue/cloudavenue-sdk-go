package clientcloudavenue

import (
	"context"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	"github.com/sethvargo/go-envconfig"
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
}

// New creates a new cloudavenue client.
func New() (*Client, error) {
	if err := c.token.RefreshToken(); err != nil {
		return nil, err
	}

	x := resty.New().
		SetDebug(c.token.debug).
		SetHeader("Accept", "application/json;version="+c.token.vcdVersion).
		SetBaseURL(c.token.GetEndpoint()).
		SetAuthToken(c.token.GetToken())

	return &Client{x}, nil
}

// GetBearerToken - Returns the bearer token
func GetBearerToken() string {
	return c.token.GetToken()
}
