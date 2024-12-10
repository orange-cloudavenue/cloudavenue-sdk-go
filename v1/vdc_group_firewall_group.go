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
