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

// * Dynamic Security Group

// Update updates the dynamic security group
func (fwgdsg *FirewallGroupDynamicSecurityGroup) Update(dynamicSecurityGroupConfig *FirewallGroupDynamicSecurityGroupModel) error {
	if dynamicSecurityGroupConfig == nil {
		return fmt.Errorf("dynamicSecurityGroupConfig is nil")
	}

	var id, name string

	if fwgdsg.edgeClient != nil {
		id = fwgdsg.edgeClient.vcdEdge.EdgeGateway.ID
		name = fwgdsg.edgeClient.vcdEdge.EdgeGateway.Name
	} else {
		id = fwgdsg.vg.GetID()
		name = fwgdsg.vg.GetName()
	}

	if err := fwgdsg.FirewallGroupDynamicSecurityGroupModel.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	v, err := fwgdsg.fwGroup.Update(dynamicSecurityGroupConfig.toGovcdtypesNsxtFirewallGroup(id, name))
	if err != nil {
		return err
	}

	fwgdsg.FirewallGroupDynamicSecurityGroupModel.fromGovcdtypesNsxtFirewallGroup(v.NsxtFirewallGroup)
	fwgdsg.fwGroup = v

	return nil
}

// Delete removes the dynamic security group
func (fwgdsg *FirewallGroupDynamicSecurityGroup) Delete() error {
	return fwgdsg.fwGroup.Delete()
}

// validateFirewallGroupModel validates the firewall group model
func (fg *FirewallGroupDynamicSecurityGroupModel) Validate() error {
	if fg.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if len(fg.Criteria) > 3 {
		return fmt.Errorf("allowed max length of Criteria is 3")
	}

	for _, criteria := range fg.Criteria {
		if len(criteria.Rules) > 4 {
			return fmt.Errorf("allowed max length of Rules is 4")
		}

		for _, rule := range criteria.Rules {
			if rule.Value == "" {
				return fmt.Errorf("value is empty")
			}

			typeOperatorMatch := false

			switch rule.RuleType {
			case FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMName:
				for _, operator := range FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMNameOperator {
					if rule.Operator == operator.Operator {
						typeOperatorMatch = true
						break
					}
				}

				if !typeOperatorMatch {
					return fmt.Errorf("rule type and operator mismatch. If rule type is %s, then operator must be one of %v", rule.RuleType, FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMNameOperator)
				}

			case FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTag:
				for _, operator := range FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTagOperator {
					if rule.Operator == operator.Operator {
						typeOperatorMatch = true
						break
					}
				}

				if !typeOperatorMatch {
					return fmt.Errorf("rule type and operator mismatch. If rule type is %s, then operator must be one of %v", rule.RuleType, FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTagOperator)
				}
			default:
				return fmt.Errorf("rule type must be one of %v", []FirewallGroupDynamicSecurityGroupModelCriteriaRuleType{FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMName, FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTag})
			}
		}
	}

	return nil
}
