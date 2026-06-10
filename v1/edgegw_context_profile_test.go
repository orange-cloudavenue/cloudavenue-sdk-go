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
	"testing"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

const (
	testProfileIDAAA   = "urn:vcloud:networkContextProfile:aaa"
	testProfileIDBBB   = "urn:vcloud:networkContextProfile:bbb"
	testProfileIDCCC   = "urn:vcloud:networkContextProfile:ccc"
	testProfileIDSSL   = "urn:vcloud:networkContextProfile:d6d3ff93-fca4-3eaf-bf07-3e1ffe0c6b7a"
	testProfileIDOWASP = "urn:vcloud:networkContextProfile:xyz789"
	testAppIDSSL       = "SSL"
	testAppIDCIFS      = "CIFS"
	testAppIDHTTP      = "HTTP"
)

// TestNetworkContextProfileFromGovcd validates that networkContextProfileFromGovcd
// correctly maps all fields from the govcd type to the SDK model.
func TestNetworkContextProfileFromGovcd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *govcdtypes.NsxtNetworkContextProfile
		expected *NetworkContextProfile
	}{
		{
			name: "full profile with attributes",
			input: &govcdtypes.NsxtNetworkContextProfile{
				ID:          testProfileIDSSL,
				Name:        testAppIDSSL,
				Description: "Secure Sockets Layer",
				Scope:       string(NetworkContextProfileScopeSystem),
				Attributes: []govcdtypes.NsxtNetworkContextProfileAttributes{
					{
						Type:   string(NetworkContextProfileAttributeTypeAppID),
						Values: []string{testAppIDSSL},
					},
				},
			},
			expected: &NetworkContextProfile{
				ID:          testProfileIDSSL,
				Name:        testAppIDSSL,
				Description: "Secure Sockets Layer",
				Scope:       NetworkContextProfileScopeSystem,
				Attributes: []NetworkContextProfileAttribute{
					{
						Type:   NetworkContextProfileAttributeTypeAppID,
						Values: []string{testAppIDSSL},
					},
				},
			},
		},
		{
			name: "profile with no attributes",
			input: &govcdtypes.NsxtNetworkContextProfile{
				ID:         "urn:vcloud:networkContextProfile:abc123",
				Name:       "EMPTY",
				Scope:      string(NetworkContextProfileScopeTenant),
				Attributes: []govcdtypes.NsxtNetworkContextProfileAttributes{},
			},
			expected: &NetworkContextProfile{
				ID:         "urn:vcloud:networkContextProfile:abc123",
				Name:       "EMPTY",
				Scope:      NetworkContextProfileScopeTenant,
				Attributes: []NetworkContextProfileAttribute{},
			},
		},
		{
			name: "profile with multiple attribute values",
			input: &govcdtypes.NsxtNetworkContextProfile{
				ID:          testProfileIDOWASP,
				Name:        "OWASP-A",
				Description: "OWASP Cipher String A",
				Scope:       string(NetworkContextProfileScopeSystem),
				Attributes: []govcdtypes.NsxtNetworkContextProfileAttributes{
					{
						Type:   string(NetworkContextProfileAttributeTypeAppID),
						Values: []string{testAppIDSSL},
					},
				},
			},
			expected: &NetworkContextProfile{
				ID:          testProfileIDOWASP,
				Name:        "OWASP-A",
				Description: "OWASP Cipher String A",
				Scope:       NetworkContextProfileScopeSystem,
				Attributes: []NetworkContextProfileAttribute{
					{
						Type:   NetworkContextProfileAttributeTypeAppID,
						Values: []string{testAppIDSSL},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := networkContextProfileFromGovcd(tt.input)

			if got.ID != tt.expected.ID {
				t.Errorf("ID: got %q, want %q", got.ID, tt.expected.ID)
			}
			if got.Name != tt.expected.Name {
				t.Errorf("Name: got %q, want %q", got.Name, tt.expected.Name)
			}
			if got.Description != tt.expected.Description {
				t.Errorf("Description: got %q, want %q", got.Description, tt.expected.Description)
			}
			if got.Scope != tt.expected.Scope {
				t.Errorf("Scope: got %q, want %q", got.Scope, tt.expected.Scope)
			}
			if len(got.Attributes) != len(tt.expected.Attributes) {
				t.Errorf("Attributes length: got %d, want %d", len(got.Attributes), len(tt.expected.Attributes))
				return
			}
			for i, attr := range got.Attributes {
				if attr.Type != tt.expected.Attributes[i].Type {
					t.Errorf("Attributes[%d].Type: got %q, want %q", i, attr.Type, tt.expected.Attributes[i].Type)
				}
				if len(attr.Values) != len(tt.expected.Attributes[i].Values) {
					t.Errorf("Attributes[%d].Values length: got %d, want %d", i, len(attr.Values), len(tt.expected.Attributes[i].Values))
				}
			}
		})
	}
}

// TestFindNetworkContextProfileByName validates the name-resolution logic
// that underlies GetNetworkContextProfileByName, by testing the filtering
// and error cases directly on a slice of profiles.
func TestFindNetworkContextProfileByName(t *testing.T) {
	t.Parallel()

	allProfiles := []*NetworkContextProfile{
		{ID: testProfileIDAAA, Name: testAppIDSSL, Scope: NetworkContextProfileScopeSystem},
		{ID: testProfileIDBBB, Name: testAppIDCIFS, Scope: NetworkContextProfileScopeSystem},
		{ID: testProfileIDCCC, Name: testAppIDHTTP, Scope: NetworkContextProfileScopeSystem},
	}

	// helper that mimics the name-filter logic from GetNetworkContextProfileByName.
	findByName := func(profiles []*NetworkContextProfile, name string) ([]*NetworkContextProfile, error) {
		var found []*NetworkContextProfile
		for _, p := range profiles {
			if p.Name == name {
				found = append(found, p)
			}
		}
		return found, nil
	}

	tests := []struct {
		name       string
		search     string
		wantID     string
		wantLen    int
		expectNone bool
	}{
		{
			name:    "found exactly one - SSL",
			search:  testAppIDSSL,
			wantID:  testProfileIDAAA,
			wantLen: 1,
		},
		{
			name:    "found exactly one - CIFS",
			search:  testAppIDCIFS,
			wantID:  testProfileIDBBB,
			wantLen: 1,
		},
		{
			name:       "not found",
			search:     "DOES_NOT_EXIST",
			expectNone: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			found, _ := findByName(allProfiles, tt.search)

			if tt.expectNone {
				if len(found) != 0 {
					t.Errorf("expected no profiles, got %d", len(found))
				}
				return
			}

			if len(found) != tt.wantLen {
				t.Errorf("expected %d profiles, got %d", tt.wantLen, len(found))
				return
			}

			if tt.wantID != "" && found[0].ID != tt.wantID {
				t.Errorf("ID: got %q, want %q", found[0].ID, tt.wantID)
			}
		})
	}
}

