package netbackup

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	commonnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
)

// VDC - Is the response structure for the GetVdc API
type VDC struct {
	ID               int    `json:"Id,omitempty"`
	Name             string `json:"Name,omitempty"`
	Identifier       string `json:"Identifier,omitempty"`
	VOrgID           int    `json:"VOrgId,omitempty"`
	ProtectionTypeID int    `json:"ProtectionTypeId,omitempty"`
}

// GetID returns the ID field of VDC
func (vdc *VDC) GetID() int {
	return vdc.ID
}

// GetIDPtr returns a pointer to the ID field of VDC
func (vdc *VDC) GetIDPtr() *int {
	return &vdc.ID
}

// GetName returns the Name field of VDC
func (vdc *VDC) GetName() string {
	return vdc.Name
}

// GetIdentifier returns the Identifier field of VDC
func (vdc *VDC) GetIdentifier() string {
	return vdc.Identifier
}

// GetVOrgID returns the VOrgID field of VDC
func (vdc *VDC) GetVOrgID() int {
	return vdc.VOrgID
}

// GetProtectionTypeID returns the ProtectionTypeID field of VDC
func (vdc *VDC) GetProtectionTypeID() int {
	return vdc.ProtectionTypeID
}

// * VDCs
// VDCs - Is the response structure for the GetVdcs API
type VDCs []VDC

// append - Append a VDC to the VDCs slice
func (v *VDCs) append(vdc VDC) {
	*v = append(*v, vdc)
}

type vdcsResponse struct {
	Data VDCs `json:"data,omitempty"`
}

// GetVdcs - Get a list of vCloud Director Virtual Data Centers
func (v *VcloudClient) GetVdcs() (resp *VDCs, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&vdcsResponse{}).
		SetError(&commonnetbackup.APIError{}).
		Get("/v6/vcloud/vdcs")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*vdcsResponse).Data, nil
}

// GetVdcsByOrgID - Get a list of vCloud Director Virtual Data Centers by Org ID
// orgID - The ID of the org in the netbackup system
func (v *VcloudClient) GetVdcsByOrgID(orgID int) (resp *VDCs, err error) {
	vdcs, err := v.GetVdcs()
	if err != nil {
		return
	}

	for _, vdc := range *vdcs {
		if vdc.VOrgID == orgID {
			resp.append(vdc)
		}
	}

	return resp, nil
}

type vdcResponse struct {
	Data VDC `json:"data,omitempty"`
}

// GetVDCByID - Get a vCloud Director Virtual Data Center by ID
// id - The ID of the vdc in the netbackup system
func (v *VcloudClient) GetVDCByID(id int) (resp *VDC, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&vdcResponse{}).
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vdcID": fmt.Sprintf("%d", id),
		}).
		Get("/v6/vcloud/vdcs/{vdcID}")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*vdcResponse).Data, nil
}

// GetVDCByIdentifier - Get a vCloud Director Virtual Data Center by Identifier
// identifier - The Identifier of the vdc in the vmware system (URN)
func (v *VcloudClient) GetVDCByIdentifier(identifier string) (resp *VDC, err error) {
	vdcs, err := v.GetVdcs()
	if err != nil {
		return
	}

	for _, vdc := range *vdcs {
		if vdc.Identifier == identifier {
			return &vdc, nil
		}
	}

	return resp, fmt.Errorf("VDC with identifier %s not found", identifier)
}

// GetVDCByName - Get a vCloud Director Virtual Data Center by Name
// name - The Name of the vdc in the vmware system
func (v *VcloudClient) GetVDCByName(name string) (resp *VDC, err error) {
	vdcs, err := v.GetVdcs()
	if err != nil {
		return
	}

	for _, vdc := range *vdcs {
		if vdc.Name == name {
			return &vdc, nil
		}
	}

	return resp, fmt.Errorf("VDC with name %s not found", name)
}

// GetVDCByNameOrIdentifier - Get a vCloud Director Virtual Data Center by Name or Identifier
// nameOrIdentifier - The Name or Identifier of the vdc in the vmware system
func (v *VcloudClient) GetVDCByNameOrIdentifier(nameOrIdentifier string) (resp *VDC, err error) {
	vdcs, err := v.GetVdcs()
	if err != nil {
		return
	}

	for _, vdc := range *vdcs {
		if vdc.Name == nameOrIdentifier || vdc.Identifier == nameOrIdentifier {
			return &vdc, nil
		}
	}

	return resp, fmt.Errorf("VDC with name or identifier %s not found", nameOrIdentifier)
}

