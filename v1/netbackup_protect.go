package v1

type protectBody struct {
	ProtectionLevelID int      `json:"ProtectionLevelId"`
	Paths             []string `json:"Paths,omitempty"`
}

// ProtectUnprotectRequest - Is the request structure for the Protects and Unprotects APIs
type ProtectUnprotectRequest struct {
	// One of the following protection level settings must be specified
	// The ID of the protection level in the netbackup system
	ProtectionLevelID *int // Optional if ProtectionLevelName is specified
	// The name of the protection level in the netbackup system
	ProtectionLevelName string // Optional if ProtectionLevelID is specified
}

type protectionLevelAppliedResponse struct {
	Data struct {
		ProtectedLevels []struct {
			EntityType      string          `json:"EntityType"`
			ProtectionLevel ProtectionLevel `json:"ProtectionLevel"`
		} `json:"ProtectedLevels"`
	} `json:"data"`
}
