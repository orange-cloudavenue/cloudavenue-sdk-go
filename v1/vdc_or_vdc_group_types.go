package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type VDCOrVDCGroupInterface interface {
	// * Global Get
	// GetName returns the name of the VDC or VDC Group
	GetName() string
	// GetID returns the ID of the VDC or VDC Group
	GetID() string
	// GetDescription returns the description of the VDC or VDC Group
	GetDescription() string

	// * Global Is
	// IsVDCGroup returns true if the object is a VDC Group
	IsVDCGroup() bool

	// * Network
	GetOpenApiOrgVdcNetworkByName(string) (*govcd.OpenApiOrgVdcNetwork, error)

	// * Security Group
	GetSecurityGroupByID(nsxtFirewallGroupID string) (*govcd.NsxtFirewallGroup, error)
	GetSecurityGroupByName(nsxtFirewallGroupName string) (*govcd.NsxtFirewallGroup, error)
	GetSecurityGroupByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error)

	// * IP Set
	GetIPSetByID(id string) (*govcd.NsxtFirewallGroup, error)
	GetIPSetByName(name string) (*govcd.NsxtFirewallGroup, error)
	GetIPSetByNameOrID(nameOrID string) (*govcd.NsxtFirewallGroup, error)
	SetIPSet(ipSetConfig *govcdtypes.NsxtFirewallGroup) (*govcd.NsxtFirewallGroup, error)

	// * NetworkConfig
	GetNetworkContextProfileByNameOrID(nameOrID string, scope VDCOrVDCGroupNetworkContextProfileScope) (*VDCOrVDCGroupNetworkContextProfile, error)
	// CreateNetworkContextProfile(name string, scope VDCOrVDCGroupNetworkContextProfileScope) (*VDCOrVDCGroupNetworkContextProfile, error)
	ListNetworkContextProfilesAttributes() any
}
