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
)

const (
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMName FirewallGroupDynamicSecurityGroupModelCriteriaRuleType = "VM_NAME"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleTypeVMTag  FirewallGroupDynamicSecurityGroupModelCriteriaRuleType = "VM_TAG"

	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorEquals   FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "EQUALS"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorContains FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "CONTAINS"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorStarts   FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "STARTS_WITH"
	FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperatorEnds     FirewallGroupDynamicSecurityGroupModelCriteriaRuleOperator = "ENDS_WITH"
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
)

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
