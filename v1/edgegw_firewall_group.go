package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

var _ FirewallGroupInterface = (*EdgeClient)(nil)

// * SecurityGroup

// CreateFirewallSecurityGroup allow creating a new security group. T
func (e *EdgeClient) CreateFirewallSecurityGroup(securityGroupConfig *FirewallGroupSecurityGroup) (*govcd.NsxtFirewallGroup, error) {
	ower := &govcdtypes.OpenApiReference{}

	if e.OwnerType.IsVDCGROUP() {
		ower.Name = e.OwnerName
	} else {
		ower.Name = e.vcdEdge.EdgeGateway.Name
		ower.ID = e.vcdEdge.EdgeGateway.ID
	}

	return e.vcdEdge.CreateNsxtFirewallGroup(&govcdtypes.NsxtFirewallGroup{
		Name:        securityGroupConfig.Name,
		Description: securityGroupConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeSecurityGroup,
		Members:     securityGroupConfig.Members,
		OwnerRef:    ower,
	})
}

// GetFirewallSecurityGroup retrieves the security group configuration for the Edge Gateway.
func (e *EdgeClient) GetFirewallSecurityGroup(nameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if urn.IsSecurityGroup(nameOrID) {
		return e.vcdEdge.GetNsxtFirewallGroupById(nameOrID)
	}

	return e.vcdEdge.GetNsxtFirewallGroupByName(nameOrID, govcdtypes.FirewallGroupTypeSecurityGroup)
}

// TODO: Add IPSet And DynamicSecurityGroup
