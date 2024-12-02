package urn

import (
	"regexp"
	"strings"
)

const (
	// Prefixes is the list of prefixes.
	VcloudPrefix      = "urn:vcloud:"
	CloudAvenuePrefix = "urn:cloudavenue:"

	// * VCD.
	Org               = URN(VcloudPrefix + "org:")
	VM                = URN(VcloudPrefix + "vm:")
	User              = URN(VcloudPrefix + "user:")
	Group             = URN(VcloudPrefix + "group:")
	Gateway           = URN(VcloudPrefix + "gateway:")
	VDC               = URN(VcloudPrefix + "vdc:")
	VDCGroup          = URN(VcloudPrefix + "vdcGroup:")
	VDCComputePolicy  = URN(VcloudPrefix + "vdcComputePolicy:")
	Network           = URN(VcloudPrefix + "network:")
	LoadBalancerPool  = URN(VcloudPrefix + "loadBalancerPool:")
	VDCStorageProfile = URN(VcloudPrefix + "vdcstorageProfile:")
	VAPP              = URN(VcloudPrefix + "vapp:")
	VAPPTemplate      = URN(VcloudPrefix + "vappTemplate:")
	Disk              = URN(VcloudPrefix + "disk:")
	SecurityGroup     = URN(VcloudPrefix + "firewallGroup:")
	Catalog           = URN(VcloudPrefix + "catalog:")
	Token             = URN(VcloudPrefix + "token:")
	AppPortProfile    = URN(VcloudPrefix + "applicationPortProfile:")

	// * CLOUDAVENUE.
	VCDA = URN(CloudAvenuePrefix + "vcda:")
)

var URNs = []URN{
	Org,
	VM,
	User,
	Group,
	Gateway,
	VDC,
	VDCGroup,
	VDCComputePolicy,
	Network,
	LoadBalancerPool,
	VDCStorageProfile,
	VAPP,
	VAPPTemplate,
	Disk,
	SecurityGroup,
	Catalog,
	Token,
	AppPortProfile,
}

type (
	URN string
)

// String returns the string representation of the URN.
func (urn URN) String() string {
	return string(urn)
}

// IsType returns true if the URN is of the specified type.
func (urn URN) IsType(prefix URN) bool {
	if urn.isEmpty() || prefix.isEmpty() {
		return false
	}

	return strings.HasPrefix(string(urn), prefix.String()) && isUUIDV4(urn.extractUUIDv4(prefix))
}

// isEmpty returns true if the URN is empty.
func (urn URN) isEmpty() bool {
	return len(urn) == 0
}

func isUUIDV4(urn string) bool {
	return regexp.MustCompile(`^([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})$`).MatchString(urn)
}

func IsUUIDV4(urn string) bool {
	return isUUIDV4(urn)
}

// ContainsPrefix returns true if the URN contains any prefix.
func (urn URN) ContainsPrefix() bool {
	// TODO add support for CloudAvenuePrefix
	return strings.Contains(string(urn), string(VcloudPrefix))
}

// extractUUIDv4 returns the UUIDv4 from the URN.
func (urn URN) extractUUIDv4(prefix URN) string {
	return extractUUIDv4(urn.String(), prefix)
}

func extractUUIDv4(urn string, prefix URN) string {
	if len(urn) == 0 || prefix.isEmpty() {
		return ""
	}

	return urn[len(prefix):]
}
