package v1

import "github.com/vmware/go-vcloud-director/v2/govcd"

type (
	FirewallGroupSecurityGroup struct {
		// vg is a unexported VDC Group Client
		// used only for the vdcGroup
		vg VDCGroup

		// edgeClient is a unexported EdgeGateway Client
		// used only for the EdgeGateway
		edgeClient *EdgeClient

		// fwGroup is a unexported NSX-T Firewall Group
		fwGroup *govcd.NsxtFirewallGroup

		*FirewallGroupSecurityGroupModel
	}

	FirewallGroupIPSet struct {
		// vg is a unexported VDC Group Client
		// used only for the vdcGroup
		vg VDCGroup

		// edgeClient is a unexported EdgeGateway Client
		// used only for the EdgeGateway
		edgeClient *EdgeClient

		// fwGroup is a unexported NSX-T Firewall Group
		fwGroup *govcd.NsxtFirewallGroup

		*FirewallGroupIPSetModel
	}

	FirewallGroupDynamicSecurityGroup struct {
		// vg is a unexported VDC Group Client
		// used only for the vdcGroup
		vg VDCGroup

		// edgeClient is a unexported EdgeGateway Client
		// used only for the EdgeGateway
		edgeClient *EdgeClient

		// fwGroup is a unexported NSX-T Firewall Group
		fwGroup *govcd.NsxtFirewallGroup

		*FirewallGroupDynamicSecurityGroupModel
	}
)
