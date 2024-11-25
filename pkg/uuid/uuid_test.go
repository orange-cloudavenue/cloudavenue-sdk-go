package uuid

import (
	"testing"
)

const (
	validUUIDv4 = "12345678-1234-1234-1234-123456789012"
)

func TestUUID_ContainsPrefix(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "ContainsPrefix",
			uuid: UUID(VM.String() + validUUIDv4),
			want: true,
		},
		{
			name: "DoesNotContainPrefix",
			uuid: UUID("urn:vm:" + validUUIDv4),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.ContainsPrefix(); got != tt.want {
				t.Errorf("UUID.ContainsPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isUUIDV4(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidUUID",
			args: args{
				uuid: validUUIDv4,
			},
			want: true,
		},
		{
			name: "InvalidUUID",
			args: args{
				uuid: "f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			args: args{
				uuid: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUUIDV4(tt.args.uuid); got != tt.want {
				t.Errorf("isUUIDV4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID_IsType(t *testing.T) {
	type args struct {
		prefix UUID
	}
	tests := []struct {
		name string
		uuid UUID
		args args
		want bool
	}{
		{
			name: "IsType",
			uuid: UUID(VM.String() + validUUIDv4),
			args: args{
				prefix: VM,
			},
			want: true,
		},
		{
			name: "IsNotType",
			uuid: UUID(VM.String() + validUUIDv4),
			args: args{
				prefix: User,
			},
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			args: args{
				prefix: VM,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsType(tt.args.prefix); got != tt.want {
				t.Errorf("UUID.IsType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractUUIDv4(t *testing.T) {
	type args struct {
		uuid   string
		prefix UUID
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ExtractUUID",
			args: args{
				uuid:   VM.String() + validUUIDv4,
				prefix: VM,
			},
			want: validUUIDv4,
		},
		{
			name: "EmptyString",
			args: args{
				uuid:   "",
				prefix: VM,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractUUIDv4(tt.args.uuid, tt.args.prefix); got != tt.want {
				t.Errorf("extractUUIDv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidUUID",
			args: args{
				uuid: VM.String() + validUUIDv4,
			},
			want: true,
		},
		{
			name: "InvalidUUID",
			args: args{
				uuid: "f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
		{
			name: "InvalidPrefix",
			args: args{
				uuid: "urn:vm:f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			args: args{
				uuid: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.args.uuid); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	type args struct {
		prefix UUID
		uuid   string
	}
	tests := []struct {
		name string
		args args
		want UUID
	}{
		{
			name: "Normalize",
			args: args{
				prefix: VM,
				uuid:   validUUIDv4,
			},
			want: UUID(VM.String() + validUUIDv4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.args.prefix, tt.args.uuid); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID_IsVM(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsVM",
			uuid: UUID(VM.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVM",
			uuid: UUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVM(); got != tt.want {
				t.Errorf("UUID.IsVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID_IsUser(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsUser",
			uuid: UUID(User.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotUser",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsUser(); got != tt.want {
				t.Errorf("UUID.IsUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGroup.
func TestUUID_IsGroup(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsGroup",
			uuid: UUID(Group.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGroup",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsGroup(); got != tt.want {
				t.Errorf("UUID.IsGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGateway.
func TestUUID_IsGateway(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsGateway",
			uuid: UUID(Gateway.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGateway",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsGateway(); got != tt.want {
				t.Errorf("UUID.IsGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDC.
func TestUUID_IsVDC(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsVDC",
			uuid: UUID(VDC.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDC",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVDC(); got != tt.want {
				t.Errorf("UUID.IsVDC() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCGroup.
func TestUUID_IsVDCGroup(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{ // IsVDCGroup
			name: "IsVDCGroup",
			uuid: UUID(VDCGroup.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotVDCGroup
			name: "IsNotVDCGroup",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVDCGroup(); got != tt.want {
				t.Errorf("UUID.IsVDCGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsNetwork.
func TestUUID_IsNetwork(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsNetwork",
			uuid: UUID(Network.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotNetwork",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsNetwork(); got != tt.want {
				t.Errorf("UUID.IsNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsLoadBalancerPool.
func TestUUID_IsLoadBalancerPool(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsLoadBalancerPool",
			uuid: UUID(LoadBalancerPool.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotLoadBalancerPool",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsLoadBalancerPool(); got != tt.want {
				t.Errorf("UUID.IsLoadBalancerPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCStorageProfile.
func TestUUID_IsVDCStorageProfile(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsVDCStorageProfile",
			uuid: UUID(VDCStorageProfile.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDCStorageProfile",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVDCStorageProfile(); got != tt.want {
				t.Errorf("UUID.IsVDCStorageProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVAPP.
func TestUUID_IsVAPP(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsVAPP",
			uuid: UUID(VAPP.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVAPP",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVAPP(); got != tt.want {
				t.Errorf("UUID.IsVAPP() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsDisk.
func TestUUID_IsDisk(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsDisk",
			uuid: UUID(Disk.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotDisk",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsDisk(); got != tt.want {
				t.Errorf("UUID.IsDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsSecurityGroup.
func TestUUID_IsSecurityGroup(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsSecurityGroup",
			uuid: UUID(SecurityGroup.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotSecurityGroup",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsSecurityGroup(); got != tt.want {
				t.Errorf("UUID.IsSecurityGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVAPPTemplate.
func TestUUID_IsVAPPTemplate(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{ // IsVAPPTemplate
			name: "IsVAPPTemplate",
			uuid: UUID(VAPPTemplate.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotVAPPTemplate
			name: "IsNotVAPPTemplate",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVAPPTemplate(); got != tt.want {
				t.Errorf("UUID.IsVAPPTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsCatalog.
func TestUUID_IsCatalog(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{ // IsCatalog
			name: "IsCatalog",
			uuid: UUID(Catalog.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotCatalog
			name: "IsNotCatalog",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsCatalog(); got != tt.want {
				t.Errorf("UUID.IsCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsToken.
func TestUUID_IsToken(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{ // IsToken
			name: "IsToken",
			uuid: UUID(Token.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotToken
			name: "IsNotToken",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsToken(); got != tt.want {
				t.Errorf("UUID.IsToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUUID_IsOrg tests the UUID.IsOrg function.
func TestUUID_IsOrg(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{ // IsOrg
			name: "IsOrg",
			uuid: UUID(Org.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotOrg
			name: "IsNotOrg",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsOrg(); got != tt.want {
				t.Errorf("UUID.IsOrg() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsType tests the TestIsType function.
func TestTestIsType(t *testing.T) {
	testCases := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{
			name:     "valid uuid",
			uuidType: VM,
			uuid:     UUID(VM.String() + validUUIDv4),
			want:     true,
		},
		{
			name:     "invalid uuid",
			uuidType: VM,
			uuid:     "invalid-uuid",
			want:     false,
		},
		{
			name:     "empty value",
			uuidType: VM,
			uuid:     "",
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := TestIsType(tc.uuidType)(tc.uuid.String())
			if tc.want && err != nil {
				t.Errorf("TestIsType() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestIsGroup(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{
			name: "IsGroup",
			uuid: UUID(Group.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGroup",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			want: false,
		},
		{
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGroup(tt.uuid.String()); got != tt.want {
				t.Errorf("IsGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsGateway.
func TestIsEdgeGateway(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{
			name: "IsGateway",
			uuid: UUID(Gateway.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotGateway",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEdgeGateway(tt.uuid.String()); got != tt.want {
				t.Errorf("IsEdgeGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDC.
func TestIsVDC(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{
			name: "IsVDC",
			uuid: UUID(VDC.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDC",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDC(tt.uuid.String()); got != tt.want {
				t.Errorf("IsVDC() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCGroup.
func TestIsVDCGroup(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{
			name: "IsVDCGroup",
			uuid: UUID(VDCGroup.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDCGroup",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDCGroup(tt.uuid.String()); got != tt.want {
				t.Errorf("IsVDCGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsNetwork.
func TestIsNetwork(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{
			name: "IsNetwork",
			uuid: UUID(Network.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotNetwork",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNetwork(tt.uuid.String()); got != tt.want {
				t.Errorf("IsNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsLoadBalancerPool.
func TestIsLoadBalancerPool(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{ // IsLoadBalancerPool
			name: "IsLoadBalancerPool",
			uuid: UUID(LoadBalancerPool.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotLoadBalancerPool
			name: "IsNotLoadBalancerPool",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLoadBalancerPool(tt.uuid.String()); got != tt.want {
				t.Errorf("IsLoadBalancerPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVDCStorageProfile.
func TestIsVDCStorageProfile(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     UUID
		want     bool
	}{
		{ // IsVDCStorageProfile
			name: "IsVDCStorageProfile",
			uuid: UUID(VDCStorageProfile.String() + validUUIDv4),
			want: true,
		},
		{ // IsNotVDCStorageProfile
			name: "IsNotVDCStorageProfile",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDCStorageProfile(tt.uuid.String()); got != tt.want {
				t.Errorf("IsVDCStorageProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsVAPP.
func TestIsVAPP(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsVAPP
			name: "IsVAPP",
			uuid: UUID(VAPP.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVAPP
			name: "IsNotVAPP",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVAPP(tt.uuid); got != tt.want {
				t.Errorf("IsVAPP() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsDisk.
func TestIsDisk(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsDisk
			name: "IsDisk",
			uuid: UUID(Disk.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotDisk
			name: "IsNotDisk",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDisk(tt.uuid); got != tt.want {
				t.Errorf("IsDisk() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsSecurityGroup.
func TestIsSecurityGroup(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsSecurityGroup
			name: "IsSecurityGroup",
			uuid: UUID(SecurityGroup.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotSecurityGroup
			name: "IsNotSecurityGroup",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSecurityGroup(tt.uuid); got != tt.want {
				t.Errorf("IsSecurityGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVCDA.
func TestIsVCDA(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsVCDA
			name: "IsVCDA",
			uuid: UUID(VCDA.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVCDA
			name: "IsNotVCDA",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVCDA(tt.uuid); got != tt.want {
				t.Errorf("IsVCDA() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVM.
func TestIsVM(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsVM
			name: "IsVM",
			uuid: UUID(VM.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVM
			name: "IsNotVM",
			uuid: UUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVM(tt.uuid); got != tt.want {
				t.Errorf("IsVM() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsUser.
func TestIsUser(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsUser
			name: "IsUser",
			uuid: UUID(User.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotUser
			name: "IsNotUser",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUser(tt.uuid); got != tt.want {
				t.Errorf("IsUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsOrg
func TestIsOrg(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsOrg
			name: "IsOrg",
			uuid: UUID(Org.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotOrg
			name: "IsNotOrg",
			uuid: UUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsOrg(tt.uuid); got != tt.want {
				t.Errorf("IsOrg() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsToken.
func TestIsToken(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsToken
			name: "IsToken",
			uuid: UUID(Token.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotToken
			name: "IsNotToken",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsToken(tt.uuid); got != tt.want {
				t.Errorf("IsToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVAPPTemplate.
func TestIsVAPPTemplate(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsVAPPTemplate
			name: "IsVAPPTemplate",
			uuid: UUID(VAPPTemplate.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVAPPTemplate
			name: "IsNotVAPPTemplate",
			uuid: UUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVAPPTemplate(tt.uuid); got != tt.want {
				t.Errorf("IsVAPPTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsCatalog
func TestIsCatalog(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsCatalog
			name: "IsCatalog",
			uuid: UUID(Catalog.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotCatalog
			name: "IsNotCatalog",
			uuid: UUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCatalog(tt.uuid); got != tt.want {
				t.Errorf("IsCatalog() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsVDCComputePolicy
func TestIsVDCComputePolicy(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsVDCComputePolicy
			name: "IsVDCComputePolicy",
			uuid: UUID(VDCComputePolicy.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotVDCComputePolicy
			name: "IsNotVDCComputePolicy",
			uuid: UUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVDCComputePolicy(tt.uuid); got != tt.want {
				t.Errorf("IsVDCComputePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsNetworkContextProfile
func TestIsNetworkContextProfile(t *testing.T) {
	tests := []struct {
		name     string
		uuidType UUID
		uuid     string
		want     bool
	}{
		{ // IsNetworkContextProfile
			name: "IsNetworkContextProfile",
			uuid: UUID(NetworkContextProfile.String() + validUUIDv4).String(),
			want: true,
		},
		{ // IsNotNetworkContextProfile
			name: "IsNotNetworkContextProfile",
			uuid: UUID("urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791").String(),
			want: false,
		},
		{ // EmptyString
			name: "EmptyString",
			uuid: UUID("").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNetworkContextProfile(tt.uuid); got != tt.want {
				t.Errorf("IsNetworkContextProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUUID_IsVDCComputePolicy(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsVDCComputePolicy",
			uuid: UUID(VDCComputePolicy.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotVDCComputePolicy",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsVDCComputePolicy(); got != tt.want {
				t.Errorf("UUID.IsVDCComputePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

// IsNetworkContextProfile
func TestUUID_IsNetworkContextProfile(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			name: "IsNetworkContextProfile",
			uuid: UUID(NetworkContextProfile.String() + validUUIDv4),
			want: true,
		},
		{
			name: "IsNotNetworkContextProfile",
			uuid: UUID("urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"),
			want: false,
		},
		{
			name: "EmptyString",
			uuid: UUID(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uuid.IsNetworkContextProfile(); got != tt.want {
				t.Errorf("UUID.IsNetworkContextProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
