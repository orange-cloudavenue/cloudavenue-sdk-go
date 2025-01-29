/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
		Bandwidth    int       `json:"rateLimit"`
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

var (
	// Source : https://wiki.cloudavenue.orange-business.com/wiki/Network
	EdgeGatewayAllowedBandwidth = map[ClassService]struct {
		T0TotalBandwidth   int
		T1AllowedBandwidth []int
	}{
		T0ClassServiceVRFStandard: {
			T0TotalBandwidth:   T0ClassesServices[T0ClassServiceVRFStandard].TotalBandwidth,
			T1AllowedBandwidth: allowedBandwidthVRFStandard,
		},
		T0ClassServiceVRFPremium: {
			T0TotalBandwidth:   T0ClassesServices[T0ClassServiceVRFPremium].TotalBandwidth,
			T1AllowedBandwidth: allowedBandwidthVRFPremium,
		},
		T0ClassServiceVRFDedicatedMedium: {
			T0TotalBandwidth:   T0ClassesServices[T0ClassServiceVRFDedicatedMedium].TotalBandwidth,
			T1AllowedBandwidth: allowedBandwidthVRFDedicatedMedium,
		},
		T0ClassServiceVRFDedicatedLarge: {
			T0TotalBandwidth:   T0ClassesServices[T0ClassServiceVRFDedicatedLarge].TotalBandwidth,
			T1AllowedBandwidth: allowedBandwidthVRFDedicatedLarge,
		},
	}

	allowedBandwidthVRFStandard        = []int{5, 25, 50, 75, 100, 150, 200, 250, 300}                                     // 5, 25, 50, 75, 100, 150, 200, 250, 300
	allowedBandwidthVRFPremium         = append(allowedBandwidthVRFStandard, []int{400, 500, 600, 700, 800, 900, 1000}...) // 5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000
	allowedBandwidthVRFDedicatedMedium = append(allowedBandwidthVRFPremium, []int{2000}...)                                // 5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000
	allowedBandwidthVRFDedicatedLarge  = append(allowedBandwidthVRFDedicatedMedium, []int{3000, 4000, 5000, 6000}...)      // 5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 6000
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
