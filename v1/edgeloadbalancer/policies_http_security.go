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

func (c *client) GetPoliciesHTTPSecurity(ctx context.Context, virtualServiceID string) (*PoliciesHTTPSecurityModel, error) {
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

	rules, err := getPoliciesHTTPSecurity(vs)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HTTP security rules: %w", err)
	}

	return (&PoliciesHTTPSecurityModel{}).fromVCD(virtualServiceID, &govcdtypes.AlbVsHttpSecurityRules{
		Values: func() (v []govcdtypes.AlbVsHttpSecurityRule) {
			if rules == nil {
				return nil
			}

			v = make([]govcdtypes.AlbVsHttpSecurityRule, len(rules))
			for i := range rules {
				v[i] = *rules[i]
			}

			return v
		}(),
	}), nil
}

var getPoliciesHTTPSecurity = func(virtualServiceClient fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
	return virtualServiceClient.GetAllHttpSecurityRules(nil)
}

func (c *client) UpdatePoliciesHTTPSecurity(ctx context.Context, policies *PoliciesHTTPSecurityModel) (*PoliciesHTTPSecurityModel, error) {
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

	policiesUpdated, err := updatePoliciesHTTPSecurity(vs, policies.toVCD())
	if err != nil {
		return nil, fmt.Errorf("error updating HTTP request rules: %w", err)
	}

	return policies.fromVCD(policies.VirtualServiceID, policiesUpdated), nil
}

var updatePoliciesHTTPSecurity = func(vs fakeVirtualServiceClient, policies *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
	policiesUpdated, err := vs.UpdateHttpSecurityRules(policies)
	if err != nil {
		return nil, fmt.Errorf("error updating HTTP request rules: %w", err)
	}
	return policiesUpdated, nil
}

func (c *client) DeletePoliciesHTTPSecurity(ctx context.Context, virtualServiceID string) error {
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
	_, err = updatePoliciesHTTPSecurity(vs, &govcdtypes.AlbVsHttpSecurityRules{})
	if err != nil {
		return fmt.Errorf("error deleting HTTP request rules: %w", err)
	}

	return nil
}
