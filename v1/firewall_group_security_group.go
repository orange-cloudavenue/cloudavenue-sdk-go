package v1

import (
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

// * Security Group

// Update updates the security group
func (fwgsg *FirewallGroupSecurityGroup) Update(securityGroupConfig *FirewallGroupSecurityGroupModel) error {
	if securityGroupConfig == nil {
		return fmt.Errorf("securityGroupConfig is nil")
	}

	var id, name string

	if fwgsg.edgeClient != nil {
		id = fwgsg.edgeClient.vcdEdge.EdgeGateway.ID
		name = fwgsg.edgeClient.vcdEdge.EdgeGateway.Name
	} else {
		id = fwgsg.vg.GetID()
		name = fwgsg.vg.GetName()
	}

	v, err := fwgsg.fwGroup.Update(&govcdtypes.NsxtFirewallGroup{
		ID:          securityGroupConfig.ID,
		Name:        securityGroupConfig.Name,
		Description: securityGroupConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeSecurityGroup,
		Members:     securityGroupConfig.Members,
		OwnerRef: &govcdtypes.OpenApiReference{
			ID:   id,
			Name: name,
		},
	})
	if err != nil {
		return err
	}

	fwgsg.FirewallGroupSecurityGroupModel.ID = v.NsxtFirewallGroup.ID
	fwgsg.FirewallGroupSecurityGroupModel.Name = v.NsxtFirewallGroup.Name
	fwgsg.FirewallGroupSecurityGroupModel.Description = v.NsxtFirewallGroup.Description
	fwgsg.FirewallGroupSecurityGroupModel.Members = v.NsxtFirewallGroup.Members

	fwgsg.fwGroup = v

	return nil
}

// Delete removes the security group
func (fwgsg *FirewallGroupSecurityGroup) Delete() error {
	return fwgsg.fwGroup.Delete()
}

// * IP Set

// Update updates the IP set
func (fwgip *FirewallGroupIPSet) Update(ipSetConfig *FirewallGroupIPSetModel) error {
	if ipSetConfig == nil {
		return fmt.Errorf("ipSetConfig is nil")
	}

	var id, name string

	if fwgip.edgeClient != nil {
		id = fwgip.edgeClient.vcdEdge.EdgeGateway.ID
		name = fwgip.edgeClient.vcdEdge.EdgeGateway.Name
	} else {
		id = fwgip.vg.GetID()
		name = fwgip.vg.GetName()
	}

	v, err := fwgip.fwGroup.Update(&govcdtypes.NsxtFirewallGroup{
		ID:          ipSetConfig.ID,
		Name:        ipSetConfig.Name,
		Description: ipSetConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeIpSet,
		IpAddresses: ipSetConfig.IPAddresses,
		OwnerRef: &govcdtypes.OpenApiReference{
			ID:   id,
			Name: name,
		},
	})
	if err != nil {
		return err
	}

	fwgip.FirewallGroupIPSetModel.ID = v.NsxtFirewallGroup.ID
	fwgip.FirewallGroupIPSetModel.Name = v.NsxtFirewallGroup.Name
	fwgip.FirewallGroupIPSetModel.Description = v.NsxtFirewallGroup.Description
	fwgip.FirewallGroupIPSetModel.IPAddresses = v.NsxtFirewallGroup.IpAddresses

	fwgip.fwGroup = v

	return nil
}

// Delete removes the IP set
func (fwgip *FirewallGroupIPSet) Delete() error {
	return fwgip.fwGroup.Delete()
}
