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

type AdminOrg struct {
	*govcd.AdminOrg
}

/*
ListCatalogs

Get the catalogs list from the admin org.
*/
func (ao *AdminOrg) ListCatalogs() *govcdtypes.CatalogsList {
	return ao.AdminOrg.AdminOrg.Catalogs
}

// GetOrgVAppLeaseSettings retrieves the lease settings for a vApp in the specified organization.
func (ao *AdminOrg) GetOrgVAppLeaseSettings() *govcdtypes.VAppLeaseSettings {
	return ao.AdminOrg.AdminOrg.OrgSettings.OrgVAppLeaseSettings
}
