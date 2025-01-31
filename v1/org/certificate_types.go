/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import "github.com/vmware/go-vcloud-director/v2/govcd"

//go:generate mockgen -source=certificate_types.go -destination=zz_generated_client_certificate_test.go -self_package github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org -package org -copyright_file "../../mock_header.txt"

type (
	internalCertificateClient interface {
		Update() (*govcd.Certificate, error)
		Delete() error
	}

	// CertificatesModel represents a certificate in the certificate library.
	CertificatesModel []*CertificateModel

	CertificateCreateRequest struct {
		// Name of the certificate
		Name string `validate:"required"`

		// Description of the certificate
		Description string `validate:"omitempty"`

		// Certificate content
		Certificate string `validate:"required"`

		// Private key content
		PrivateKey string `validate:"omitempty"`

		// Passphrase for the private key
		Passphrase string `validate:"omitempty"`
	}

	CertificateUpdateRequest struct {
		// Name of the certificate
		Name string `validate:"required"`

		// Description of the certificate
		Description string `validate:"omitempty"`
	}

	// CertificateModel represents a certificate in the certificate library.
	CertificateModel struct {
		ID string `validate:"required"`

		// Name of the certificate
		Name string `validate:"required"`

		// Description of the certificate
		Description string `validate:"omitempty"`

		// Certificate content
		Certificate string `validate:"required"`
	}
)
