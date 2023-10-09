package clientnetbackup

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/sethvargo/go-envconfig"
)

var c = &internalClient{}

// Opts - Is a struct that contains the options for the netbackup client
type Opts struct {
	Endpoint string `env:"ENDPOINT,default=https://backup1.cloudavenue.orange-business.com/NetBackupSelfServiceNetBackupPanels/Api"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Debug    bool   `env:"DEBUG,default=false"`
}

type internalClient struct {
	token token
}

// Init - Initializes the netbackup client
func Init(opts Opts) (err error) {
	l := envconfig.PrefixLookuper("NETBACKUP_", envconfig.OsLookuper())
	if err := envconfig.ProcessWith(context.Background(), &opts, l); err != nil {
		return err
	}

	c.token.username = opts.Username
	c.token.password = opts.Password
	c.token.endpoint = opts.Endpoint
	c.token.debug = opts.Debug

	return
}

type Client struct {
	*resty.Client
}

// new creates a new netbackup client.
func New() (*Client, error) {
	if err := c.token.RefreshToken(); err != nil {
		return nil, err
	}

	x := resty.New().
		SetDebug(c.token.debug).
		SetBaseURL(c.token.endpoint).
		SetAuthToken(c.token.GetToken())

	return &Client{x}, nil
}
