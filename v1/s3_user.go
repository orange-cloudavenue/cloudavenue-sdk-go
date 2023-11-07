package v1

import (
	"fmt"

	clients3 "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/s3"
)

type (
	S3Users []S3User
	S3User  struct {
		Name             string   `json:"name"`
		ID               string   `json:"id"`
		FullName         string   `json:"fullName"`
		Role             string   `json:"role"`
		SubordinateRoles []string `json:"subordinateRoles"`
		Active           bool     `json:"active"`
		CurrentRegion    string   `json:"currentRegion"`
		Site             struct {
			ID           string `json:"id"`
			Name         string `json:"name"`
			RestEndpoint string `json:"restEndpoint"`
		} `json:"site"`
		Usage struct {
			Scope          string `json:"scope"`
			TotalBytes     int    `json:"totalBytes"`
			UsedBytes      int    `json:"usedBytes"`
			AvailableBytes int    `json:"availableBytes"`
			ObjectCount    int    `json:"objectCount"`
			BucketCount    int    `json:"bucketCount"`
		} `json:"usage"`
		ExistedOnStorage bool     `json:"existedOnStorage"`
		CanoncialID      string   `json:"storageCanonicalId"`
		PolicyArns       []string `json:"policyArns"`
		Remote           bool     `json:"remote"`
		PoseAsUser       bool     `json:"poseAsUser"`
		SourceTenant     string   `json:"sourceTenant"`
	}
)

func (s *S3User) GetName() string {
	return s.Name
}

func (s *S3User) GetID() string {
	return s.ID
}

func (s *S3User) GetFullName() string {
	return s.FullName
}

func (s *S3User) GetRole() string {
	return s.Role
}

func (s *S3User) GetSubordinateRoles() []string {
	return s.SubordinateRoles
}

func (s *S3User) IsActive() bool {
	return s.Active
}

func (s *S3User) GetCurrentRegion() string {
	return s.CurrentRegion
}

func (s *S3User) GetSiteID() string {
	return s.Site.ID
}

func (s *S3User) GetSiteName() string {
	return s.Site.Name
}

func (s *S3User) GetSiteRestEndpoint() string {
	return s.Site.RestEndpoint
}

func (s *S3User) GetUsageScope() string {
	return s.Usage.Scope
}

func (s *S3User) GetTotalBytes() int {
	return s.Usage.TotalBytes
}

func (s *S3User) GetUsedBytes() int {
	return s.Usage.UsedBytes
}

func (s *S3User) GetAvailableBytes() int {
	return s.Usage.AvailableBytes
}

func (s *S3User) GetObjectCount() int {
	return s.Usage.ObjectCount
}

func (s *S3User) GetBucketCount() int {
	return s.Usage.BucketCount
}

func (s *S3User) GetExistedOnStorage() bool {
	return s.ExistedOnStorage
}

func (s *S3User) GetPolicyArns() []string {
	return s.PolicyArns
}

func (s *S3User) IsRemote() bool {
	return s.Remote
}

func (s *S3User) IsPoseAsUser() bool {
	return s.PoseAsUser
}

func (s *S3User) GetSourceTenant() string {
	return s.SourceTenant
}

// * user

// GetUser - Get a S3 user by username
func (s S3Client) GetUser(username string) (resp *S3User, err error) {
	r, err := clients3.NewOSE().R().
		SetResult(&S3User{}).
		SetPathParams(map[string]string{
			"orgID":    clients3.GetOrganizationID(),
			"userName": username,
		}).
		Get("/api/v1/core/tenants/{orgID}/users/{userName}")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, fmt.Errorf("error getting user: %s", r.Error())
	}

	return r.Result().(*S3User), nil
}

// GetCanonicalID - Get a S3 user canonical ID by username
func (s *S3User) GetCanonicalID() (resp string, err error) {
	if s.CanoncialID == "" {
		r, err := clients3.NewOSE().R().
			SetResult(&resp).
			SetPathParams(map[string]string{
				"orgID":    clients3.GetOrganizationID(),
				"userName": s.GetName(),
			}).
			Get("/api/v1/core/tenants/{orgID}/users/{userName}/canonical-id")
		if err != nil {
			return resp, err
		}

		if r.IsError() {
			return resp, fmt.Errorf("error getting canonical ID: %s", r.Error())
		}

		s.CanoncialID = *r.Result().(*string)
	}

	return s.CanoncialID, nil
}

// * users

// GetUsers - Get all S3 users
func (s S3Client) GetUsers() (resp *S3Users, err error) {
	type allUsers struct {
		Items S3Users `json:"items"`
	}

	r, err := clients3.NewOSE().R().
		SetResult(&allUsers{}).
		SetPathParams(map[string]string{
			"orgID": clients3.GetOrganizationID(),
		}).
		Get("/api/v1/core/tenants/{orgID}/users")
	if err != nil {
		return
	}

	if r.IsError() {
		return resp, fmt.Errorf("error getting users: %s", r.Error())
	}

	return &r.Result().(*allUsers).Items, nil
}

// UserExists - Check if a user exists
func (s *S3Users) UserExists(username string) (exist bool, user *S3User) {
	for _, user := range *s {
		if user.Name == username {
			return true, &user
		}
	}

	return false, nil
}
