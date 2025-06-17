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

	return (&PoliciesHTTPSecurityModel{}).fromVCD(virtualServiceID, rules), nil
}

var getPoliciesHTTPSecurity = func(virtualServiceClient fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
	return virtualServiceClient.GetAllHttpSecurityRules(nil)
}

func (c *client) UpdatePoliciesHTTPSecurity(ctx context.Context, policies *PoliciesHTTPSecurityModel) (*PoliciesHTTPSecurityModel, error) {
	if err := validators.New().StructCtx(ctx, policies); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, fmt.Errorf("error refreshing client: %w", err)
	}

	// * Get the virtual service
	vs, err := c.getVirtualService(ctx, "", policies.VirtualServiceID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving virtual service: %w", err)
	}

	policiesUpdated, err := updatePoliciesHTTPSecurity(vs, policies.toVCD())
	if err != nil {
		return nil, fmt.Errorf("error updating HTTP security rules: %w", err)
	}

	// Convert the updated policies to a slice of pointers
	var rulesUpdated []*govcdtypes.AlbVsHttpSecurityRule
	for i := range policiesUpdated.Values {
		rulesUpdated = append(rulesUpdated, &policiesUpdated.Values[i])
	}

	return policies.fromVCD(policies.VirtualServiceID, rulesUpdated), nil
}

var updatePoliciesHTTPSecurity = func(vs fakeVirtualServiceClient, policies *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
	return vs.UpdateHttpSecurityRules(policies)
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
		return fmt.Errorf("error deleting HTTP security rules: %w", err)
	}

	return nil
}
