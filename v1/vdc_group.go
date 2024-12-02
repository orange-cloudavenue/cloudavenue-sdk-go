package v1

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// GetVDCGroup retrieves the VDC Group by its name.
// It returns a pointer to the VDC Group and an error if any.
func (v *CAVVdc) GetVDCGroup(vdcGroupName string) (*VDCGroup, error) {
	if vdcGroupName == "" {
		return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	x, err := c.AdminOrg.GetVdcGroupByName(vdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDC, vdcGroupName, err)
	}

	return &VDCGroup{
		VdcGroup: x,
	}, nil
}

// GetName returns the name of the VDC Group.
func (g VDCGroup) GetName() string {
	return g.VdcGroup.VdcGroup.Name
}

// GetID returns the ID of the VDC Group.
func (g VDCGroup) GetID() string {
	return g.VdcGroup.VdcGroup.Id
}

// GetDescription returns the description of the VDC Group.
func (g VDCGroup) GetDescription() string {
	return g.VdcGroup.VdcGroup.Description
}

// IsVDCGroup return true if the object is a VDC Group.
func (g VDCGroup) IsVDCGroup() bool {
	return govcd.OwnerIsVdcGroup(g.GetID())
}

// GetSecurityGroupByID return the NSX-T security group using the ID provided in the argument.
func (g VDCGroup) GetSecurityGroupByID(nsxtFirewallGroupID string) (*govcd.NsxtFirewallGroup, error) {
	return g.VdcGroup.GetNsxtFirewallGroupById(nsxtFirewallGroupID)
}

// GetSecurityGroupByName return the NSX-T security group using the name provided in the argument.
func (g VDCGroup) GetSecurityGroupByName(nsxtFirewallGroupName string) (*govcd.NsxtFirewallGroup, error) {
	return g.VdcGroup.GetNsxtFirewallGroupByName(nsxtFirewallGroupName, govcdtypes.FirewallGroupTypeSecurityGroup)
}

// GetSecurityGroupByNameOrID return the NSX-T security group using the name or ID provided in the argument.
func (g VDCGroup) GetSecurityGroupByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if urn.IsValid(nsxtFirewallGroupNameOrID) {
		return g.GetSecurityGroupByID(nsxtFirewallGroupNameOrID)
	}

	return g.GetSecurityGroupByName(nsxtFirewallGroupNameOrID)
}

// GetIPSetByID return the NSX-T firewall group using the ID provided in the argument.
func (g VDCGroup) GetIPSetByID(id string) (*govcd.NsxtFirewallGroup, error) {
	return g.VdcGroup.GetNsxtFirewallGroupById(id)
}

// GetIPSetByName return the NSX-T firewall group using the name provided in the argument.
func (g VDCGroup) GetIPSetByName(name string) (*govcd.NsxtFirewallGroup, error) {
	return g.VdcGroup.GetNsxtFirewallGroupByName(name, govcdtypes.FirewallGroupTypeIpSet)
}

// GetIPSetByNameOrID return the NSX-T firewall group using the name or ID provided in the argument.
func (g VDCGroup) GetIPSetByNameOrID(nameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if urn.IsValid(nameOrID) {
		return g.GetIPSetByID(nameOrID)
	}

	return g.GetIPSetByName(nameOrID)
}

// SetIPSet set the NSX-T firewall group using the name provided in the argument.
func (g VDCGroup) SetIPSet(ipSetConfig *govcdtypes.NsxtFirewallGroup) (*govcd.NsxtFirewallGroup, error) {
	return g.VdcGroup.CreateNsxtFirewallGroup(ipSetConfig)
}
