package netbackup

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	commonnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
)

type MachineClient struct{}

// Machine - Is the response structure for the Machines APIs.
type Machine struct {
	IsVisibleToAllUsers       bool      `json:"IsVisibleToAllUsers,omitempty"`
	NetBackupClientName       string    `json:"NetBackupClientName,omitempty"`
	LastUpdatedDateTime       time.Time `json:"LastUpdatedDateTime,omitempty"`
	MachineCode               string    `json:"MachineCode,omitempty"` // Is a URN in the form of urn:vcloud:vm:xxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	DisplayName               string    `json:"DisplayName,omitempty"`
	ProtectionTypeName        string    `json:"ProtectionTypeName,omitempty"`
	Hardware                  string    `json:"Hardware,omitempty"`
	TrafficLightStatus        string    `json:"TrafficLightStatus,omitempty"`
	IsSyncingNetBackupData    bool      `json:"IsSyncingNetBackupData,omitempty"`
	PolicyAppendix            string    `json:"PolicyAppendix,omitempty"`
	AdditionalData            string    `json:"AdditionalData,omitempty"`
	Os                        string    `json:"OS,omitempty"`
	ImportSource              string    `json:"ImportSource,omitempty"`
	CustomerCode              string    `json:"CustomerCode,omitempty"` // Is a Org Name
	ProviderAssetType         string    `json:"ProviderAssetType,omitempty"`
	IsDeletedFromImportSource bool      `json:"IsDeletedFromImportSource,omitempty"`
	VAppID                    int       `json:"VAppId,omitempty"`
	LastSuccessfulBackupDate  time.Time `json:"LastSuccessfulBackupDate,omitempty"`
	SyncLastError             string    `json:"SyncLastError,omitempty"`
	ProtectionTypeID          int       `json:"ProtectionTypeId,omitempty"`
	IsInVCloud                bool      `json:"IsInVCloud,omitempty"`
	CatalogName               string    `json:"CatalogName,omitempty"`
	ID                        int       `json:"Id,omitempty"`
	VMDisplayName             string    `json:"VMDisplayName,omitempty"` // Is a VM Name with suffix "-Random4Letters/Numbers" (e.g. "demo-4f5a")
	CreatedDateTime           time.Time `json:"CreatedDateTime,omitempty"`
	Location                  string    `json:"Location,omitempty"`
}

// GetID returns the ID field of Machine.
func (m *Machine) GetID() int {
	return m.ID
}

// GetIDPtr returns a pointer to the ID field of Machine.
func (m *Machine) GetIDPtr() *int {
	return &m.ID
}

// GetIsVisibleToAllUsers returns the IsVisibleToAllUsers field of Machine.
func (m *Machine) GetIsVisibleToAllUsers() bool {
	return m.IsVisibleToAllUsers
}

// GetNetBackupClientName returns the NetBackupClientName field of Machine.
func (m *Machine) GetNetBackupClientName() string {
	return m.NetBackupClientName
}

// GetLastUpdatedDateTime returns the LastUpdatedDateTime field of Machine.
func (m *Machine) GetLastUpdatedDateTime() time.Time {
	return m.LastUpdatedDateTime
}

// GetMachineCode returns the MachineCode field of Machine.
func (m *Machine) GetMachineCode() string {
	return m.MachineCode
}

// GetDisplayName returns the DisplayName field of Machine.
func (m *Machine) GetDisplayName() string {
	return m.DisplayName
}

// GetName returns the DisplayName field of Machine.
func (m *Machine) GetName() string {
	return m.DisplayName
}

// GetProtectionTypeName returns the ProtectionTypeName field of Machine.
func (m *Machine) GetProtectionTypeName() string {
	return m.ProtectionTypeName
}

// GetHardware returns the Hardware field of Machine.
func (m *Machine) GetHardware() string {
	return m.Hardware
}

// GetTrafficLightStatus returns the TrafficLightStatus field of Machine.
func (m *Machine) GetTrafficLightStatus() string {
	return m.TrafficLightStatus
}

// GetIsSyncingNetBackupData returns the IsSyncingNetBackupData field of Machine.
func (m *Machine) GetIsSyncingNetBackupData() bool {
	return m.IsSyncingNetBackupData
}

// GetPolicyAppendix returns the PolicyAppendix field of Machine.
func (m *Machine) GetPolicyAppendix() string {
	return m.PolicyAppendix
}

