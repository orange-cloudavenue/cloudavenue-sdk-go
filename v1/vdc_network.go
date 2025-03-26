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
	"github.com/avast/retry-go/v4"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// * Isolated Network

// GetNetworkIsolated returns the isolated network by its name or ID.
func (v *VDC) GetNetworkIsolated(nameOrID string) (*VDCNetworkIsolated, error) {
	net, err := v.genericGetNetwork(nameOrID)
	if err != nil {
		return nil, err
	}

	x := &VDCNetworkIsolated{
		VDCNetwork: VDCNetwork[*VDCNetworkIsolatedModel]{
			v:    v,
			net:  net,
			data: &VDCNetworkIsolatedModel{},
		},
	}

	x.data.fromVDCNetworkModel(net.OpenApiOrgVdcNetwork)
	x.VDCNetworkIsolatedModel = x.data
	return x, nil
}

// CreateNetworkIsolated creates an isolated network.
func (v *VDC) CreateNetworkIsolated(model *VDCNetworkIsolatedModel) (*VDCNetworkIsolated, error) {
	net, err := v.createVDCNetwork(model.toVDCNetworkModel(v, govcdtypes.OrgVdcNetworkTypeIsolated))
	if err != nil {
		return nil, err
	}

	x := &VDCNetworkIsolated{
		VDCNetwork: VDCNetwork[*VDCNetworkIsolatedModel]{
			v:    v,
			net:  net,
			data: &VDCNetworkIsolatedModel{},
		},
	}

	x.data.fromVDCNetworkModel(net.OpenApiOrgVdcNetwork)
	x.VDCNetworkIsolatedModel = x.data
	return x, nil
}

// * Routed Network

// GetNetworkRouted returns the routed network by its name or ID.
func (v *VDC) GetNetworkRouted(nameOrID string) (*VDCNetworkRouted, error) {
	net, err := v.genericGetNetwork(nameOrID)
	if err != nil {
		return nil, err
	}

	x := &VDCNetworkRouted{
		VDCNetwork: VDCNetwork[*VDCNetworkRoutedModel]{
			v:    v,
			net:  net,
			data: &VDCNetworkRoutedModel{},
		},
	}

	x.data.fromVDCNetworkModel(net.OpenApiOrgVdcNetwork)
	x.VDCNetworkRoutedModel = x.data
	return x, nil
}

// CreateNetworkRouted creates a routed network.
func (v *VDC) CreateNetworkRouted(model *VDCNetworkRoutedModel) (*VDCNetworkRouted, error) {
	net, err := v.createVDCNetwork(model.toVDCNetworkModel(v, govcdtypes.OrgVdcNetworkTypeRouted))
	if err != nil {
		return nil, err
	}

	x := &VDCNetworkRouted{
		VDCNetwork: VDCNetwork[*VDCNetworkRoutedModel]{
			v:    v,
			net:  net,
			data: &VDCNetworkRoutedModel{},
		},
	}

	x.data.fromVDCNetworkModel(net.OpenApiOrgVdcNetwork)
	x.VDCNetworkRoutedModel = x.data
	return x, nil
}

func (v VDC) genericGetNetwork(nameOrID string) (*govcd.OpenApiOrgVdcNetwork, error) {
	var values *govcd.OpenApiOrgVdcNetwork

	err := retry.Do(
		func() error {
			var err error
			if urn.IsNetwork(nameOrID) {
				values, err = v.getVDCNetworkByID(nameOrID)
			} else {
				values, err = v.getVDCNetworkByName(nameOrID)
			}

			return err
		},
		retry.RetryIf(govcd.ContainsNotFound),
		retry.Attempts(5),
	)

	return values, err
}

// IsRouted returns true if the network is routed.
func (v VDCNetwork[T]) IsRouted() bool {
	return v.net.IsRouted()
}

// IsIsolated returns true if the network is isolated.
func (v VDCNetwork[T]) IsIsolated() bool {
	return v.net.IsIsolated()
}

// IsShared returns true if the network is shared.
func (v VDCNetwork[T]) IsShared() bool {
	if v.net.OpenApiOrgVdcNetwork.Shared == nil {
		return false
	}
	return *v.net.OpenApiOrgVdcNetwork.Shared
}
