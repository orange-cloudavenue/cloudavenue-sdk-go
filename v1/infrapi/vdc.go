package infrapi

import (
	"fmt"
	"regexp"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	rules "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi/rules"
)

type (
	CAVVDC               struct{}
	VDCs                 []CAVVirtualDataCenter
	CAVVirtualDataCenter struct {
		VDCGroup string                  `json:"vdcGroup,omitempty"`
		VDC      CAVVirtualDataCenterVDC `json:"vdc"`
	}
	CAVVirtualDataCenterVDC struct {
		Name                string             `json:"name"`
		Description         string             `json:"description"`
		ServiceClass        ServiceClass       `json:"vdcServiceClass"`
		DisponibilityClass  DisponibilityClass `json:"vdcDisponibilityClass"`
		BillingModel        BillingModel       `json:"vdcBillingModel"`
		VCPUInMhz           int                `json:"vcpuInMhz2"`
		CPUAllocated        int                `json:"cpuAllocated"`
		MemoryAllocated     int                `json:"memoryAllocated"`
		StorageBillingModel BillingModel       `json:"vdcStorageBillingModel"`
		StorageProfiles     []StorageProfile   `json:"vdcStorageProfiles,omitempty"`
	}
	StorageProfile struct {
		Class   StorageProfileClass `json:"class"`
		Limit   int                 `json:"limit"`
		Default bool                `json:"default"`
	}

	BillingModel        = rules.BillingModel
	ServiceClass        = rules.ServiceClass
	DisponibilityClass  = rules.DisponibilityClass
	StorageBillingModel = rules.BillingModel
	StorageProfileClass = rules.StorageProfileClass
)

// GetName - Return the VDC name
func (v *CAVVirtualDataCenter) GetName() string {
	return v.VDC.Name
}

// GetDescription - Return the VDC description
func (v *CAVVirtualDataCenter) GetDescription() string {
	return v.VDC.Description
}

// GetServiceClass - Return the VDC service class
func (v *CAVVirtualDataCenter) GetServiceClass() ServiceClass {
	return v.VDC.ServiceClass
}

// GetDisponibilityClass - Return the VDC disponibility class
func (v *CAVVirtualDataCenter) GetDisponibilityClass() DisponibilityClass {
	return v.VDC.DisponibilityClass
}

// GetBillingModel - Return the VDC billing model
func (v *CAVVirtualDataCenter) GetBillingModel() BillingModel {
	return v.VDC.BillingModel
}

// GetVCPUInMhz - Return the VDC vcpu in mhz2
func (v *CAVVirtualDataCenter) GetVCPUInMhz() int {
	return v.VDC.VCPUInMhz
}

// GetCPUAllocated - Return the VDC cpu allocated
func (v *CAVVirtualDataCenter) GetCPUAllocated() int {
	return v.VDC.CPUAllocated
}

// GetMemoryAllocated - Return the VDC memory allocated
func (v *CAVVirtualDataCenter) GetMemoryAllocated() int {
	return v.VDC.MemoryAllocated
}

// GetStorageBillingModel - Return the VDC storage billing model
func (v *CAVVirtualDataCenter) GetStorageBillingModel() BillingModel {
	return v.VDC.StorageBillingModel
}

// GetStorageProfiles - Return the VDC storage profiles
func (v *CAVVirtualDataCenter) GetStorageProfiles() []StorageProfile {
	return v.VDC.StorageProfiles
}

// GetVdcGroup - Return the VDC vdc group
func (v *CAVVirtualDataCenter) GetVDCGroup() string {
	return v.VDCGroup
}

// SetName - Set the VDC name
// Name respects the following regex: ^[a-zA-Z0-9-_]{1,64}$
func (v *CAVVirtualDataCenter) SetName(name string) error {
	// check if the name respects the regex
	re := regexp.MustCompile(`^[a-zA-Z0-9-_]{1,64}$`)
	if !re.MatchString(name) {
		return fmt.Errorf("invalid name value: %s (Accepted values: ^[a-zA-Z0-9-_]{1,64}$)", name)
	}
	v.VDC.Name = name
	return nil
}

// SetDescription - Set the VDC description
func (v *CAVVirtualDataCenter) SetDescription(description string) {
	v.VDC.Description = description
}

// SetCPUAllocated - Set the VDC cpu allocated
func (v *CAVVirtualDataCenter) SetCPUAllocated(cpuAllocated int) {
	v.VDC.CPUAllocated = cpuAllocated
}

// SetMemoryAllocated - Set the VDC memory allocated
func (v *CAVVirtualDataCenter) SetMemoryAllocated(memoryAllocated int) {
	v.VDC.MemoryAllocated = memoryAllocated
}

// AddStorageProfile - Add a storage profile
func (v *CAVVirtualDataCenter) AddStorageProfile(storageProfile StorageProfile) {
	v.VDC.StorageProfiles = append(v.VDC.StorageProfiles, storageProfile)
}

// RemoveStorageProfile - Remove a storage profile
func (v *CAVVirtualDataCenter) RemoveStorageProfile(storageProfile StorageProfile) {
	for i, storageProfileToRemove := range v.VDC.StorageProfiles {
		if storageProfileToRemove.Class == storageProfile.Class {
			v.VDC.StorageProfiles = append(v.VDC.StorageProfiles[:i], v.VDC.StorageProfiles[i+1:]...)
		}
	}
}

// SetStorageProfiles - Set the VDC storage profiles
func (v *CAVVirtualDataCenter) SetStorageProfiles(storageProfiles []StorageProfile) {
	v.VDC.StorageProfiles = storageProfiles
}

// SetVCPUInMhz - Set the VDC vcpu in mhz
func (v *CAVVirtualDataCenter) SetVCPUInMhz(vcpuInMhz int) {
	v.VDC.VCPUInMhz = vcpuInMhz
}

// Set - Set the VDC
func (v *CAVVirtualDataCenter) Set(vdc *CAVVirtualDataCenter) {
	v.VDC = vdc.VDC
}

// IsValid - Check if everythings is valid
func (v *CAVVirtualDataCenter) IsValid(isUpdate bool) error {
	return rules.Validate(rules.ValidateData{
		ServiceClass:        v.VDC.ServiceClass,
		DisponibilityClass:  v.VDC.DisponibilityClass,
		BillingModel:        v.VDC.BillingModel,
		VCPUInMhz:           v.VDC.VCPUInMhz,
		CPUAllocated:        v.VDC.CPUAllocated,
		MemoryAllocated:     v.VDC.MemoryAllocated,
		StorageBillingModel: v.VDC.StorageBillingModel,
		StorageProfiles: func() map[rules.StorageProfileClass]struct {
			Limit   int
			Default bool
		} {
			storageProfiles := make(map[rules.StorageProfileClass]struct {
				Limit   int
				Default bool
			})
			for _, storageProfile := range v.VDC.StorageProfiles {
				storageProfiles[storageProfile.Class] = struct {
					Limit   int
					Default bool
				}{Limit: storageProfile.Limit, Default: storageProfile.Default}
			}
			return storageProfiles
		}(),
	}, isUpdate)
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
		VDCName string `json:"vdc_name"`
		VDCUUID string `json:"vdc_uuid"`
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
		response, err := v.Get(vdc.VDCName)
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
		SetPathParam("vdcName", v.VDC.Name).
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
		SetPathParam("vdcName", v.VDC.Name).
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
	if err := value.IsValid(false); err != nil {
		return nil, fmt.Errorf("error on create VDC: %w", err)
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

	return v.Get(value.VDC.Name)
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