// GetAdditionalData returns the AdditionalData field of Machine.
func (m *Machine) GetAdditionalData() string {
	return m.AdditionalData
}

// GetOs returns the Os field of Machine.
func (m *Machine) GetOs() string {
	return m.Os
}

// GetImportSource returns the ImportSource field of Machine.
func (m *Machine) GetImportSource() string {
	return m.ImportSource
}

// GetCustomerCode returns the CustomerCode field of Machine.
func (m *Machine) GetCustomerCode() string {
	return m.CustomerCode
}

// GetProviderAssetType returns the ProviderAssetType field of Machine.
func (m *Machine) GetProviderAssetType() string {
	return m.ProviderAssetType
}

// GetIsDeletedFromImportSource returns the IsDeletedFromImportSource field of Machine.
func (m *Machine) GetIsDeletedFromImportSource() bool {
	return m.IsDeletedFromImportSource
}

// GetVAppID returns the VAppID field of Machine.
func (m *Machine) GetVAppID() int {
	return m.VAppID
}

// GetLastSuccessfulBackupDate returns the LastSuccessfulBackupDate field of Machine.
func (m *Machine) GetLastSuccessfulBackupDate() time.Time {
	return m.LastSuccessfulBackupDate
}

// GetSyncLastError returns the SyncLastError field of Machine.
func (m *Machine) GetSyncLastError() string {
	return m.SyncLastError
}

// GetProtectionTypeID returns the ProtectionTypeID field of Machine.
func (m *Machine) GetProtectionTypeID() int {
	return m.ProtectionTypeID
}

// GetIsInVCloud returns the IsInVCloud field of Machine.
func (m *Machine) GetIsInVCloud() bool {
	return m.IsInVCloud
}

// GetCatalogName returns the CatalogName field of Machine.
func (m *Machine) GetCatalogName() string {
	return m.CatalogName
}

// GetVMDisplayName returns the VMDisplayName field of Machine.
func (m *Machine) GetVMDisplayName() string {
	return m.VMDisplayName
}

// GetCreatedDateTime returns the CreatedDateTime field of Machine.
func (m *Machine) GetCreatedDateTime() time.Time {
	return m.CreatedDateTime
}

// GetLocation returns the Location field of Machine.
func (m *Machine) GetLocation() string {
	return m.Location
}

// * Machines
// Machines - Is the response structure for the GetMachines API.
type Machines []Machine

// append - Append a Machine to the Machines slice.
func (m *Machines) append(machine Machine) { //nolint:unused
	*m = append(*m, machine)
}

type machinesResponse struct {
	Data Machines `json:"data,omitempty"`
}

// GetMachines - Get a list of NetBackup Machines.
func (m *MachineClient) GetMachines() (resp *Machines, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&machinesResponse{}).
		SetError(&commonnetbackup.APIError{}).
		Get("/v6/machines")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*machinesResponse).Data, nil
}

type machineResponse struct {
	Data Machine `json:"data,omitempty"`
}

// GetMachineByID - Get a NetBackup Machine by ID.
func (m *MachineClient) GetMachineByID(id int) (resp *Machine, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&machineResponse{}).
		SetError(&commonnetbackup.APIError{}).
		SetPathParam("id", fmt.Sprintf("%d", id)).
		Get("/v6/machines/{id}")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*machineResponse).Data, nil
}

// GetMachineByName - Get a NetBackup Machine by Name.
func (m *MachineClient) GetMachineByName(name string) (resp *Machine, err error) {
	machines, err := m.GetMachines()
	if err != nil {
		return
	}

	for _, machine := range *machines {
		if machine.DisplayName == name {
			return &machine, nil
		}
	}

	return resp, fmt.Errorf("Machine with name %s not found", name)
}

// GetMachineByIdentifier - Get a NetBackup Machine by Identifier.
func (m *MachineClient) GetMachineByIdentifier(identifier string) (resp *Machine, err error) {
	machines, err := m.GetMachines()
	if err != nil {
		return
	}

	for _, machine := range *machines {
		if machine.MachineCode == identifier {
			return &machine, nil
		}
	}

	return resp, fmt.Errorf("Machine with identifier %s not found", identifier)
}

