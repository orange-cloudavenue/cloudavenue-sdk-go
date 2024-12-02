package clientcloudavenue

import (
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// getOrg returns the org object from the vCloud Director API
func (v *Client) getOrg() (err error) {
	if !urn.IsOrg(v.GetOrganizationID()) {
		return fmt.Errorf("invalid organization ID format: %s", v.GetOrganizationID())
	}

	v.Org, err = v.Vmware.GetOrgByName(v.GetOrganization())

	return err
}

// getAdminOrg returns the admin org object from the vCloud Director API
func (v *Client) getAdminOrg() (err error) {
	if !urn.IsOrg(v.GetOrganizationID()) {
		return fmt.Errorf("invalid organization ID format: %s", v.GetOrganizationID())
	}

	v.AdminOrg, err = v.Vmware.GetAdminOrgByName(v.GetOrganization())

	return err
}
