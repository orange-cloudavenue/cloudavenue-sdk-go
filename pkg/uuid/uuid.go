package uuid

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// Prefixes is the list of prefixes.
	VcloudPrefix      = "urn:vcloud:"
	CloudAvenuePrefix = "urn:cloudavenue:"

	// * VCD.
	Org                   = UUID(VcloudPrefix + "org:")
	VM                    = UUID(VcloudPrefix + "vm:")
	User                  = UUID(VcloudPrefix + "user:")
	Group                 = UUID(VcloudPrefix + "group:")
	Gateway               = UUID(VcloudPrefix + "gateway:")
	VDC                   = UUID(VcloudPrefix + "vdc:")
	VDCGroup              = UUID(VcloudPrefix + "vdcGroup:")
	VDCComputePolicy      = UUID(VcloudPrefix + "vdcComputePolicy:")
	Network               = UUID(VcloudPrefix + "network:")
	LoadBalancerPool      = UUID(VcloudPrefix + "loadBalancerPool:")
	VDCStorageProfile     = UUID(VcloudPrefix + "vdcstorageProfile:")
	VAPP                  = UUID(VcloudPrefix + "vapp:")
	VAPPTemplate          = UUID(VcloudPrefix + "vappTemplate:")
	Disk                  = UUID(VcloudPrefix + "disk:")
	SecurityGroup         = UUID(VcloudPrefix + "firewallGroup:")
	Catalog               = UUID(VcloudPrefix + "catalog:")
	Token                 = UUID(VcloudPrefix + "token:")
	NetworkContextProfile = UUID(VcloudPrefix + "networkContextProfile:")

	// * CLOUDAVENUE.
	VCDA = UUID(CloudAvenuePrefix + "vcda:")
)

var UUIDs = []UUID{
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
	NetworkContextProfile,
}

type (
	UUID string
)

// String returns the string representation of the UUID.
func (uuid UUID) String() string {
	return string(uuid)
}

// IsType returns true if the UUID is of the specified type.
func (uuid UUID) IsType(prefix UUID) bool {
	if uuid.isEmpty() || prefix.isEmpty() {
		return false
	}

	return strings.HasPrefix(string(uuid), prefix.String()) && isUUIDV4(uuid.extractUUIDv4(prefix))
}

// isNotEmpty returns true if the UUID is not empty.
func (uuid UUID) isEmpty() bool {
	return len(uuid) == 0
}

func isUUIDV4(uuid string) bool {
	return regexp.MustCompile(`(?m)^\w{8}-\w{4}-\w{4}-\w{4}-\w{12}$`).MatchString(uuid)
}

func IsUUIDV4(uuid string) bool {
	return isUUIDV4(uuid)
}

// ContainsPrefix returns true if the UUID contains any prefix.
func (uuid UUID) ContainsPrefix() bool {
	return strings.Contains(string(uuid), string(VcloudPrefix))
}

// extractUUIDv4 returns the UUIDv4 from the UUID.
func (uuid UUID) extractUUIDv4(prefix UUID) string {
	return extractUUIDv4(uuid.String(), prefix)
}

func extractUUIDv4(uuid string, prefix UUID) string {
	if len(uuid) == 0 || prefix.isEmpty() {
		return ""
	}

	return uuid[len(prefix):]
}

func IsValid(uuid string) bool {
	if len(uuid) == 0 {
		return false
	}

	u := UUID(uuid)

	for _, prefix := range UUIDs {
		if u.IsType(prefix) {
			return isUUIDV4(extractUUIDv4(uuid, prefix))
		}
	}
	return false
}

// Normalize returns the UUID with the prefix if prefix is missing.
func Normalize(prefix UUID, uuid string) UUID {
	if len(uuid) == 0 || prefix.isEmpty() {
		return ""
	}

	u := UUID(uuid)
	if u.ContainsPrefix() {
		return u
	}

	return prefix + u
}

// IsOrg returns true if the UUID is an Org UUID.
func (uuid UUID) IsOrg() bool {
	return uuid.IsType(Org)
}

// IsVM returns true if the UUID is a VM UUID.
func (uuid UUID) IsVM() bool {
	return uuid.IsType(VM)
}

// IsUser returns true if the UUID is a User UUID.
func (uuid UUID) IsUser() bool {
	return uuid.IsType(User)
}

// IsGroup returns true if the UUID is a Group UUID.
func (uuid UUID) IsGroup() bool {
	return uuid.IsType(Group)
}

// IsGateway returns true if the UUID is a Gateway UUID.
func (uuid UUID) IsGateway() bool {
	return uuid.IsType(Gateway)
}

// IsVDC returns true if the UUID is a VDC UUID.
func (uuid UUID) IsVDC() bool {
	return uuid.IsType(VDC)
}

// IsVDCGroup returns true if the UUID is a VDCGroup UUID.
func (uuid UUID) IsVDCGroup() bool {
	return uuid.IsType(VDCGroup)
}

// IsNetwork returns true if the UUID is a Network UUID.
func (uuid UUID) IsNetwork() bool {
	return uuid.IsType(Network)
}

// IsLoadBalancerPool returns true if the UUID is a LoadBalancerPool UUID.
func (uuid UUID) IsLoadBalancerPool() bool {
	return uuid.IsType(LoadBalancerPool)
}

