/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package iam

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

const (
	// UserTypeLocal is the type of the user.
	UserTypeLocal UserType = govcd.OrgUserProviderIntegrated
	// UserTypeSAML is the type of the user.
	UserTypeSAML UserType = govcd.OrgUserProviderSAML
)
