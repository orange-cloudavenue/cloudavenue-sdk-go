package v1

import (
	"context"
	"fmt"
	"sync"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi"
)

type (
	CAVVdc struct{}
)

// ! Errors
var (
	ErrEmptyVDCNameProvided    = fmt.Errorf("empty VDC name provided")
	ErrRetrievingOrg           = fmt.Errorf("error retrieving org")
	ErrRetrievingOrgAdmin      = fmt.Errorf("error retrieving org admin")
	ErrRetrievingVDC           = fmt.Errorf("error retrieving VDC")
	ErrRetrievingVDCOrVDCGroup = fmt.Errorf("error retrieving VDC or VDC Group")
)

// *
// * VDC
// *

type (
	VDC struct {
		vmware  *govcd.Vdc
		infrapi *infrapi.CAVVirtualDataCenter
		// netbackup
	}
)

// Get retrieves the VDC (Virtual Data Center) by its name.
// It returns a pointer to the VDC and an error if any.
// The function performs concurrent requests to retrieve the VDC from both the VMware and the infrapi.
// It uses goroutines and channels to handle the concurrent requests and waits for all goroutines to finish using a WaitGroup.
// The function returns the VDC that was successfully retrieved from either the VMware or the infrapi.
func (v *CAVVdc) GetVDC(vdcName string) (*VDC, error) {
	if vdcName == "" {
		return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	org, err := getOrg()
	if err != nil {
		return nil, err
	}

	// wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// channels
	var (
		errChan = make(chan error, 2)
		vdcChan = make(chan interface{}, 2)
		done    = make(chan bool)
	)

	defer close(errChan)
	defer close(vdcChan)

	getVDC := new(VDC)

	// goroutine to get the VDC from the vmware
	wg.Add(1)
	go func() {
		defer wg.Done()
		vdc, err := org.GetVDCByNameOrId(vdcName, false)
		if err != nil {
			errChan <- err
			return
		}
		vdcChan <- vdc
	}()

	// goroutine to get the VDC from the infrapi
	wg.Add(1)
	go func() {
		defer wg.Done()
		infraPIVDC := infrapi.CAVVDC{}
		vdc, err := infraPIVDC.Get(vdcName)
		if err != nil {
			errChan <- err
			return
		}
		vdcChan <- vdc
	}()

	go func() {
		wg.Wait()
		done <- true
	}()

	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			return nil, err
		case vdc := <-vdcChan:
			switch x := vdc.(type) {
			case *govcd.Vdc:
				getVDC.vmware = x
			case *infrapi.CAVVirtualDataCenter:
				getVDC.infrapi = x
			default:
				return nil, fmt.Errorf("unknown type %T", x)
			}
		case <-done:
			break
		}
	}

	return getVDC, nil
}

func (v *VDC) Vmware() *govcd.Vdc {
	return v.vmware
}

// New creates a new VDC.
// For the context use contet.withTimeout to set a timeout.
func (v *CAVVdc) New(ctx context.Context, object *infrapi.CAVVirtualDataCenter) (*VDC, error) {
	infraPIVDC := infrapi.CAVVDC{}
	vdcCreated, err := infraPIVDC.New(ctx, object)
	if err != nil {
		return nil, err
	}

	return v.GetVDC(vdcCreated.GetName())
}

// List returns the list of VDCs.
// TODO - refacto to return a slice of VDC
func (v *CAVVdc) List() (*infrapi.VDCs, error) {
	infraPIVDC := infrapi.CAVVDC{}
	return infraPIVDC.List()
}

// ? VMware

// GetName returns the name of the VDC.
func (v *VDC) GetName() string {
	return v.vmware.Vdc.Name
}

// GetID returns the ID of the VDC.
func (v *VDC) GetID() string {
	return v.vmware.Vdc.ID
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
// Name respects the following regex: ^[a-zA-Z0-9-_]{1,64}$
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
func (v *VDC) Delete(ctx context.Context) (err error) {
	return v.infrapi.Delete(ctx)
}

// Update updates the VDC.
func (v *VDC) Update(ctx context.Context) (err error) {
	return v.infrapi.Update(ctx)
}

// *
// * VDCGroup
// *

type (
	VDCGroup struct {
		vmware *govcd.VdcGroup
	}
)

// GetVDCGroup retrieves the VDC Group by its name.
// It returns a pointer to the VDC Group and an error if any.
func (v *CAVVdc) GetVDCGroup(vdcGroupName string) (*VDCGroup, error) {
	if vdcGroupName == "" {
		return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	adminOrg, err := getAdminOrg()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRetrievingOrgAdmin, err)
	}

	x, err := adminOrg.GetVdcGroupByName(vdcGroupName)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDC, vdcGroupName, err)
	}

	return &VDCGroup{
		vmware: x,
	}, nil
}

// GetName returns the name of the VDC Group.
func (v *VDCGroup) GetName() string {
	return v.vmware.VdcGroup.Name
}

// GetID returns the ID of the VDC Group.
func (v *VDCGroup) GetID() string {
	return v.vmware.VdcGroup.Id
}

// GetDescription returns the description of the VDC Group.
func (v *VDCGroup) GetDescription() string {
	return v.vmware.VdcGroup.Description
}

// *
// * VDCOrVDCGroup
// *

type VDCOrVDCGroup interface {
	// GetName returns the name of the VDC or VDC Group
	GetName() string
	// GetID returns the ID of the VDC or VDC Group
	GetID() string
	// GetDescription returns the description of the VDC or VDC Group
	GetDescription() string
}

// GetVDCOrVDCGroup returns the VDC or VDC Group by its name.
// It returns a pointer to the VDC or VDC Group and an error if any.
func (v *CAVVdc) GetVDCOrVDCGroup(vdcOrVDCGroupName string) (VDCOrVDCGroup, error) {
	xVDCGroup, err := v.GetVDCGroup(vdcOrVDCGroupName)
	if err == nil {
		return xVDCGroup, nil
	}

	xVDC, err := v.GetVDC(vdcOrVDCGroupName)
	if err == nil {
		return xVDC, nil
	}

	return nil, ErrRetrievingVDCOrVDCGroup
}
