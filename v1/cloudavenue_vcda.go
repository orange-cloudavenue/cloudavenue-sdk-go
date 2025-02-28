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
	"fmt"

	"golang.org/x/exp/slices"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type (
	VCDA    struct{}
	VDCAIps []string
	VDCAIP  string
)

// List of on premise IP addresses allowed for this organization's draas offer.
func (v *VCDA) List() (VDCAIps, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	r, err := c.R().
		SetResult(&VDCAIps{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/vcda/ips")
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on list VDCA IPs: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return *r.Result().(*VDCAIps), nil
}

// IsIPExists - Returns true if the IP exists.
func (v *VDCAIps) IsIPExists(ip string) bool {
	return slices.Contains(*v, ip)
}

// RegisterIP - Registers a new IP to the list.
func (v *VCDA) RegisterIP(ip string) error {
	c, err := clientcloudavenue.New()
	if err != nil {
		return err
	}

	r, err := c.R().
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("ip", ip).
		Post("/api/customers/v2.0/vcda/ips/{ip}/")
	if err != nil {
		return err
	}

	if r.IsError() {
		return fmt.Errorf("error on add new VDCA IP: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return nil
}

// DeleteIP - Deletes an IP from the list.
func (v *VDCAIps) DeleteIP(ip string) error {
	c, err := clientcloudavenue.New()
	if err != nil {
		return err
	}

	r, err := c.R().
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("ip", ip).
		Delete("/api/customers/v2.0/vcda/ips/{ip}/")
	if err != nil {
		return err
	}

	if r.IsError() {
		return fmt.Errorf("error on delete VDCA IP: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	for i, vIP := range *v {
		if vIP == ip {
			*v = append((*v)[:i], (*v)[i+1:]...)
		}
	}

	return nil
}

// DeleteAllIPs - Deletes all IPs from the list.
func (v *VDCAIps) DeleteAllIPs() error {
	for _, ip := range *v {
		err := v.DeleteIP(ip)
		if err != nil {
			return err
		}
	}

	return nil
}
