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
	"fmt"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func (c *client) ListCertificatesInLibrary(_ context.Context) (CertificatesModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	certificates, err := c.clientGoVCDAdminOrg.GetAllCertificatesFromLibrary(nil)
	if err != nil {
		return nil, err
	}
	if len(certificates) == 0 {
		return nil, fmt.Errorf("no certificates found in the library")
	}

	x := make(CertificatesModel, 0)
	for _, certificate := range certificates {
		x = append(x, &CertificateModel{
			ID:          certificate.CertificateLibrary.Id,
			Name:        certificate.CertificateLibrary.Alias,
			Description: certificate.CertificateLibrary.Description,
			Certificate: certificate.CertificateLibrary.Certificate,
		})
	}

	return x, nil
}

func (c *client) GetCertificateFromLibrary(ctx context.Context, nameOrID string) (*CertificateModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	certificate, err := c.getCertificateFromLibrary(ctx, nameOrID)
	if err != nil {
		return nil, err
	}

	return &CertificateModel{
		ID:          certificate.CertificateLibrary.Id,
		Name:        certificate.CertificateLibrary.Alias,
		Description: certificate.CertificateLibrary.Description,
		Certificate: certificate.CertificateLibrary.Certificate,
	}, nil
}

func (c *client) getCertificateFromLibrary(_ context.Context, nameOrID string) (*govcd.Certificate, error) {
	var (
		certificate *govcd.Certificate
		err         error
	)

	if urn.IsCertificateLibraryItem(nameOrID) {
		certificate, err = c.clientGoVCDAdminOrg.GetCertificateFromLibraryById(nameOrID)
	} else {
		certificate, err = c.clientGoVCDAdminOrg.GetCertificateFromLibraryByName(nameOrID)
	}

	return certificate, err
}

// CreateCertificateLibrary creates a new certificate library.
func (c *client) CreateCertificateInLibrary(ctx context.Context, cert *CertificateCreateRequest) (*CertificateModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	if err := validators.New().StructCtx(ctx, cert); err != nil {
		return nil, err
	}

	govcdCert := &govcdtypes.CertificateLibraryItem{
		Alias:                cert.Name,
		Certificate:          cert.Certificate,
		Description:          cert.Description,
		PrivateKey:           cert.PrivateKey,
		PrivateKeyPassphrase: cert.Passphrase,
	}

	// Create the certificate library
	certCreated, err := c.clientGoVCDAdminOrg.AddCertificateToLibrary(govcdCert)
	if err != nil {
		return nil, fmt.Errorf("error while creating certificate library: %s", err.Error())
	}

	return &CertificateModel{
		ID:          certCreated.CertificateLibrary.Id,
		Name:        certCreated.CertificateLibrary.Alias,
		Description: certCreated.CertificateLibrary.Description,
		Certificate: certCreated.CertificateLibrary.Certificate,
	}, nil
}

// UpdateCertificateInLibrary updates a certificate in the library.
func (c *client) UpdateCertificateInLibrary(ctx context.Context, certificateID string, cert *CertificateUpdateRequest) (*CertificateModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	if err := validators.New().StructCtx(ctx, cert); err != nil {
		return nil, err
	}

	certificate, err := c.getCertificateFromLibrary(ctx, certificateID)
	if err != nil {
		return nil, err
	}

	certificate.CertificateLibrary.Alias = cert.Name
	certificate.CertificateLibrary.Description = cert.Description

	certUpdated, err := updateCertificateInLibrary(certificate)
	if err != nil {
		return nil, fmt.Errorf("error while updating certificate library: %s", err.Error())
	}

	return &CertificateModel{
		ID:          certUpdated.CertificateLibrary.Id,
		Name:        certUpdated.CertificateLibrary.Alias,
		Description: certUpdated.CertificateLibrary.Description,
		Certificate: certUpdated.CertificateLibrary.Certificate,
	}, nil
}

var updateCertificateInLibrary = func(cert internalCertificateClient) (*govcd.Certificate, error) {
	return cert.Update()
}

// DeleteCertificateFromLibrary deletes a certificate from the library.
func (c *client) DeleteCertificateFromLibrary(ctx context.Context, certificateID string) error {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	certificate, err := c.getCertificateFromLibrary(ctx, certificateID)
	if err != nil {
		return err
	}

	return deleteCertificateFromLibrary(certificate)
}

var deleteCertificateFromLibrary = func(cert internalCertificateClient) error {
	return cert.Delete()
}