// IsVDCStorageProfile returns true if the UUID is a VDCStorageProfile UUID.
func (uuid UUID) IsVDCStorageProfile() bool {
	return uuid.IsType(VDCStorageProfile)
}

// IsVAPP returns true if the UUID is a VAPP UUID.
func (uuid UUID) IsVAPP() bool {
	return uuid.IsType(VAPP)
}

// IsVAPPTemplate returns true if the UUID is a VAPPTemplate UUID.
func (uuid UUID) IsVAPPTemplate() bool {
	return uuid.IsType(VAPPTemplate)
}

// IsDisk returns true if the UUID is a Disk UUID.
func (uuid UUID) IsDisk() bool {
	return uuid.IsType(Disk)
}

// IsSecurityGroup returns true if the UUID is a SecurityGroup UUID.
func (uuid UUID) IsSecurityGroup() bool {
	return uuid.IsType(SecurityGroup)
}

// IsCatalog returns true if the UUID is a Catalog UUID.
func (uuid UUID) IsCatalog() bool {
	return uuid.IsType(Catalog)
}

// IsToken returns true if the UUID is a Token UUID.
func (uuid UUID) IsToken() bool {
	return uuid.IsType(Token)
}

// IsVDCComputePolicy returns true if the UUID is a VDCComputePolicy UUID.
func (uuid UUID) IsVDCComputePolicy() bool {
	return uuid.IsType(VDCComputePolicy)
}

// IsNetworkContextProfile returns true if the UUID is a NetworkContextProfile UUID.
func (uuid UUID) IsNetworkContextProfile() bool {
	return uuid.IsType(NetworkContextProfile)
}

// * End Methods

// IsOrg returns true if the UUID is an Org UUID.
func IsOrg(uuid string) bool {
	return UUID(uuid).IsType(Org)
}

// IsEdgeGateway returns true if the UUID is a EdgeGateway UUID.
func IsEdgeGateway(uuid string) bool {
	return UUID(uuid).IsType(Gateway)
}

// IsVDC returns true if the UUID is a VDC UUID.
func IsVDC(uuid string) bool {
	return UUID(uuid).IsType(VDC)
}

// IsVDCGroup returns true if the UUID is a VDCGroup UUID.
func IsVDCGroup(uuid string) bool {
	return UUID(uuid).IsType(VDCGroup)
}

// IsNetwork returns true if the UUID is a Network UUID.
func IsNetwork(uuid string) bool {
	return UUID(uuid).IsType(Network)
}

// IsLoadBalancerPool returns true if the UUID is a LoadBalancerPool UUID.
func IsLoadBalancerPool(uuid string) bool {
	return UUID(uuid).IsType(LoadBalancerPool)
}

// IsVDCStorageProfile returns true if the UUID is a VDCStorageProfile UUID.
func IsVDCStorageProfile(uuid string) bool {
	return UUID(uuid).IsType(VDCStorageProfile)
}

// IsVAPP returns true if the UUID is a VAPP UUID.
func IsVAPP(uuid string) bool {
	return UUID(uuid).IsType(VAPP)
}

// IsVAPPTemplate returns true if the UUID is a VAPPTemplate UUID.
func IsVAPPTemplate(uuid string) bool {
	return UUID(uuid).IsType(VAPPTemplate)
}

// IsDisk returns true if the UUID is a Disk UUID.
func IsDisk(uuid string) bool {
	return UUID(uuid).IsType(Disk)
}

// IsSecurityGroup returns true if the UUID is a SecurityGroup UUID.
func IsSecurityGroup(uuid string) bool {
	return UUID(uuid).IsType(SecurityGroup)
}

// IsVCDA returns true if the UUID is a VCDA UUID.
func IsVCDA(uuid string) bool {
	return UUID(uuid).IsType(VCDA)
}

// IsVM returns true if the UUID is a VM UUID.
func IsVM(uuid string) bool {
	return UUID(uuid).IsType(VM)
}

// IsUser returns true if the UUID is a User UUID.
func IsUser(uuid string) bool {
	return UUID(uuid).IsType(User)
}

// IsGroup returns true if the UUID is a Group UUID.
func IsGroup(uuid string) bool {
	return UUID(uuid).IsType(Group)
}

// IsCatalog returns true if the UUID is a Catalog UUID.
func IsCatalog(uuid string) bool {
	return UUID(uuid).IsType(Catalog)
}

// IsToken returns true if the UUID is a Token UUID.
func IsToken(uuid string) bool {
	return UUID(uuid).IsType(Token)
}

// IsVDCComputePolicy returns true if the UUID is a VDCComputePolicy UUID.
func IsVDCComputePolicy(uuid string) bool {
	return UUID(uuid).IsType(VDCComputePolicy)
}

// IsNetworkContextProfile returns true if the UUID is a NetworkContextProfile UUID.
func IsNetworkContextProfile(uuid string) bool {
	return UUID(uuid).IsType(NetworkContextProfile)
}

// * End Functions

// Special functions for the terraform provider test.
// TestIsType returns true if the UUID is of the specified type.
func TestIsType(uuidType UUID) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if value == "" {
			return nil
		}

		if !UUID(value).IsType(uuidType) {
			return fmt.Errorf("uuid %s is not of type %s", value, uuidType)
		}
		return nil
	}
}
