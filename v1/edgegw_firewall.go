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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

// GetFirewall retrieves the firewall configuration for the Edge Gateway.
// It first fetches the VMware Edge Gateway associated with the EdgeGw instance,
// and then retrieves the firewall configuration from the Edge Gateway.
// Returns an EdgeGatewayFirewall instance containing the firewall configuration,
// or an error if any step in the process fails.
func (e *EdgeClient) GetFirewall() (*EdgeGatewayFirewall, error) {
	edge, err := e.GetVmwareEdgeGateway()
	if err != nil {
		return nil, err
	}

	fw, err := edge.GetNsxtFirewall()
	if err != nil {
		return nil, err
	}

	return &EdgeGatewayFirewall{
		NsxtFirewall: fw,
	}, nil
}

// UpdateFirewall updates the firewall configuration for the Edge Gateway.
// It first fetches the VMware Edge Gateway associated with the EdgeGw instance,
// and then updates the firewall configuration on the Edge Gateway.
// Returns an error if any step in the process fails.
func (e *EdgeGatewayFirewall) UpdateFirewall(fwRules []*govcdtypes.NsxtFirewallRule) (err error) {
	e.NsxtFirewall, err = e.client.UpdateNsxtFirewall(&govcdtypes.NsxtFirewallRuleContainer{
		UserDefinedRules: fwRules,
	})

	return err
}

// DeleteFirewall deletes the firewall configuration for the Edge Gateway.
func (e *EdgeGatewayFirewall) DeleteFirewall() error {
	fw, err := e.client.GetNsxtFirewall()
	if err != nil {
		return err
	}

	return fw.DeleteAllRules()
}
