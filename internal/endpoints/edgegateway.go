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
	EdgeGatewayCreateFromVDC      = "/infrapicustomerproxy/v2.0/vdcs/{vdc-name}/edges"
	EdgeGatewayCreateFromVDCGroup = "/infrapicustomerproxy/v2.0/vdc-groups/{vdc-group-name}/edges"
	EdgeGatewayGet                = "/infrapicustomerproxy/v2.0/edges/{edge-id}"
	EdgeGatewayList               = "/infrapicustomerproxy/v2.0/edges"
	EdgeGatewayDelete             = EdgeGatewayGet
	EdgeGatewayUpdate             = EdgeGatewayGet

	NetworkServiceGet    = "/infrapicustomerproxy/v2.0/network"
	NetworkServiceCreate = "/infrapicustomerproxy/v2.0/services"
	NetworkServiceDelete = "/infrapicustomerproxy/v2.0/services/{service-id}"
)
