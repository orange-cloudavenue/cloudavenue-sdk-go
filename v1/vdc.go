// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0

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
	"context"
	"errors"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi"
)

type (
	CAVVdc struct{}
)

// ! Errors.
var (
	ErrEmptyVDCNameProvided    = errors.New("empty VDC name provided")
	ErrEmptyVDCIDProvided      = errors.New("empty VDC ID provided")
	ErrRetrievingOrg           = errors.New("error retrieving org")
	ErrRetrievingOrgAdmin      = errors.New("error retrieving org admin")
	ErrRetrievingVDC           = errors.New("error retrieving VDC")
	ErrRetrievingVDCGroup      = errors.New("error retrieving VDC Group")
	ErrRetrievingVDCOrVDCGroup = errors.New("error retrieving VDC or VDC Group")
)

// GetVDC retrieves the VDC (Virtual Data Center) by its name.
// It returns a pointer to the VDC and an error if any.
// The function performs sequential lookups from three sources: the infrapi Get lookup,
// the VMware GetVDCByNameOrId lookup, and an infrapi List() name-scan.
// A successful infrapi Get is enough to return the VDC, and also triggers the VMware
// lookup to populate the VMware side of the object. The function returns an error only
// when all three lookups fail, with the error chosen by priority Get, GetVmware, List.
func (v *CAVVdc) GetVDC(vdcName string) (*VDC, error) {
	if vdcName == "" {
		return nil, ErrEmptyVDCNameProvided
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	getVDC := new(VDC)

	// First lookup: infrapi Get(name).
	infraPIVDC := infrapi.CAVVDC{}
	vdc, errGet := infraPIVDC.Get(vdcName)
	if errGet == nil && vdc != nil {
		getVDC.infrapi = vdc

		// Also resolve the VMware side so VDC methods backed by the VMware API
		// (e.g. CreateVAPP, GetSecurityGroupByID) keep working. This is
		// best-effort: a failure here only loses the VMware side.
		vdcVmware, errGetVmware := c.Org.GetVDCByNameOrId(vdcName, true)
		if errGetVmware == nil && vdcVmware != nil {
			getVDC.Vdc = vdcVmware
		}

		return getVDC, nil
	}

	// Second lookup: vmware GetVDCByNameOrId.
	vdcVmware, errGetVmware := c.Org.GetVDCByNameOrId(vdcName, true)
	if errGetVmware == nil && vdcVmware != nil {
		getVDC.Vdc = vdcVmware
		return getVDC, nil
	}

	// Third lookup: infrapi List() name-scan.
	vdcs, errList := infraPIVDC.List()
	if errList == nil && vdcs != nil {
		for _, vdc := range *vdcs {
			if vdc.VDC.Name == vdcName {
				getVDC.infrapi = &vdc
				return getVDC, nil
			}
		}
	}

	// All three lookups failed: return the single error by priority Get, GetVmware, List.
	err = errGet
	if err == nil {
		err = errGetVmware
	}
	if err == nil {
		err = errList
	}
	if err == nil {
		// Defensive guard: all lookups failed without returning an error.
		return nil, fmt.Errorf("%w: %s: no VDC data retrieved", ErrRetrievingVDC, vdcName)
	}

	return nil, fmt.Errorf("%w: %w", ErrRetrievingVDC, err)
}

func (v *VDC) Vmware() *govcd.Vdc {
	return v.Vdc
}

// New creates a new VDC.
// For the context use context.WithTimeout to set a timeout.
func (v *CAVVdc) New(ctx context.Context, object *infrapi.CAVVirtualDataCenter) (*VDC, error) {
	if object == nil {
		return nil, fmt.Errorf("error on create VDC: object is nil")
	}

	infraPIVDC := infrapi.CAVVDC{}
	vdcCreated, err := infraPIVDC.New(ctx, object)
	if err != nil {
		return nil, fmt.Errorf("error on create VDC: %w", err)
	}

	return v.GetVDC(vdcCreated.GetName())
}

// List returns the list of VDCs.
// TODO - refacto to return a slice of VDC.
func (v *CAVVdc) List() (*infrapi.VDCs, error) {
	infraPIVDC := infrapi.CAVVDC{}
	return infraPIVDC.List()
}

// ? VMware

// GetName returns the name of the VDC.
func (v *VDC) GetName() string {
	if v.Vdc == nil || v.Vdc.Vdc == nil {
		return ""
	}

	return v.Vdc.Vdc.Name
}

// GetID returns the ID of the VDC.
func (v *VDC) GetID() string {
	if v.Vdc == nil || v.Vdc.Vdc == nil {
		return ""
	}

	return v.Vdc.Vdc.ID
}

// ? Infrapi

// GetDescription returns the description of the VDC.
func (v *VDC) GetDescription() string {
	return v.infrapi.GetDescription()
}

// GetServiceClass returns the service class of the VDC.
func (v *VDC) GetServiceClass() infrapi.ServiceClass {
	return v.infrapi.GetServiceClass()
}

// GetDisponibilityClass returns the disponibility class of the VDC.
func (v *VDC) GetDisponibilityClass() infrapi.DisponibilityClass {
	return v.infrapi.GetDisponibilityClass()
}

// GetBillingModel returns the billing model of the VDC.
func (v *VDC) GetBillingModel() infrapi.BillingModel {
	return v.infrapi.GetBillingModel()
}

// GetVCPUInMhz returns the VCPU in MHz of the VDC.
func (v *VDC) GetVCPUInMhz() int {
	return v.infrapi.GetVCPUInMhz()
}

// GetCPUAllocated returns the CPU allocated of the VDC.
func (v *VDC) GetCPUAllocated() int {
	return v.infrapi.GetCPUAllocated()
}

// GetMemoryAllocated returns the memory allocated of the VDC.
func (v *VDC) GetMemoryAllocated() int {
	return v.infrapi.GetMemoryAllocated()
}

// GetStorageBillingModel returns the storage billing model of the VDC.
func (v *VDC) GetStorageBillingModel() infrapi.BillingModel {
	return v.infrapi.GetStorageBillingModel()
}

// GetStorageProfiles returns the storage profiles of the VDC.
func (v *VDC) GetStorageProfiles() []infrapi.StorageProfile {
	return v.infrapi.GetStorageProfiles()
}

// SetName set the name of the VDC.
// Name respects the following regex: ^[a-zA-Z0-9-_]{1,64}$.
func (v *VDC) SetName(name string) error {
	return v.infrapi.SetName(name)
}

// SetDescription set the description of the VDC.
func (v *VDC) SetDescription(description string) {
	v.infrapi.SetDescription(description)
}

// SetCPUAllocated set the CPU allocated of the VDC.
func (v *VDC) SetCPUAllocated(cpuAllocated int) {
	v.infrapi.SetCPUAllocated(cpuAllocated)
}

// SetMemoryAllocated set the memory allocated of the VDC.
func (v *VDC) SetMemoryAllocated(memoryAllocated int) {
	v.infrapi.SetMemoryAllocated(memoryAllocated)
}

// AddStorageProfile add a storage profile to the VDC.
func (v *VDC) AddStorageProfile(storageProfile infrapi.StorageProfile) {
	v.infrapi.AddStorageProfile(storageProfile)
}

// RemoveStorageProfile remove a storage profile from the VDC.
func (v *VDC) RemoveStorageProfile(storageProfileName infrapi.StorageProfile) {
	v.infrapi.RemoveStorageProfile(storageProfileName)
}

// SetStorageProfiles set the storage profiles of the VDC.
func (v *VDC) SetStorageProfiles(storageProfiles []infrapi.StorageProfile) {
	v.infrapi.SetStorageProfiles(storageProfiles)
}

// SetVCPUInMhz set the VCPU in MHz of the VDC.
func (v *VDC) SetVCPUInMhz(vcpuInMhz int) {
	v.infrapi.SetVCPUInMhz(vcpuInMhz)
}

// Set set the VDC.
func (v *VDC) Set(vdc *infrapi.CAVVirtualDataCenter) {
	v.infrapi.Set(vdc)
}

// IsValid returns true if the VDC is valid.
func (v *VDC) IsValid(isUpdate bool) error {
	return v.infrapi.IsValid(isUpdate)
}

// Delete deletes the VDC.
func (v *VDC) Delete(ctx context.Context) (job *commoncloudavenue.JobCreatedAPIResponse, err error) {
	return v.infrapi.Delete(ctx)
}

// Update updates the VDC.
func (v *VDC) Update(ctx context.Context) (err error) {
	return v.infrapi.Update(ctx)
}

// IsVDCGroup return true if the object is a VDC Group.
func (v VDC) IsVDCGroup() bool {
	return govcd.OwnerIsVdcGroup(v.GetID())
}

// GetSecurityGroupByID return the NSX-T security group using the ID provided in the argument.
func (v VDC) GetSecurityGroupByID(nsxtFirewallGroupID string) (*govcd.NsxtFirewallGroup, error) {
	return v.Vdc.GetNsxtFirewallGroupById(nsxtFirewallGroupID)
}

// GetSecurityGroupByName return the NSX-T security group using the name provided in the argument.
func (v VDC) GetSecurityGroupByName(nsxtFirewallGroupName string) (*govcd.NsxtFirewallGroup, error) {
	return v.Vdc.GetNsxtFirewallGroupByName(nsxtFirewallGroupName, govcdtypes.FirewallGroupTypeSecurityGroup)
}

// GetSecurityGroupByNameOrID return the NSX-T security group using the name or ID provided in the argument.
func (v VDC) GetSecurityGroupByNameOrID(nsxtFirewallGroupNameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if urn.IsValid(nsxtFirewallGroupNameOrID) {
		return v.GetSecurityGroupByID(nsxtFirewallGroupNameOrID)
	}

	return v.GetSecurityGroupByName(nsxtFirewallGroupNameOrID)
}

// GetIPSetByID return the NSX-T firewall group using the ID provided in the argument.
func (v VDC) GetIPSetByID(id string) (*govcd.NsxtFirewallGroup, error) {
	return v.Vdc.GetNsxtFirewallGroupById(id)
}

// GetIPSetByName return the NSX-T firewall group using the name provided in the argument.
func (v VDC) GetIPSetByName(name string) (*govcd.NsxtFirewallGroup, error) {
	return v.Vdc.GetNsxtFirewallGroupByName(name, govcdtypes.FirewallGroupTypeIpSet)
}

// GetIPSetByNameOrId return the NSX-T firewall group using the name or ID provided in the argument.
func (v VDC) GetIPSetByNameOrID(nameOrID string) (*govcd.NsxtFirewallGroup, error) {
	if urn.IsValid(nameOrID) {
		return v.GetIPSetByID(nameOrID)
	}

	return v.GetIPSetByName(nameOrID)
}

// SetIPSet set the NSX-T firewall group using the name provided in the argument.
func (v VDC) SetIPSet(ipSetConfig *govcdtypes.NsxtFirewallGroup) (*govcd.NsxtFirewallGroup, error) {
	return v.Vdc.CreateNsxtFirewallGroup(ipSetConfig)
}

// GetDefaultPlacementPolicyID give you the ID of the default placement policy.
func (v VDC) GetDefaultPlacementPolicyID() string {
	if v.Vdc == nil || v.Vdc.Vdc == nil || v.Vdc.Vdc.DefaultComputePolicy == nil {
		return ""
	}

	return v.Vdc.Vdc.DefaultComputePolicy.ID
}

// GetVAPP give you the vApp using the name provided in the argument.
func (v VDC) GetVAPP(nameOrID string, refresh bool) (*VAPP, error) {
	vapp, err := v.Vdc.GetVAppByNameOrId(nameOrID, refresh)
	if err != nil {
		return nil, err
	}

	return &VAPP{vapp}, nil
}

// CreateVAPP create new vApp.
func (v VDC) CreateVAPP(name, description string) (*VAPP, error) {
	vapp, err := v.Vdc.CreateRawVApp(name, description)
	if err != nil {
		return nil, err
	}

	return &VAPP{vapp}, nil
}

// getVDCNetworkById returns the VDC Network by its ID.
func (v VDC) getVDCNetworkByID(id string) (*govcd.OpenApiOrgVdcNetwork, error) {
	return v.GetOpenApiOrgVdcNetworkById(id)
}

// getVDCNetworkByName returns the VDC Network by its name.
func (v VDC) getVDCNetworkByName(name string) (*govcd.OpenApiOrgVdcNetwork, error) {
	return v.GetOpenApiOrgVdcNetworkByName(name)
}

// createVDCNetwork creates a VDC Network.
func (v VDC) createVDCNetwork(networkConfig *govcdtypes.OpenApiOrgVdcNetwork) (*govcd.OpenApiOrgVdcNetwork, error) {
	return v.CreateOpenApiOrgVdcNetwork(networkConfig)
}
