/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package v1 provides CloudAvenue SDK client implementations
package v1

// NetworkContextProfile types define Layer 7 traffic classification and control for VMware NSX Edge Gateways.
// These types enable fine-grained application-level firewall policies by classifying traffic based on:
//   - Application protocols and signatures (APP_ID attributes)
//   - Domain/FQDN patterns (DOMAIN_NAME attributes)
//   - Security attributes like TLS versions, cipher suites, and SMB protocol versions
//
// Network Context Profiles are used in NSX Distributed Firewall and Edge Gateway firewall rules to enforce
// security policies beyond traditional port-based rules, supporting use cases such as:
//   - Application-aware access control
//   - Protocol version compliance (e.g., blocking legacy TLS 1.0)
//   - Ransomware protection (isolating SMB v1)
//   - Zero Trust Network Access policies
//
// For more information, see:
//   - VMware NSX Administration Guide: Network Context Profiles
//   - VMware NSX Security: Firewall Rules and Application Identification
//   - CloudAvenue Documentation: Edge Gateway Security Policies
//   - Broadcom NSX Documentation: https://techdocs.broadcom.com (search for "Network Context Profile")

// ============================================================================
// TYPE DEFINITIONS & CONSTANTS
// ============================================================================

// NetworkContextProfileScope represents the scope of a Network Context Profile.
// SYSTEM scopes apply cluster-wide, PROVIDER enables service provider settings, TENANT restricts to specific organizations.
type NetworkContextProfileScope string

const (
	NetworkContextProfileScopeSystem   NetworkContextProfileScope = "SYSTEM"
	NetworkContextProfileScopeProvider NetworkContextProfileScope = "PROVIDER"
	NetworkContextProfileScopeTenant   NetworkContextProfileScope = "TENANT"
)

// NetworkContextProfileAttributeType represents the type of a profile attribute.
// These are used to classify and control network traffic at Layer 7 (application layer) in NSX edge gateways.
type NetworkContextProfileAttributeType string

const (
	// NetworkContextProfileAttributeTypeAppID identifies Layer 7 applications using application signatures and protocols.
	// Enables fine-grained traffic filtering by specific application types rather than just ports.
	NetworkContextProfileAttributeTypeAppID NetworkContextProfileAttributeType = "APP_ID"

	// NetworkContextProfileAttributeTypeDomainName matches traffic by FQDN or wildcard domain name patterns.
	// Supports wildcard patterns like "*.example.com" or "myhost.corp.local" for hostname-based classification.
	NetworkContextProfileAttributeTypeDomainName NetworkContextProfileAttributeType = "DOMAIN_NAME"
)

// NetworkContextProfileSubAttributeType represents sub-attributes of network profiles.
// These enable granular control over Layer 7 attributes like TLS versions and protocol specifics.
type NetworkContextProfileSubAttributeType string

const (
	NetworkContextProfileSubAttributeTypeTLSVersion     NetworkContextProfileSubAttributeType = "TLS_VERSION"
	NetworkContextProfileSubAttributeTypeTLSCipherSuite NetworkContextProfileSubAttributeType = "TLS_CIPHER_SUITE"
	NetworkContextProfileSubAttributeTypeCIFSSMBVersion NetworkContextProfileSubAttributeType = "CIFS_SMB_VERSION"
)

// ============================================================================
// ATTRIBUTE DEFINITION TYPES
// ============================================================================

// NetworkContextProfileValueDefinition represents a single value in an attribute definition.
type NetworkContextProfileValueDefinition struct {
	Value       string
	Description string
}

// NetworkContextProfileAttributeDefinition describes a network context profile attribute or sub-attribute.
type NetworkContextProfileAttributeDefinition struct {
	Name        NetworkContextProfileAttributeType
	Description string
	Values      []NetworkContextProfileValueDefinition
}

// ValuesOnly returns a slice containing just the values (not descriptions).
func (a NetworkContextProfileAttributeDefinition) ValuesOnly() []string {
	values := make([]string, len(a.Values))
	for i, v := range a.Values {
		values[i] = v.Value
	}
	return values
}

// ============================================================================
// ATTRIBUTE DEFINITIONS (ORDERED BY DEPENDENCY)
// ============================================================================

// NetworkContextProfileTLSVersionDefinition defines valid TLS protocol versions.
var NetworkContextProfileTLSVersionDefinition = NetworkContextProfileAttributeDefinition{
	Name:        NetworkContextProfileAttributeTypeAppID,
	Description: "TLS protocol version classifications. Restricts encrypted traffic by protocol version to enforce minimum security standards (e.g., block TLS 1.0/1.1 for compliance).",
	Values: []NetworkContextProfileValueDefinition{
		{Value: "TLS_V10", Description: "TLS 1.0"},
		{Value: "TLS_V11", Description: "TLS 1.1"},
		{Value: "TLS_V12", Description: "TLS 1.2"},
		{Value: "TLS_V13", Description: "TLS 1.3"},
	},
}

