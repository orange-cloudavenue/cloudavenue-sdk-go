/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"context"
	"fmt"

	"github.com/orange-cloudavenue/common-go/validators"

	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

// GetProperties retrieves the properties of the client from the Cloudavenue API.
// It refreshes the client session before making the request to ensure the session is valid.
//
// Returns:
// - client: The properties client.
// - error: An error if there was an issue with the request or response.
func (c *client) GetProperties(ctx context.Context) (values *PropertiesModel, err error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	r, err := c.clientCloudavenue.R().
		SetContext(ctx).
		SetResult(&propertiesResponse{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/configurations")
	if err != nil {
		return
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on get properties: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return &PropertiesModel{
		FullName:    r.Result().(*propertiesResponse).FullName,
		Description: r.Result().(*propertiesResponse).Description,
		Email:       r.Result().(*propertiesResponse).Email,

		BillingModel: r.Result().(*propertiesResponse).InternetBillingMode,
	}, nil
}

// UpdateProperties updates the properties of the client in the Cloudavenue API.
func (c *client) UpdateProperties(ctx context.Context, properties *PropertiesRequest) (job *commoncloudavenue.JobCreatedAPIResponse, err error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	if err := validators.New().StructCtx(ctx, properties); err != nil {
		return nil, err
	}

	r, err := c.clientCloudavenue.R().
		SetContext(ctx).
		SetBody(properties).
		SetResult(&commoncloudavenue.JobCreatedAPIResponse{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Put("/api/customers/v2.0/configurations")
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on update properties: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobCreatedAPIResponse), nil
}
