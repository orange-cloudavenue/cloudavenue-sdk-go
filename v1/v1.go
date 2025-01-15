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

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/iam"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/netbackup"
)

type V1 struct {
	Netbackup   netbackup.Netbackup
	PublicIP    PublicIP
	EdgeGateway EdgeGateway
	T0          Tier0
	VCDA        VCDA
	BMS         BMS
	// VDC         VDC is a method of the V1 struct that returns a pointer to the CAVVdc struct
	// S3          *s3.S3 - S3 is a method of the V1 struct that returns a pointer to the AWS S3 client preconfigured
}

func (v *V1) AdminVDC() *CAVAdminVDC {
	return &CAVAdminVDC{}
}

func (v *V1) VDC() *CAVVdc {
	return &CAVVdc{}
}

func (v *V1) Vmware() (*govcd.VCDClient, error) {
	client, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}
	return client.Vmware, nil
}

func (v *V1) Org() (*Org, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return &Org{
		Org: c.Org,
	}, nil
}

func (v *V1) AdminOrg() (*AdminOrg, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	if c.AdminOrg == nil {
		return nil, fmt.Errorf("admin org is nil")
	}

	return &AdminOrg{
		AdminOrg: c.AdminOrg,
	}, nil
}

func (v *V1) IAM() (*iam.Client, error) {
	return iam.NewClient()
}