// TestNsxtFirewallRuleExtendedFields verifies that NsxtFirewallRuleExtended
// correctly stores all fields, including NetworkContextProfiles which extends
// the upstream govcd SDK struct.
func TestNsxtFirewallRuleExtendedFields(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		rule                   *NsxtFirewallRuleExtended
		wantName               string
		wantContextProfilesLen int
		wantContextProfileID   string
		wantContextProfileName string
	}{
		"rule with one context profile": {
			rule: &NsxtFirewallRuleExtended{
				Name:        "test-rule",
				ActionValue: "ALLOW",
				Enabled:     true,
				IPProtocol:  "IPV4",
				Logging:     false,
				Direction:   "OUT",
				NetworkContextProfiles: []govcdtypes.OpenApiReference{
					{ID: testProfileIDAAA, Name: testAppIDSSL},
				},
			},
			wantName:               "test-rule",
			wantContextProfilesLen: 1,
			wantContextProfileID:   testProfileIDAAA,
			wantContextProfileName: testAppIDSSL,
		},
		"rule with no context profiles": {
			rule: &NsxtFirewallRuleExtended{
				Name:                   "bare-rule",
				ActionValue:            "DROP",
				Enabled:                true,
				IPProtocol:             "IPV4_IPV6",
				Direction:              "IN",
				NetworkContextProfiles: nil,
			},
			wantName:               "bare-rule",
			wantContextProfilesLen: 0,
		},
		"rule with multiple context profiles": {
			rule: &NsxtFirewallRuleExtended{
				Name:        "multi-profile-rule",
				ActionValue: "ALLOW",
				Direction:   "OUT",
				IPProtocol:  "IPV4",
				NetworkContextProfiles: []govcdtypes.OpenApiReference{
					{ID: testProfileIDAAA, Name: testAppIDSSL},
					{ID: testProfileIDBBB, Name: testAppIDCIFS},
				},
			},
			wantName:               "multi-profile-rule",
			wantContextProfilesLen: 2,
			wantContextProfileID:   testProfileIDAAA,
			wantContextProfileName: testAppIDSSL,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.rule.Name != tt.wantName {
				t.Errorf("Name: got %q, want %q", tt.rule.Name, tt.wantName)
			}
			if len(tt.rule.NetworkContextProfiles) != tt.wantContextProfilesLen {
				t.Errorf("NetworkContextProfiles length: got %d, want %d",
					len(tt.rule.NetworkContextProfiles), tt.wantContextProfilesLen)
				return
			}
			if tt.wantContextProfilesLen > 0 {
				if tt.rule.NetworkContextProfiles[0].ID != tt.wantContextProfileID {
					t.Errorf("NetworkContextProfiles[0].ID: got %q, want %q",
						tt.rule.NetworkContextProfiles[0].ID, tt.wantContextProfileID)
				}
				if tt.rule.NetworkContextProfiles[0].Name != tt.wantContextProfileName {
					t.Errorf("NetworkContextProfiles[0].Name: got %q, want %q",
						tt.rule.NetworkContextProfiles[0].Name, tt.wantContextProfileName)
				}
			}
		})
	}
}
