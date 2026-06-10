/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

// NetworkContextProfile represents a Network Context Profile (Layer 7 context profile)
// available on an Edge Gateway. Profiles are always SYSTEM-scoped and read-only for tenants.
type NetworkContextProfile struct {
	// ID is the URN of the profile (e.g. urn:vcloud:networkContextProfile:...).
	ID string

	// Name is the human-readable name (e.g. "SSL", "CIFS", "HTTP").
	Name string

	// Description provides a human-readable description of the profile.
	Description string

	// Scope is always "SYSTEM" for predefined VMware profiles.
	Scope string

	// Attributes describes the Layer 7 characteristics of the profile.
	Attributes []NetworkContextProfileAttribute
}

// NetworkContextProfileAttribute is a single attribute of a Network Context Profile.
type NetworkContextProfileAttribute struct {
	// Type is the attribute type, e.g. "APP_ID" or "DOMAIN_NAME".
	Type string

	// Values is the list of values for this attribute (e.g. ["SSL"], ["office365.com"]).
	Values []string
}
