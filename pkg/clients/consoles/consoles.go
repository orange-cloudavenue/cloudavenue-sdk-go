package consoles

import "regexp"

type (
	Console      string
	LocationCode string

	console struct {
		SiteName            string
		LocationCode        LocationCode
		SiteID              Console
		URL                 string
		Services            services
		OrganizationPattern *regexp.Regexp
	}

	services struct {
		S3   service
		VCDA service
	}

	service struct {
		Enabled  bool
		Endpoint string
	}
)

const (
	Console1 Console = "console1" // Externe VDR
	Console2 Console = "console2" // Internal
	Console4 Console = "console4" // Externe CHA
	// Console5 Console = "console5" // Internal

	LocationVDR LocationCode = "vdr"
	LocationCHR LocationCode = "chr"
)

var consoles = map[Console]console{
	Console1: {
		SiteName:            "Console Externe",
		LocationCode:        LocationVDR,
		SiteID:              Console1,
		URL:                 "https://console1.cloudavenue.orange-business.com",
		OrganizationPattern: regexp.MustCompile(`^cav01ev01ocb\d{7}$`),
		Services: services{
			S3: service{
				Enabled:  true,
				Endpoint: "https://s3console1.cloudavenue.orange-business.com",
			},
		},
	},
	Console2: {
		SiteName:            "Console Interne",
		LocationCode:        LocationVDR,
		SiteID:              Console2,
		URL:                 "https://console2.cloudavenue.orange-business.com",
		OrganizationPattern: regexp.MustCompile(`^cav01iv02ocb\d{7}$`),
		Services: services{
			S3: service{
				Enabled:  true,
				Endpoint: "https://s3console2.cloudavenue.orange-business.com",
			},
		},
	},
	Console4: {
		SiteName:            "Console Externe",
		LocationCode:        LocationCHR,
		SiteID:              Console4,
		URL:                 "https://console4.cloudavenue.orange-business.com",
		OrganizationPattern: regexp.MustCompile(`^cav02ev04ocb\d{7}$`),
	},
}

// FindBySiteID - Returns the console by its siteID
func FindBySiteID(siteID string) (Console, bool) {
	for c, console := range consoles {
		if console.SiteID == Console(siteID) {
			return c, true
		}
	}

	return "", false
}

// FindByURL - Returns the console by its URL
func FindByURL(url string) (Console, bool) {
	for c, console := range consoles {
		if console.URL == url {
			return c, true
		}
	}

	return "", false
}

// FingByOrganizationName - Returns the console by its organization name
func FingByOrganizationName(organizationName string) (Console, error) {
	for c, console := range consoles {
		if console.OrganizationPattern.MatchString(organizationName) {
			return c, nil
		}
	}

	return "", ErrOrganizationFormatIsInvalid
}

// GetSiteName - Returns the site name
func (c Console) GetSiteName() string {
	return consoles[c].SiteName
}

// GetLocationCode - Returns the location code
func (c Console) GetLocationCode() LocationCode {
	return consoles[c].LocationCode
}

// GetSiteID - Returns the site ID
func (c Console) GetSiteID() Console {
	return consoles[c].SiteID
}

// GetURL - Returns the URL
func (c Console) GetURL() string {
	return consoles[c].URL
}

// S3IsEnabled - Returns true if the S3 service is enabled
func (c Console) S3IsEnabled() bool {
	return consoles[c].Services.S3.Enabled
}

// S3GetEndpoint - Returns the S3 endpoint
func (c Console) GetS3Endpoint() string {
	return consoles[c].Services.S3.Endpoint
}

// VCDAIsEnabled - Returns true if the VCDA service is enabled
func (c Console) IsVCDAEnabled() bool {
	return consoles[c].Services.VCDA.Enabled
}

// VCDAGetEndpoint - Returns the VCDA endpoint
func (c Console) GetVCDAEndpoint() string {
	return consoles[c].Services.VCDA.Endpoint
}
