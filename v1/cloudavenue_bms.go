package v1

import (
	"fmt"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type BMS struct{}

/*
	type BMSs struct {
		Networks []Network `json:"BMSnetworkDetails"`
		BMSs     []BMS     `json:"BMSDetails"`
	}

	type Network struct {
		VLANID string `json:"vlanId"`
		Subnet string `json:"subnet"`
		Prefix string `json:"prefix"`
	}

	type BMS struct {
		BMSType           string    `json:"bmsType"`
		Hostname          string    `json:"hostname"`
		OS                string    `json:"os"`
		BiosConfiguration string    `json:"biosConfiguration"`
		Storage           []Storage `json:"storage"`
	}
	type Storage struct {
		Local  []StorageDetails `json:"local"`
		Shared []StorageDetails `json:"shared"`
		System []StorageDetails `json:"system"`
		Data   []StorageDetails `json:"data"`
	}

	type StorageDetails struct {
		StorageClass string `json:"storageClass"`
		Size         string `json:"size"`
	}
*/

func (v *BMS) List() (BMS, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return BMS{}, err
	}

	r, err := c.R().
		SetResult(&BMS{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/bms")
	if err != nil {
		return BMS{}, err
	}

	if r.IsError() {
		return BMS{}, fmt.Errorf("error on list BMS(s): %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return *r.Result().(*BMS), nil
}
