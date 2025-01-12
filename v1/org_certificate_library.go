package v1

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// ListCertificateLibrary returns an array of CertificateLibraryModel
func (o *Org) ListCertificateLibrary() ([]*CertificateLibraryModel, error) {
	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	// Find the certificate library
	certificates, err := c.AdminOrg.GetAllCertificatesFromLibrary(url.Values{})
	if err != nil {
		return nil, fmt.Errorf("error while fetching certificate library: %s", err.Error())
	}
	if len(certificates) == 0 {
		return nil, fmt.Errorf("no certificate library found")
	}

	// For x make it in []*CertificateLibraryModel
	x := make([]*CertificateLibraryModel, 0)
	for _, certificate := range certificates {
		x = append(x, &CertificateLibraryModel{
			ID:          certificate.CertificateLibrary.Id,
			Name:        certificate.CertificateLibrary.Alias,
			Description: certificate.CertificateLibrary.Description,
			Certificate: certificate.CertificateLibrary.Certificate,
		})
	}

	return x, nil
}

// GetCertificateLibrary returns a CertificateLibraryModel
func (o *Org) GetCertificateLibrary(nameOrID string) (*CertificateLibraryModel, error) {
	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	var certificate *govcd.Certificate
	// Find the certificate library
	if urn.IsCertificateLibraryItem(nameOrID) {
		certificate, err = c.AdminOrg.GetCertificateFromLibraryById(nameOrID)
	} else {
		certificate, err = c.AdminOrg.GetCertificateFromLibraryByName(nameOrID)
	}
	if err != nil {
		return nil, fmt.Errorf("error while fetching certificate library: %s", err.Error())
	}

	return &CertificateLibraryModel{
		ID:          certificate.CertificateLibrary.Id,
		Name:        certificate.CertificateLibrary.Alias,
		Description: certificate.CertificateLibrary.Description,
		Certificate: certificate.CertificateLibrary.Certificate,
	}, nil
}

// CreateCertificateLibrary creates a CertificateLibraryModel
func (o *Org) CreateCertificateLibrary(certificate *CertificateLibraryModel) (*CertificateLibraryModel, error) {
	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	govcdCert := &govcdtypes.CertificateLibraryItem{
		Alias:       certificate.Name,
		Certificate: certificate.Certificate,
	}
	// Add Description if it is not empty
	if certificate.Description != "" {
		govcdCert.Description = certificate.Description
	}
	// Add PrivateKey if it is not empty
	if certificate.PrivateKey != "" {
		govcdCert.PrivateKey = certificate.PrivateKey
	}
	// Add Passphrase if it is not empty
	if certificate.Passphrase != "" {
		govcdCert.PrivateKeyPassphrase = certificate.Passphrase
	}

	// Create the certificate library
	cert, err := c.AdminOrg.AddCertificateToLibrary(govcdCert)
	if err != nil {
		return nil, fmt.Errorf("error while creating certificate library: %s", err.Error())
	}

	return &CertificateLibraryModel{
		ID:          cert.CertificateLibrary.Id,
		Name:        cert.CertificateLibrary.Alias,
		Description: cert.CertificateLibrary.Description,
		Certificate: cert.CertificateLibrary.Certificate,
		PrivateKey:  cert.CertificateLibrary.PrivateKey,
		Passphrase:  cert.CertificateLibrary.PrivateKeyPassphrase,
	}, nil
}

// Update updates a CertificateLibraryModel
func (cl *CertificateLibraryModel) Update() (*CertificateLibraryModel, error) {
	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	var certificate *govcd.Certificate
	// Find the certificate library
	if urn.IsCertificateLibraryItem(cl.ID) {
		certificate, err = c.AdminOrg.GetCertificateFromLibraryById(cl.ID)
	} else {
		certificate, err = c.AdminOrg.GetCertificateFromLibraryByName(cl.Name)
	}
	if err != nil {
		return nil, err
	}

	// Set the new values
	certificate.CertificateLibrary.Alias = cl.Name
	certificate.CertificateLibrary.Description = cl.Description

	// Update the certificate library
	certificateUpdated, err := certificate.Update()
	if err != nil {
		return nil, err
	}

	return &CertificateLibraryModel{
		ID:          certificateUpdated.CertificateLibrary.Id,
		Name:        certificateUpdated.CertificateLibrary.Alias,
		Description: certificateUpdated.CertificateLibrary.Description,
		Certificate: certificateUpdated.CertificateLibrary.Certificate,
		PrivateKey:  certificateUpdated.CertificateLibrary.PrivateKey,
		Passphrase:  certificateUpdated.CertificateLibrary.PrivateKeyPassphrase,
	}, nil
}

// Delete deletes a CertificateLibraryModel
func (cl *CertificateLibraryModel) Delete() error {
	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return err
	}

	var certificate *govcd.Certificate
	// Find the certificate library
	if urn.IsCertificateLibraryItem(cl.ID) {
		certificate, err = c.AdminOrg.GetCertificateFromLibraryById(cl.ID)
	} else {
		certificate, err = c.AdminOrg.GetCertificateFromLibraryByName(cl.Name)
	}
	if err != nil {
		return err
	}

	// Delete the certificate library
	err = certificate.Delete()
	if err != nil {
		return err
	}

	return nil
}
