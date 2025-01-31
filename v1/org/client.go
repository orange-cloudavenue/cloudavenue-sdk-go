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
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

//go:generate mockgen -source=client.go -destination=zz_generated_client_test.go -self_package github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org -package org -copyright_file "../../mock_header.txt"

type (
	Client interface {
		// * Properties
		GetProperties(ctx context.Context) (client *PropertiesModel, err error)
		UpdateProperties(ctx context.Context, properties *PropertiesRequest) (job *commoncloudavenue.JobCreatedAPIResponse, err error)

		// * Certificates
		ListCertificatesInLibrary(ctx context.Context) (CertificatesModel, error)
		GetCertificateFromLibrary(ctx context.Context, certificateNameOrID string) (*CertificateModel, error)
		CreateCertificateInLibrary(ctx context.Context, cert *CertificateCreateRequest) (*CertificateModel, error)
		UpdateCertificateInLibrary(ctx context.Context, certificateID string, cert *CertificateUpdateRequest) (*CertificateModel, error)
		DeleteCertificateFromLibrary(ctx context.Context, certificateID string) error
	}

	internalClient interface {
		clientGoVCDAdminOrg
		clientCloudavenue
	}

	client struct {
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
		R() *resty.Request
	}
)

// NewClient creates a new Org client.
func NewClient() (Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return &client{
		clientCloudavenue:   c,
		clientGoVCDAdminOrg: c.AdminOrg,
	}, nil
}

// NesFakeClient creates a new fake Org client used for testing.
func NewFakeClient(i internalClient) (Client, error) {
	return &client{
		clientCloudavenue:   i,
		clientGoVCDAdminOrg: i,
	}, nil
}
