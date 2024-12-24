package v1

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	FirewallGroupInterface interface {
		// * Security Group
		CreateFirewallSecurityGroup(*FirewallGroupSecurityGroupModel) (*FirewallGroupSecurityGroup, error)
		GetFirewallSecurityGroup(nameOrID string) (*FirewallGroupSecurityGroup, error)

		// * IP Set
		CreateFirewallIPSet(*FirewallGroupIPSetModel) (*FirewallGroupIPSet, error)
		GetFirewallIPSet(nameOrID string) (*FirewallGroupIPSet, error)

		// * App Port Profile
		CreateFirewallAppPortProfile(*FirewallGroupAppPortProfileModel) (*FirewallGroupAppPortProfile, error)
		GetFirewallAppPortProfile(nameOrID string) (*FirewallGroupAppPortProfile, error)
		FindFirewallAppPortProfile(name string) (*FirewallGroupAppPortProfiles, error)
	}
)

type (
	// FirewallGroup contains the basic information of a firewall group.
	FirewallGroupModel struct {
		// ID contains Firewall Group ID (URN format)
		// e.g. urn:vcloud:firewallGroup:d7f4e0b4-b83f-4a07-9f22-d242c9c0987a
		ID string `json:"id,omitempty"`
		// Name contains Firewall Group name
		// Name must be unique.
		Name string `json:"name"`
		// Description contains Firewall Group description
		Description string `json:"description,omitempty"`
	}

	// FirewallGroupIPSet contains the information of an IPSet firewall group.
	FirewallGroupIPSetModel struct {
		FirewallGroupModel `json:",inline"`
		// IP Addresses included in the group. This
		// can support IPv4 and IPv6 addresses in single, range, and CIDR formats.
		// E.g [
		//     "12.12.12.1",
		//     "10.10.10.0/24",
		//     "11.11.11.1-11.11.11.2",
		//     "2001:db8::/48",
		//	   "2001:db6:0:0:0:0:0:0-2001:db6:0:ffff:ffff:ffff:ffff:ffff",
		// ],
		IPAddresses []string `json:"ipAddresses"`
	}

	// FirewallGroupSecurityGroup contains the information of a SecurityGroup firewall group.
	FirewallGroupSecurityGroupModel struct {
		FirewallGroupModel `json:",inline"`
		// Members define list of Org VDC networks belonging to this Firewall Group
		Members []govcdtypes.OpenApiReference `json:"members"`
	}

	// FirewallGroupDynamicSecurityGroup contains the information of a DynamicSecurityGroup firewall group.
	FirewallGroupDynamicSecurityGroupModel struct {
		FirewallGroupModel `json:",inline"`
		// VmCriteria defines list of dynamic criteria that determines whether a VM belongs
		// to a dynamic firewall group. A VM needs to meet at least one criteria to belong to the
		// firewall group. In other words, the logical AND is used for rules within a single criteria
		// and the logical OR is used in between each criteria.
		Criteria FirewallGroupDynamicSecurityGroupModelCriterias `json:"Criteria"`
	}

	// FirewallGroupDynamicSecurityGroupModelCriterias defines list of dynamic criteria that determines whether a VM belongs to a dynamic firewall group.
	// A VM needs to meet at least one criteria to belong to the firewall group.
	// In other words, the logical AND is used for rules within a single criteria and the logical OR is used in between each criteria.
	// Allowed max length of Criteria is 3.
	FirewallGroupDynamicSecurityGroupModelCriterias []FirewallGroupDynamicSecurityGroupModelCriteria
	FirewallGroupDynamicSecurityGroupModelCriteria  struct {
		Rules FirewallGroupDynamicSecurityGroupModelCriteriaRules `json:"rules"`
	}

	// Allowed max length of Rules is 4.
	FirewallGroupDynamicSecurityGroupModelCriteriaRules []FirewallGroupDynamicSecurityGroupModelCriteriaRule
	FirewallGroupDynamicSecurityGroupModelCriteriaRule  struct {
		RuleType FirewallGroupDynamicSecurityGroupModelCriteriaRuleType
		Operator FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator
		Value    string
	}

	FirewallGroupDynamicSecurityGroupModelCriteriaRuleType     string
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator string

	FirewallGroupAppPortProfileModel struct {
		// ID contains App Port Profile ID (URN format)
		// e.g. urn:vcloud:applicationPortProfile::d7f4e0b4-b83f-4a07-9f22-d242c9c0987a
		ID string `json:"id,omitempty"`

		// Name contains App Port Profile name
		// Name must be unique by scope.
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`

		// ApplicationPorts contains one or more protocol and port definitions
		ApplicationPorts FirewallGroupAppPortProfileModelPorts `json:"applicationPorts"`
	}

	FirewallGroupAppPortProfileModelResponse struct {
		FirewallGroupAppPortProfileModel

		// Scope can be one of the following:
		//  * "TENANT" - v1.FirewallGroupAppPortProfileModelScopeTenant
		//  * "PROVIDER" - v1.FirewallGroupAppPortProfileModelScopeProvider
		//  * "SYSTEM" - v1.FirewallGroupAppPortProfileModelScopeSystem
		Scope FirewallGroupAppPortProfileModelScope `json:"scope"`
	}

	FirewallGroupAppPortProfileModelPorts []FirewallGroupAppPortProfileModelPort
	FirewallGroupAppPortProfileModelPort  struct {
		// Protocol can be one of the following:
		//  * "ICMPv4" - v1.FirewallGroupAppPortProfileModelPortProtocolICMPv4
		//  * "ICMPv6" - v1.FirewallGroupAppPortProfileModelPortProtocolICMPv6
		//  * "TCP" - v1.FirewallGroupAppPortProfileModelPortProtocolTCP
		//  * "UDP" - v1.FirewallGroupAppPortProfileModelPortProtocolUDP
		Protocol FirewallGroupAppPortProfileModelPortProtocol `json:"protocol"`

		// DestinationPorts is required when protocol is TCP or UDP , but can define list of ports ("1000", "1500") or port ranges ("1200-1400")
		DestinationPorts []string `json:"destinationPorts"`
	}

	// Protocol can be one of the following:
	//  * "ICMPv4" - v1.FirewallGroupAppPortProfileModelPortProtocolICMPv4
	//  * "ICMPv6" - v1.FirewallGroupAppPortProfileModelPortProtocolICMPv6
	//  * "TCP" - v1.FirewallGroupAppPortProfileModelPortProtocolTCP
	//  * "UDP" - v1.FirewallGroupAppPortProfileModelPortProtocolUDP
	FirewallGroupAppPortProfileModelPortProtocol string

	// Scope can be one of the following:
	//  * "TENANT" - v1.FirewallGroupAppPortProfileModelScopeTenant
	//  * "PROVIDER" - v1.FirewallGroupAppPortProfileModelScopeProvider
	//  * "SYSTEM" - v1.FirewallGroupAppPortProfileModelScopeSystem
	FirewallGroupAppPortProfileModelScope string
)

const (
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMName FirewallGroupDynamicSecurityGroupModelCriteriaRuleType = "VM_NAME"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTag  FirewallGroupDynamicSecurityGroupModelCriteriaRuleType = "VM_TAG"

	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorEquals   FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "EQUALS"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorContains FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "CONTAINS"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorStarts   FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "STARTS_WITH"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorEnds     FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "ENDS_WITH"

	FirewallGroupAppPortProfileModelPortProtocolICMPv4 FirewallGroupAppPortProfileModelPortProtocol = "ICMPv4"
	FirewallGroupAppPortProfileModelPortProtocolICMPv6 FirewallGroupAppPortProfileModelPortProtocol = "ICMPv6"
	FirewallGroupAppPortProfileModelPortProtocolTCP    FirewallGroupAppPortProfileModelPortProtocol = "TCP"
	FirewallGroupAppPortProfileModelPortProtocolUDP    FirewallGroupAppPortProfileModelPortProtocol = "UDP"

	FirewallGroupAppPortProfileModelScopeTenant   FirewallGroupAppPortProfileModelScope = govcdtypes.ApplicationPortProfileScopeTenant
	FirewallGroupAppPortProfileModelScopeProvider FirewallGroupAppPortProfileModelScope = govcdtypes.ApplicationPortProfileScopeProvider
	FirewallGroupAppPortProfileModelScopeSystem   FirewallGroupAppPortProfileModelScope = govcdtypes.ApplicationPortProfileScopeSystem
)

var (
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypes = []struct {
		Type        FirewallGroupDynamicSecurityGroupModelCriteriaRuleType
		Description string
	}{
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMName, "The criteria is based on the VM name."},
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTag, "The criteria is based on the VM tag."},
	}

	FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTagOperator = []struct {
		Operator    FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator
		Description string
	}{
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorEquals, "The VM name must be equal to the `value`."},
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorContains, "The `value` must be contained in the VM name."},
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorStarts, "The VM name must start with the `value`."},
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorEnds, "The VM name must end with the `value`."},
	}

	FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMNameOperator = []struct {
		Operator    FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator
		Description string
	}{
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorStarts, "The VM tag must start with the `value`."},
		{FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorContains, "The `value` must be contained in the VM tag."},
	}

	FirewallGroupAppPortProfileModelPortProtocols = []FirewallGroupAppPortProfileModelPortProtocol{
		FirewallGroupAppPortProfileModelPortProtocolICMPv4,
		FirewallGroupAppPortProfileModelPortProtocolICMPv6,
		FirewallGroupAppPortProfileModelPortProtocolTCP,
		FirewallGroupAppPortProfileModelPortProtocolUDP,
	}

	FirewallGroupAppPortProfileModelScopes = []FirewallGroupAppPortProfileModelScope{
		FirewallGroupAppPortProfileModelScopeTenant,
		FirewallGroupAppPortProfileModelScopeProvider,
		FirewallGroupAppPortProfileModelScopeSystem,
	}
)

// * Dynamic Security Group

// toGovcdtypesNsxtFirewallGroup is a method to convert a FirewallGroupDynamicSecurityGroupModel to a go-vcd firewall group.
func (fg *FirewallGroupDynamicSecurityGroupModel) toGovcdtypesNsxtFirewallGroup(ownerID, ownerName string) *govcdtypes.NsxtFirewallGroup {
	nsxtFirewallGroup := &govcdtypes.NsxtFirewallGroup{
		ID:          fg.ID,
		Name:        fg.Name,
		Description: fg.Description,
		OwnerRef: &govcdtypes.OpenApiReference{
			ID:   ownerID,
			Name: ownerName,
		},
		TypeValue:  govcdtypes.FirewallGroupTypeVmCriteria,
		VmCriteria: make([]govcdtypes.NsxtFirewallGroupVmCriteria, 0, len(fg.Criteria)),
	}

	for _, criteria := range fg.Criteria {
		vmCriteria := govcdtypes.NsxtFirewallGroupVmCriteria{
			VmCriteriaRule: make([]govcdtypes.NsxtFirewallGroupVmCriteriaRule, 0, len(criteria.Rules)),
		}

		for _, rule := range criteria.Rules {
			vmCriteria.VmCriteriaRule = append(vmCriteria.VmCriteriaRule, govcdtypes.NsxtFirewallGroupVmCriteriaRule{
				AttributeType:  string(rule.RuleType),
				Operator:       string(rule.Operator),
				AttributeValue: rule.Value,
			})
		}

		nsxtFirewallGroup.VmCriteria = append(nsxtFirewallGroup.VmCriteria, vmCriteria)
	}

	return nsxtFirewallGroup
}

// fromGovcdtypesNsxtFirewallGroup is a method to convert a go-vcd firewall group to a FirewallGroupDynamicSecurityGroupModel.
func (fg *FirewallGroupDynamicSecurityGroupModel) fromGovcdtypesNsxtFirewallGroup(nsxtFirewallGroup *govcdtypes.NsxtFirewallGroup) {
	fg.ID = nsxtFirewallGroup.ID
	fg.Name = nsxtFirewallGroup.Name
	fg.Description = nsxtFirewallGroup.Description

	fg.Criteria = make(FirewallGroupDynamicSecurityGroupModelCriterias, 0, len(nsxtFirewallGroup.VmCriteria))
	for _, criteria := range nsxtFirewallGroup.VmCriteria {
		vmCriteria := FirewallGroupDynamicSecurityGroupModelCriteria{
			Rules: make(FirewallGroupDynamicSecurityGroupModelCriteriaRules, 0, len(criteria.VmCriteriaRule)),
		}

		for _, rule := range criteria.VmCriteriaRule {
			vmCriteria.Rules = append(vmCriteria.Rules, FirewallGroupDynamicSecurityGroupModelCriteriaRule{
				RuleType: FirewallGroupDynamicSecurityGroupModelCriteriaRuleType(rule.AttributeType),
				Operator: FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator(rule.Operator),
				Value:    rule.AttributeValue,
			})
		}

		fg.Criteria = append(fg.Criteria, vmCriteria)
	}
}

// * App Port Profile

// toGovcdtypesNsxtAppPortProfile is a method to convert a FirewallGroupAppPortProfileModel to a go-vcd app port profile.
func (appPortProfile *FirewallGroupAppPortProfileModel) toGovcdtypesNsxtAppPortProfile(orgID, vDCOrVDCGroupID string) *govcdtypes.NsxtAppPortProfile {
	nsxtAppPortProfile := &govcdtypes.NsxtAppPortProfile{
		ID:               appPortProfile.ID,
		Name:             appPortProfile.Name,
		Description:      appPortProfile.Description,
		ApplicationPorts: make([]govcdtypes.NsxtAppPortProfilePort, 0, len(appPortProfile.ApplicationPorts)),
		OrgRef: &govcdtypes.OpenApiReference{
			ID: orgID,
		},
		ContextEntityId: vDCOrVDCGroupID,
		Scope:           govcdtypes.ApplicationPortProfileScopeTenant,
	}

	for _, port := range appPortProfile.ApplicationPorts {
		nsxtAppPortProfile.ApplicationPorts = append(nsxtAppPortProfile.ApplicationPorts, govcdtypes.NsxtAppPortProfilePort{
			Protocol:         string(port.Protocol),
			DestinationPorts: port.DestinationPorts,
		})
	}

	return nsxtAppPortProfile
}

// fromGovcdtypesNsxtAppPortProfile is a method to convert a go-vcd app port profile to a FirewallGroupAppPortProfileModel.
func (appPortProfile *FirewallGroupAppPortProfileModel) fromGovcdtypesNsxtAppPortProfile(nsxtAppPortProfile *govcdtypes.NsxtAppPortProfile) {
	appPortProfile.ID = nsxtAppPortProfile.ID
	appPortProfile.Name = nsxtAppPortProfile.Name
	appPortProfile.Description = nsxtAppPortProfile.Description

	appPortProfile.ApplicationPorts = make(FirewallGroupAppPortProfileModelPorts, 0, len(nsxtAppPortProfile.ApplicationPorts))
	for _, port := range nsxtAppPortProfile.ApplicationPorts {
		appPortProfile.ApplicationPorts = append(appPortProfile.ApplicationPorts, FirewallGroupAppPortProfileModelPort{
			Protocol:         FirewallGroupAppPortProfileModelPortProtocol(port.Protocol),
			DestinationPorts: port.DestinationPorts,
		})
	}
}
