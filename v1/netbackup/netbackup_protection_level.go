package netbackup

type ProtectionLevelClient struct{}

type ProtectionLevel struct {
	ID                            int    `json:"Id,omitempty"`
	ProtectionTypeID              int    `json:"ProtectionTypeId,omitempty"`
	Name                          string `json:"Name,omitempty"`
	Code                          string `json:"Code,omitempty"`
	Description                   string `json:"Description,omitempty"`
	Sequence                      int    `json:"Sequence,omitempty"`
	Color                         string `json:"Color,omitempty"`
	RequestTypeCode               string `json:"RequestTypeCode,omitempty"`
	IsVisible                     bool   `json:"IsVisible,omitempty"`
	IsBackupNow                   bool   `json:"IsBackupNow,omitempty"`
	IsManaged                     bool   `json:"IsManaged,omitempty"`
	SupportsFileProtect           bool   `json:"SupportsFileProtect,omitempty"`
	SupportsSingleClientBackupNow bool   `json:"SupportsSingleClientBackupNow,omitempty"`
}

// GetID returns the ID field of GetProtectionLevelResponse
func (r *ProtectionLevel) GetID() int {
	return r.ID
}

// GetProtectionTypeID returns the ProtectionTypeID field of ProtectionLevel
func (r *ProtectionLevel) GetProtectionTypeID() int {
	return r.ProtectionTypeID
}

// GetName returns the Name field of ProtectionLevel
func (r *ProtectionLevel) GetName() string {
	return r.Name
}

// GetCode returns the Code field of ProtectionLevel
func (r *ProtectionLevel) GetCode() string {
	return r.Code
}

// GetDescription returns the Description field of ProtectionLevel
func (r *ProtectionLevel) GetDescription() string {
	return r.Description
}

// GetSequence returns the Sequence field of ProtectionLevel
func (r *ProtectionLevel) GetSequence() int {
	return r.Sequence
}

// GetColor returns the Color field of ProtectionLevel
func (r *ProtectionLevel) GetColor() string {
	return r.Color
}

// GetRequestTypeCode returns the RequestTypeCode field of ProtectionLevel
func (r *ProtectionLevel) GetRequestTypeCode() string {
	return r.RequestTypeCode
}

// GetIsVisible returns the IsVisible field of ProtectionLevel
func (r *ProtectionLevel) GetIsVisible() bool {
	return r.IsVisible
}

// GetIsBackupNow returns the IsBackupNow field of ProtectionLevel
func (r *ProtectionLevel) GetIsBackupNow() bool {
	return r.IsBackupNow
}

// GetIsManaged returns the IsManaged field of ProtectionLevel
func (r *ProtectionLevel) GetIsManaged() bool {
	return r.IsManaged
}

// GetSupportsFileProtect returns the SupportsFileProtect field of ProtectionLevel
func (r *ProtectionLevel) GetSupportsFileProtect() bool {
	return r.SupportsFileProtect
}

// GetSupportsSingleClientBackupNow returns the SupportsSingleClientBackupNow field of ProtectionLevel
func (r *ProtectionLevel) GetSupportsSingleClientBackupNow() bool {
	return r.SupportsSingleClientBackupNow
}
