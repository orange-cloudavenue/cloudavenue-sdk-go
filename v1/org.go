package v1

import (
	"fmt"
	"strings"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/uuid"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type (
	CAVOrg struct {
		*govcd.Org
	}
)

// getOrg returns the org object from the vCloud Director API
func getOrg() (*CAVOrg, error) {
	client, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	if !uuid.IsOrg(client.GetOrganizationID()) {
		return nil, fmt.Errorf("invalid organization ID format: %s", client.GetOrganizationID())
	}

	org := govcd.NewOrg(&client.Vmware.Client)
	org.TenantContext = &govcd.TenantContext{
		OrgId:   strings.TrimPrefix(client.GetOrganizationID(), uuid.Org.String()),
		OrgName: client.GetOrganization(),
	}

	return &CAVOrg{
		Org: org,
	}, nil
}

// getAdminOrg returns the admin org object from the vCloud Director API
func getAdminOrg() (*govcd.AdminOrg, error) {
	client, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	if !uuid.IsOrg(client.GetOrganizationID()) {
		return nil, fmt.Errorf("invalid organization ID format: %s", client.GetOrganizationID())
	}

	adminOrg := govcd.NewAdminOrg(&client.Vmware.Client)
	adminOrg.TenantContext = &govcd.TenantContext{
		OrgId:   strings.TrimPrefix(client.GetOrganizationID(), uuid.Org.String()),
		OrgName: adminOrg.AdminOrg.Name,
	}

	return adminOrg, nil
}
