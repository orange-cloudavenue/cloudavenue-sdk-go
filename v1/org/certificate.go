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
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers/validators"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func (c *Client) ListCertificatesInLibrary() (CertificatesModel, error) {
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

func (c *Client) GetCertificateFromLibrary(nameOrID string) (*CertificateClient, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	var (
		certificate *govcd.Certificate
		err         error
	)

	if urn.IsCertificateLibraryItem(nameOrID) {
		certificate, err = c.clientGoVCDAdminOrg.GetCertificateFromLibraryById(nameOrID)
	} else {
		certificate, err = c.clientGoVCDAdminOrg.GetCertificateFromLibraryByName(nameOrID)
	}
	if err != nil {
		return nil, err
	}

	return &CertificateClient{
		govcdAdminOrg: c.clientGoVCDAdminOrg,
		certVCD:       certificate,

		Certificate: CertificateModel{
			ID:          certificate.CertificateLibrary.Id,
			Name:        certificate.CertificateLibrary.Alias,
			Description: certificate.CertificateLibrary.Description,
			Certificate: certificate.CertificateLibrary.Certificate,
		},
	}, nil
}

// CreateCertificateLibrary creates a new certificate library.
func (c *Client) CreateCertificateInLibrary(cert CertificateModel) (*CertificateClient, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	if err := validators.New().Struct(cert); err != nil {
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

	return &CertificateClient{
		govcdAdminOrg: c.clientGoVCDAdminOrg,
		certVCD:       certCreated,

		Certificate: CertificateModel{
			ID:          certCreated.CertificateLibrary.Id,
			Name:        certCreated.CertificateLibrary.Alias,
			Description: certCreated.CertificateLibrary.Description,
			Certificate: certCreated.CertificateLibrary.Certificate,
		},
	}, nil
}

// Update updates a certificate in the library.
// Only the Name and Description can be updated.
func (c *CertificateClient) Update() error {
	if err := validators.New().Struct(c.Certificate); err != nil {
		return err
	}

	c.certVCD.CertificateLibrary.Alias = c.Certificate.Name
	c.certVCD.CertificateLibrary.Description = c.Certificate.Description

	certUpdated, err := c.certVCD.Update()
	if err != nil {
		return fmt.Errorf("error while updating certificate library: %s", err.Error())
	}

	c.Certificate.Name = certUpdated.CertificateLibrary.Alias
	c.Certificate.Description = certUpdated.CertificateLibrary.Description

	return nil
}

// Delete deletes a certificate from the library.
func (c *CertificateClient) Delete() error {
	return c.certVCD.Delete()
}
