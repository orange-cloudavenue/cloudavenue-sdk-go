package v1

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

// CreateFirewall creates a Distributed Firewall for the VDC Group.
func (g VDCGroup) CreateFirewall(value VDCGroupFirewallType) (*VDCGroupFirewall, error) {
	var err error

	vdcgfw := &VDCGroupFirewall{
		vdcGroup: g.vg,
	}

	// Create the Distributed Firewall
	g.vg, err = g.vg.ActivateDfw()
	if err != nil {
		return nil, err
	}

	// Enable the Distributed Firewall
	if err := vdcgfw.enableOrDisable(value.Enabled); err != nil {
		return nil, err
	}

	fw, err := g.GetFirewall()
	if err != nil {
		return nil, err
	}

	// Update the Distributed Firewall rules
	if err := fw.updateRules(value.Rules); err != nil {
		return nil, err
	}

	return fw, nil
}

// GetFirewall returns the Distributed Firewall for the VDC Group.
func (g VDCGroup) GetFirewall() (*VDCGroupFirewall, error) {
	df, err := g.vg.GetDistributedFirewall()
	if err != nil {
		return nil, err
	}

	return &VDCGroupFirewall{
		vdcGroup: df.VdcGroup,
		vgf:      df,
	}, nil
}

// GetRules returns the Distributed Firewall rules for the VDC Group.
func (g VDCGroupFirewall) GetRules() VDCGroupFirewallTypeRules {
	return g.vcdRulesToRules(g.vgf.DistributedFirewallRuleContainer.Values)
}

// IsEnabled returns true if the Distributed Firewall is enabled for the VDC Group.
func (g VDCGroupFirewall) IsEnabled() (bool, error) {
	policies, err := g.vdcGroup.GetDfwPolicies()
	if err != nil {
		return false, err
	}

	if policies.DefaultPolicy == nil || policies.DefaultPolicy.Enabled == nil {
		return false, nil
	}

	return *policies.DefaultPolicy.Enabled, nil
}

// enableOrDisable enables or disables the Distributed Firewall for the VDC Group.
func (g *VDCGroupFirewall) enableOrDisable(enable bool) (err error) {
	// Enable the Distributed Firewall
	isEnabled, err := g.IsEnabled()
	if err != nil {
		return err
	}

	if enable && !isEnabled {
		g.vdcGroup, err = g.vdcGroup.EnableDefaultPolicy()
	} else if !enable && isEnabled {
		g.vdcGroup, err = g.vdcGroup.DisableDefaultPolicy()
	}

	return err
}

// UpdateFirewall updates the Distributed Firewall for the VDC Group.
func (g *VDCGroupFirewall) UpdateFirewall(value VDCGroupFirewallType) (err error) {
	if err := g.enableOrDisable(value.Enabled); err != nil {
		return err
	}

	return g.updateRules(value.Rules)
}

// UpdateRules updates the Distributed Firewall rules for the VDC Group.
func (g *VDCGroupFirewall) updateRules(rules VDCGroupFirewallTypeRules) error {
	df, err := g.vgf.VdcGroup.UpdateDistributedFirewall(&govcdtypes.DistributedFirewallRules{
		Values: g.rulesToVCDRules(rules),
	})
	if err != nil {
		return err
	}

	// Update the VDC Group Firewall object
	g.vgf = df

	return nil
}

// Delete() deletes the Distributed Firewall for the VDC Group.
func (g *VDCGroupFirewall) Delete() error {
	if err := g.vgf.VdcGroup.DeleteAllDistributedFirewallRules(); err != nil {
		return err
	}
	if err := g.enableOrDisable(false); err != nil {
		return err
	}

	_, err := g.vdcGroup.DeactivateDfw()
	return err
}

func (g *VDCGroupFirewall) rulesToVCDRules(rules VDCGroupFirewallTypeRules) []*govcdtypes.DistributedFirewallRule {
	var vcdRules []*govcdtypes.DistributedFirewallRule
	for _, r := range rules {
		vcdRules = append(vcdRules, &govcdtypes.DistributedFirewallRule{
			Name:        r.Name,
			Enabled:     r.Enabled,
			Direction:   string(r.Direction),
			IpProtocol:  string(r.IPProtocol),
			ActionValue: string(r.Action),

			ID:                        r.ID,
			Logging:                   r.Logging,
			Description:               r.Description,
			Comments:                  r.Description,
			ApplicationPortProfiles:   r.ApplicationPortProfiles,
			SourceFirewallGroups:      r.SourceFirewallGroups,
			DestinationFirewallGroups: r.DestinationFirewallGroups,
			SourceGroupsExcluded:      r.SourceGroupsExcluded,
			DestinationGroupsExcluded: r.DestinationGroupsExcluded,
		})
	}
	return vcdRules
}

func (g *VDCGroupFirewall) vcdRulesToRules(rules []*govcdtypes.DistributedFirewallRule) VDCGroupFirewallTypeRules {
	var vcdRules VDCGroupFirewallTypeRules
	for _, r := range rules {
		vcdRules = append(vcdRules, &VDCGroupFirewallTypeRule{
			Name:       r.Name,
			Enabled:    r.Enabled,
			Direction:  VDCGroupFirewallTypeRuleDirection(r.Direction),
			IPProtocol: VDCGroupFirewallTypeRuleIPProtocol(r.IpProtocol),
			Action:     VDCGroupFirewallTypeRuleAction(r.ActionValue),

			ID:                        r.ID,
			Logging:                   r.Logging,
			Description:               r.Description,
			ApplicationPortProfiles:   r.ApplicationPortProfiles,
			SourceFirewallGroups:      r.SourceFirewallGroups,
			DestinationFirewallGroups: r.DestinationFirewallGroups,
			SourceGroupsExcluded:      r.SourceGroupsExcluded,
			DestinationGroupsExcluded: r.DestinationGroupsExcluded,
		})
	}
	return vcdRules
}
