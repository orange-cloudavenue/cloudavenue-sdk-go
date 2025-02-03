package edgeloadbalancer

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

//go:generate mockgen -source=client.go -destination=zz_generated_client_test.go -self_package github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer -package edgeloadbalancer -copyright_file "../../mock_header.txt"

type (
	Client interface {
		ListServiceEngineGroups(ctx context.Context, edgeGatewayID string) ([]*ServiceEngineGroupModel, error)
		GetServiceEngineGroup(ctx context.Context, edgeGatewayID, nameOrID string) (*ServiceEngineGroupModel, error)
	}

	internalClient interface {
		clientGoVCD
		clientCloudavenue
	}

	client struct {
		clientGoVCD       clientGoVCD
		clientCloudavenue clientCloudavenue
	}

	clientGoVCD interface {
		GetAllAlbServiceEngineGroupAssignments(queryParameters url.Values) ([]*govcd.NsxtAlbServiceEngineGroupAssignment, error)
	}

	clientCloudavenue interface {
		Refresh() error
		R() *resty.Request
	}
)

// NewClient creates a new edgegateway load balancer client.
func NewClient() (Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return &client{
		clientCloudavenue: c,
		clientGoVCD:       c.Vmware,
	}, nil
}

// NewFakeClient creates a new fake Org client used for testing.
func NewFakeClient(i internalClient) (Client, error) {
	return &client{
		clientCloudavenue: i,
		clientGoVCD:       i,
	}, nil
}
