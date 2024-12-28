package v1

import (
	"fmt"
	"strings"

	"github.com/avast/retry-go/v4"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

var _ FirewallGroupInterface = (*VDCGroup)(nil)

// * Security Group

// CreateFirewallSecurityGroup allow creating a new security group for the VDC Group.
func (g VDCGroup) CreateFirewallSecurityGroup(securityGroupConfig *FirewallGroupSecurityGroupModel) (*FirewallGroupSecurityGroup, error) {
	v, err := g.vg.CreateNsxtFirewallGroup(&govcdtypes.NsxtFirewallGroup{
		Name:        securityGroupConfig.Name,
		Description: securityGroupConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeSecurityGroup,
		Members:     securityGroupConfig.Members,
		OwnerRef: &govcdtypes.OpenApiReference{
			ID:   g.GetID(),
			Name: g.GetName(),
		},
	})
	if err != nil {
		return nil, err
	}

	securityGroupConfig.ID = v.NsxtFirewallGroup.ID

	return &FirewallGroupSecurityGroup{
		fwGroup:                         v,
		FirewallGroupSecurityGroupModel: securityGroupConfig,
		vg:                              g,
	}, nil
}

func (g VDCGroup) genericGetFirewallSecurityGroup(nameOrID string) (*govcd.NsxtFirewallGroup, error) {
	var values *govcd.NsxtFirewallGroup

	if nameOrID == "" {
		return nil, fmt.Errorf("the name or ID must be provided")
	}

	err := retry.Do(
		func() error {
			var err error
			if urn.IsSecurityGroup(nameOrID) {
				values, err = g.vg.GetNsxtFirewallGroupById(nameOrID)
			} else {
				values, err = g.vg.GetNsxtFirewallGroupByName(nameOrID, govcdtypes.FirewallGroupTypeSecurityGroup)
			}

			return err
		},
		retry.RetryIf(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "could not find NSX-T Firewall Group")
		}),
		retry.Attempts(5),
	)

	return values, err
}

// GetFirewallSecurityGroup retrieves the security group configuration for the VDC Group.
func (g VDCGroup) GetFirewallSecurityGroup(nameOrID string) (*FirewallGroupSecurityGroup, error) {
	vv, err := g.genericGetFirewallSecurityGroup(nameOrID)
	if err != nil {
		return nil, err
	}

	return &FirewallGroupSecurityGroup{
		fwGroup: vv,
		FirewallGroupSecurityGroupModel: &FirewallGroupSecurityGroupModel{
			FirewallGroupModel: FirewallGroupModel{
				ID:          vv.NsxtFirewallGroup.ID,
				Name:        vv.NsxtFirewallGroup.Name,
				Description: vv.NsxtFirewallGroup.Description,
			},
			Members: vv.NsxtFirewallGroup.Members,
		},
		vg: g,
	}, nil
}

// * IP Set

// CreateFirewallIPSet allow creating a new IP set for the VDC Group.
func (g VDCGroup) CreateFirewallIPSet(ipSetConfig *FirewallGroupIPSetModel) (*FirewallGroupIPSet, error) {
	v, err := g.vg.CreateNsxtFirewallGroup(&govcdtypes.NsxtFirewallGroup{
		Name:        ipSetConfig.Name,
		Description: ipSetConfig.Description,
		TypeValue:   govcdtypes.FirewallGroupTypeIpSet,
		IpAddresses: ipSetConfig.IPAddresses,
		OwnerRef: &govcdtypes.OpenApiReference{
			ID:   g.GetID(),
			Name: g.GetName(),
		},
	})
	if err != nil {
		return nil, err
	}

	ipSetConfig.ID = v.NsxtFirewallGroup.ID

	return &FirewallGroupIPSet{
		fwGroup:                 v,
		FirewallGroupIPSetModel: ipSetConfig,
		vg:                      g,
	}, nil
}

func (g VDCGroup) genericGetFirewallIPSet(nameOrID string) (*govcd.NsxtFirewallGroup, error) {
	var values *govcd.NsxtFirewallGroup

	if nameOrID == "" {
		return nil, fmt.Errorf("the name or ID must be provided")
	}

	err := retry.Do(
		func() error {
			var err error
			if urn.IsSecurityGroup(nameOrID) {
				values, err = g.vg.GetNsxtFirewallGroupById(nameOrID)
			} else {
				values, err = g.vg.GetNsxtFirewallGroupByName(nameOrID, govcdtypes.FirewallGroupTypeIpSet)
			}

			return err
		},
		retry.RetryIf(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "could not find NSX-T Firewall Group")
		}),
		retry.Attempts(5),
	)

	return values, err
}

