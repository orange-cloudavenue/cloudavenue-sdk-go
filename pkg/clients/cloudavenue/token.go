package clientcloudavenue

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

type token struct {
	baererToken string
	expiresAt   time.Time

	username   string
	password   string
	org        string
	orgID      string
	vdc        string
	vcdVersion string

	endpoint string

	debug bool
}

// GetOrganization - Returns the organization
func (t *token) GetOrganization() string {
	return t.org
}

// GetVCDVersion - Returns the VCD version
func (t *token) GetVCDVersion() string {
	return t.vcdVersion
}

// GetEndpoint - Returns the API endpoint
func (t *token) GetEndpoint() string {
	return t.endpoint
}

// GetEndpointURL - Returns the API endpoint URL
func (t *token) GetEndpointURL() url.URL {
	u, _ := url.Parse(t.endpoint)
	return *u
}

// IsExpired - Returns true if the token is expired
func (t *token) IsExpired() bool {
	return t.expiresAt.Before(time.Now())
}

// IsSet - Returns true if the token is set
func (t *token) IsSet() bool {
	return t.baererToken != ""
}

// GetToken - Returns the token
func (t *token) GetToken() string {
	return t.baererToken
}

// GetOrgID - Returns the organization ID
func (t *token) GetOrgID() string {
	return t.orgID
}

// RefreshToken - Refreshes the token
func (t *token) RefreshToken() error {
	if !t.IsSet() || t.IsExpired() {
		c := resty.New().SetBaseURL(t.endpoint)

		r, err := c.R().
			SetDebug(t.debug).
			// SetHeader("Content-Type", "application/x-www-form-url-encoded").
			SetHeader("Accept", "application/json;version="+t.vcdVersion).
			SetResult(&authTokenResponse{}).
			SetError(&APIErrorResponse{}).
			SetBasicAuth(t.username+"@"+t.org, t.password).
			Post("/cloudapi/1.0.0/sessions")
		if err != nil {
			return err
		}

		if r.IsError() {
			return fmt.Errorf("authentification failed : HTTPCode:%s - %s ", r.Status(), r.Error().(*APIErrorResponse).FormatError())
		}
		// Set the token
		t.baererToken = r.Header().Get("X-VMWARE-VCLOUD-ACCESS-TOKEN")

		// Calculate the expiration date
		refreshedToken := r.Result().(*authTokenResponse)
		// Set OrganizationID
		t.orgID = refreshedToken.Org.ID
		t.expiresAt = time.Now().Add(time.Duration(refreshedToken.SessionIdleTimeoutMinutes) * time.Minute)
	}
	return nil
}

type authTokenResponse struct {
	ID   string `json:"id"`
	User struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"user"`
	Org struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"org"`
	Location string   `json:"location"`
	Roles    []string `json:"roles"`
	RoleRefs []struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"roleRefs"`
	SessionIdleTimeoutMinutes int `json:"sessionIdleTimeoutMinutes"`
}

type APIErrorResponse struct {
	Code    string `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// FormatError - Formats the error
func (e *APIErrorResponse) FormatError() string {
	return fmt.Sprintf("ErrorCode:%s - ErrorReason:%s - ErrorMessage:%s", e.Code, e.Reason, e.Message)
}

// ToError - Converts an APIErrorResponse to an error
func ToError(e *APIErrorResponse) error {
	return fmt.Errorf("error on API call: %s", e.FormatError())
}
