package v1

import (
	"fmt"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type (
	BMS struct {
		BMSNetworks []BMSNetwork `json:"network"`
		BMSDetails  []BMSDetail  `json:"bms"`
	}

	BMSNetwork struct {
		VLANID string `json:"vlanid"`
		Subnet string `json:"subnet"`
		Prefix string `json:"prefix"`
	}

	BMSDetail struct {
		BMSType           string     `json:"bmsType"`
		Hostname          string     `json:"hostname"`
		OS                string     `json:"os"`
		BiosConfiguration string     `json:"biosConfiguration"`
		Storages          BMSStorage `json:"storage"`
	}

	BMSStorage struct {
		Local  []BMSStorageDetail `json:"local,omitempty"`
		Shared []BMSStorageDetail `json:"shared,omitempty"`
		System []BMSStorageDetail `json:"system,omitempty"`
		Data   []BMSStorageDetail `json:"data,omitempty"`
	}

	BMSStorageDetail struct {
		StorageClass string `json:"storageClass"`
		Size         string `json:"size"`
	}
)

// ! BMS
// Return a Slice of BMS struct.
func (v *BMS) List() (response *[]BMS, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult([]BMS{}). //- because the response is a slice of struct
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/bms")
	if err != nil {
		return
	}

	if r.IsError() {
		return response, fmt.Errorf("error on list BMS(s): %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*[]BMS), nil
}

func (v *BMS) GetNetworks() []BMSNetwork {
	return v.BMSNetworks
}

func (v *BMS) GetBMS() []BMSDetail {
	return v.BMSDetails
}

func (v *BMS) GetBMSByHostname(hostname string) (response *BMSDetail, err error) {
	// For each BMS find the one with the same hostname
	for _, bms := range v.BMSDetails {
		if bms.Hostname == hostname {
			return &bms, nil
		}
	}
	return nil, fmt.Errorf("BMS with hostname %s not found", hostname)
}

func (v *BMSDetail) GetStorages() BMSStorage {
	return v.Storages
}

func (v *BMSStorage) GetLocal() []BMSStorageDetail {
	return v.Local
}

func (v *BMSStorage) GetShared() []BMSStorageDetail {
	return v.Shared
}

func (v *BMSStorage) GetSystem() []BMSStorageDetail {
	return v.System
}

func (v *BMSStorage) GetData() []BMSStorageDetail {
	return v.Data
}

func (v *BMSStorageDetail) GetStorageClass() string {
	return v.StorageClass
}

func (v *BMSStorageDetail) GetSize() string {
	return v.Size
}

func (v *BMSNetwork) GetVLANID() string {
	return v.VLANID
}

func (v *BMSNetwork) GetSubnet() string {
	return v.Subnet
}

func (v *BMSNetwork) GetPrefix() string {
	return v.Prefix
}