// GetMachineByNameOrIdentifier - Get a NetBackup Machine by Name or Identifier.
func (m *MachineClient) GetMachineByNameOrIdentifier(nameOrIdentifier string) (resp *Machine, err error) {
	machines, err := m.GetMachines()
	if err != nil {
		return
	}

	for _, machine := range *machines {
		if machine.DisplayName == nameOrIdentifier || machine.MachineCode == nameOrIdentifier {
			return &machine, nil
		}
	}

	return resp, fmt.Errorf("Machine with name or identifier %s not found", nameOrIdentifier)
}

// * Protection Level

// ListProtectionLevelsAvailable - List the protection levels available for a Machine.
func (m *Machine) ListProtectionLevelsAvailable() (resp *ProtectionLevels, err error) {
	pL := ProtectionLevelClient{}
	return pL.ListProtectionLevels(listProtectionLevelsRequest{
		MachineID: m.GetIDPtr(),
	})
}

// GetProtectionLevelAvailableByName - Get a protection level available for a Machine by Name.
func (m *Machine) GetProtectionLevelAvailableByName(name string) (resp *ProtectionLevel, err error) {
	pL := ProtectionLevelClient{}
	return pL.getProtectionLevelByName(getProtectionLevelByNameRequest{
		MachineID:           m.GetIDPtr(),
		ProtectionLevelName: &name,
	})
}

// GetProtectionLevelAvailableByID - Get a protection level available for a Machine by ID.
func (m *Machine) GetProtectionLevelAvailableByID(id int) (resp *ProtectionLevel, err error) {
	pL := ProtectionLevelClient{}
	return pL.getProtectionLevelByID(getProtectionLevelByIDRequest{
		MachineID:         m.GetIDPtr(),
		ProtectionLevelID: &id,
	})
}

// ListProtectionLevels - List the protection levels applied to a Machine.
func (m *Machine) ListProtectionLevels() (resp *ProtectionLevels, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return resp, err
	}

	r, err := c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"machineID": fmt.Sprintf("%d", m.GetID()),
		}).
		SetResult(&protectionLevelAppliedResponse{}).
		Get("/v6/machines/{machineID}/protected")
	if err != nil {
		return resp, err
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	resp = &ProtectionLevels{}
	for _, pL := range r.Result().(*protectionLevelAppliedResponse).Data.ProtectedLevels {
		if strings.ToLower(pL.EntityType) == "machine" {
			resp.append(pL.ProtectionLevel)
		}
	}

	return resp, nil
}

// * Protect Machine

// Protect - Protect a Machine.
func (m *Machine) Protect(req ProtectUnprotectRequest) (job *commonnetbackup.JobAPIResponse, err error) {
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
		protectionLevel, err = m.GetProtectionLevelAvailableByName(req.ProtectionLevelName)
		if err != nil {
			return job, err
		}
	} else {
		protectionLevel, err = m.GetProtectionLevelAvailableByID(*req.ProtectionLevelID)
		if err != nil {
			return job, err
		}
	}

	r, err = c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"machineID": fmt.Sprintf("%d", m.GetID()),
		}).
		SetBody(protectBody{
			ProtectionLevelID: protectionLevel.GetID(),
			Paths:             []string{},
		}).
		SetResult(&commonnetbackup.JobAPIResponse{}).
		Post("/v6/machines/{machineID}/protect")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*commonnetbackup.JobAPIResponse), nil
}

// * Unprotect Machine

// Unprotect - Unprotect a Machine.
func (m *Machine) Unprotect(req ProtectUnprotectRequest) (job *commonnetbackup.JobAPIResponse, err error) {
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
		protectionLevel, err = m.GetProtectionLevelAvailableByName(req.ProtectionLevelName)
		if err != nil {
			return job, err
		}
	} else {
		protectionLevel, err = m.GetProtectionLevelAvailableByID(*req.ProtectionLevelID)
		if err != nil {
			return job, err
		}
	}

	r, err = c.R().
		SetError(&commonnetbackup.APIError{}).
		SetPathParams(map[string]string{
			"machineID": fmt.Sprintf("%d", m.GetID()),
		}).
		SetBody(protectBody{
			ProtectionLevelID: protectionLevel.GetID(),
			Paths:             []string{},
		}).
		SetResult(&commonnetbackup.JobAPIResponse{}).
		Post("/v6/machines/{machineID}/unprotect")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*commonnetbackup.JobAPIResponse), nil
}
