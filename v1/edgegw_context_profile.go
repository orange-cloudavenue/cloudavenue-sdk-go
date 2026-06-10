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
	IpProtocol                string                           `json:"ipProtocol"`
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

// nsxtFirewallRulesMinAPIVersion is the minimum VCD API version required for Edge Gateway firewall rules.
// See: govcd/openapi_endpoints.go - OpenApiEndpointNsxtFirewallRules requires API version 34.0.
const nsxtFirewallRulesMinAPIVersion = "34.0"

// getGovcdClient returns the underlying govcd.Client from the cloudavenue client pool.
func getGovcdClient() (*govcd.Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, fmt.Errorf("error initialising CloudAvenue client: %w", err)
	}
	return &c.Vmware.Client, nil
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

	result := &NsxtFirewallRuleContainerExtended{}
	if err := client.OpenApiGetItem(nsxtFirewallRulesMinAPIVersion, urlRef, nil, result, nil); err != nil {
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

	result := &NsxtFirewallRuleContainerExtended{}
	if err := client.OpenApiPutItem(nsxtFirewallRulesMinAPIVersion, urlRef, nil, container, result, nil); err != nil {
		return nil, fmt.Errorf("error updating NSX-T Firewall with context profiles: %w", err)
	}

	return result, nil
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

	queryParams := url.Values{}
	queryParams.Set("filter", fmt.Sprintf("_context==%s", vcdEdge.EdgeGateway.ID))

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
		return nil, fmt.Errorf("%s: network context profile with name %q not found", govcd.ErrorEntityNotFound, name)
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
		Scope:       raw.Scope,
		Attributes:  make([]NetworkContextProfileAttribute, len(raw.Attributes)),
	}

	for i, a := range raw.Attributes {
		p.Attributes[i] = NetworkContextProfileAttribute{
			Type:   a.Type,
			Values: a.Values,
		}
	}

	return p
}
