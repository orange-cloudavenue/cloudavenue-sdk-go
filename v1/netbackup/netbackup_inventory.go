/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package netbackup

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

	type jobAPIResponse struct {
		Data []struct {
			ID     int    `json:"Id,omitempty"`
			Status string `json:"Status,omitempty"`
		} `json:"data,omitempty"`
	}

	r, err := c.R().
		SetError(&commonnetbackup.APIError{}).
		SetResult(&jobAPIResponse{}).
		SetHeader("Content-Length", "0").
		Post("/v6/assetimport/vcloud/tenants/import")
	if err != nil {
		return job, err
	}

	if r.IsError() {
		return job, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	jAPIResponse := &commonnetbackup.JobAPIResponse{
		Data: struct {
			ID     int    `json:"Id,omitempty"`
			Status string `json:"Status,omitempty"`
		}{
			ID:     r.Result().(*jobAPIResponse).Data[0].ID,
			Status: r.Result().(*jobAPIResponse).Data[0].Status,
		},
	}

	return jAPIResponse, nil
}
