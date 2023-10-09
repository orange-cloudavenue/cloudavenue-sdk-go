package v1

import (
	"fmt"

	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	commonnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
)

type Org struct {
	ID           int    `json:"Id"`
	Name         string `json:"Name"`
	DisplayName  string `json:"DisplayName"`
	CustomerCode string `json:"CustomerCode"`
	ImportSource string `json:"ImportSource"`
	Location     string `json:"Location"`
}

// GetID returns the ID field of Org
func (r *Org) GetID() int {
	return r.ID
}

// GetIDPtr returns a pointer to the ID field of Org
func (r *Org) GetIDPtr() *int {
	return &r.ID
}

// GetName returns the Name field of Org
func (r *Org) GetName() string {
	return r.Name
}

// GetDisplayName returns the DisplayName field of Org
func (r *Org) GetDisplayName() string {
	return r.DisplayName
}

// GetCustomerCode returns the CustomerCode field of Org
func (r *Org) GetCustomerCode() string {
	return r.CustomerCode
}

// GetImportSource returns the ImportSource field of Org
func (r *Org) GetImportSource() string {
	return r.ImportSource
}

// GetLocation returns the Location field of Org
func (r *Org) GetLocation() string {
	return r.Location
}

// * Orgs
// Orgs - Is the response structure for the GetOrgs API
type Orgs []Org

type orgsResponse struct {
	Data Orgs `json:"data"`
}

// GetOrgs - Get a list of vCloud Director Organizations
func (v *VCloudClient) GetOrgs() (resp *Orgs, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&orgsResponse{}).
		SetError(&commonnetbackup.APIError{}).
		Get("/v6/vcloud/orgs")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*orgsResponse).Data, nil
}

// * Org

type orgResponse struct {
	Data Org `json:"data"`
}

// GetOrg - Get a vCloud Director Organization
func (v *VCloudClient) GetOrg(id int) (resp *Org, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&orgResponse{}).
		SetError(commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vAppID": fmt.Sprintf("%d", id),
		}).
		Get("/v6/vcloud/orgs/{orgID}")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*orgResponse).Data, nil
}

// GetOrgByName - Get a vCloud Director Organization by name
func (v *VCloudClient) GetOrgByName(name string) (resp *Org, err error) {
	orgs, err := v.GetOrgs()
	if err != nil {
		return
	}

	for _, org := range *orgs {
		if org.Name == name {
			return &org, nil
		}
	}

	return resp, fmt.Errorf("org %s not found", name)
}
