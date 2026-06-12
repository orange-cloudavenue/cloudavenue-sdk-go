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

	// NetworkContextProfileAttributeTypeDomainName matches traffic by FQDN or wildcard domain name.
	// Values are wildcard-capable domain strings, e.g. "*.example.com" or "myhost.corp.local".
	NetworkContextProfileAttributeTypeDomainName NetworkContextProfileAttributeType = "DOMAIN_NAME"
)

// NetworkContextProfileSubAttributeType represents the type of a sub-attribute.
type NetworkContextProfileSubAttributeType string

const (
	NetworkContextProfileSubAttributeTypeTLSVersion     NetworkContextProfileSubAttributeType = "TLS_VERSION"
	NetworkContextProfileSubAttributeTypeTLSCipherSuite NetworkContextProfileSubAttributeType = "TLS_CIPHER_SUITE"
	NetworkContextProfileSubAttributeTypeCIFSSMBVersion NetworkContextProfileSubAttributeType = "CIFS_SMB_VERSION"
)

// NetworkContextProfileKnownSubAttributeTypes lists all known sub-attribute types.
var NetworkContextProfileKnownSubAttributeTypes = []string{
	string(NetworkContextProfileSubAttributeTypeTLSVersion),
	string(NetworkContextProfileSubAttributeTypeTLSCipherSuite),
	string(NetworkContextProfileSubAttributeTypeCIFSSMBVersion),
}

// NetworkContextProfileKnownTLSVersions lists valid TLS_VERSION sub-attribute values.
var NetworkContextProfileKnownTLSVersions = []string{
	"TLS_V10", "TLS_V11", "TLS_V12", "TLS_V13",
}

// NetworkContextProfileKnownTLSCipherSuites lists valid TLS_CIPHER_SUITE sub-attribute values.
var NetworkContextProfileKnownTLSCipherSuites = []string{
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA",
	"TLS_DHE_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_DHE_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA",
	"TLS_DHE_RSA_WITH_AES_256_CBC_SHA256",
	"TLS_DHE_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA",
	"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384",
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
	"TLS_RSA_WITH_3DES_EDE_CBC_SHA",
	"TLS_RSA_WITH_AES_128_CBC_SHA",
	"TLS_RSA_WITH_AES_128_CBC_SHA256",
	"TLS_RSA_WITH_AES_128_GCM_SHA256",
	"TLS_RSA_WITH_AES_256_CBC_SHA",
	"TLS_RSA_WITH_AES_256_CBC_SHA256",
	"TLS_RSA_WITH_AES_256_GCM_SHA384",
}

// NetworkContextProfileKnownCIFSSMBVersions lists valid CIFS_SMB_VERSION sub-attribute values.
var NetworkContextProfileKnownCIFSSMBVersions = []string{
	"CIFS_SMB_V1", "CIFS_SMB_V2", "CIFS_SMB_V3",
}

// NetworkContextProfileKnownSubAttributeValues is the union of all valid values across
// all sub-attribute types. Use this for schema validation when the sub-attribute type
// is not known at validation time.
var NetworkContextProfileKnownSubAttributeValues = func() []string {
	all := make([]string, 0,
		len(NetworkContextProfileKnownTLSVersions)+
			len(NetworkContextProfileKnownTLSCipherSuites)+
			len(NetworkContextProfileKnownCIFSSMBVersions))
	all = append(all, NetworkContextProfileKnownTLSVersions...)
	all = append(all, NetworkContextProfileKnownTLSCipherSuites...)
	all = append(all, NetworkContextProfileKnownCIFSSMBVersions...)
	return all
}()

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

	// OrgID is the URN of the owning organisation, populated from the API response
	// for TENANT profiles.
	OrgID string

	// Attributes describes the Layer 7 characteristics of the profile.
	Attributes []NetworkContextProfileAttribute
}

// NetworkContextProfileAttribute is a single attribute of a Network Context Profile.
// The Type field determines whether this is an APP_ID or DOMAIN_NAME attribute.
type NetworkContextProfileAttribute struct {
	// Type is the attribute type — always APP_ID for user-managed profiles.
	Type NetworkContextProfileAttributeType

	// Values is the list of app identifiers for this attribute (e.g. ["SSL"]).
	Values []string

	// SubAttributes provides optional refinements (e.g. TLS version, cipher suites).
	SubAttributes []NetworkContextProfileSubAttribute
}

// NetworkContextProfileSubAttribute is a typed sub-attribute within an APP_ID attribute.
type NetworkContextProfileSubAttribute struct {
	// Type identifies the sub-attribute (e.g. TLS_VERSION, TLS_CIPHER_SUITE, CIFS_SMB_VERSION).
	Type NetworkContextProfileSubAttributeType

	// Values is the list of allowed values for this sub-attribute.
	Values []string
}

// NetworkContextProfileAttributesCatalog holds the server-side allowlist of valid
// attribute values returned by GET /cloudapi/1.0.0/networkContextProfiles/attributes.
type NetworkContextProfileAttributesCatalog struct {
	// DomainNameValues is the list of valid DOMAIN_NAME values on this platform.
	// Only values in this list can be used in a domain_name attribute block.
	DomainNameValues []string

	// AppIDValues is the list of valid APP_ID values on this platform.
	AppIDValues []string
}

// networkContextProfileAttributesAPIResponse is the raw API response shape.
type networkContextProfileAttributesAPIResponse struct {
	Attributes []struct {
		Type          string   `json:"type"`
		Values        []string `json:"values"`
		SubAttributes []struct {
			Type   string   `json:"type"`
			Values []string `json:"values"`
		} `json:"subAttributes"`
	} `json:"attributes"`
}
