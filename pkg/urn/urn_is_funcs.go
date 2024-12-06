package urn

// IsAppPortProfile returns true if the URN is a AppPortProfile URN.
func IsAppPortProfile(urn string) bool {
	return URN(urn).IsType(AppPortProfile)
}

// IsOrg returns true if the URN is an Org URN.
func IsOrg(urn string) bool {
	return URN(urn).IsType(Org)
}

// IsEdgeGateway returns true if the URN is a EdgeGateway URN.
func IsEdgeGateway(urn string) bool {
	return URN(urn).IsType(Gateway)
}

// IsVDC returns true if the URN is a VDC URN.
func IsVDC(urn string) bool {
	return URN(urn).IsType(VDC)
}

// IsVDCGroup returns true if the URN is a VDCGroup URN.
func IsVDCGroup(urn string) bool {
	return URN(urn).IsType(VDCGroup)
}

// IsNetwork returns true if the URN is a Network URN.
func IsNetwork(urn string) bool {
	return URN(urn).IsType(Network)
}

// IsLoadBalancerPool returns true if the URN is a LoadBalancerPool URN.
func IsLoadBalancerPool(urn string) bool {
	return URN(urn).IsType(LoadBalancerPool)
}

// IsVDCStorageProfile returns true if the URN is a VDCStorageProfile URN.
func IsVDCStorageProfile(urn string) bool {
	return URN(urn).IsType(VDCStorageProfile)
}

// IsVAPP returns true if the URN is a VAPP URN.
func IsVAPP(urn string) bool {
	return URN(urn).IsType(VAPP)
}

// IsVAPPTemplate returns true if the URN is a VAPPTemplate URN.
func IsVAPPTemplate(urn string) bool {
	return URN(urn).IsType(VAPPTemplate)
}

// IsDisk returns true if the URN is a Disk URN.
func IsDisk(urn string) bool {
	return URN(urn).IsType(Disk)
}

// IsSecurityGroup returns true if the URN is a SecurityGroup URN.
func IsSecurityGroup(urn string) bool {
	return URN(urn).IsType(SecurityGroup)
}

// IsVCDA returns true if the URN is a VCDA URN.
func IsVCDA(urn string) bool {
	return URN(urn).IsType(VCDA)
}

// IsVM returns true if the URN is a VM URN.
func IsVM(urn string) bool {
	return URN(urn).IsType(VM)
}

// IsUser returns true if the URN is a User URN.
func IsUser(urn string) bool {
	return URN(urn).IsType(User)
}

// IsGroup returns true if the URN is a Group URN.
func IsGroup(urn string) bool {
	return URN(urn).IsType(Group)
}

// IsCatalog returns true if the URN is a Catalog URN.
func IsCatalog(urn string) bool {
	return URN(urn).IsType(Catalog)
}

// IsToken returns true if the URN is a Token URN.
func IsToken(urn string) bool {
	return URN(urn).IsType(Token)
}

// IsVDCComputePolicy returns true if the URN is a VDCComputePolicy URN.
func IsVDCComputePolicy(urn string) bool {
	return URN(urn).IsType(VDCComputePolicy)
}

// IsCertificateLibraryItem returns true if the URN is a CertificateLibraryItem URN.
func IsCertificateLibraryItem(urn string) bool {
	return URN(urn).IsType(CertificateLibraryItem)
}

// IsLoadBalancerVirtualService returns true if the URN is a LoadBalancerVirtualService URN.
func IsLoadBalancerVirtualService(urn string) bool {
	return URN(urn).IsType(LoadBalancerVirtualService)
}

// IsServiceEngineGroup returns true if the URN is a ServiceEngineGroup URN.
func IsServiceEngineGroup(urn string) bool {
	return URN(urn).IsType(ServiceEngineGroup)
}