// ListProtectionLevelsAvailable - List the protection levels available for a vCloud Director Virtual Application
func (vdc *VDC) ListProtectionLevelsAvailable() (resp *ProtectionLevels, err error) {
	pL := ProtectionLevelClient{}
	return pL.ListProtectionLevels(listProtectionLevelsRequest{
		VDCID: vdc.GetIDPtr(),
	})
}

// GetProtectionLevelAvailableByName - Get a protection level by name for a vCloud Director Virtual Application
func (vdc *VDC) GetProtectionLevelAvailableByName(name string) (resp *ProtectionLevel, err error) {
	pL := ProtectionLevelClient{}
	return pL.getProtectionLevelByName(getProtectionLevelByNameRequest{
		VDCID:               vdc.GetIDPtr(),
		ProtectionLevelName: &name,
	})
}

// GetProtectionLevelAvailableByID - Get a protection level by ID for a vCloud Director Virtual Application
func (vdc *VDC) GetProtectionLevelAvailableByID(id int) (resp *ProtectionLevel, err error) {
	pL := ProtectionLevelClient{}
	return pL.getProtectionLevelByID(getProtectionLevelByIDRequest{
		VDCID:             vdc.GetIDPtr(),
		ProtectionLevelID: &id,
	})
}

// ListProtectionLevels - List the protection levels applied to a vCloud Director Virtual Application
func (vdc *VDC) ListProtectionLevels() (resp *ProtectionLevels, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return resp, err
	}

	r, err := c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vdcID": fmt.Sprintf("%d", vdc.GetID()),
		}).
		SetResult(&protectionLevelAppliedResponse{}).
		Get("/v6/vcloud/vdcs/{vdcID}/protected")
	if err != nil {
		return resp, err
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	resp = &ProtectionLevels{}
	for _, pL := range r.Result().(*protectionLevelAppliedResponse).Data.ProtectedLevels {
		if strings.ToLower(pL.EntityType) == "vdc" {
			resp.append(pL.ProtectionLevel)
		}
	}

	return resp, nil
}

// * ProtectVdc

// ProtectVdc - Protect a vCloud Director Virtual Data Center
func (vdc *VDC) Protect(req ProtectUnprotectRequest) (job *commonnetbackup.JobAPIResponse, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return job, err
	}

	if req.ProtectionLevelID == nil && req.ProtectionLevelName == "" {
		return job, fmt.Errorf("you must specify a ProtectionLevelID or ProtectionLevelName")
	}

	var (
		r               *resty.Response
		protectionLevel *ProtectionLevel
	)

	if req.ProtectionLevelID == nil {
		protectionLevel, err = vdc.GetProtectionLevelAvailableByName(req.ProtectionLevelName)
		if err != nil {
			return job, err
		}
	} else {
		protectionLevel, err = vdc.GetProtectionLevelAvailableByID(*req.ProtectionLevelID)
		if err != nil {
			return job, err
		}
	}

	r, err = c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vdcID": fmt.Sprintf("%d", vdc.GetID()),
		}).
		SetBody(protectBody{
			ProtectionLevelID: protectionLevel.GetID(),
			Paths:             []string{},
		}).
		SetResult(&commonnetbackup.JobAPIResponse{}).
		Post("/v6/vcloud/vdcs/{vdcID}/protect")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*commonnetbackup.JobAPIResponse), nil
}

// * UnprotectVdc

// UnprotectVdc - Unprotect a vCloud Director Virtual Data Center
func (vdc *VDC) Unprotect(req ProtectUnprotectRequest) (job *commonnetbackup.JobAPIResponse, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return job, err
	}

	if req.ProtectionLevelID == nil && req.ProtectionLevelName == "" {
		return job, fmt.Errorf("you must specify a ProtectionLevelID or ProtectionLevelName")
	}

	var (
		r               *resty.Response
		protectionLevel *ProtectionLevel
	)

	if req.ProtectionLevelID == nil {
		protectionLevel, err = vdc.GetProtectionLevelAvailableByName(req.ProtectionLevelName)
		if err != nil {
			return job, err
		}
	} else {
		protectionLevel, err = vdc.GetProtectionLevelAvailableByID(*req.ProtectionLevelID)
		if err != nil {
			return job, err
		}
	}

	r, err = c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vdcID": fmt.Sprintf("%d", vdc.GetID()),
		}).
		SetBody(protectBody{
			ProtectionLevelID: protectionLevel.GetID(),
		}).
		SetResult(&commonnetbackup.JobAPIResponse{}).
		Post("/v6/vcloud/vdcs/{vdcID}/unprotect")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*commonnetbackup.JobAPIResponse), nil
}
