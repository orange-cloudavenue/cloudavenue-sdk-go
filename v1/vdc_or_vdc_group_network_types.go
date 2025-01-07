/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/utils"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

type (
	vdcNetworkInterface interface {
		GetID() string
		GetName() string

		getVDCNetworkByID(id string) (*govcd.OpenApiOrgVdcNetwork, error)
		getVDCNetworkByName(name string) (*govcd.OpenApiOrgVdcNetwork, error)
		createVDCNetwork(networkConfig *govcdtypes.OpenApiOrgVdcNetwork) (*govcd.OpenApiOrgVdcNetwork, error)
	}

	vdcNetworkModelInterface interface {
		toVDCNetworkModel(v vdcNetworkInterface) *govcdtypes.OpenApiOrgVdcNetwork
		fromVDCNetworkModel(*govcdtypes.OpenApiOrgVdcNetwork)
	}

	VDCNetwork[T vdcNetworkModelInterface] struct {
		v    vdcNetworkInterface
		net  *govcd.OpenApiOrgVdcNetwork
		data T
	}

	// * VDC Or VDCGroup Isolated Network.

	VDCNetworkIsolated struct {
		VDCNetwork[*VDCNetworkIsolatedModel]
		*VDCNetworkIsolatedModel
	}

	VDCNetworkIsolatedModel = VDCNetworkModel

	// * Common Data structs.
	VDCNetworkModel struct {
		ID                      string                `json:"id,omitempty"`
		Name                    string                `json:"name"`
		Description             string                `json:"description"`
		Status                  string                `json:"status,omitempty"`
		Subnet                  VDCNetworkModelSubnet `json:"subnet"`
		GuestVLANTaggingAllowed *bool                 `json:"guestVlanTaggingAllowed"`
	}

	VDCNetworkModelSubnet struct {
		Gateway      string                        `json:"gateway"`
		PrefixLength int                           `json:"prefixLength"`
		DNSServer1   string                        `json:"dnsServer1"`
		DNSServer2   string                        `json:"dnsServer2"`
		DNSSuffix    string                        `json:"dnsSuffix"`
		IPRanges     VDCNetworkModelSubnetIPRanges `json:"ipRanges"`
	}

	VDCNetworkModelSubnetIPRanges []VDCNetworkModelSubnetIPRange

	VDCNetworkModelSubnetIPRange struct {
		StartAddress string `json:"startAddress"`
		EndAddress   string `json:"endAddress"`
	}
)

func (ipr VDCNetworkModelSubnetIPRanges) ToVcdIPRanges() govcdtypes.OrgVdcNetworkSubnetIPRanges {
	var ipRanges govcdtypes.OrgVdcNetworkSubnetIPRanges
	for _, ipRange := range ipr {
		ipRanges.Values = append(ipRanges.Values, govcdtypes.OpenApiIPRangeValues{
			StartAddress: ipRange.StartAddress,
			EndAddress:   ipRange.EndAddress,
		})
	}
	return ipRanges
}

// Update updates the network.
func (n *VDCNetwork[T]) Update(model T) error {
	net, err := n.net.Update(model.toVDCNetworkModel(n.v))
	if err != nil {
		return err
	}

	n.net = net
	return nil
}

// Delete deletes the network.
func (n *VDCNetwork[T]) Delete() error {
	return n.net.Delete()
}

// * Isolated Network

// ToVDCNetworkModel converts the VDCNetworkIsolated to govcd.OpenApiOrgVdcNetwork.
func (n *VDCNetworkIsolatedModel) toVDCNetworkModel(v vdcNetworkInterface) *govcdtypes.OpenApiOrgVdcNetwork {
	return &govcdtypes.OpenApiOrgVdcNetwork{
		ID:          n.ID,
		Name:        n.Name,
		Description: n.Description,
		NetworkType: govcdtypes.OrgVdcNetworkTypeIsolated,
		OwnerRef: &govcdtypes.OpenApiReference{
			ID:   v.GetID(),
			Name: v.GetName(),
		},
		Subnets: govcdtypes.OrgVdcNetworkSubnets{
			Values: []govcdtypes.OrgVdcNetworkSubnetValues{
				{
					Gateway:      n.Subnet.Gateway,
					PrefixLength: n.Subnet.PrefixLength,
					IPRanges:     n.Subnet.IPRanges.ToVcdIPRanges(),
					DNSServer1:   n.Subnet.DNSServer1,
					DNSServer2:   n.Subnet.DNSServer2,
					DNSSuffix:    n.Subnet.DNSSuffix,
				},
			},
		},
		GuestVlanTaggingAllowed: n.GuestVLANTaggingAllowed,
		Shared: func() *bool {
			if urn.IsVDCGroup(v.GetID()) {
				return utils.ToPTR(true)
			}
			return nil
		}(),
	}
}

// fromVDCNetworkModel converts the govcd.OpenApiOrgVdcNetwork to VDCNetworkIsolated.
func (n *VDCNetworkIsolatedModel) fromVDCNetworkModel(net *govcdtypes.OpenApiOrgVdcNetwork) {
	n.ID = net.ID
	n.Name = net.Name
	n.Description = net.Description
	n.Status = net.Status
	n.GuestVLANTaggingAllowed = net.GuestVlanTaggingAllowed
	n.Subnet = VDCNetworkModelSubnet{
		Gateway:      net.Subnets.Values[0].Gateway,
		PrefixLength: net.Subnets.Values[0].PrefixLength,
		DNSServer1:   net.Subnets.Values[0].DNSServer1,
		DNSServer2:   net.Subnets.Values[0].DNSServer2,
		DNSSuffix:    net.Subnets.Values[0].DNSSuffix,
		IPRanges: func() []VDCNetworkModelSubnetIPRange {
			var ipRanges []VDCNetworkModelSubnetIPRange
			for _, ipRange := range net.Subnets.Values[0].IPRanges.Values {
				ipRanges = append(ipRanges, VDCNetworkModelSubnetIPRange{
					StartAddress: ipRange.StartAddress,
					EndAddress:   ipRange.EndAddress,
				})
			}
			return ipRanges
		}(),
	}
}
