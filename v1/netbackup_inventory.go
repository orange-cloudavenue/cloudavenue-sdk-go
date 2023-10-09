package v1

import (
	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	commonnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
)

type InventoryClient struct{}

// Refresh refreshes the inventory.
func (i *InventoryClient) Refresh() (job *commonnetbackup.JobAPIResponse, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return job, err
	}

	r, err := c.R().
		SetError(&commonnetbackup.APIError{}).
		SetResult(&commonnetbackup.JobAPIResponse{}).
		SetHeader("Content-Length", "0").
		Post("/v6/assetimport/vcloud/tenants/import")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return r.Result().(*commonnetbackup.JobAPIResponse), nil
}
