package v1

import (
	"context"
	"fmt"
	"sync"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/uuid"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi"
)

// Contrain the VDC object.
var _ VDCOrVDCGroupInterface = (*VDC)(nil)

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

// Get retrieves the VDC (Virtual Data Center) by its name.
// It returns a pointer to the VDC and an error if any.
// The function performs concurrent requests to retrieve the VDC from both the VMware and the infrapi.
// It uses goroutines and channels to handle the concurrent requests and waits for all goroutines to finish using a WaitGroup.
// The function returns the VDC that was successfully retrieved from either the VMware or the infrapi.
func (v *CAVVdc) GetVDC(vdcName string) (*VDC, error) {
	if vdcName == "" {
		return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	c, err := clientcloudavenue.New()
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
		vdc, err := c.Org.GetVDCByNameOrId(vdcName, true)
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
				getVDC.Vdc = x
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
	return v.Vdc
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
	return v.Vdc.Vdc.Name
}

// GetID returns the ID of the VDC.
func (v *VDC) GetID() string {
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
	if uuid.IsValid(nsxtFirewallGroupNameOrID) {
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
	if uuid.IsValid(nameOrID) {
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

// GetNetworkContextProfileByName retrieves a network context profile by name or ID
func (v VDC) GetNetworkContextProfileByNameOrID(name string, scope VDCOrVDCGroupNetworkContextProfileScope) (*VDCOrVDCGroupNetworkContextProfile, error) {
	return getNetworkContextProfile(name, v.GetID(), scope)
}

// ListNetworkContextProfilesAttributes retrieves all network context profiles attributes
func (v VDC) ListNetworkContextProfilesAttributes() interface{} {
	return listNetworkContextProfileAttributes()
}
