/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

// NetworkContextProfileScope represents the scope of a Network Context Profile.
type NetworkContextProfileScope string

const (
	NetworkContextProfileScopeSystem   NetworkContextProfileScope = "SYSTEM"
	NetworkContextProfileScopeProvider NetworkContextProfileScope = "PROVIDER"
	NetworkContextProfileScopeTenant   NetworkContextProfileScope = "TENANT"
)

// NetworkContextProfileAttributeType represents the type of a profile attribute.
type NetworkContextProfileAttributeType string

const (
	// NetworkContextProfileAttributeTypeAppID identifies Layer 7 applications.
	NetworkContextProfileAttributeTypeAppID NetworkContextProfileAttributeType = "APP_ID"
)

// NetworkContextProfileKnownAppIDs lists all APP_ID values known to be available
// on Cloud Avenue (sourced from SYSTEM profiles via the VCD API).
var NetworkContextProfileKnownAppIDs = []string{
	"360ANTIV", "ACTIVDIR", "AMQP", "AVAST", "AVG", "AVIRA", "BDEFNDER",
	"BLAST", "CA_CERT", "CIFS", "CLDAP", "CTRXCGP", "CTRXGOTO", "CTRXICA",
	"DCERPC", "DHCP", "DIAMETER", "DNS", "EPIC", "ESET", "FPROT", "FTP",
	"GITHUB", "HTTP", "HTTP2", "IMAP", "KASPRSKY", "KERBEROS", "LDAP",
	"MAXDB", "MCAFEE", "MSSQL", "MYSQL", "NFS", "NNTP", "NTBIOSNS", "NTP",
	"OCSP", "ORACLE", "PANDA", "PCOIP", "POP3", "RADIUS", "RDP", "RTCP",
	"RTP", "RTSP", "SIP", "SMTP", "SNMP", "SSH", "SSL", "SYMUPDAT",
	"SYSLOG", "TELNET", "TFTP", "VNC", "WINS",
}

// NetworkContextProfile represents a Network Context Profile (Layer 7 context profile).
// SYSTEM and PROVIDER profiles are read-only; TENANT profiles can be created/updated/deleted.
type NetworkContextProfile struct {
	// ID is the URN of the profile (e.g. urn:vcloud:networkContextProfile:...).
	ID string

	// Name is the human-readable name (e.g. "SSL", "CIFS", "my-custom-profile").
	Name string

	// Description provides a human-readable description of the profile.
	Description string

	// Scope is one of SYSTEM, PROVIDER, TENANT.
	Scope NetworkContextProfileScope

	// OrgID is the URN of the owning organisation (required for TENANT profiles).
	OrgID string

	// OwnerVDCID is the URN of the owning VDC (used as contextEntityId on create).
	OwnerVDCID string

	// Attributes describes the Layer 7 characteristics of the profile.
	Attributes []NetworkContextProfileAttribute
}

// NetworkContextProfileAttribute is a single attribute of a Network Context Profile.
type NetworkContextProfileAttribute struct {
	// Type is the attribute type, e.g. "APP_ID".
	Type NetworkContextProfileAttributeType

	// Values is the list of values for this attribute (e.g. ["SSL"]).
	Values []string
}
