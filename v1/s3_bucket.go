package v1

import (
	"fmt"
	"time"

	clients3 "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/s3"
)

type (
	SyncBucketResponse struct {
		ID               string    `json:"id"`
		VcdID            string    `json:"vcdId"`
		VcdAssociationID string    `json:"vcdAssociationId"`
		Description      string    `json:"description"`
		Status           string    `json:"status"`
		ResourceType     string    `json:"resourceType"`
		ResourceKey      string    `json:"resourceKey"`
		Progress         int       `json:"progress"`
		Tenant           string    `json:"tenant"`
		Owner            string    `json:"owner"`
		StartDate        time.Time `json:"startDate"`
		EndDate          time.Time `json:"endDate"`
		Reason           string    `json:"reason"`
		Metadata         struct {
			AdditionalProp1 string `json:"additionalProp1"`
			AdditionalProp2 string `json:"additionalProp2"`
			AdditionalProp3 string `json:"additionalProp3"`
		} `json:"metadata"`
	}
)

// SyncBucket - Syncs a bucket.
func (s S3Client) SyncBucket(bucketName string) (err error) {
	r, err := clients3.NewOSE().R().
		SetResult(&SyncBucketResponse{}).
		SetPathParams(map[string]string{
			"bucketName": bucketName,
		}).
		SetQueryParam("sync", "").
		Get("/api/v1/s3/{bucketName}")
	if err != nil {
		return
	}

	if r.IsError() {
		return fmt.Errorf("error syncing bucket: %s", r.Error())
	}

	return nil
}
