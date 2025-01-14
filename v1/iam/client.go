package iam

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

//go:generate mockgen -source=client.go -destination=mock/zz_generated_client.go

type (
	Client struct {
		clientGoVCDAdminOrg
		clientCloudavenue
	}

	clientGoVCDAdminOrg interface {
		CreateUser(*govcdtypes.User) (*govcd.OrgUser, error)
		GetUserByNameOrId(identifier string, refresh bool) (*govcd.OrgUser, error)
		GetRoleReference(roleName string) (*govcdtypes.Reference, error)
	}

	clientCloudavenue interface {
		Refresh() error
	}
)

// NewClient creates a new IAM client.
func NewClient() (*Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return &Client{
		clientCloudavenue:   c,
		clientGoVCDAdminOrg: c.AdminOrg,
	}, nil
}
