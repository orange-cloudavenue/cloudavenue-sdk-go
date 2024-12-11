package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type (
	VDCGroup struct {
		// vg is a unexported VDC Group Client
		vg *govcd.VdcGroup

		// VdcGroup is a exported client for VDC Group
		VDCGroupInterface
	}

	VDCGroupInterface interface {
		GetOpenApiOrgVdcNetworkByName(string) (*govcd.OpenApiOrgVdcNetwork, error)
	}
)
