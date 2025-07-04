/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package endpoints

// List of API endpoints.
const (
	EdgeGatewayCreateFromVDC      = CloudavenueV2 + "vdcs/{vdc-name}/edges"
	EdgeGatewayCreateFromVDCGroup = CloudavenueV2 + "vdc-groups/{vdc-group-name}/edges"
	EdgeGatewayGet                = CloudavenueV2 + "edges/{edge-id}"
	EdgeGatewayList               = CloudavenueV2 + "edges"
	EdgeGatewayDelete             = EdgeGatewayGet
	EdgeGatewayUpdate             = EdgeGatewayGet

	vmwareEdgeGateway = VmwareV2 + "edgeGateways/"
	FirewallRules     = vmwareEdgeGateway + "{edge-id}/firewall/rules"
	FirewallRule      = vmwareEdgeGateway + "{edge-id}/firewall/rules/{rule-id}"

	NetworkServiceGet    = CloudavenueV2 + "network"
	NetworkServiceCreate = CloudavenueV2 + "services"
	NetworkServiceDelete = CloudavenueV2 + "services/{service-id}"
)
