package clientnetbackup

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/sethvargo/go-envconfig"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

var c = &internalClient{}

// Opts - Is a struct that contains the options for the netbackup client
type Opts struct {
	org      string
	Endpoint string `env:"ENDPOINT,overwrite"` // Deprecated - use URL instead
	URL      string `env:"URL,overwrite"`
	Username string `env:"USERNAME,overwrite"`
	Password string `env:"PASSWORD,overwrite"`
	Debug    bool   `env:"DEBUG,overwrite"`
}

type internalClient struct {
	token token
}

// Init - Initializes the netbackup client
func Init(opts *Opts, organizationName string) error {
	opts.org = organizationName
	if err := opts.Validate(); err != nil {
		return err
	}

	if opts.Username != "" && opts.Password != "" {
		c.token.username = opts.Username
		c.token.password = opts.Password
		c.token.endpoint = opts.URL
		c.token.debug = opts.Debug
	}

	return nil
}

func (o *Opts) Validate() error {
	l := envconfig.PrefixLookuper("NETBACKUP_", envconfig.OsLookuper())
	config := &envconfig.Config{
		Target:   o,
		Lookuper: l,
	}
	if err := envconfig.ProcessWith(context.Background(), config); err != nil {
		return err
	}

	if o.org == "" && (o.Endpoint == "" && o.URL == "") {
		return fmt.Errorf("failed to retrieve the netbackup URL. Because the organization and the URL are %w", errors.ErrEmpty)
	}

	if o.org != "" && (o.Endpoint == "" && o.URL == "") {
		console, err := consoles.FingByOrganizationName(o.org)
		if err != nil {
			return err
		}

		if !console.Services().Netbackup.IsEnabled() {
			return fmt.Errorf("the netbackup service is not enabled for the location %s", console.GetSiteID())
		}

		if o.Debug {
			log.Default().Printf("Found netbackup console %s with URL %s", console.GetSiteID(), console.Services().Netbackup.GetEndpoint())
		}
		o.URL = console.Services().Netbackup.GetEndpoint()
		o.Endpoint = o.URL
	}

	if o.URL == "" && o.Endpoint != "" {
		o.URL = o.Endpoint
	}

	if (o.Username == "" && o.Password != "") || (o.Username != "" && o.Password == "") {
		return fmt.Errorf("the username or password are %w", errors.ErrEmpty)
	}

	// username and password are not checked because they can be empty (NetBackupClient not used)

	return nil
}

type Client struct {
	*resty.Client
}

// new creates a new netbackup client.
func New() (*Client, error) {
	if !isCredentialProvider() {
		return nil, fmt.Errorf("the netbackup client is not configured")
	}

	if err := c.token.RefreshToken(); err != nil {
		return nil, err
	}

	x := resty.New().
		SetDebug(c.token.debug).
		SetBaseURL(c.token.endpoint).
		SetAuthToken(c.token.GetToken())

	return &Client{x}, nil
}

// isCredentialProvider - Returns true if the client is a credential provider
func isCredentialProvider() bool {
	return c.token.username != "" && c.token.password != "" && c.token.endpoint != ""
}
