package urn

import "testing"

func TestIsGroup(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{
			name: "IsGroup",
			urn:  URN(Group.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGroup",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGroup(tt.urn.String()); got != tt.want {
				t.Errorf("IsGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGateway.
func TestIsEdgeGateway(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{
			name: "IsGateway",
			urn:  URN(Gateway.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGateway",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEdgeGateway(tt.urn.String()); got != tt.want {
				t.Errorf("IsEdgeGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDC.
func TestIsVDC(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{
			name: "IsVDC",
			urn:  URN(VDC.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDC",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDC(tt.urn.String()); got != tt.want {
				t.Errorf("IsVDC() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCGroup.
func TestIsVDCGroup(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{
			name: "IsVDCGroup",
			urn:  URN(VDCGroup.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDCGroup",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDCGroup(tt.urn.String()); got != tt.want {
				t.Errorf("IsVDCGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsNetwork.
func TestIsNetwork(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{
			name: "IsNetwork",
			urn:  URN(Network.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotNetwork",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNetwork(tt.urn.String()); got != tt.want {
				t.Errorf("IsNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsLoadBalancerPool.
func TestIsLoadBalancerPool(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{ // IsLoadBalancerPool
			name: "IsLoadBalancerPool",
			urn:  URN(LoadBalancerPool.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotLoadBalancerPool
			name: "IsNotLoadBalancerPool",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLoadBalancerPool(tt.urn.String()); got != tt.want {
				t.Errorf("IsLoadBalancerPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCStorageProfile.
func TestIsVDCStorageProfile(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{ // IsVDCStorageProfile
			name: "IsVDCStorageProfile",
			urn:  URN(VDCStorageProfile.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotVDCStorageProfile
			name: "IsNotVDCStorageProfile",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDCStorageProfile(tt.urn.String()); got != tt.want {
				t.Errorf("IsVDCStorageProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVAPP.
func TestIsVAPP(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsVAPP
			name: "IsVAPP",
			urn:  URN(VAPP.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVAPP
			name: "IsNotVAPP",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVAPP(tt.urn); got != tt.want {
				t.Errorf("IsVAPP() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsDisk.
func TestIsDisk(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsDisk
			name: "IsDisk",
			urn:  URN(Disk.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotDisk
			name: "IsNotDisk",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDisk(tt.urn); got != tt.want {
				t.Errorf("IsDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsSecurityGroup.
func TestIsSecurityGroup(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsSecurityGroup
			name: "IsSecurityGroup",
			urn:  URN(SecurityGroup.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotSecurityGroup
			name: "IsNotSecurityGroup",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSecurityGroup(tt.urn); got != tt.want {
				t.Errorf("IsSecurityGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVCDA.
func TestIsVCDA(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsVCDA
			name: "IsVCDA",
			urn:  URN(VCDA.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVCDA
			name: "IsNotVCDA",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVCDA(tt.urn); got != tt.want {
				t.Errorf("IsVCDA() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVM.
func TestIsVM(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsVM
			name: "IsVM",
			urn:  URN(VM.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVM
			name: "IsNotVM",
			urn:  URN("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVM(tt.urn); got != tt.want {
				t.Errorf("IsVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsUser.
func TestIsUser(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsUser
			name: "IsUser",
			urn:  URN(User.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotUser
			name: "IsNotUser",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUser(tt.urn); got != tt.want {
				t.Errorf("IsUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsOrg.
func TestIsOrg(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsOrg
			name: "IsOrg",
			urn:  URN(Org.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotOrg
			name: "IsNotOrg",
			urn:  URN("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOrg(tt.urn); got != tt.want {
				t.Errorf("IsOrg() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsToken.
func TestIsToken(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsToken
			name: "IsToken",
			urn:  URN(Token.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotToken
			name: "IsNotToken",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsToken(tt.urn); got != tt.want {
				t.Errorf("IsToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsAppPortProfile.
func TestIsAppPortProfile(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsAppPortProfile
			name: "IsAppPortProfile",
			urn:  URN(AppPortProfile.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotAppPortProfile
			name: "IsNotAppPortProfile",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAppPortProfile(tt.urn); got != tt.want {
				t.Errorf("IsAppPortProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVDCComputePolicy.
func TestIsVDCComputePolicy(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsVDCComputePolicy
			name: "IsVDCComputePolicy",
			urn:  URN(VDCComputePolicy.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVDCComputePolicy
			name: "IsNotVDCComputePolicy",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDCComputePolicy(tt.urn); got != tt.want {
				t.Errorf("IsAppPortProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsCatalog.
func TestIsCatalog(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsCatalog
			name: "IsCatalog",
			urn:  URN(Catalog.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotCatalog
			name: "IsNotCatalog",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCatalog(tt.urn); got != tt.want {
				t.Errorf("IsCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVAPPTemplate.
func TestIsVAPPTemplate(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsVAPPTemplate
			name: "IsVAPPTemplate",
			urn:  URN(VAPPTemplate.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVAPPTemplate
			name: "IsNotVAPPTemplate",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVAPPTemplate(tt.urn); got != tt.want {
				t.Errorf("IsVAPPTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsCertificateLibraryItem.
func TestIsCertificateLibraryItem(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsCertificateLibraryItem
			name: "IsCertificateLibraryItem",
			urn:  URN(CertificateLibraryItem.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotCertificateLibraryItem
			name: "IsNotCertificateLibraryItem",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCertificateLibraryItem(tt.urn); got != tt.want {
				t.Errorf("IsCertificateLibraryItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsLoadBalancerVirtualService.
func TestIsLoadBalancerVirtualService(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsLoadBalancerVirtualService
			name: "IsLoadBalancerVirtualService",
			urn:  URN(LoadBalancerVirtualService.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotLoadBalancerVirtualService
			name: "IsNotLoadBalancerVirtualService",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLoadBalancerVirtualService(tt.urn); got != tt.want {
				t.Errorf("IsLoadBalancerVirtualService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsServiceEngineGroup.
func TestIsServiceEngineGroup(t *testing.T) {
	tests := []struct {
		name    string
		urnType URN
		urn     string
		want    bool
	}{
		{ // IsServiceEngineGroup
			name: "IsServiceEngineGroup",
			urn:  URN(ServiceEngineGroup.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotServiceEngineGroup
			name: "IsNotServiceEngineGroup",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			urn:  URN("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServiceEngineGroup(tt.urn); got != tt.want {
				t.Errorf("IsServiceEngineGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
