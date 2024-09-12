package v1

import (
	"encoding/json"
	"fmt"
	"io"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type (
	// BMS struct{}

	BMS struct {
		BMSNetworks []BMSNetwork `json:"BMSnetworks"`
		BMSDetails  []BMSDetail  `json:"BMSDetails"`
	}

	BMSNetwork struct {
		VLANID string `json:"vlanId"`
		Subnet string `json:"subnet"`
		Prefix string `json:"prefix"`
	}

	BMSDetail struct {
		BMSType           string     `json:"bmsType"`
		Hostname          string     `json:"hostname"`
		OS                string     `json:"os"`
		BiosConfiguration string     `json:"biosConfiguration"`
		Storages          BMSStorage `json:"storages"`
	}

	BMSStorage struct {
		Local  BMSStorageDetail `json:"local"`
		Shared BMSStorageDetail `json:"shared"`
		System BMSStorageDetail `json:"system"`
		Data   BMSStorageDetail `json:"data"`
	}

	BMSStorageDetail struct {
		StorageClass string `json:"storageClass"`
		Size         string `json:"size"`
	}
)

// ! BMS
// function List BMSolution
func (v *BMS) List() (response *BMS, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	r, err := c.R().
		SetDoNotParseResponse(true).
		// SetResult(&BMS{}). - because the response is not a correct json format
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/bms")
	if err != nil {
		return nil, err
	}
	rep, err := io.ReadAll(r.RawBody())
	if err != nil {
		return nil, err
	}
	// ! convert []byte to string to fix the json format response
	stringRep := string(rep)
	// delete last character "]"
	stringRep = stringRep[:len(stringRep)-1]
	// delete first character "["
	stringRep = stringRep[1:]

	// put the json format y into struct *BMS
	err = json.Unmarshal([]byte(stringRep), &response)
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return response, fmt.Errorf("error on list BMS(s): %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return response, nil
	// return r.Result().(*BMS), nil
}

func (v *BMS) GetNetworks() []BMSNetwork {
	return v.BMSNetworks
}

func (v *BMS) GetBMSDetails() []BMSDetail {
	return v.BMSDetails
}

func (v *BMSDetail) GetBMSStorage() BMSStorage {
	return v.Storages
}

func (v *BMSStorage) GetLocal() BMSStorageDetail {
	return v.Local
}

func (v *BMSStorage) GetShared() BMSStorageDetail {
	return v.Shared
}

func (v *BMSStorage) GetSystem() BMSStorageDetail {
	return v.System
}

func (v *BMSStorage) GetData() BMSStorageDetail {
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
