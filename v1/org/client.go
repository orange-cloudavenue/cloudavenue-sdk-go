/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

//go:generate mockgen -source=client.go -destination=mock/zz_generated_client.go

type (
	Client struct {
		clientGoVCDAdminOrg clientGoVCDAdminOrg
		clientCloudavenue   clientCloudavenue
	}

	clientGoVCDAdminOrg interface {
		GetAllCertificatesFromLibrary(queryParameters url.Values) ([]*govcd.Certificate, error)
		GetCertificateFromLibraryById(id string) (*govcd.Certificate, error)
		GetCertificateFromLibraryByName(name string) (*govcd.Certificate, error)
		AddCertificateToLibrary(certificateConfig *govcdtypes.CertificateLibraryItem) (*govcd.Certificate, error)
	}

	clientCloudavenue interface {
		Refresh() error
	}
)

// NewClient creates a new Org client.
func NewClient() (*Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return &Client{
		clientCloudavenue:   c,
		clientGoVCDAdminOrg: c.AdminOrg,
	}, nil
}
