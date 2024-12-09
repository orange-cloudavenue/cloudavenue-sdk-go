package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

var _ FirewallGroupInterface = (*VDCGroup)(nil)

// CreateFirewallSecurityGroup allow creating a new security group for the VDC Group.
func (g VDCGroup) CreateFirewallSecurityGroup(securityGroupConfig *FirewallGroupSecurityGroup) (*govcd.NsxtFirewallGroup, error) {
	return g.vg.CreateNsxtFirewallGroup(&govcdtypes.NsxtFirewallGroup{
		Name:        securityGroupConfig.Name,
		Description: securityGroupConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeSecurityGroup,
		Members:     securityGroupConfig.Members,
		OwnerRef: &govcdtypes.OpenApiReference{
			ID:   g.GetID(),
			Name: g.GetName(),
		},
	})
}

// GetFirewallSecurityGroup retrieves the security group configuration for the VDC Group.
func (g VDCGroup) GetFirewallSecurityGroup(nameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if urn.IsSecurityGroup(nameOrID) {
		return g.vg.GetNsxtFirewallGroupById(nameOrID)
	}

	return g.vg.GetNsxtFirewallGroupByName(nameOrID, govcdtypes.FirewallGroupTypeSecurityGroup)
}
