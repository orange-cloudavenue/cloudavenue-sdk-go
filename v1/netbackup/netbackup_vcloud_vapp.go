package netbackup

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	commonnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
)

// VApp - Is the response structure for the GetVApp API
type VApp struct {
	ID               int    `json:"Id,omitempty"`
	Identifier       string `json:"Identifier,omitempty"`
	ProtectionTypeID int    `json:"ProtectionTypeId,omitempty"`
	VDCID            int    `json:"VdcId,omitempty"`
	Name             string `json:"Name,omitempty"`
}

// GetID returns the ID field of VApp
func (vApp *VApp) GetID() int {
	return vApp.ID
}

// GetIDPtr returns a pointer to the ID field of VApp
func (vApp *VApp) GetIDPtr() *int {
	return &vApp.ID
}

// GetName returns the Name field of VApp
func (vApp *VApp) GetName() string {
	return vApp.Name
}

// GetIdentifier returns the Identifier field of VApp
func (vApp *VApp) GetIdentifier() string {
	return vApp.Identifier
}

// GetProtectionTypeID returns the ProtectionTypeID field of VApp
func (vApp *VApp) GetProtectionTypeID() int {
	return vApp.ProtectionTypeID
}

// * VApps
// VApps - Is the response structure for the GetVApps API
type VApps []VApp

// append - Append a VApp to the VApps slice
func (v *VApps) append(vapp VApp) {
	*v = append(*v, vapp)
}

type VAppsResponse struct {
	Data VApps `json:"data,omitempty"`
}

// GetVApps - Get a list of vCloud Director Virtual Applications
func (v *VcloudClient) GetVApps() (resp *VApps, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&VAppsResponse{}).
		SetError(&commonnetbackup.APIError{}).
		Get("/v6/vcloud/vapps")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*VAppsResponse).Data, nil
}

// * VApp

type VAppResponse struct {
	Data VApp `json:"data,omitempty"`
}

// GetVAppByID - Get a vCloud Director Virtual Application by ID
// id - The ID of the vapp in the netbackup system
func (v *VcloudClient) GetVAppByID(id int) (resp *VApp, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&VAppResponse{}).
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vAppID": fmt.Sprintf("%d", id),
		}).
		Get("/v6/vcloud/vapps/{vAppID}")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*VAppResponse).Data, nil
}

// GetVAppByName - Get a vCloud Director Virtual Application by Name
// name - The name of the vapp in the netbackup system
func (v *VcloudClient) GetVAppByName(name string) (resp *VApp, err error) {
	vapps, err := v.GetVApps()
	if err != nil {
		return
	}

	for _, vapp := range *vapps {
		if vapp.Name == name {
			return &vapp, nil
		}
	}

	return resp, fmt.Errorf("VApp with name %s not found", name)
}

// GetVAppByIdentifier - Get a vCloud Director Virtual Application by Identifier
// identifier - The Identifier of the vapp in the vmware system (URN)
func (v *VcloudClient) GetVAppByIdentifier(identifier string) (resp *VApp, err error) {
	vapps, err := v.GetVApps()
	if err != nil {
		return
	}

	for _, vapp := range *vapps {
		if vapp.Identifier == identifier {
			return &vapp, nil
		}
	}

	return resp, fmt.Errorf("VApp with identifier %s not found", identifier)
}

// GetVdcByNameOrIdentifier - Get a vCloud Director Virtual Application by Name or Identifier
// nameOrIdentifier - The Name or Identifier of the vapp in the vmware system
func (v *VcloudClient) GetVAppByNameOrIdentifier(nameOrIdentifier string) (resp *VApp, err error) {
	vapps, err := v.GetVApps()
	if err != nil {
		return
	}

	for _, vapp := range *vapps {
		if vapp.Name == nameOrIdentifier || vapp.Identifier == nameOrIdentifier {
			return &vapp, nil
		}
	}

	return resp, fmt.Errorf("VApp with name or identifier %s not found", nameOrIdentifier)
}

// * VApp Machines

// GetVAppMachinesResponse - Is the response structure for the GetVAppMachines API
type GetVAppMachinesResponse struct {
	Data []struct {
		IsVisibleToAllUsers       bool      `json:"IsVisibleToAllUsers,omitempty"`
		NetBackupClientName       string    `json:"NetBackupClientName,omitempty"`
		LastUpdatedDateTime       time.Time `json:"LastUpdatedDateTime,omitempty"`
		MachineCode               string    `json:"MachineCode,omitempty"`
		DisplayName               string    `json:"DisplayName,omitempty"`
		ProtectionTypeName        string    `json:"ProtectionTypeName,omitempty"`
		Hardware                  string    `json:"Hardware,omitempty"`
		TrafficLightStatus        string    `json:"TrafficLightStatus,omitempty"`
		IsSyncingNetBackupData    bool      `json:"IsSyncingNetBackupData,omitempty"`
		PolicyAppendix            string    `json:"PolicyAppendix,omitempty"`
		AdditionalData            string    `json:"AdditionalData,omitempty"`
		Os                        string    `json:"OS,omitempty"`
		ImportSource              string    `json:"ImportSource,omitempty"`
		CustomerCode              string    `json:"CustomerCode,omitempty"`
		ProviderAssetType         string    `json:"ProviderAssetType,omitempty"`
		IsDeletedFromImportSource bool      `json:"IsDeletedFromImportSource,omitempty"`
		VAppID                    int       `json:"VAppId,omitempty"`
		LastSuccessfulBackupDate  time.Time `json:"LastSuccessfulBackupDate,omitempty"`
		SyncLastError             string    `json:"SyncLastError,omitempty"`
		ProtectionTypeID          int       `json:"ProtectionTypeId,omitempty"`
		IsInVCloud                bool      `json:"IsInVCloud,omitempty"`
		Links                     []struct {
			Rel    string `json:"Rel,omitempty"`
			Href   string `json:"Href,omitempty"`
			Method string `json:"Method,omitempty"`
		} `json:"Links,omitempty"`
		CatalogName     string    `json:"CatalogName,omitempty"`
		ID              int       `json:"Id,omitempty"`
		VMDisplayName   string    `json:"VMDisplayName,omitempty"`
		CreatedDateTime time.Time `json:"CreatedDateTime,omitempty"`
		Location        string    `json:"Location,omitempty"`
	} `json:"Data,omitempty"`
}

