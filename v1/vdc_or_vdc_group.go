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
	"errors"

	"github.com/vmware/go-vcloud-director/v2/govcd"
)

// *
// * VDCOrVDCGroup
// *

// GetVDCOrVDCGroup returns the VDC or VDC Group by its name.
// It returns a pointer to the VDC or VDC Group and an error if any.
func (v *CAVVdc) GetVDCOrVDCGroup(vdcOrVDCGroupName string) (VDCOrVDCGroupInterface, error) {
	errs := []error{}

	xVDCGroup, err := v.GetVDCGroup(vdcOrVDCGroupName)
	if err != nil {
		if !govcd.ContainsNotFound(err) {
			errs = append(errs, err)
		}
	} else {
		return xVDCGroup, nil
	}

	xVDC, err := v.GetVDC(vdcOrVDCGroupName)
	if err != nil {
		if !govcd.ContainsNotFound(err) {
			errs = append(errs, err)
		}
	} else {
		return xVDC, nil
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return nil, ErrRetrievingVDCOrVDCGroup
}
