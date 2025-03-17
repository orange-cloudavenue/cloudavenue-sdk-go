/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

import (
	"context"
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers/validators"
)

func (c *client) GetPoliciesHTTPResponse(ctx context.Context, virtualServiceID string) (*PoliciesHTTPResponseModel, error) {
	if err := c.virtualServiceIDValidator(virtualServiceID); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	// * Get the virtual service
	vs, err := c.getVirtualService(ctx, "", virtualServiceID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving virtual service: %w", err)
	}

	rules, err := getPoliciesHTTPResponse(vs)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HTTP response rules: %w", err)
	}

	return (&PoliciesHTTPResponseModel{}).fromVCD(virtualServiceID, &govcdtypes.AlbVsHttpResponseRules{
		Values: func() (v []govcdtypes.AlbVsHttpResponseRule) {
			if rules == nil {
				return nil
			}

			v = make([]govcdtypes.AlbVsHttpResponseRule, len(rules))
			for i := range rules {
				v[i] = *rules[i]
			}

			return v
		}(),
	}), nil
}

var getPoliciesHTTPResponse = func(virtualServiceClient fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
	return virtualServiceClient.GetAllHttpResponseRules(nil)
}

func (c *client) UpdatePoliciesHTTPResponse(ctx context.Context, policies *PoliciesHTTPResponseModel) (*PoliciesHTTPResponseModel, error) {
	if err := validators.New().Struct(policies); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	// * Get the virtual service
	vs, err := c.getVirtualService(ctx, "", policies.VirtualServiceID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving virtual service: %w", err)
	}

	policiesUpdated, err := updatePoliciesHTTPResponse(vs, policies.toVCD())
	if err != nil {
		return nil, fmt.Errorf("error updating HTTP response rules: %w", err)
	}

	return policies.fromVCD(policies.VirtualServiceID, policiesUpdated), nil
}

var updatePoliciesHTTPResponse = func(vs fakeVirtualServiceClient, policies *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
	return vs.UpdateHttpResponseRules(policies)
}

func (c *client) DeletePoliciesHTTPResponse(ctx context.Context, virtualServiceID string) error {
	if err := c.virtualServiceIDValidator(virtualServiceID); err != nil {
		return err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	// * Get the virtual service
	vs, err := c.getVirtualService(ctx, "", virtualServiceID)
	if err != nil {
		return fmt.Errorf("error retrieving virtual service: %w", err)
	}
	_, err = updatePoliciesHTTPResponse(vs, &govcdtypes.AlbVsHttpResponseRules{})
	if err != nil {
		return fmt.Errorf("error deleting HTTP response rules: %w", err)
	}

	return nil
}
