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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

// NsxtFirewallRuleContainerExtended mirrors govcdtypes.NsxtFirewallRuleContainer
// but adds NetworkContextProfiles per firewall rule, which the upstream govcd SDK omits.
// See: https://techdocs.broadcom.com/us/en/vmware-cis/nsx/nsxt-dc/3-0/administration-guide/inventory/profiles.html
type NsxtFirewallRuleContainerExtended struct {
	UserDefinedRules []*NsxtFirewallRuleExtended `json:"userDefinedRules"`
}

// NsxtFirewallRuleExtended extends govcdtypes.NsxtFirewallRule with NetworkContextProfiles.
type NsxtFirewallRuleExtended struct {
	ID                        string                           `json:"id,omitempty"`
	Name                      string                           `json:"name"`
	ActionValue               string                           `json:"actionValue,omitempty"`
	Enabled                   bool                             `json:"enabled"`
	IPProtocol                string                           `json:"ipProtocol"`
	Logging                   bool                             `json:"logging"`
	Direction                 string                           `json:"direction"`
	SourceFirewallGroups      []govcdtypes.OpenApiReference    `json:"sourceFirewallGroups,omitempty"`
	DestinationFirewallGroups []govcdtypes.OpenApiReference    `json:"destinationFirewallGroups,omitempty"`
	ApplicationPortProfiles   []govcdtypes.OpenApiReference    `json:"applicationPortProfiles,omitempty"`
	NetworkContextProfiles    []govcdtypes.OpenApiReference    `json:"networkContextProfiles,omitempty"`
	Version                   *NsxtFirewallRuleExtendedVersion `json:"version,omitempty"`
}

// NsxtFirewallRuleExtendedVersion holds the version of a firewall rule.
type NsxtFirewallRuleExtendedVersion struct {
	Version *int `json:"version,omitempty"`
}

// nsxtFirewallRulesMinAPIVersion is the minimum API version for the Edge Gateway firewall rules endpoint.
// See: govcd/openapi_endpoints.go - OpenApiEndpointNsxtFirewallRules requires API version 34.0.
const nsxtFirewallRulesMinAPIVersion = "34.0"

// networkContextProfileAPIPayload is the JSON payload sent to the VCD API on
// POST and PUT. It mirrors NsxtNetworkContextProfile but adds contextEntityId
// which is required on create and omitted in read responses.
type networkContextProfileAPIPayload struct {
	ID              string                                     `json:"id,omitempty"`
	Name            string                                     `json:"name"`
	Description     string                                     `json:"description,omitempty"`
	Scope           string                                     `json:"scope"`
	ContextEntityID string                                     `json:"contextEntityId,omitempty"`
	OrgRef          *networkContextProfileOrgRef               `json:"orgRef,omitempty"`
	Attributes      []networkContextProfileAttributeAPIPayload `json:"attributes"`
}

type networkContextProfileOrgRef struct {
	ID string `json:"id"`
}

type networkContextProfileAttributeAPIPayload struct {
	Type          string   `json:"type"`
	Values        []string `json:"values"`
	SubAttributes any      `json:"subAttributes"`
}

// getGovcdClient returns the underlying govcd.Client from the cloudavenue client pool.
func getGovcdClient() (*govcd.Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	return &c.Vmware.Client, nil
}

// resolveFirewallAPIVersion replicates the logic of govcd's unexported
// checkOpenApiEndpointCompatibility using only exported Client methods:
//
//   - If the VCD instance maximum version is below the endpoint minimum → error
//   - If the client's configured APIVersion is higher than the minimum → use client's version
//   - Otherwise → fall back to the endpoint minimum
//
// This ensures we never send a version the server doesn't support, while still
// honouring the minimum contract for the endpoint.
func resolveFirewallAPIVersion(client *govcd.Client) (string, error) {
	// Check the VCD instance can support this endpoint at all.
	if client.APIVCDMaxVersionIs("< " + nsxtFirewallRulesMinAPIVersion) {
		maxVer, err := client.MaxSupportedVersion()
		if err != nil {
			maxVer = "unknown"
		}
		return "", fmt.Errorf(
			"edge gateway firewall endpoint requires API >= %s; this VCD instance supports up to %s",
			nsxtFirewallRulesMinAPIVersion, maxVer,
		)
	}

	// Use the client's configured version only when it is strictly above the minimum
	// AND does not exceed the VCD instance's maximum supported version.
	// This prevents sending a version the server does not understand.
	if client.APIClientVersionIs("> "+nsxtFirewallRulesMinAPIVersion) &&
		!client.APIVCDMaxVersionIs("< "+client.APIVersion) {
		return client.APIVersion, nil
	}

	return nsxtFirewallRulesMinAPIVersion, nil
}

