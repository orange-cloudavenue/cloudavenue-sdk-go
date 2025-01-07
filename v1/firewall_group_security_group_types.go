/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

import "github.com/vmware/go-vcloud-director/v2/govcd"

type (
	FirewallGroupSecurityGroup struct {
		// vg is a unexported VDC Group Client
		// used only for the vdcGroup
		vg VDCGroup

		// edgeClient is a unexported EdgeGateway Client
		// used only for the EdgeGateway
		edgeClient *EdgeClient

		// fwGroup is a unexported NSX-T Firewall Group
		fwGroup *govcd.NsxtFirewallGroup

		*FirewallGroupSecurityGroupModel
	}

	FirewallGroupIPSet struct {
		// vg is a unexported VDC Group Client
		// used only for the vdcGroup
		vg VDCGroup

		// edgeClient is a unexported EdgeGateway Client
		// used only for the EdgeGateway
		edgeClient *EdgeClient

		// fwGroup is a unexported NSX-T Firewall Group
		fwGroup *govcd.NsxtFirewallGroup

		*FirewallGroupIPSetModel
	}

	FirewallGroupDynamicSecurityGroup struct {
		// vg is a unexported VDC Group Client
		// used only for the vdcGroup
		vg VDCGroup

		// edgeClient is a unexported EdgeGateway Client
		// used only for the EdgeGateway
		edgeClient *EdgeClient

		// fwGroup is a unexported NSX-T Firewall Group
		fwGroup *govcd.NsxtFirewallGroup

		*FirewallGroupDynamicSecurityGroupModel
	}

	FirewallGroupAppPortProfile struct {
		// vdcOrVDCGroup is a unexported EdgeGateway Or VDC Group Interface
		vdcOrVDCGroup idOrNameInterface

		// org is a unexported Org Client
		org *govcd.Org

		// appProfile is a unexported NSX-T Application Port Profile
		appProfile *govcd.NsxtAppPortProfile

		*FirewallGroupAppPortProfileModelResponse
	}

	FirewallGroupAppPortProfiles struct {
		// vdcOrVDCGroup is a unexported EdgeGateway Or VDC Group Interface
		vdcOrVDCGroup idOrNameInterface

		// org is a unexported Org Client
		org *govcd.Org

		AppPortProfiles []*FirewallGroupAppPortProfileModelResponse
	}
)