// GetVAppMachines - Get a list of vCloud Director Virtual Application Machines
func (v *VcloudClient) GetVAppMachines(vAppID int) (resp *GetVAppMachinesResponse, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&GetVAppMachinesResponse{}).
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vAppID": fmt.Sprintf("%d", vAppID),
		}).
		Get("/v6/vcloud/vapps/{vAppID}/machines")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*GetVAppMachinesResponse), nil
}

// * Protection Level

// ListProtectionLevelsAvailable - List the protection levels available for a vCloud Director Virtual Application
func (vApp *VApp) ListProtectionLevelsAvailable() (resp *ProtectionLevels, err error) {
	pL := ProtectionLevelClient{}
	return pL.ListProtectionLevels(listProtectionLevelsRequest{
		VAppID: vApp.GetIDPtr(),
	})
}

// GetProtectionLevelAvailableByName - Get a protection level by name for a vCloud Director Virtual Application
func (vApp *VApp) GetProtectionLevelAvailableByName(name string) (resp *ProtectionLevel, err error) {
	pL := ProtectionLevelClient{}
	return pL.getProtectionLevelByName(getProtectionLevelByNameRequest{
		VAppID:              vApp.GetIDPtr(),
		ProtectionLevelName: &name,
	})
}

// GetProtectionLevelAvailableByID - Get a protection level by ID for a vCloud Director Virtual Application
func (vApp *VApp) GetProtectionLevelAvailableByID(id int) (resp *ProtectionLevel, err error) {
	pL := ProtectionLevelClient{}
	return pL.getProtectionLevelByID(getProtectionLevelByIDRequest{
		VAppID:            vApp.GetIDPtr(),
		ProtectionLevelID: &id,
	})
}

// ListProtectionLevels - List the protection levels applied to a vCloud Director Virtual Application
func (vApp *VApp) ListProtectionLevels() (resp *ProtectionLevels, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return resp, err
	}

	r, err := c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vAppID": fmt.Sprintf("%d", vApp.GetID()),
		}).
		SetResult(&protectionLevelAppliedResponse{}).
		Get("/v6/vcloud/vapps/{vAppID}/protected")
	if err != nil {
		return resp, err
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	resp = &ProtectionLevels{}
	for _, pL := range r.Result().(*protectionLevelAppliedResponse).Data.ProtectedLevels {
		if strings.ToLower(pL.EntityType) == "vapp" {
			resp.append(pL.ProtectionLevel)
		}
	}

	return resp, nil
}

// * ProtectVApp

// ProtectVApp - Protect a vCloud Director Virtual Application
func (vApp *VApp) Protect(req ProtectUnprotectRequest) (job *commonnetbackup.JobAPIResponse, err error) {
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
		protectionLevel, err = vApp.GetProtectionLevelAvailableByName(req.ProtectionLevelName)
		if err != nil {
			return job, err
		}
	} else {
		protectionLevel, err = vApp.GetProtectionLevelAvailableByID(*req.ProtectionLevelID)
		if err != nil {
			return job, err
		}
	}

	r, err = c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vAppID": fmt.Sprintf("%d", vApp.GetID()),
		}).
		SetBody(protectBody{
			ProtectionLevelID: protectionLevel.GetID(),
			Paths:             []string{},
		}).
		SetResult(&commonnetbackup.JobAPIResponse{}).
		Post("/v6/vcloud/vapps/{vAppID}/protect")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*commonnetbackup.JobAPIResponse), nil
}

// * UnprotectVApp

// UnprotectVApp - Unprotect a vCloud Director Virtual Application
func (vApp *VApp) Unprotect(req ProtectUnprotectRequest) (job *commonnetbackup.JobAPIResponse, err error) {
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
		protectionLevel, err = vApp.GetProtectionLevelAvailableByName(req.ProtectionLevelName)
		if err != nil {
			return job, err
		}
	} else {
		protectionLevel, err = vApp.GetProtectionLevelAvailableByID(*req.ProtectionLevelID)
		if err != nil {
			return job, err
		}
	}

	r, err = c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"vAppID": fmt.Sprintf("%d", vApp.GetID()),
		}).
		SetBody(protectBody{
			ProtectionLevelID: protectionLevel.GetID(),
		}).
		SetResult(&commonnetbackup.JobAPIResponse{}).
		Post("/v6/vcloud/vapps/{vAppID}/unprotect")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*commonnetbackup.JobAPIResponse), nil
}
