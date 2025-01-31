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
	"fmt"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type (
	Tier0 struct{}
	T0s   []T0
	T0    struct {
		Tier0Vrf          string       `json:"tier0_vrf"`
		Tier0Provider     string       `json:"tier0_provider"`
		Tier0ClassService string       `json:"tier0_class_service"`
		ClassService      ClassService `json:"class_service"`
		Services          T0Services   `json:"services"`
	}
	T0Services []T0Service
	T0Service  struct {
		Service string `json:"service"`
		VLANID  any    `json:"vlanId"`
	}

	ClassService string
)

const (
	// TOClassServiceVRFStandard - VRF Standard.
	T0ClassServiceVRFStandard ClassService = "VRF_STANDARD"
	// T0ClassServiceVRFPremium - VRF Premium.
	T0ClassServiceVRFPremium ClassService = "VRF_PREMIUM"
	// T0ClassServiceVRFDedicatedMedium - VRF Dedicated Medium.
	T0ClassServiceVRFDedicatedMedium ClassService = "VRF_DEDICATED_MEDIUM"
	// T0ClassServiceVRFDedicatedLarge - VRF Dedicated Large.
	T0ClassServiceVRFDedicatedLarge ClassService = "VRF_DEDICATED_LARGE"
)

var T0ClassesServices = map[ClassService]struct {
	TotalBandwidth int
}{
	T0ClassServiceVRFStandard: {
		TotalBandwidth: 300,
	},
	T0ClassServiceVRFPremium: {
		TotalBandwidth: 1000,
	},
	T0ClassServiceVRFDedicatedMedium: {
		TotalBandwidth: 3500,
	},
	T0ClassServiceVRFDedicatedLarge: {
		TotalBandwidth: 10000,
	},
}

// * T0

// GetTier0ClassService - Returns the Tier0ClassService.
func (t *T0) GetTier0ClassService() string {
	return t.Tier0ClassService
}

// GetName - Returns the Tier0Vrf.
func (t *T0) GetName() string {
	return t.Tier0Vrf
}

// GetTier0Vrf - Returns the Tier0Vrf.
func (t *T0) GetTier0Vrf() string {
	return t.Tier0Vrf
}

// GetTier0Provider - Returns the Tier0Provider.
func (t *T0) GetTier0Provider() string {
	return t.Tier0Provider
}

// GetClassService - Returns the ClassService.
func (t *T0) GetClassService() ClassService {
	return t.ClassService
}

// GetServices - Returns the Services.
func (t *T0) GetServices() T0Services {
	return t.Services
}

// * T0Service

// GetService - Returns the Service.
func (t *T0Service) GetService() string {
	return t.Service
}

// GetVLANID - Returns the VLANID.
func (t *T0Service) GetVLANID() any {
	return t.VLANID
}

// * ClassService

// IsVRFStandard - Returns true if the ClassService is VRFStandard.
func (c ClassService) IsVRFStandard() bool {
	return c == T0ClassServiceVRFStandard
}

// IsVRFPremium - Returns true if the ClassService is VRFPremium.
func (c ClassService) IsVRFPremium() bool {
	return c == T0ClassServiceVRFPremium
}

// IsVRFDedicatedMedium - Returns true if the ClassService is VRFDedicatedMedium.
func (c ClassService) IsVRFDedicatedMedium() bool {
	return c == T0ClassServiceVRFDedicatedMedium
}

// IsVRFDedicatedLarge - Returns true if the ClassService is VRFDedicatedLarge.
func (c ClassService) IsVRFDedicatedLarge() bool {
	return c == T0ClassServiceVRFDedicatedLarge
}

// * List

// GetT0s - Returns the list of T0s.
func (t *Tier0) GetT0s() (listOfT0s *T0s, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&[]string{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/tier-0-vrfs")
	if err != nil {
		return
	}

	if r.IsError() {
		return listOfT0s, fmt.Errorf("error on list T0s: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	listOfT0s = &T0s{}

	for _, t0 := range *r.Result().(*[]string) {
		response, err := t.GetT0(t0)
		if err != nil {
			return listOfT0s, err
		}

		*listOfT0s = append(*listOfT0s, *response)
	}

	return listOfT0s, nil
}

// GetT0 - Returns the T0.
func (t *Tier0) GetT0(t0 string) (response *T0, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&T0{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("t0Name", t0).
		Get("/api/customers/v2.0/tier-0-vrfs/{t0Name}")
	if err != nil {
		return
	}

	if r.IsError() {
		return response, fmt.Errorf("error on get T0: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*T0), nil
}

// GetBandwidthCapacity - Returns the Bandwidth Capacity of the T0 in Mbps.
func (t *T0) GetBandwidthCapacity() (bandwidthCapacity int, err error) {
	if v, ok := T0ClassesServices[t.GetClassService()]; ok {
		return v.TotalBandwidth, nil
	}

	return 0, fmt.Errorf("unknown class service: %s", t.GetClassService())
}
