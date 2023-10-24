package s3

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type token struct {
	userName             string
	cavToken             string
	organizationName     string
	organizationID       string
	accessKey, secretKey string

	oseEndpoint string
	s3Endpoint  string
	debug       bool
}

// GetEndpointOSE - Returns the OSE endpoint
func (t *token) GetEndpointOSE() string {
	return t.oseEndpoint
}

// GetEndpointS3 - Returns the S3 endpoint
func (t *token) GetEndpointS3() string {
	return t.s3Endpoint
}

// IsSet - Returns true if the accessKey and secretKey are set
func (t *token) IsSet() bool {
	return t.accessKey != "" && t.secretKey != ""
}

// GetToken - Returns the token
func (t *token) GetToken() string {
	return t.cavToken
}

// GetAccessKey - Returns the accessKey
func (t *token) GetAccessKey() string {
	return t.accessKey
}

// GetSecretKey - Returns the secretKey
func (t *token) GetSecretKey() string {
	return t.secretKey
}

// RefreshAccessKey - Refreshes the accessKey and secretKey
func (t *token) RefreshAccessKey() error {
	if !t.IsSet() {
		c := resty.New().
			SetDebug(t.debug).
			SetAuthToken(t.GetToken()).
			SetBaseURL(t.GetEndpointOSE())

		if t.organizationID == "" {
			type tenantsResponse struct {
				Items []struct {
					Name          string `json:"name"`
					OrgID         string `json:"orgId"`
					SiteName      string `json:"siteName"`
					SiteID        string `json:"siteId"`
					FullID        string `json:"fullId"`
					SitePortalURL string `json:"sitePortalUrl"`
					Accessible    bool   `json:"accessible"`
					Local         bool   `json:"local"`
				} `json:"items"`
				PageInfo struct {
					Offset int `json:"offset"`
					Limit  int `json:"limit"`
					Total  int `json:"total"`
				} `json:"pageInfo"`
			}

			// Get Organization ID
			r, err := c.R().
				SetQueryParam("accessible-only", "true").
				SetResult(&tenantsResponse{}).
				Get("/api/v1/core/associated-tenants")
			if err != nil {
				return err
			}

			if r.IsError() {
				return fmt.Errorf("error getting organization ID: %s", r.Error())
			}

			if len(r.Result().(*tenantsResponse).Items) == 0 {
				return fmt.Errorf("no accessible tenants found")
			}

			// find first tenant match with organizationName
			for _, item := range r.Result().(*tenantsResponse).Items {
				if item.Name == t.organizationName {
					t.organizationID = item.OrgID
					break
				}
			}
		}

		type credentialsResponse struct {
			Items []struct {
				TenantID        string    `json:"tenantId"`
				StorageTenantID string    `json:"storageTenantId"`
				Owner           string    `json:"owner,omitempty"`
				OwnerID         string    `json:"ownerId"`
				StorageUserID   string    `json:"storageUserId"`
				Immutable       bool      `json:"immutable"`
				AccessKey       string    `json:"accessKey"`
				SecretKey       string    `json:"secretKey"`
				Active          bool      `json:"active"`
				CreatedDate     time.Time `json:"createdDate"`
				AllowedBuckets  []any     `json:"allowedBuckets"`
				Extension       struct{}  `json:"extension"`
				ProviderOwner   bool      `json:"providerOwner"`
			} `json:"items"`
		}

		// Get Access Key / Secret Key
		r, err := c.R().
			SetHeader("Accept", "application/json").
			SetPathParams(map[string]string{
				"organizationID": t.organizationID,
				"userName":       t.userName,
			}).
			SetResult(&credentialsResponse{}).
			Get("/api/v1/core/tenants/{organizationID}/users/{userName}/credentials")
		if err != nil {
			return err
		}

		if r.IsError() {
			return fmt.Errorf("error getting access token: %s", r.Error())
		}

		if len(r.Result().(*credentialsResponse).Items) == 0 {
			return fmt.Errorf("no access token found for user %s", t.userName)
		}

		// find first immutable credentials
		for _, item := range r.Result().(*credentialsResponse).Items {
			if item.Immutable {
				t.accessKey = item.AccessKey
				t.secretKey = item.SecretKey
				break
			}
		}
	}
	return nil
}
