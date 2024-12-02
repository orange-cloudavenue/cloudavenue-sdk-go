package urn

// IsAppPortProfile returns true if the URN is a AppPortProfile URN.
func (urn URN) IsAppPortProfile() bool {
	return urn.IsType(AppPortProfile)
}

// IsOrg returns true if the URN is an Org URN.
func (urn URN) IsOrg() bool {
	return urn.IsType(Org)
}

// IsVM returns true if the URN is a VM URN.
func (urn URN) IsVM() bool {
	return urn.IsType(VM)
}

// IsUser returns true if the URN is a User URN.
func (urn URN) IsUser() bool {
	return urn.IsType(User)
}

// IsGroup returns true if the URN is a Group URN.
func (urn URN) IsGroup() bool {
	return urn.IsType(Group)
}

// IsGateway returns true if the URN is a Gateway URN.
func (urn URN) IsGateway() bool {
	return urn.IsType(Gateway)
}

// IsVDC returns true if the URN is a VDC URN.
func (urn URN) IsVDC() bool {
	return urn.IsType(VDC)
}

// IsVDCGroup returns true if the URN is a VDCGroup URN.
func (urn URN) IsVDCGroup() bool {
	return urn.IsType(VDCGroup)
}

// IsNetwork returns true if the URN is a Network URN.
func (urn URN) IsNetwork() bool {
	return urn.IsType(Network)
}

// IsLoadBalancerPool returns true if the URN is a LoadBalancerPool URN.
func (urn URN) IsLoadBalancerPool() bool {
	return urn.IsType(LoadBalancerPool)
}

// IsVDCStorageProfile returns true if the URN is a VDCStorageProfile URN.
func (urn URN) IsVDCStorageProfile() bool {
	return urn.IsType(VDCStorageProfile)
}

// IsVAPP returns true if the URN is a VAPP URN.
func (urn URN) IsVAPP() bool {
	return urn.IsType(VAPP)
}

// IsVAPPTemplate returns true if the URN is a VAPPTemplate URN.
func (urn URN) IsVAPPTemplate() bool {
	return urn.IsType(VAPPTemplate)
}

// IsDisk returns true if the URN is a Disk URN.
func (urn URN) IsDisk() bool {
	return urn.IsType(Disk)
}

// IsSecurityGroup returns true if the URN is a SecurityGroup URN.
func (urn URN) IsSecurityGroup() bool {
	return urn.IsType(SecurityGroup)
}

// IsCatalog returns true if the URN is a Catalog URN.
func (urn URN) IsCatalog() bool {
	return urn.IsType(Catalog)
}

// IsToken returns true if the URN is a Token URN.
func (urn URN) IsToken() bool {
	return urn.IsType(Token)
}

// IsVDCComputePolicy returns true if the URN is a VDCComputePolicy URN.
func (urn URN) IsVDCComputePolicy() bool {
	return urn.IsType(VDCComputePolicy)
}
