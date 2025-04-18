package edgegateway

type (
	ServiceModel struct {
		IPAllocated string

		ServiceDetails []ServiceModelDetails
	}

	ServiceModelDetails struct {
		// Category is the category of the service
		Category string
		// Network is the network of the service
		Network string
		// Services is the list of services
		Services []ServiceModelDetail
	}

	ServiceModelDetail struct {
		// Name is the name of the service
		Name string
		// Description
		Description string
		// DocumentationURL is the URL of the documentation
		DocumentationURL string
		// IP is the IP address of the service
		IP []string
		// FQDN is the FQDN of the service
		FQDN  []string
		Ports []NetworkServicesModelSvcServiceDetailsPorts
	}
)

var listOfServices = []ServiceModelDetails{
	{
		Category: "administration",
		Network:  "57.199.209.192/27",
		Services: []ServiceModelDetail{
			{
				Name:        "linux-repository",
				Description: "Linux (Debian, Ubuntu, CentOS) package repository",
				IP:          []string{"57.199.209.214"},
				FQDN:        []string{"repo.service.cav"},
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     3142,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "rhui-repository",
				Description: "Red Hat (RHUI) package repository",
				IP:          []string{"57.199.209.197"},
				FQDN:        []string{"rhui.service.cav"},
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     80,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "windows-repository",
				Description: "Windows (WSUS) package repository",
				IP:          []string{"57.199.209.212"},
				FQDN:        []string{"wsus.service.cav"},
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     80,
						Protocol: "tcp",
					},
					{
						Port:     443,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "windows-kms",
				Description: "Windows (KMS) license server",
				IP: []string{
					"57.199.209.210",
				},
				FQDN: []string{"kms.service.cav"},
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     1688,
						Protocol: "tcp",
					},
				},
			},
			{
				Name:        "ntp",
				Description: "Network Time Protocol (NTP) server",
				IP: []string{
					"57.199.209.217",
					"57.199.209.218",
				},
				FQDN: []string{
					"ntp1.service.cav",
					"ntp2.service.cav",
				},
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     123,
						Protocol: "udp",
					},
				},
			},
			{
				Name:             "dns-authoritative",
				Description:      "DNS authoritative server. Use for resolving cloudavenue services names",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/service-zone-dns/",
				IP: []string{
					"57.199.209.207",
					"57.199.209.208",
				},
				FQDN: nil,
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     53,
						Protocol: "tcp",
					},
					{
						Port:     53,
						Protocol: "udp",
					},
				},
			},
			{
				Name:             "dns-resolver",
				Description:      "DNS resolver. Use for resolving cloudavenue services names and public names",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/service-zone-dns/",
				IP: []string{
					"57.199.209.220",
					"57.199.209.221",
				},
				FQDN: nil,
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     53,
						Protocol: "tcp",
					},
					{
						Port:     53,
						Protocol: "udp",
					},
				},
			},
			{
				Name:             "smtp",
				Description:      "SMTP relay. Use for sending emails",
				DocumentationURL: "https://cloud.orange-business.com/en/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/practical-sheets/services-area/services-en/smtp-service-2/",
				IP: []string{
					"57.199.209.206",
				},
				FQDN: []string{"smtp.service.cav"},
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     25,
						Protocol: "tcp",
					},
				},
			},
		},
	},
	{
		Category: "s3",
		Network:  "194.206.55.5/32",
		Services: []ServiceModelDetail{
			{
				Name:             "s3-internal",
				Description:      "S3 internal service. Use for accessing S3 directly from the organization",
				DocumentationURL: "https://cloud.orange-business.com/offres/infrastructure-iaas/cloud-avenue/wiki-cloud-avenue/fiches-pratiques/stockage/stockage-objet-s3/guide-de-demarrage/premiere-utilisation-stockage-objet/",
				IP: []string{
					"194.206.55.5",
				},
				FQDN: []string{"s3-region01-priv.cloudavenue.orange-business.com"},
				Ports: []NetworkServicesModelSvcServiceDetailsPorts{
					{
						Port:     443,
						Protocol: "tcp",
					},
				},
			},
		},
	},
}
