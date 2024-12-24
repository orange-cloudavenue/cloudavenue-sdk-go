package v1

import (
	"fmt"
	"regexp"

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

// validate validates the firewall group model
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

// * Application Port Profile

// validate validates the application port profile model
func (appPortProfile *FirewallGroupAppPortProfileModel) Validate() error {
	// This regex validate the conformity of a TCP Port (1-65535) or/and a TCP port range (8080-8090)
	re := regexp.MustCompile(`(?m)^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])(-([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]))?$`)

	if appPortProfile == nil {
		return fmt.Errorf("appPortProfile is nil")
	}

	if appPortProfile.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if len(appPortProfile.ApplicationPorts) == 0 {
		return fmt.Errorf("ApplicationPorts is empty")
	}

	for _, appPort := range appPortProfile.ApplicationPorts {
		switch appPort.Protocol {
		case FirewallGroupAppPortProfileModelPortProtocolICMPv4, FirewallGroupAppPortProfileModelPortProtocolICMPv6:
			if len(appPort.DestinationPorts) != 0 {
				return fmt.Errorf("port must be empty for protocol %s", appPort.Protocol)
			}
		case FirewallGroupAppPortProfileModelPortProtocolTCP, FirewallGroupAppPortProfileModelPortProtocolUDP:
			if len(appPort.DestinationPorts) == 0 {
				return fmt.Errorf("port is required for protocol %s", appPort.Protocol)
			}

			for _, port := range appPort.DestinationPorts {
				if !re.Match([]byte(port)) {
					return fmt.Errorf("port %s is invalid", port)
				}
			}

		default:
			return fmt.Errorf("protocol must be one of %v", []FirewallGroupAppPortProfileModelPortProtocol{FirewallGroupAppPortProfileModelPortProtocolICMPv4, FirewallGroupAppPortProfileModelPortProtocolICMPv6, FirewallGroupAppPortProfileModelPortProtocolTCP, FirewallGroupAppPortProfileModelPortProtocolUDP})
		}
	}

	return nil
}

// Update updates the application port profile
func (fwgap *FirewallGroupAppPortProfile) Update(appPortProfileConfig *FirewallGroupAppPortProfileModel) error {
	if appPortProfileConfig == nil {
		return fmt.Errorf("appPortProfileConfig is nil")
	}

	if err := appPortProfileConfig.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	v, err := fwgap.appProfile.Update(appPortProfileConfig.toGovcdtypesNsxtAppPortProfile(fwgap.org.Org.ID, fwgap.vdcOrVDCGroup.GetID()))
	if err != nil {
		return err
	}

	fwgap.fromGovcdtypesNsxtAppPortProfile(v.NsxtAppPortProfile)
	fwgap.Scope = FirewallGroupAppPortProfileModelScope(v.NsxtAppPortProfile.Scope)

	fwgap.appProfile = v

	return nil
}

// Delete removes the application port profile
func (fwgap *FirewallGroupAppPortProfile) Delete() error {
	return fwgap.appProfile.Delete()
}
