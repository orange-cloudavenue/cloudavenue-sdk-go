package infrapi

import (
	"encoding/json"
	"fmt"
	"regexp"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type (
	CAVVDC               struct{}
	VDCs                 []CAVVirtualDataCenter
	CAVVirtualDataCenter struct {
		VdcGroup string                  `json:"vdcGroup,omitempty"`
		Vdc      CAVVirtualDataCenterVDC `json:"vdc"`
	}
	CAVVirtualDataCenterVDC struct {
		Name                string                 `json:"name"`
		Description         string                 `json:"description"`
		ServiceClass        VDCServiceClass        `json:"vdcServiceClass"`
		DisponibilityClass  VDCDisponibilityClass  `json:"vdcDisponibilityClass"`
		BillingModel        VDCBillingModel        `json:"vdcBillingModel"`
		VcpuInMhz2          int                    `json:"vcpuInMhz2"`
		CPUAllocated        int                    `json:"cpuAllocated"`
		MemoryAllocated     int                    `json:"memoryAllocated"`
		StorageBillingModel VDCStorageBillingModel `json:"vdcStorageBillingModel"`
		StorageProfiles     []VDCStrorageProfile   `json:"vdcStorageProfiles,omitempty"`
	}
	VDCStrorageProfile struct {
		Class   VDCStrorageProfileClass `json:"class"`
		Limit   int                     `json:"limit"`
		Default bool                    `json:"default"`
	}

	// VDCBillingModel - VDC Billing Model
	VDCBillingModel string
	// VDCServiceClass - VDC Service Class
	VDCServiceClass string
	// VDCDisponibilityClass - VDC Disponibility Class
	VDCDisponibilityClass string
	// VDCStorageBillingModel - VDC Storage Billing Model
	VDCStorageBillingModel string
	// VDCStrorageProfileClass - VDC Storage Profile Class
	VDCStrorageProfileClass string
)

const (
	// VDCBillingModelPayAsYouGo - Pay as you go
	VDCBillingModelPayAsYouGo VDCBillingModel = "PAYG"
	// VDCBillingModelReserved - Reserved
	VDCBillingModelReserved VDCBillingModel = "RESERVED"
	// VDCBillingModelDraas - DRaaS
	VDCBillingModelDraas VDCBillingModel = "DRAAS"

	// VDCServiceClassStandard - Standard
	VDCServiceClassStandard VDCServiceClass = "STD"
	// VDCServiceClassEco - Eco
	VDCServiceClassEco VDCServiceClass = "ECO"
	// VDCServiceClassHP - HP
	VDCServiceClassHP VDCServiceClass = "HP"
	// VDCServiceClassVoIP - VoIP
	VDCServiceClassVoIP VDCServiceClass = "VOIP"

	// VDCDisponibilityClassOneRoom - One Room
	VDCDisponibilityClassOneRoom VDCDisponibilityClass = "ONE-ROOM"
	// VDCDisponibilityClassDualRoom - Dual Room
	VDCDisponibilityClassDualRoom VDCDisponibilityClass = "DUAL-ROOM"
	// VDCDisponibilityClassHADualRoom - HA Dual Room
	VDCDisponibilityClassHADualRoom VDCDisponibilityClass = "HA-DUAL-ROOM"

	// VDCStorageBillingModelPayAsYouGo - Pay as you go
	VDCStorageBillingModelPayAsYouGo VDCStorageBillingModel = "PAYG"
	// VDCStorageBillingModelReserved - Reserved
	VDCStorageBillingModelReserved VDCStorageBillingModel = "RESERVED"

	// VDCStrorageProfileClassSilver - Silver
	VDCStrorageProfileClassSilver VDCStrorageProfileClass = "silver"
	// VDCStrorageProfileClassSilverR1 - Silver R1
	VDCStrorageProfileClassSilverR1 VDCStrorageProfileClass = "silver_r1"
	// VDCStrorageProfileClassSilverR2 - Silver R2
	VDCStrorageProfileClassSilverR2 VDCStrorageProfileClass = "silver_r2"
	// VDCStrorageProfileClassGold - Gold
	VDCStrorageProfileClassGold VDCStrorageProfileClass = "gold"
	// VDCStrorageProfileClassGoldR1 - Gold R1
	VDCStrorageProfileClassGoldR1 VDCStrorageProfileClass = "gold_r1"
	// VDCStrorageProfileClassGoldR2 - Gold R2
	VDCStrorageProfileClassGoldR2 VDCStrorageProfileClass = "gold_r2"
	// VDCStrorageProfileClassGoldHm - Gold HM
	VDCStrorageProfileClassGoldHm VDCStrorageProfileClass = "gold_hm"
	// VDCStrorageProfileClassPlatinum3k - Platinum 3k
	VDCStrorageProfileClassPlatinum3k VDCStrorageProfileClass = "platinum3k"
	// VDCStrorageProfileClassPlatinum3kR1 - Platinum 3k R1
	VDCStrorageProfileClassPlatinum3kR1 VDCStrorageProfileClass = "platinum3k_r1"
	// VDCStrorageProfileClassPlatinum3kR2 - Platinum 3k R2
	VDCStrorageProfileClassPlatinum3kR2 VDCStrorageProfileClass = "platinum3k_r2"
	// VDCStrorageProfileClassPlatinum3kHm - Platinum 3k HM
	VDCStrorageProfileClassPlatinum3kHm VDCStrorageProfileClass = "platinum3k_hm"
	// VDCStrorageProfileClassPlatinum7k - Platinum 7k
	VDCStrorageProfileClassPlatinum7k VDCStrorageProfileClass = "platinum7k"
	// VDCStrorageProfileClassPlatinum7kR1 - Platinum 7k R1
	VDCStrorageProfileClassPlatinum7kR1 VDCStrorageProfileClass = "platinum7k_r1"
	// VDCStrorageProfileClassPlatinum7kR2 - Platinum 7k R2
	VDCStrorageProfileClassPlatinum7kR2 VDCStrorageProfileClass = "platinum7k_r2"
	// VDCStrorageProfileClassPlatinum7kHm - Platinum 7k HM
	VDCStrorageProfileClassPlatinum7kHm VDCStrorageProfileClass = "platinum7k_hm"
)

var (
	ErrVDCServiceClass         = fmt.Errorf("invalid service class value (Allowed values: %s, %s, %s, %s)", VDCServiceClassStandard, VDCServiceClassEco, VDCServiceClassHP, VDCServiceClassVoIP)
	ErrVDCBillingModel         = fmt.Errorf("invalid billing model value (Allowed values: %s, %s, %s)", VDCBillingModelPayAsYouGo, VDCBillingModelReserved, VDCBillingModelDraas)
	ErrVDCDisponibilityClass   = fmt.Errorf("invalid disponibility class value (Allowed values: %s, %s, %s)", VDCDisponibilityClassOneRoom, VDCDisponibilityClassDualRoom, VDCDisponibilityClassHADualRoom)
	ErrVDCStorageBillingModel  = fmt.Errorf("invalid storage billing model value (Allowed values: %s, %s)", VDCStorageBillingModelPayAsYouGo, VDCStorageBillingModelReserved)
	ErrVDCStrorageProfileClass = fmt.Errorf("invalid storage profile class value (Allowed values: %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)", VDCStrorageProfileClassSilver, VDCStrorageProfileClassSilverR1, VDCStrorageProfileClassSilverR2, VDCStrorageProfileClassGold, VDCStrorageProfileClassGoldR1, VDCStrorageProfileClassGoldR2, VDCStrorageProfileClassGoldHm, VDCStrorageProfileClassPlatinum3k, VDCStrorageProfileClassPlatinum3kR1, VDCStrorageProfileClassPlatinum3kR2, VDCStrorageProfileClassPlatinum3kHm, VDCStrorageProfileClassPlatinum7k, VDCStrorageProfileClassPlatinum7kR1, VDCStrorageProfileClassPlatinum7kR2, VDCStrorageProfileClassPlatinum7kHm)
)

// GetName - Return the VDC name
func (v *CAVVirtualDataCenter) GetName() string {
	return v.Vdc.Name
}

// GetDescription - Return the VDC description
func (v *CAVVirtualDataCenter) GetDescription() string {
	return v.Vdc.Description
}

// GetServiceClass - Return the VDC service class
func (v *CAVVirtualDataCenter) GetServiceClass() VDCServiceClass {
	return v.Vdc.ServiceClass
}

// GetDisponibilityClass - Return the VDC disponibility class
func (v *CAVVirtualDataCenter) GetDisponibilityClass() VDCDisponibilityClass {
	return v.Vdc.DisponibilityClass
}

// GetBillingModel - Return the VDC billing model
func (v *CAVVirtualDataCenter) GetBillingModel() VDCBillingModel {
	return v.Vdc.BillingModel
}

// GetVcpuInMhz2 - Return the VDC vcpu in mhz2
func (v *CAVVirtualDataCenter) GetVcpuInMhz2() int {
	return v.Vdc.VcpuInMhz2
}

// GetCPUAllocated - Return the VDC cpu allocated
func (v *CAVVirtualDataCenter) GetCPUAllocated() int {
	return v.Vdc.CPUAllocated
}

// GetMemoryAllocated - Return the VDC memory allocated
func (v *CAVVirtualDataCenter) GetMemoryAllocated() int {
	return v.Vdc.MemoryAllocated
}

// GetStorageBillingModel - Return the VDC storage billing model
func (v *CAVVirtualDataCenter) GetStorageBillingModel() VDCStorageBillingModel {
	return v.Vdc.StorageBillingModel
}

// GetStorageProfiles - Return the VDC storage profiles
func (v *CAVVirtualDataCenter) GetStorageProfiles() []VDCStrorageProfile {
	return v.Vdc.StorageProfiles
}

// GetVdcGroup - Return the VDC vdc group
func (v *CAVVirtualDataCenter) GetVdcGroup() string {
	return v.VdcGroup
}

// SetName - Set the VDC name
// Name respects the following regex: ^[a-zA-Z0-9-_]{1,64}$
func (v *CAVVirtualDataCenter) SetName(name string) error {
	// check if the name respects the regex
	re := regexp.MustCompile(`^[a-zA-Z0-9-_]{1,64}$`)
	if !re.MatchString(name) {
		return fmt.Errorf("invalid name value: %s (Accepted values: ^[a-zA-Z0-9-_]{1,64}$)", name)
	}
	v.Vdc.Name = name
	return nil
}

// SetDescription - Set the VDC description
func (v *CAVVirtualDataCenter) SetDescription(description string) {
	v.Vdc.Description = description
}

// SetCPUAllocated - Set the VDC cpu allocated
func (v *CAVVirtualDataCenter) SetCPUAllocated(cpuAllocated int) {
	v.Vdc.CPUAllocated = cpuAllocated
}

// SetMemoryAllocated - Set the VDC memory allocated
func (v *CAVVirtualDataCenter) SetMemoryAllocated(memoryAllocated int) {
	v.Vdc.MemoryAllocated = memoryAllocated
}

// AddStorageProfile - Add a storage profile
func (v *CAVVirtualDataCenter) AddStorageProfile(storageProfile VDCStrorageProfile) {
	v.Vdc.StorageProfiles = append(v.Vdc.StorageProfiles, storageProfile)
}

// RemoveStorageProfile - Remove a storage profile
func (v *CAVVirtualDataCenter) RemoveStorageProfile(storageProfile VDCStrorageProfile) {
	for i, storageProfileToRemove := range v.Vdc.StorageProfiles {
		if storageProfileToRemove.Class == storageProfile.Class {
			v.Vdc.StorageProfiles = append(v.Vdc.StorageProfiles[:i], v.Vdc.StorageProfiles[i+1:]...)
		}
	}
}

// SetStorageProfiles - Set the VDC storage profiles
func (v *CAVVirtualDataCenter) SetStorageProfiles(storageProfiles []VDCStrorageProfile) {
	v.Vdc.StorageProfiles = storageProfiles
}

// SetVCPUInMhz2 - Set the VDC vcpu in mhz2
func (v *CAVVirtualDataCenter) SetVCPUInMhz2(vcpuInMhz2 int) {
	v.Vdc.VcpuInMhz2 = vcpuInMhz2
}

// Set - Set the VDC
func (v *CAVVirtualDataCenter) Set(vdc *CAVVirtualDataCenter) {
	v.Vdc = vdc.Vdc
}

// MarshalJSON - Marshal the VDC
func (v *CAVVirtualDataCenter) MarshalJSON() ([]byte, error) {
	if ok, err := v.IsValid(); !ok {
		return nil, err
	}

	return json.Marshal(*v)
}

// Is

// IsValid - Check if the billing model is valid
func (v VDCBillingModel) IsValid() bool {
	switch v {
	case VDCBillingModelPayAsYouGo, VDCBillingModelReserved, VDCBillingModelDraas:
		return true
	}

	return false
}

// IsValid - Check if the service class is valid
func (v VDCServiceClass) IsValid() bool {
	switch v {
	case VDCServiceClassStandard, VDCServiceClassEco, VDCServiceClassHP, VDCServiceClassVoIP:
		return true
	}

	return false
}

// IsValid - Check if the storage billing model is valid
func (v VDCStorageBillingModel) IsValid() bool {
	switch v {
	case VDCStorageBillingModelPayAsYouGo, VDCStorageBillingModelReserved:
		return true
	}

	return false
}

// IsValid - Check if the storage profile class is valid
func (s VDCStrorageProfileClass) IsValid() bool {
	switch s {
	case VDCStrorageProfileClassSilver, VDCStrorageProfileClassSilverR1, VDCStrorageProfileClassSilverR2, VDCStrorageProfileClassGold, VDCStrorageProfileClassGoldR1, VDCStrorageProfileClassGoldR2, VDCStrorageProfileClassGoldHm, VDCStrorageProfileClassPlatinum3k, VDCStrorageProfileClassPlatinum3kR1, VDCStrorageProfileClassPlatinum3kR2, VDCStrorageProfileClassPlatinum3kHm, VDCStrorageProfileClassPlatinum7k, VDCStrorageProfileClassPlatinum7kR1, VDCStrorageProfileClassPlatinum7kR2, VDCStrorageProfileClassPlatinum7kHm:
		return true
	}

	return false
}

// IsValid - Check if the disponibility class is valid
func (v VDCDisponibilityClass) IsValid() bool {
	switch v {
	case VDCDisponibilityClassOneRoom, VDCDisponibilityClassDualRoom, VDCDisponibilityClassHADualRoom:
		return true
	}

	return false
}

// IsValid - Check if everythings is valid
func (v *CAVVirtualDataCenter) IsValid() (bool, error) {
	if !v.Vdc.StorageBillingModel.IsValid() {
		return false, ErrVDCStorageBillingModel
	}

	if !v.Vdc.ServiceClass.IsValid() {
		return false, ErrVDCServiceClass
	}

	if !v.Vdc.DisponibilityClass.IsValid() {
		return false, ErrVDCDisponibilityClass
	}

	if !v.Vdc.BillingModel.IsValid() {
		return false, ErrVDCBillingModel
	}

	for _, storageProfile := range v.Vdc.StorageProfiles {
		if !storageProfile.Class.IsValid() {
			return false, ErrVDCStrorageProfileClass
		}
	}

	return true, nil
}

// Get VDC - Return the VDC Object
func (v *CAVVDC) Get(vdcName string) (*CAVVirtualDataCenter, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	r, err := c.R().
		SetResult(&CAVVirtualDataCenter{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("vdcName", vdcName).
		Get("/api/customers/v2.0/vdcs/{vdcName}")
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on get VDC: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*CAVVirtualDataCenter), nil
}

// List - Return the list of VDCs
func (v *CAVVDC) List() (*VDCs, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	type listOfVDCs []struct {
		VdcName string `json:"vdc_name"`
		VdcUUID string `json:"vdc_uuid"`
	}

	r, err := c.R().
		SetResult(&listOfVDCs{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/vdcs")
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on list VDCs: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	vdcS := &VDCs{}

	// TODO : Use waitgroup to get all VDCs
	for _, vdc := range *r.Result().(*listOfVDCs) {
		response, err := v.Get(vdc.VdcName)
		if err != nil {
			return vdcS, err
		}

		*vdcS = append(*vdcS, *response)
	}

	return vdcS, nil
}

// Delete - Delete the VDC
func (v *CAVVirtualDataCenter) Delete() (job *commoncloudavenue.JobStatus, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("vdcName", v.Vdc.Name).
		Delete("/api/customers/v2.0/vdcs/{vdcName}")
	if err != nil {
		return
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on delete VDC: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus), nil
}

// Update - Update the VDC
func (v *CAVVirtualDataCenter) Update() (err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return err
	}

	r, err := c.R().
		SetBody(v).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("vdcName", v.Vdc.Name).
		SetResult(&commoncloudavenue.JobStatus{}).
		Put("/api/customers/v2.0/vdcs/{vdcName}")
	if err != nil {
		return err
	}

	if r.IsError() {
		return fmt.Errorf("error on update VDC: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus).Wait(1, 90)
}

// New - Create a new VDC
func (v *CAVVDC) New(value *CAVVirtualDataCenter) (vdc *CAVVirtualDataCenter, err error) {
	if ok, err := value.IsValid(); !ok {
		return nil, err
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetBody(value).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetResult(&commoncloudavenue.JobStatus{}).
		Post("/api/customers/v2.0/vdcs")
	if err != nil {
		return
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on create VDC: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	if err := r.Result().(*commoncloudavenue.JobStatus).Wait(1, 90); err != nil {
		return nil, err
	}

	return v.Get(value.Vdc.Name)
}

// GetVMwareObject - Return the VMware object
func (v *CAVVirtualDataCenter) GetVMwareObject() (*govcd.Vdc, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	org, err := c.Vmware.GetOrgByName(c.GetOrganization())
	if err != nil {
		return nil, err
	}

	return org.GetVDCByName(v.GetName(), true)
}
