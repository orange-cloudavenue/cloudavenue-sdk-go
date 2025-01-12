package v1

type (
	CertificateLibrary struct{}

	// CertificatesLibraryModel represents a certificate library in Cloud Avenue
	CertificatesLibraryModel struct {
		ID string `json:"id,omitempty"` // urn format of the service engine group

		// Name of the certificate
		Name string `json:"name"`

		// Description of the certificate
		Description string `json:"description,omitempty"`

		// Certificate content
		Certificate string `json:"certificate,omitempty"`

		// Private key content
		PrivateKey string `json:"privateKey,omitempty"`

		// Passphrase for the private key
		Passphrase string `json:"passphrase,omitempty"`
	}
)
