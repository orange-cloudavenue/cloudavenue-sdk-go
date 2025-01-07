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
	"fmt"

	"github.com/go-resty/resty/v2"

	clientnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/netbackup"
	commonnetbackup "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/netbackup"
)

type ProtectionLevels []ProtectionLevel

// append - Append a ProtectionLevel to the ProtectionLevels slice.
func (p *ProtectionLevels) append(protectionLevel ProtectionLevel) {
	*p = append(*p, protectionLevel)
}

type protectionLevelsResponse struct {
	Data ProtectionLevels `json:"data,omitempty"`
}

type listProtectionLevelsRequest struct {
	VAppID    *int
	VDCID     *int
	MachineID *int
}

// ListProtectionLevels - Get a list of protection levels
// Use listProtectionLevelsRequest to specify the VAppID, VDCID or MachineID.
func (p *ProtectionLevelClient) ListProtectionLevels(req listProtectionLevelsRequest) (resp *ProtectionLevels, err error) {
	c, err := clientnetbackup.New()
	if err != nil {
		return resp, err
	}

	if req.VAppID == nil && req.VDCID == nil && req.MachineID == nil {
		return resp, fmt.Errorf("you must specify a VAppID, VDCID or MachineID")
	}

	var r *resty.Response

	hReq := c.R().
		SetResult(&protectionLevelsResponse{}).
		SetError(&commonnetbackup.APIError{})

	switch {
	case req.VAppID != nil:
		r, err = hReq.
			SetPathParams(map[string]string{
				"VAppId": fmt.Sprintf("%d", *req.VAppID),
			}).
			Get("/v6/vcloud/vapps/{VAppId}/protection/levels")
	case req.VDCID != nil:
		r, err = hReq.
			SetPathParams(map[string]string{
				"VdcId": fmt.Sprintf("%d", *req.VDCID),
			}).
			Get("/v6/vcloud/vdcs/{VdcId}/protection/levels")
	case req.MachineID != nil:
		r, err = hReq.
			SetPathParams(map[string]string{
				"MachineId": fmt.Sprintf("%d", *req.MachineID),
			}).
			SetQueryParam("MachineId", fmt.Sprintf("%d", *req.MachineID)).
			Get("/v6/machines/{MachineId}/protection/levels")
	}
	if err != nil {
		return resp, err
	}

	if r.IsError() {
		return resp, commonnetbackup.ToError(r.Error().(*commonnetbackup.APIError))
	}

	return &r.Result().(*protectionLevelsResponse).Data, nil
}
