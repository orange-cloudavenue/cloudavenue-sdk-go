package org

import "github.com/vmware/go-vcloud-director/v2/govcd"

type (
	CertificateClient struct {
		govcdAdminOrg clientGoVCDAdminOrg
		certVCD       *govcd.Certificate

		// Data
		Certificate CertificateModel
	}

	// CertificatesModel represents a certificate in the certificate library.
	CertificatesModel []*CertificateModel

	// CertificateModel represents a certificate in the certificate library.
	CertificateModel struct {
		ID string

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
)
