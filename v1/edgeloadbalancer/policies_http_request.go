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

	"github.com/orange-cloudavenue/common-go/validators"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

func (c *client) GetPoliciesHTTPRequest(ctx context.Context, virtualServiceID string) (*PoliciesHTTPRequestModel, error) {
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

	rules, err := getPoliciesHTTPRequest(vs)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HTTP request rules: %w", err)
	}

	return (&PoliciesHTTPRequestModel{}).fromVCD(virtualServiceID, &govcdtypes.AlbVsHttpRequestRules{
		Values: func() (v []govcdtypes.AlbVsHttpRequestRule) {
			if rules == nil {
				return nil
			}

			v = make([]govcdtypes.AlbVsHttpRequestRule, len(rules))
			for i := range rules {
				v[i] = *rules[i]
			}

			return v
		}(),
	}), nil
}

var getPoliciesHTTPRequest = func(virtualServiceClient fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
	return virtualServiceClient.GetAllHttpRequestRules(nil)
}

func (c *client) UpdatePoliciesHTTPRequest(ctx context.Context, policies *PoliciesHTTPRequestModel) (*PoliciesHTTPRequestModel, error) {
	if err := validators.New().StructCtx(ctx, policies); err != nil {
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

	policiesUpdated, err := updatePoliciesHTTPRequest(vs, policies.toVCD())
	if err != nil {
		return nil, fmt.Errorf("error updating HTTP request rules: %w", err)
	}

	return policies.fromVCD(policies.VirtualServiceID, policiesUpdated), nil
}

var updatePoliciesHTTPRequest = func(vs fakeVirtualServiceClient, policies *govcdtypes.AlbVsHttpRequestRules) (*govcdtypes.AlbVsHttpRequestRules, error) {
	policiesUpdated, err := vs.UpdateHttpRequestRules(policies)
	if err != nil {
		return nil, fmt.Errorf("error updating HTTP request rules: %w", err)
	}
	return policiesUpdated, nil
}

func (c *client) DeletePoliciesHTTPRequest(ctx context.Context, virtualServiceID string) error {
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
	_, err = updatePoliciesHTTPRequest(vs, &govcdtypes.AlbVsHttpRequestRules{})
	if err != nil {
		return fmt.Errorf("error deleting HTTP request rules: %w", err)
	}

	return nil
}
