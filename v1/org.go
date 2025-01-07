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
)

type Org struct {
	*govcd.Org
}

// Returns the name of the organization.
func (o *Org) GetName() string {
	return o.Org.Org.Name
}

// Returns the ID of the organization.
func (o *Org) GetID() string {
	return o.Org.Org.ID
}

// TODO refacto
// GetNetworkDHCP returns the DHCP object for the org network provided in parameter.
func (o *Org) GetNetworkDHCP(orgNetworkID string) (*govcd.OpenApiOrgVdcNetworkDhcp, error) {
	if err := o.Refresh(); err != nil {
		return nil, err
	}

	orgNetwork, err := o.GetOpenApiOrgVdcNetworkById(orgNetworkID)
	if err != nil {
		return nil, err
	}

	return orgNetwork.GetOpenApiOrgVdcNetworkDhcp()
}

// UpdateNetworkDHCP updates the DHCP object for the org network provided in parameter.
func (o *Org) UpdateNetworkDHCP(orgNetworkID string, dhcpParameters *govcdtypes.OpenApiOrgVdcNetworkDhcp) error {
	orgNetwork, err := o.GetOpenApiOrgVdcNetworkById(orgNetworkID)
	if err != nil {
		return err
	}

	_, err = orgNetwork.UpdateDhcp(dhcpParameters)
	return err
}

// DeleteNetworkDHCP deletes the DHCP object for the org network provided in parameter.
func (o *Org) DeleteNetworkDHCP(orgNetworkID string) error {
	orgNetwork, err := o.GetOpenApiOrgVdcNetworkById(orgNetworkID)
	if err != nil {
		return err
	}

	return orgNetwork.DeletNetworkDhcp()
}
