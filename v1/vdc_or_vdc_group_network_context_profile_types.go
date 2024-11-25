package v1

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	VDCOrVDCGroupNetworkContextProfile struct {
		*govcdtypes.NsxtNetworkContextProfile
	}

	VDCOrVDCGroupNetworkContextProfileScope string
)

const (
	VDCOrVDCGroupNetworkContextProfileScopeSystem   VDCOrVDCGroupNetworkContextProfileScope = "system"
	VDCOrVDCGroupNetworkContextProfileScopeProvider VDCOrVDCGroupNetworkContextProfileScope = "provider"
	VDCOrVDCGroupNetworkContextProfileScopeTenant   VDCOrVDCGroupNetworkContextProfileScope = "tenant"
)