// GetFirewallIPSet retrieves the IP set configuration for the VDC Group.
func (g VDCGroup) GetFirewallIPSet(nameOrID string) (*FirewallGroupIPSet, error) {
	v, err := g.genericGetFirewallIPSet(nameOrID)
	if err != nil {
		return nil, err
	}

	return &FirewallGroupIPSet{
		fwGroup: v,
		FirewallGroupIPSetModel: &FirewallGroupIPSetModel{
			FirewallGroupModel: FirewallGroupModel{
				ID:          v.NsxtFirewallGroup.ID,
				Name:        v.NsxtFirewallGroup.Name,
				Description: v.NsxtFirewallGroup.Description,
			},
			IPAddresses: v.NsxtFirewallGroup.IpAddresses,
		},
		vg: g,
	}, nil
}

// * Dynamic Security Group

// CreateFirewallDynamicSecurityGroup allow creating a new dynamic security group for the VDC Group.
func (g VDCGroup) CreateFirewallDynamicSecurityGroup(dynamicSecurityGroupConfig *FirewallGroupDynamicSecurityGroupModel) (*FirewallGroupDynamicSecurityGroup, error) {
	if dynamicSecurityGroupConfig == nil {
		return nil, fmt.Errorf("dynamicSecurityGroupConfig is nil")
	}

	if err := dynamicSecurityGroupConfig.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	v, err := g.vg.CreateNsxtFirewallGroup(dynamicSecurityGroupConfig.toGovcdtypesNsxtFirewallGroup(g.GetID(), g.GetName()))
	if err != nil {
		return nil, err
	}

	dynamicSecurityGroupConfig.fromGovcdtypesNsxtFirewallGroup(v.NsxtFirewallGroup)

	return &FirewallGroupDynamicSecurityGroup{
		fwGroup:                                v,
		FirewallGroupDynamicSecurityGroupModel: dynamicSecurityGroupConfig,
		vg:                                     g,
	}, nil
}

func (g VDCGroup) genericGetFirewallDynamicSecurityGroup(nameOrID string) (*govcd.NsxtFirewallGroup, error) {
	var values *govcd.NsxtFirewallGroup

	if nameOrID == "" {
		return nil, fmt.Errorf("the name or ID must be provided")
	}

	err := retry.Do(
		func() error {
			var err error
			if urn.IsSecurityGroup(nameOrID) {
				values, err = g.vg.GetNsxtFirewallGroupById(nameOrID)
			} else {
				values, err = g.vg.GetNsxtFirewallGroupByName(nameOrID, govcdtypes.FirewallGroupTypeVmCriteria)
			}

			return err
		},
		retry.RetryIf(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "could not find NSX-T Firewall Group")
		}),
		retry.Attempts(5),
	)

	return values, err
}

// GetFirewallDynamicSecurityGroup retrieves the dynamic security group configuration for the VDC Group.
func (g VDCGroup) GetFirewallDynamicSecurityGroup(nameOrID string) (*FirewallGroupDynamicSecurityGroup, error) {
	v, err := g.genericGetFirewallDynamicSecurityGroup(nameOrID)
	if err != nil {
		return nil, err
	}

	x := &FirewallGroupDynamicSecurityGroupModel{}
	x.fromGovcdtypesNsxtFirewallGroup(v.NsxtFirewallGroup)

	return &FirewallGroupDynamicSecurityGroup{
		fwGroup:                                v,
		FirewallGroupDynamicSecurityGroupModel: x,
		vg:                                     g,
	}, nil
}

// * App Port Profile

// CreateFirewallAppPortProfile allow creating a new application port profile for the VDC Group.
func (g *VDCGroup) CreateFirewallAppPortProfile(appPortProfileConfig *FirewallGroupAppPortProfileModel) (*FirewallGroupAppPortProfile, error) {
	return createFirewallAppPortProfile(appPortProfileConfig, g)
}

// GetFirewallAppPortProfile retrieves the application port profile configuration for the VDC Group.
// This function retrieves the application port profile created by the user.
// For retrieving the application port profile created by the system, use FindFirewallAppPortProfile.
func (g *VDCGroup) GetFirewallAppPortProfile(nameOrID string) (*FirewallGroupAppPortProfile, error) {
	return getFirewallAppPortProfile(nameOrID, g)
}

// FindFirewallAppPortProfile retrieves the application port profile configuration for the VDC Group.
// This function retrieves the application port profile created by the user, cloudavenue provider or the system.
func (g *VDCGroup) FindFirewallAppPortProfile(nameOrID string) (*FirewallGroupAppPortProfiles, error) {
	return findFirewallAppPortProfile(nameOrID, g)
}
