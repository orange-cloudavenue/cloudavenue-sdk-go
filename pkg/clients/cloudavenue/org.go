package clientcloudavenue

import (
	"fmt"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/uuid"
)

// getOrg returns the org object from the vCloud Director API
func (v *Client) getOrg() error {
	if !uuid.IsOrg(v.GetOrganizationID()) {
		return fmt.Errorf("invalid organization ID format: %s", v.GetOrganizationID())
	}

	v.Org = govcd.NewOrg(&v.Vmware.Client)
	v.Org.TenantContext = &govcd.TenantContext{
		OrgId:   strings.TrimPrefix(v.GetOrganizationID(), uuid.Org.String()),
		OrgName: v.GetOrganization(),
	}

	return nil
}

// getAdminOrg returns the admin org object from the vCloud Director API
func (v *Client) getAdminOrg() error {
	if !uuid.IsOrg(v.GetOrganizationID()) {
		return fmt.Errorf("invalid organization ID format: %s", v.GetOrganizationID())
	}

	v.AdminOrg = govcd.NewAdminOrg(&v.Vmware.Client)
	v.AdminOrg.TenantContext = &govcd.TenantContext{
		OrgId:   strings.TrimPrefix(v.GetOrganizationID(), uuid.Org.String()),
		OrgName: v.AdminOrg.AdminOrg.Name,
	}

	return nil
}