// GetFirewallExtended retrieves the Edge Gateway firewall rules using the extended struct
// that includes NetworkContextProfiles, bypassing the govcd SDK limitation.
func (e *EdgeClient) GetFirewallExtended() (*NsxtFirewallRuleContainerExtended, error) {
	vcdEdge, err := e.GetVmwareEdgeGateway()
	if err != nil {
		return nil, fmt.Errorf("error retrieving VMware Edge Gateway: %w", err)
	}

	client, err := getGovcdClient()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(govcdtypes.OpenApiPathVersion1_0_0+govcdtypes.OpenApiEndpointNsxtFirewallRules, vcdEdge.EdgeGateway.ID)

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error building firewall endpoint URL: %w", err)
	}

	apiVersion, err := resolveFirewallAPIVersion(client)
	if err != nil {
		return nil, err
	}

	result := &NsxtFirewallRuleContainerExtended{}
	if err := client.OpenApiGetItem(apiVersion, urlRef, nil, result, nil); err != nil {
		return nil, fmt.Errorf("error retrieving NSX-T Firewall with context profiles: %w", err)
	}

	return result, nil
}

// UpdateFirewallExtended updates the Edge Gateway firewall rules using the extended struct
// that includes NetworkContextProfiles, bypassing the govcd SDK limitation.
func (e *EdgeClient) UpdateFirewallExtended(container *NsxtFirewallRuleContainerExtended) (*NsxtFirewallRuleContainerExtended, error) {
	vcdEdge, err := e.GetVmwareEdgeGateway()
	if err != nil {
		return nil, fmt.Errorf("error retrieving VMware Edge Gateway: %w", err)
	}

	client, err := getGovcdClient()
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(govcdtypes.OpenApiPathVersion1_0_0+govcdtypes.OpenApiEndpointNsxtFirewallRules, vcdEdge.EdgeGateway.ID)

	urlRef, err := client.OpenApiBuildEndpoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error building firewall endpoint URL: %w", err)
	}

	apiVersion, err := resolveFirewallAPIVersion(client)
	if err != nil {
		return nil, err
	}

	result := &NsxtFirewallRuleContainerExtended{}
	if err := client.OpenApiPutItem(apiVersion, urlRef, nil, container, result, nil); err != nil {
		return nil, fmt.Errorf("error updating NSX-T Firewall with context profiles: %w", err)
	}

	return result, nil
}

// networkContextProfilesEndpoint returns the base API endpoint for network context profiles.
func networkContextProfilesEndpoint() string {
	return govcdtypes.OpenApiPathVersion1_0_0 + "networkContextProfiles"
}

// resolveNetworkContextProfilesAPIVersion resolves the best API version for
// networkContextProfiles using the same strategy as resolveFirewallAPIVersion.
const nsxtNetworkContextProfilesMinAPIVersion = "35.0"

func resolveNetworkContextProfilesAPIVersion(client *govcd.Client) (string, error) {
	if client.APIVCDMaxVersionIs("< " + nsxtNetworkContextProfilesMinAPIVersion) {
		maxVer, err := client.MaxSupportedVersion()
		if err != nil {
			maxVer = "unknown"
		}
		return "", fmt.Errorf(
			"networkContextProfiles endpoint requires API >= %s; this VCD instance supports up to %s",
			nsxtNetworkContextProfilesMinAPIVersion, maxVer,
		)
	}

	// Use the client's configured version only when it is strictly above the minimum
	// AND does not exceed the VCD instance's maximum supported version.
	if client.APIClientVersionIs("> "+nsxtNetworkContextProfilesMinAPIVersion) &&
		!client.APIVCDMaxVersionIs("< "+client.APIVersion) {
		return client.APIVersion, nil
	}

	return nsxtNetworkContextProfilesMinAPIVersion, nil
}

