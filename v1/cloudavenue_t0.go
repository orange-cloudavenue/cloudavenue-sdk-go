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
		VlanID  any    `json:"vlanId"`
	}

	ClassService string
)

const (
	// ClassServiceVRFStandard - VRF Standard
	ClassServiceVRFStandard ClassService = "VRF_STANDARD"
	// ClassServiceVRFPremium - VRF Premium
	ClassServiceVRFPremium ClassService = "VRF_PREMIUM"
)

// * T0

// GetTier0ClassService - Returns the Tier0ClassService
func (t *T0) GetTier0ClassService() string {
	return t.Tier0ClassService
}

// GetName - Returns the Tier0Vrf
func (t *T0) GetName() string {
	return t.Tier0Vrf
}

// GetTier0Vrf - Returns the Tier0Vrf
func (t *T0) GetTier0Vrf() string {
	return t.Tier0Vrf
}

// GetTier0Provider - Returns the Tier0Provider
func (t *T0) GetTier0Provider() string {
	return t.Tier0Provider
}

// GetClassService - Returns the ClassService
func (t *T0) GetClassService() ClassService {
	return t.ClassService
}

// GetServices - Returns the Services
func (t *T0) GetServices() T0Services {
	return t.Services
}

// * T0Service

// GetService - Returns the Service
func (t *T0Service) GetService() string {
	return t.Service
}

// GetVlanID - Returns the VlanID
func (t *T0Service) GetVlanID() any {
	return t.VlanID
}

// * ClassService

// IsVRFStandard - Returns true if the ClassService is VRFStandard
func (c ClassService) IsVRFStandard() bool {
	return c == ClassServiceVRFStandard
}

// IsVRFPremium - Returns true if the ClassService is VRFPremium
func (c ClassService) IsVRFPremium() bool {
	return c == ClassServiceVRFPremium
}

// * List

// GetT0s - Returns the list of T0s
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

// GetT0 - Returns the T0
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
