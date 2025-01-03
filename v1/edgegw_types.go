package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

const (
	OwnerVDC      OwnerType = "vdc"
	ownerVDCGROUP OwnerType = "vdc-group"
)

type (
	EdgeGateways    []EdgeGatewayType
	EdgeGatewayType struct {
		Tier0VrfName string    `json:"tier0VrfId"`
		EdgeID       string    `json:"edgeId"`
		EdgeName     string    `json:"edgeName"`
		OwnerType    OwnerType `json:"ownerType"`
		OwnerName    string    `json:"ownerName"`
		Description  string    `json:"description"`
		Bandwidth    Bandwidth `json:"rateLimit"`
	}

	EdgeClient struct {
		EdgeVCDInterface
		*EdgeGatewayType
		vcdEdge *govcd.NsxtEdgeGateway
	}

	// This interface contains all methods for the edge gateway in the CloudAvenue environment.
	// This list of methods are directly inherited from the go-vcloud-director/v2/govcd package.
	EdgeVCDInterface interface {
		GetNsxtFirewall() (*govcd.NsxtFirewall, error)
		UpdateNsxtFirewall(firewallRules *govcdtypes.NsxtFirewallRuleContainer) (*govcd.NsxtFirewall, error)
	}
)

// * Getters

// GetTier0VrfID - Returns the Tier0VrfID.
func (e *EdgeGatewayType) GetTier0VrfID() string {
	return e.Tier0VrfName
}

// GetT0 - Returns the Tier0VrfID (alias).
func (e *EdgeGatewayType) GetT0() string {
	return e.Tier0VrfName
}

// GetID - Returns the EdgeID.
func (e *EdgeGatewayType) GetID() string {
	return e.EdgeID
}

// GetName - Returns the EdgeName.
func (e *EdgeGatewayType) GetName() string {
	return e.EdgeName
}

// GetOwnerType - Returns the OwnerType.
func (e *EdgeGatewayType) GetOwnerType() OwnerType {
	return e.OwnerType
}

// GetOwnerName - Returns the OwnerName.
func (e *EdgeGatewayType) GetOwnerName() string {
	return e.OwnerName
}

// GetDescription - Returns the Description.
func (e *EdgeGatewayType) GetDescription() string {
	return e.Description
}