// CreateNetworkContextProfile creates a new TENANT-scoped Network Context Profile.
// The VCD API responds with 202 Accepted + a task; we wait for the task then
// fetch the profile by name to obtain its assigned ID.
func (e *EdgeClient) CreateNetworkContextProfile(profile *NetworkContextProfile) (*NetworkContextProfile, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile must not be nil")
	}
	if profile.Name == "" {
		return nil, fmt.Errorf("profile.Name must not be empty")
	}

	vcdEdge, err := e.GetVmwareEdgeGateway()
	if err != nil {
		return nil, fmt.Errorf("error retrieving VMware Edge Gateway: %w", err)
	}

	cavc, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	client := &cavc.Vmware.Client

	apiVersion, err := resolveNetworkContextProfilesAPIVersion(client)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(networkContextProfilesEndpoint())
	if err != nil {
		return nil, fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	payload := networkContextProfileToAPIPayload(profile, vcdEdge.EdgeGateway.OwnerRef.ID, cavc.Org.Org.ID)

	// POST returns 202 Accepted with a task — use async variant and wait.
	task, err := client.OpenApiPostItemAsync(apiVersion, urlRef, nil, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating Network Context Profile: %w", err)
	}

	if err := task.WaitTaskCompletion(); err != nil {
		return nil, fmt.Errorf("error waiting for Network Context Profile creation task: %w", err)
	}

	// Fetch the created profile by name (POST doesn't return the ID).
	return e.GetNetworkContextProfileByName(profile.Name)
}

