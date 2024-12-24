package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

var _ FirewallGroupInterface = (*EdgeClient)(nil)

// * SecurityGroup

// CreateFirewallSecurityGroup allow creating a new security group. T
func (e *EdgeClient) CreateFirewallSecurityGroup(securityGroupConfig *FirewallGroupSecurityGroupModel) (*FirewallGroupSecurityGroup, error) {
	ower := &govcdtypes.OpenApiReference{}

	if e.OwnerType.IsVDCGROUP() {
		ower.Name = e.OwnerName
	} else {
		ower.Name = e.vcdEdge.EdgeGateway.Name
		ower.ID = e.vcdEdge.EdgeGateway.ID
	}

	securityGroup, err := e.vcdEdge.CreateNsxtFirewallGroup(&govcdtypes.NsxtFirewallGroup{
		Name:        securityGroupConfig.Name,
		Description: securityGroupConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeSecurityGroup,
		Members:     securityGroupConfig.Members,
		OwnerRef:    ower,
	})
	if err != nil {
		return nil, err
	}

	securityGroupConfig.ID = securityGroup.NsxtFirewallGroup.ID

	return &FirewallGroupSecurityGroup{
		edgeClient:                      e,
		FirewallGroupSecurityGroupModel: securityGroupConfig,
		fwGroup:                         securityGroup,
	}, nil
}

// GetFirewallSecurityGroup retrieves the security group configuration for the Edge Gateway.
func (e *EdgeClient) GetFirewallSecurityGroup(nameOrID string) (*FirewallGroupSecurityGroup, error) {
	var (
		values *govcd.NsxtFirewallGroup
		err    error
	)

	if urn.IsSecurityGroup(nameOrID) {
		values, err = e.vcdEdge.GetNsxtFirewallGroupById(nameOrID)
	} else {
		values, err = e.vcdEdge.GetNsxtFirewallGroupByName(nameOrID, govcdtypes.FirewallGroupTypeSecurityGroup)
	}

	if err != nil {
		return nil, err
	}

	return &FirewallGroupSecurityGroup{
		fwGroup: values,
		FirewallGroupSecurityGroupModel: &FirewallGroupSecurityGroupModel{
			FirewallGroupModel: FirewallGroupModel{
				ID:          values.NsxtFirewallGroup.ID,
				Name:        values.NsxtFirewallGroup.Name,
				Description: values.NsxtFirewallGroup.Description,
			},
			Members: values.NsxtFirewallGroup.Members,
		},
		edgeClient: e,
	}, nil
}

// * IPSet

// CreateFirewallIPSet allow creating a new IPSet group.
func (e *EdgeClient) CreateFirewallIPSet(ipSetConfig *FirewallGroupIPSetModel) (*FirewallGroupIPSet, error) {
	ower := &govcdtypes.OpenApiReference{}

	if e.OwnerType.IsVDCGROUP() {
		ower.Name = e.OwnerName
	} else {
		ower.Name = e.vcdEdge.EdgeGateway.Name
		ower.ID = e.vcdEdge.EdgeGateway.ID
	}

	ipSet, err := e.vcdEdge.CreateNsxtFirewallGroup(&govcdtypes.NsxtFirewallGroup{
		Name:        ipSetConfig.Name,
		Description: ipSetConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeIpSet,
		IpAddresses: ipSetConfig.IPAddresses,
		OwnerRef:    ower,
	})
	if err != nil {
		return nil, err
	}

	ipSetConfig.ID = ipSet.NsxtFirewallGroup.ID

	return &FirewallGroupIPSet{
		edgeClient:              e,
		FirewallGroupIPSetModel: ipSetConfig,
		fwGroup:                 ipSet,
	}, nil
}

// GetFirewallIPSet retrieves the IPSet configuration for the Edge Gateway.
func (e *EdgeClient) GetFirewallIPSet(nameOrID string) (*FirewallGroupIPSet, error) {
	var (
		values *govcd.NsxtFirewallGroup
		err    error
	)

	if urn.IsSecurityGroup(nameOrID) {
		values, err = e.vcdEdge.GetNsxtFirewallGroupById(nameOrID)
	} else {
		values, err = e.vcdEdge.GetNsxtFirewallGroupByName(nameOrID, govcdtypes.FirewallGroupTypeIpSet)
	}

	if err != nil {
		return nil, err
	}

	return &FirewallGroupIPSet{
		fwGroup: values,
		FirewallGroupIPSetModel: &FirewallGroupIPSetModel{
			FirewallGroupModel: FirewallGroupModel{
				ID:          values.NsxtFirewallGroup.ID,
				Name:        values.NsxtFirewallGroup.Name,
				Description: values.NsxtFirewallGroup.Description,
			},
			IPAddresses: values.NsxtFirewallGroup.IpAddresses,
		},
		edgeClient: e,
	}, nil
}

// * App Port Profile

// CreateFirewallAppPortProfile allow creating a new application port profile for the Edge Gateway.
func (e *EdgeClient) CreateFirewallAppPortProfile(appPortProfileConfig *FirewallGroupAppPortProfileModel) (*FirewallGroupAppPortProfile, error) {
	vdcOrVDCGroup, err := (&CAVVdc{}).GetVDCOrVDCGroup(e.vcdEdge.EdgeGateway.OwnerRef.Name)
	if err != nil {
		return nil, err
	}

	return createFirewallAppPortProfile(appPortProfileConfig, vdcOrVDCGroup)
}

// GetFirewallAppPortProfile retrieves the application port profile configuration for the VDC Group.
// This function retrieves the application port profile created by the user.
// For retrieving the application port profile created by the system, use FindFirewallAppPortProfile.
func (e *EdgeClient) GetFirewallAppPortProfile(nameOrID string) (*FirewallGroupAppPortProfile, error) {
	vdcOrVDCGroup, err := (&CAVVdc{}).GetVDCOrVDCGroup(e.vcdEdge.EdgeGateway.OwnerRef.Name)
	if err != nil {
		return nil, err
	}

	return getFirewallAppPortProfile(nameOrID, vdcOrVDCGroup)
}

// FindFirewallAppPortProfile retrieves the application port profile configuration for the VDC Group.
// This function retrieves the application port profile created by the user, cloudavenue provider or the system.
func (e *EdgeClient) FindFirewallAppPortProfile(nameOrID string) (*FirewallGroupAppPortProfiles, error) {
	vdcOrVDCGroup, err := (&CAVVdc{}).GetVDCOrVDCGroup(e.vcdEdge.EdgeGateway.OwnerRef.Name)
	if err != nil {
		return nil, err
	}

	return findFirewallAppPortProfile(nameOrID, vdcOrVDCGroup)
}