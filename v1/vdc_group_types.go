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
)

type (
	VDCGroup struct {
		// vg is a unexported VDC Group Client
		vg *govcd.VdcGroup

		// VdcGroup is a exported client for VDC Group
		VDCGroupInterface
	}

	VDCGroupInterface interface {
		GetOpenApiOrgVdcNetworkByName(string) (*govcd.OpenApiOrgVdcNetwork, error)
	}
)