// NetworkContextProfileTLSCipherSuiteDefinition defines valid TLS cipher suites.
var NetworkContextProfileTLSCipherSuiteDefinition = NetworkContextProfileAttributeDefinition{
	Name:        NetworkContextProfileAttributeTypeAppID,
	Description: "TLS cipher suite classifications. Enables security policies based on specific cipher algorithms (e.g., reject weak ciphers, enforce FIPS-140 compliance).",
	Values: []NetworkContextProfileValueDefinition{
		{Value: "TLS_DHE_RSA_WITH_AES_128_CBC_SHA"},
		{Value: "TLS_DHE_RSA_WITH_AES_128_CBC_SHA256"},
		{Value: "TLS_DHE_RSA_WITH_AES_128_GCM_SHA256"},
		{Value: "TLS_DHE_RSA_WITH_AES_256_CBC_SHA"},
		{Value: "TLS_DHE_RSA_WITH_AES_256_CBC_SHA256"},
		{Value: "TLS_DHE_RSA_WITH_AES_256_GCM_SHA384"},
		{Value: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"},
		{Value: "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256"},
		{Value: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
		{Value: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"},
		{Value: "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384"},
		{Value: "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
		{Value: "TLS_RSA_WITH_3DES_EDE_CBC_SHA"},
		{Value: "TLS_RSA_WITH_AES_128_CBC_SHA"},
		{Value: "TLS_RSA_WITH_AES_128_CBC_SHA256"},
		{Value: "TLS_RSA_WITH_AES_128_GCM_SHA256"},
		{Value: "TLS_RSA_WITH_AES_256_CBC_SHA"},
		{Value: "TLS_RSA_WITH_AES_256_CBC_SHA256"},
		{Value: "TLS_RSA_WITH_AES_256_GCM_SHA384"},
	},
}

// NetworkContextProfileCIFSSMBVersionDefinition defines valid CIFS/SMB protocol versions.
var NetworkContextProfileCIFSSMBVersionDefinition = NetworkContextProfileAttributeDefinition{
	Name:        NetworkContextProfileAttributeTypeAppID,
	Description: "CIFS/SMB protocol version classifications. Controls Windows file sharing traffic by SMB version for security (e.g., isolate legacy SMB v1 for ransomware protection).",
	Values: []NetworkContextProfileValueDefinition{
		{
			Value:       "SMB_V1",
			Description: "Legacy SMB 1.0 protocol. Deprecated and vulnerable to multiple known exploits.",
		},
		{
			Value:       "SMB_V2",
			Description: "Server Message Block version 2.0.",
		},
		{
			Value:       "SMB_V3",
			Description: "Server Message Block version 3.0.",
		},
	},
}

// NetworkContextProfileSubAttributeDefinitions maps sub-attribute type names to their definitions.
var NetworkContextProfileSubAttributeDefinitions = map[string]NetworkContextProfileAttributeDefinition{
	string(NetworkContextProfileSubAttributeTypeTLSVersion):     NetworkContextProfileTLSVersionDefinition,
	string(NetworkContextProfileSubAttributeTypeTLSCipherSuite): NetworkContextProfileTLSCipherSuiteDefinition,
	string(NetworkContextProfileSubAttributeTypeCIFSSMBVersion): NetworkContextProfileCIFSSMBVersionDefinition,
}

// NetworkContextProfileAppIDSubAttributeDefinition describes the sub-attributes available for APP_ID attributes.
var NetworkContextProfileAppIDSubAttributeDefinition = NetworkContextProfileAttributeDefinition{
	Name:        NetworkContextProfileAttributeTypeAppID,
	Description: "Application-level sub-attributes (TLS versions, cipher suites, SMB versions). Enables advanced filtering rules within APP_ID classifications.",
	Values: []NetworkContextProfileValueDefinition{
		{
			Value:       string(NetworkContextProfileSubAttributeTypeTLSVersion),
			Description: "Restrict traffic by TLS protocol version.",
		},
		{
			Value:       string(NetworkContextProfileSubAttributeTypeTLSCipherSuite),
			Description: "Restrict traffic by TLS cipher suite.",
		},
		{
			Value:       string(NetworkContextProfileSubAttributeTypeCIFSSMBVersion),
			Description: "Restrict traffic by CIFS/SMB protocol version.",
		},
	},
}

// NetworkContextProfileAppIDDefinition defines all valid Layer 7 application identifiers (APP_ID values).
var NetworkContextProfileAppIDDefinition = NetworkContextProfileAttributeDefinition{
	Name:        NetworkContextProfileAttributeTypeAppID,
	Description: "Layer 7 application classifications for network traffic. Enables policies based on specific applications and protocols rather than just IP/port, used in NSX firewall policies.",
	Values: []NetworkContextProfileValueDefinition{
		{Value: "360ANTIV", Description: "360 Total Security antivirus traffic"},
		{Value: "ACTIVDIR", Description: "Microsoft Active Directory / LDAP authentication"},
		{Value: "AMQP", Description: "Advanced Message Queuing Protocol (message broker)"},
		{Value: "AVAST", Description: "Avast antivirus update traffic"},
		{Value: "AVG", Description: "AVG antivirus update traffic"},
		{Value: "AVIRA", Description: "Avira antivirus update traffic"},
		{Value: "BDEFNDER", Description: "Bitdefender antivirus update traffic"},
		{Value: "BLAST", Description: "VMware Blast display protocol (Horizon)"},
		{Value: "CA_CERT", Description: "Certificate Authority certificate validation (OCSP/CRL)"},
		{Value: "CIFS", Description: "Common Internet File System / SMB file sharing"},
		{Value: "CLDAP", Description: "Connectionless LDAP (UDP-based directory queries)"},
		{Value: "CTRXCGP", Description: "Citrix Common Gateway Protocol"},
		{Value: "CTRXGOTO", Description: "Citrix GoTo conferencing traffic"},
		{Value: "CTRXICA", Description: "Citrix Independent Computing Architecture (ICA/HDX)"},
		{Value: "DCERPC", Description: "Microsoft DCE/RPC distributed computing traffic"},
		{Value: "DHCP", Description: "Dynamic Host Configuration Protocol"},
		{Value: "DIAMETER", Description: "Diameter AAA protocol (successor to RADIUS, used in LTE/telecom)"},
		{Value: "DNS", Description: "Domain Name System"},
		{Value: "EPIC", Description: "Epic Systems healthcare application traffic"},
		{Value: "ESET", Description: "ESET antivirus update traffic"},
		{Value: "FPROT", Description: "F-Prot antivirus update traffic"},
		{Value: "FTP", Description: "File Transfer Protocol"},
		{Value: "GITHUB", Description: "GitHub version control and CI/CD traffic"},
		{Value: "HTTP", Description: "Hypertext Transfer Protocol (plain HTTP)"},
		{Value: "HTTP2", Description: "HTTP/2 protocol"},
		{Value: "IMAP", Description: "Internet Message Access Protocol (email)"},
		{Value: "KASPRSKY", Description: "Kaspersky antivirus update traffic"},
		{Value: "KERBEROS", Description: "Kerberos network authentication protocol"},
		{Value: "LDAP", Description: "Lightweight Directory Access Protocol"},
		{Value: "MAXDB", Description: "SAP MaxDB database traffic"},
		{Value: "MCAFEE", Description: "McAfee antivirus update traffic"},
		{Value: "MSSQL", Description: "Microsoft SQL Server database traffic"},
		{Value: "MYSQL", Description: "MySQL / MariaDB database traffic"},
		{Value: "NFS", Description: "Network File System"},
		{Value: "NNTP", Description: "Network News Transfer Protocol"},
		{Value: "NTBIOSNS", Description: "NetBIOS Name Service (Windows naming)"},
		{Value: "NTP", Description: "Network Time Protocol"},
		{Value: "OCSP", Description: "Online Certificate Status Protocol"},
		{Value: "ORACLE", Description: "Oracle database traffic"},
		{Value: "PANDA", Description: "Panda Security antivirus update traffic"},
		{Value: "PCOIP", Description: "PC-over-IP display protocol (VMware Horizon / Teradici)"},
		{Value: "POP3", Description: "Post Office Protocol v3 (email retrieval)"},
		{Value: "RADIUS", Description: "Remote Authentication Dial-In User Service"},
		{Value: "RDP", Description: "Microsoft Remote Desktop Protocol"},
		{Value: "RTCP", Description: "Real-Time Control Protocol (companion to RTP)"},
		{Value: "RTP", Description: "Real-Time Transport Protocol (audio/video streams)"},
		{Value: "RTSP", Description: "Real Time Streaming Protocol"},
		{Value: "SIP", Description: "Session Initiation Protocol (VoIP signalling)"},
		{Value: "SMTP", Description: "Simple Mail Transfer Protocol"},
		{Value: "SNMP", Description: "Simple Network Management Protocol"},
		{Value: "SSH", Description: "Secure Shell"},
		{Value: "SSL", Description: "SSL/TLS encrypted traffic"},
		{Value: "SYMUPDAT", Description: "Symantec / Norton antivirus update traffic"},
		{Value: "SYSLOG", Description: "Syslog event logging"},
		{Value: "TELNET", Description: "Telnet remote terminal protocol (unencrypted)"},
		{Value: "TFTP", Description: "Trivial File Transfer Protocol"},
		{Value: "VNC", Description: "Virtual Network Computing remote desktop"},
		{Value: "WINS", Description: "Windows Internet Name Service"},
	},
}

// ============================================================================
// MODEL STRUCTS
// ============================================================================

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
