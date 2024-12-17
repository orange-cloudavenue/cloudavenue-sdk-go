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
		VMCriteria []govcdtypes.NsxtFirewallGroupVmCriteria `json:"vmCriteria"`
	}
)
