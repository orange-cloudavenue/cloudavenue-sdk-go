/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

// GetAllNetworkContextProfiles returns all Network Context Profiles available
// in the context of the VDC Group (SYSTEM + PROVIDER + TENANT scopes).
func (g VDCGroup) GetAllNetworkContextProfiles() ([]*NetworkContextProfile, error) {
	cavc, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	client := &cavc.Vmware.Client

	// Use the named filter parameter introduced in API 38.0,
	// replacing the deprecated _context== format.
	queryParams := url.Values{}
	queryParams.Set("filter", fmt.Sprintf("vdcGroupId==%s", g.vg.VdcGroup.Id))

	raw, err := govcd.GetAllNetworkContextProfiles(client, queryParams)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Network Context Profiles for VDC Group: %w", err)
	}

	profiles := make([]*NetworkContextProfile, len(raw))
	for i, r := range raw {
		profiles[i] = networkContextProfileFromGovcd(r)
	}

	return profiles, nil
}

// GetNetworkContextProfileByName returns a single Network Context Profile by name
// within the VDC Group context.
func (g VDCGroup) GetNetworkContextProfileByName(name string) (*NetworkContextProfile, error) {
	all, err := g.GetAllNetworkContextProfiles()
	if err != nil {
		return nil, err
	}

	var found []*NetworkContextProfile
	for _, p := range all {
		if p.Name == name {
			found = append(found, p)
		}
	}

	switch len(found) {
	case 0:
		return nil, fmt.Errorf("%w: network context profile with name %q not found", govcd.ErrorEntityNotFound, name)
	case 1:
		return found[0], nil
	default:
		return nil, fmt.Errorf("found %d network context profiles with name %q, please use ID to disambiguate", len(found), name)
	}
}

// GetNetworkContextProfileByID returns a single Network Context Profile by its URN ID.
func (g VDCGroup) GetNetworkContextProfileByID(id string) (*NetworkContextProfile, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
	}

	cavc, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	c := &cavc.Vmware.Client

	urlRef, err := c.OpenApiBuildEndpoint(networkContextProfilesEndpoint() + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	result := &networkContextProfileAPIPayload{}
	if err := c.OpenApiGetItem(c.APIVersion, urlRef, nil, result, nil); err != nil {
		return nil, fmt.Errorf("error retrieving Network Context Profile %q: %w", id, err)
	}

	return networkContextProfileFromAPIPayload(result), nil
}

// CreateNetworkContextProfile creates a new TENANT-scoped Network Context Profile
// within the VDC Group context.
func (g VDCGroup) CreateNetworkContextProfile(profile *NetworkContextProfile) (*NetworkContextProfile, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile must not be nil")
	}
	if profile.Name == "" {
		return nil, fmt.Errorf("profile.Name must not be empty")
	}

	cavc, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	client := &cavc.Vmware.Client

	urlRef, err := client.OpenApiBuildEndpoint(networkContextProfilesEndpoint())
	if err != nil {
		return nil, fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	payload := networkContextProfileToAPIPayload(profile, g.vg.VdcGroup.Id, cavc.Org.Org.ID)

	task, err := client.OpenApiPostItemAsync(client.APIVersion, urlRef, nil, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating Network Context Profile: %w", err)
	}

	if err := task.WaitTaskCompletion(); err != nil {
		return nil, fmt.Errorf("error waiting for Network Context Profile creation task: %w", err)
	}

	return g.GetNetworkContextProfileByName(profile.Name)
}

// UpdateNetworkContextProfile updates an existing TENANT-scoped Network Context Profile.
func (g VDCGroup) UpdateNetworkContextProfile(profile *NetworkContextProfile) (*NetworkContextProfile, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile must not be nil")
	}
	if profile.ID == "" {
		return nil, fmt.Errorf("profile.ID must not be empty")
	}

	cavc, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	client := &cavc.Vmware.Client

	urlRef, err := client.OpenApiBuildEndpoint(networkContextProfilesEndpoint() + "/" + profile.ID)
	if err != nil {
		return nil, fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	payload := networkContextProfileToAPIPayload(profile, g.vg.VdcGroup.Id, cavc.Org.Org.ID)

	task, err := client.OpenApiPutItemAsync(client.APIVersion, urlRef, nil, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error updating Network Context Profile %q: %w", profile.ID, err)
	}

	if err := task.WaitTaskCompletion(); err != nil {
		return nil, fmt.Errorf("error waiting for Network Context Profile update task: %w", err)
	}

	return g.GetNetworkContextProfileByID(profile.ID)
}

// DeleteNetworkContextProfile deletes a TENANT-scoped Network Context Profile by ID.
func (g VDCGroup) DeleteNetworkContextProfile(id string) error {
	if id == "" {
		return fmt.Errorf("id must not be empty")
	}

	cavc, err := clientcloudavenue.New()
	if err != nil {
		return fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	c := &cavc.Vmware.Client

	urlRef, err := c.OpenApiBuildEndpoint(networkContextProfilesEndpoint() + "/" + id)
	if err != nil {
		return fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	if err := c.OpenApiDeleteItem(c.APIVersion, urlRef, nil, nil); err != nil {
		return fmt.Errorf("error deleting Network Context Profile %q: %w", id, err)
	}

	return nil
}

// GetNetworkContextProfileAttributes returns the server-side catalog of valid attribute
// values for Network Context Profiles, scoped to this VDC Group.
// The returned catalog lists valid DOMAIN_NAME and APP_ID values for this platform instance.
// Only values present in this catalog can be used when creating or updating a profile.
func (g VDCGroup) GetNetworkContextProfileAttributes() (*NetworkContextProfileAttributesCatalog, error) {
	cavc, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	client := &cavc.Vmware.Client

	queryParams := url.Values{}
	queryParams.Set("filter", fmt.Sprintf("vdcGroupId==%s", g.vg.VdcGroup.Id))

	urlRef, err := client.OpenApiBuildEndpoint(networkContextProfilesEndpoint() + "/attributes")
	if err != nil {
		return nil, fmt.Errorf("error building networkContextProfiles/attributes endpoint: %w", err)
	}

	result := &networkContextProfileAttributesAPIResponse{}
	if err := client.OpenApiGetItem(client.APIVersion, urlRef, queryParams, result, nil); err != nil {
		return nil, fmt.Errorf("error retrieving Network Context Profile attributes: %w", err)
	}

	return networkContextProfileAttributesFromAPIResponse(result), nil
}
