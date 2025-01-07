/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package urn

import "testing"

func TestURN_IsVM(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsVM",
			urn:  URN(VM.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVM",
			urn:  URN("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsVM(); got != tt.want {
				t.Errorf("URN.IsVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURN_IsUser(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsUser",
			urn:  URN(User.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotUser",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsUser(); got != tt.want {
				t.Errorf("URN.IsUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGroup.
func TestURN_IsGroup(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
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
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsGroup(); got != tt.want {
				t.Errorf("URN.IsGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGateway.
func TestURN_IsGateway(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsGateway",
			urn:  URN(Gateway.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGateway",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsGateway(); got != tt.want {
				t.Errorf("URN.IsGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsAppPortProfile tests the IsAppProfile function.
func TestURN_IsAppPortProfile(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsAppProfile
			name: "IsAppProfile",
			urn:  URN(AppPortProfile.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotAppProfile
			name: "IsNotAppProfile",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsAppPortProfile(); got != tt.want {
				t.Errorf("VcloudURN.IsAppPortProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDC.
func TestURN_IsVDC(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsVDC",
			urn:  URN(VDC.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDC",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsVDC(); got != tt.want {
				t.Errorf("URN.IsVDC() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCGroup.
func TestURN_IsVDCGroup(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsVDCGroup
			name: "IsVDCGroup",
			urn:  URN(VDCGroup.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotVDCGroup
			name: "IsNotVDCGroup",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsVDCGroup(); got != tt.want {
				t.Errorf("URN.IsVDCGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsNetwork.
func TestURN_IsNetwork(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsNetwork",
			urn:  URN(Network.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotNetwork",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsNetwork(); got != tt.want {
				t.Errorf("URN.IsNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsLoadBalancerPool.
func TestURN_IsLoadBalancerPool(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsLoadBalancerPool",
			urn:  URN(LoadBalancerPool.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotLoadBalancerPool",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsLoadBalancerPool(); got != tt.want {
				t.Errorf("URN.IsLoadBalancerPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCStorageProfile.
func TestURN_IsVDCStorageProfile(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsVDCStorageProfile",
			urn:  URN(VDCStorageProfile.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDCStorageProfile",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsVDCStorageProfile(); got != tt.want {
				t.Errorf("URN.IsVDCStorageProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVAPP.
func TestURN_IsVAPP(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsVAPP",
			urn:  URN(VAPP.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVAPP",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsVAPP(); got != tt.want {
				t.Errorf("URN.IsVAPP() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsDisk.
func TestURN_IsDisk(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsDisk",
			urn:  URN(Disk.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotDisk",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsDisk(); got != tt.want {
				t.Errorf("URN.IsDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsSecurityGroup.
func TestURN_IsSecurityGroup(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "IsSecurityGroup",
			urn:  URN(SecurityGroup.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotSecurityGroup",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsSecurityGroup(); got != tt.want {
				t.Errorf("URN.IsSecurityGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVAPPTemplate.
func TestURN_IsVAPPTemplate(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsVAPPTemplate
			name: "IsVAPPTemplate",
			urn:  URN(VAPPTemplate.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotVAPPTemplate
			name: "IsNotVAPPTemplate",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsVAPPTemplate(); got != tt.want {
				t.Errorf("URN.IsVAPPTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsCatalog.
func TestURN_IsCatalog(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsCatalog
			name: "IsCatalog",
			urn:  URN(Catalog.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotCatalog
			name: "IsNotCatalog",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsCatalog(); got != tt.want {
				t.Errorf("URN.IsCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsToken.
func TestURN_IsToken(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsToken
			name: "IsToken",
			urn:  URN(Token.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotToken
			name: "IsNotToken",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsToken(); got != tt.want {
				t.Errorf("URN.IsToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsCatalog.
func TestVcloudURN_IsCatalog(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsCatalog
			name: "IsCatalog",
			urn:  URN(Catalog.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotCatalog
			name: "IsNotCatalog",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsCatalog(); got != tt.want {
				t.Errorf("VcloudURN.IsCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestURN_IsOrg tests the URN.IsOrg function.
func TestURN_IsOrg(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsOrg
			name: "IsOrg",
			urn:  URN(Org.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotOrg
			name: "IsNotOrg",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsOrg(); got != tt.want {
				t.Errorf("URN.IsOrg() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestURN_IsVDCComputePolicy tests the URN.IsVDCComputePolicy function.
func TestURN_IsVDCComputePolicy(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsVDCComputePolicy
			name: "IsVDCComputePolicy",
			urn:  URN(VDCComputePolicy.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotVDCComputePolicy
			name: "IsNotVDCComputePolicy",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsVDCComputePolicy(); got != tt.want {
				t.Errorf("URN.IsVDCComputePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestURN_IsLoadBalancerVirtualService tests the URN.IsLoadBalancerVirtualService function.
func TestURN_IsLoadBalancerVirtualService(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsLoadBalancerVirtualService
			name: "IsLoadBalancerVirtualService",
			urn:  URN(LoadBalancerVirtualService.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotLoadBalancerVirtualService
			name: "IsNotLoadBalancerVirtualService",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsLoadBalancerVirtualService(); got != tt.want {
				t.Errorf("URN.IsLoadBalancerVirtualService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestURN_IsCertificateLibraryItem tests the URN.IsCertificateLibraryItem function.
func TestURN_IsCertificateLibraryItem(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsCertificateLibraryItem
			name: "IsCertificateLibraryItem",
			urn:  URN(CertificateLibraryItem.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotCertificateLibraryItem
			name: "IsNotCertificateLibraryItem",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsCertificateLibraryItem(); got != tt.want {
				t.Errorf("URN.IsCertificateLibraryItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestURN_IsServiceEngineGroup tests the URN.IsServiceEngineGroup function.
func TestURN_IsServiceEngineGroup(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{ // IsServiceEngineGroup
			name: "IsServiceEngineGroup",
			urn:  URN(ServiceEngineGroup.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotServiceEngineGroup
			name: "IsNotServiceEngineGroup",
			urn:  URN("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsServiceEngineGroup(); got != tt.want {
				t.Errorf("URN.IsServiceEngineGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
