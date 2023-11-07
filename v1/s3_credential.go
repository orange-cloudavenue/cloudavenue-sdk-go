package v1

import (
	"fmt"
	"time"

	clients3 "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/s3"
)

type (
	S3Credentials []S3Credential
	S3Credential  struct {
		TenantID        string    `json:"tenantId"`
		StorageTenantID string    `json:"storageTenantId"`
		Owner           string    `json:"owner"`
		OwnerID         string    `json:"ownerId"`
		StorageUserID   string    `json:"storageUserId"`
		Type            string    `json:"type"`
		Immutable       bool      `json:"immutable"`
		AccessKey       string    `json:"accessKey"`
		SecretKey       string    `json:"secretKey"`
		Active          bool      `json:"active"`
		CreatedDate     time.Time `json:"createdDate"`
		LastAccessDate  time.Time `json:"lastAccessDate"`
		AppName         string    `json:"appName"`
		AllowedBuckets  []string  `json:"allowedBuckets"`
		UsedK8SClusters []string  `json:"usedK8sClusters"`
		ProviderOwner   bool      `json:"providerOwner"`
	}
)

func (c *S3Credential) GetTenantID() string {
	return c.TenantID
}

func (c *S3Credential) GetStorageTenantID() string {
	return c.StorageTenantID
}

func (c *S3Credential) GetOwner() string {
	return c.Owner
}

func (c *S3Credential) GetOwnerID() string {
	return c.OwnerID
}

func (c *S3Credential) GetStorageUserID() string {
	return c.StorageUserID
}

func (c *S3Credential) GetType() string {
	return c.Type
}

func (c *S3Credential) IsImmutable() bool {
	return c.Immutable
}

func (c *S3Credential) GetAccessKey() string {
	return c.AccessKey
}

func (c *S3Credential) GetSecretKey() string {
	return c.SecretKey
}

func (c *S3Credential) IsActive() bool {
	return c.Active
}

func (c *S3Credential) GetCreatedDate() time.Time {
	return c.CreatedDate
}

func (c *S3Credential) GetLastAccessDate() time.Time {
	return c.LastAccessDate
}

func (c *S3Credential) GetAppName() string {
	return c.AppName
}

func (c *S3Credential) GetAllowedBuckets() []string {
	return c.AllowedBuckets
}

func (c *S3Credential) GetUsedK8SClusters() []string {
	return c.UsedK8SClusters
}

func (c *S3Credential) IsProviderOwner() bool {
	return c.ProviderOwner
}

// GetCredentials - Get a list of credentials
func (s *S3User) GetCredentials() (resp *S3Credentials, err error) {
	type allCredentials struct {
		Items S3Credentials `json:"items"`
	}

	r, err := clients3.NewOSE().R().
		SetResult(&allCredentials{}).
		SetPathParams(map[string]string{
			"orgID":    clients3.GetOrganizationID(),
			"userName": s.GetName(),
		}).
		Get("/api/v1/core/tenants/{orgID}/users/{userName}/credentials")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, fmt.Errorf("error getting credential: %s", r.Error())
	}

	return &r.Result().(*allCredentials).Items, nil
}

// GetCredential - Get a credential by access key
func (s *S3User) GetCredential(accessKey string) (resp *S3Credential, err error) {
	r, err := clients3.NewOSE().R().
		SetResult(&S3Credential{}).
		SetPathParams(map[string]string{
			"orgID":     clients3.GetOrganizationID(),
			"userName":  s.GetName(),
			"accessKey": accessKey,
		}).
		Get("/api/v1/core/tenants/{orgID}/users/{userName}/credentials/{accessKey}")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, fmt.Errorf("error getting credential: %s", r.Error())
	}

	return r.Result().(*S3Credential), nil
}

// NewCredential - Create a new credential
func (s *S3User) NewCredential(username string) (resp *S3Credential, err error) {
	r, err := clients3.NewOSE().R().
		SetResult(&S3Credential{}).
		SetPathParams(map[string]string{
			"orgID":    clients3.GetOrganizationID(),
			"userName": username,
		}).
		Post("/api/v1/core/tenants/{orgID}/users/{userName}/credentials")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, fmt.Errorf("error creating credential: %s", r.Error())
	}

	return r.Result().(*S3Credential), nil
}

// DeleteCredential - Delete a credential
func (c *S3Credential) Delete() (err error) {
	r, err := clients3.NewOSE().R().
		SetPathParams(map[string]string{
			"orgID":     clients3.GetOrganizationID(),
			"userName":  c.GetOwner(),
			"accessKey": c.GetAccessKey(),
		}).
		Delete("/api/v1/core/tenants/{orgID}/users/{userName}/credentials/{accessKey}")
	if err != nil {
		return
	}

	if r.IsError() {
		return fmt.Errorf("error deleting credential: %s", r.Error())
	}

	return nil
}