// GetNetworkContextProfileByID retrieves a Network Context Profile by its URN ID.
func (e *EdgeClient) GetNetworkContextProfileByID(id string) (*NetworkContextProfile, error) {
	if id == "" {
		return nil, fmt.Errorf("id must not be empty")
	}

	c, err := getGovcdClient()
	if err != nil {
		return nil, err
	}

	apiVersion, err := resolveNetworkContextProfilesAPIVersion(c)
	if err != nil {
		return nil, err
	}

	urlRef, err := c.OpenApiBuildEndpoint(networkContextProfilesEndpoint() + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	result := &networkContextProfileAPIPayload{}
	if err := c.OpenApiGetItem(apiVersion, urlRef, nil, result, nil); err != nil {
		return nil, fmt.Errorf("error retrieving Network Context Profile %q: %w", id, err)
	}

	return networkContextProfileFromAPIPayload(result), nil
}

// UpdateNetworkContextProfile updates an existing TENANT-scoped Network Context Profile.
func (e *EdgeClient) UpdateNetworkContextProfile(profile *NetworkContextProfile) (*NetworkContextProfile, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile must not be nil")
	}
	if profile.ID == "" {
		return nil, fmt.Errorf("profile.ID must not be empty")
	}

	vcdEdge, err := e.GetVmwareEdgeGateway()
	if err != nil {
		return nil, fmt.Errorf("error retrieving VMware Edge Gateway: %w", err)
	}

	cavc, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	client := &cavc.Vmware.Client

	apiVersion, err := resolveNetworkContextProfilesAPIVersion(client)
	if err != nil {
		return nil, err
	}

	urlRef, err := client.OpenApiBuildEndpoint(networkContextProfilesEndpoint() + "/" + profile.ID)
	if err != nil {
		return nil, fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	payload := networkContextProfileToAPIPayload(profile, vcdEdge.EdgeGateway.OwnerRef.ID, cavc.Org.Org.ID)

	// PUT also returns 202 Accepted — use async variant and wait.
	task, err := client.OpenApiPutItemAsync(apiVersion, urlRef, nil, payload, nil)
	if err != nil {
		return nil, fmt.Errorf("error updating Network Context Profile %q: %w", profile.ID, err)
	}

	if err := task.WaitTaskCompletion(); err != nil {
		return nil, fmt.Errorf("error waiting for Network Context Profile update task: %w", err)
	}

	// Re-read by ID to get updated state.
	return e.GetNetworkContextProfileByID(profile.ID)
}

// DeleteNetworkContextProfile deletes a TENANT-scoped Network Context Profile by ID.
func (e *EdgeClient) DeleteNetworkContextProfile(id string) error {
	if id == "" {
		return fmt.Errorf("id must not be empty")
	}

	c, err := getGovcdClient()
	if err != nil {
		return err
	}

	apiVersion, err := resolveNetworkContextProfilesAPIVersion(c)
	if err != nil {
		return err
	}

	urlRef, err := c.OpenApiBuildEndpoint(networkContextProfilesEndpoint() + "/" + id)
	if err != nil {
		return fmt.Errorf("error building networkContextProfiles endpoint: %w", err)
	}

	if err := c.OpenApiDeleteItem(apiVersion, urlRef, nil, nil); err != nil {
		return fmt.Errorf("error deleting Network Context Profile %q: %w", id, err)
	}

	return nil
}

// networkContextProfileToAPIPayload converts an SDK model to the API JSON payload.
func networkContextProfileToAPIPayload(p *NetworkContextProfile, ownerVDCID, orgID string) *networkContextProfileAPIPayload {
	attrs := make([]networkContextProfileAttributeAPIPayload, len(p.Attributes))
	for i, a := range p.Attributes {
		attrs[i] = networkContextProfileAttributeAPIPayload{
			Type:          string(a.Type),
			Values:        a.Values,
			SubAttributes: nil,
		}
	}

	payload := &networkContextProfileAPIPayload{
		ID:              p.ID,
		Name:            p.Name,
		Description:     p.Description,
		Scope:           string(NetworkContextProfileScopeTenant),
		ContextEntityID: ownerVDCID,
		Attributes:      attrs,
	}

	if orgID != "" {
		payload.OrgRef = &networkContextProfileOrgRef{ID: orgID}
	}

	return payload
}

// networkContextProfileFromAPIPayload converts an API JSON payload to the SDK model.
func networkContextProfileFromAPIPayload(p *networkContextProfileAPIPayload) *NetworkContextProfile {
	attrs := make([]NetworkContextProfileAttribute, len(p.Attributes))
	for i, a := range p.Attributes {
		attrs[i] = NetworkContextProfileAttribute{
			Type:   NetworkContextProfileAttributeType(a.Type),
			Values: a.Values,
		}
	}

	profile := &NetworkContextProfile{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Scope:       NetworkContextProfileScope(p.Scope),
		Attributes:  attrs,
	}
	if p.OrgRef != nil {
		profile.OrgID = p.OrgRef.ID
	}
	return profile
}

// GetAllNetworkContextProfiles returns all Network Context Profiles available
// in the context of the Edge Gateway (SYSTEM + PROVIDER + TENANT scopes).
func (e *EdgeClient) GetAllNetworkContextProfiles() ([]*NetworkContextProfile, error) {
	vcdEdge, err := e.GetVmwareEdgeGateway()
	if err != nil {
		return nil, fmt.Errorf("error retrieving VMware Edge Gateway: %w", err)
	}

	client, err := getGovcdClient()
	if err != nil {
		return nil, err
	}

	// Network Context Profiles are scoped to VDC/VDC Group, not to the edge gateway directly.
	// The _context filter must be the owning VDC or VDC Group ID.
	queryParams := url.Values{}
	queryParams.Set("filter", fmt.Sprintf("_context==%s", vcdEdge.EdgeGateway.OwnerRef.ID))

	raw, err := govcd.GetAllNetworkContextProfiles(client, queryParams)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Network Context Profiles: %w", err)
	}

	profiles := make([]*NetworkContextProfile, len(raw))
	for i, r := range raw {
		profiles[i] = networkContextProfileFromGovcd(r)
	}

	return profiles, nil
}

// GetNetworkContextProfileByName returns a single Network Context Profile by name.
func (e *EdgeClient) GetNetworkContextProfileByName(name string) (*NetworkContextProfile, error) {
	all, err := e.GetAllNetworkContextProfiles()
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

// networkContextProfileFromGovcd converts a govcd NsxtNetworkContextProfile to the SDK model.
func networkContextProfileFromGovcd(raw *govcdtypes.NsxtNetworkContextProfile) *NetworkContextProfile {
	p := &NetworkContextProfile{
		ID:          raw.ID,
		Name:        raw.Name,
		Description: raw.Description,
		Scope:       NetworkContextProfileScope(raw.Scope),
		Attributes:  make([]NetworkContextProfileAttribute, len(raw.Attributes)),
	}

	if raw.OrgRef != nil {
		p.OrgID = raw.OrgRef.ID
	}

	for i, a := range raw.Attributes {
		p.Attributes[i] = NetworkContextProfileAttribute{
			Type:   NetworkContextProfileAttributeType(a.Type),
			Values: a.Values,
		}
	}

	return p
}
