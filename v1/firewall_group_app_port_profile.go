package v1

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// findFirewallAppPortProfile searches for a firewall application port profile by name or ID within a specified edge gateway or VDC group.
// It returns a FirewallGroupAppPortProfiles struct containing the found profiles or an error if no profile is found.
//
// Parameters:
// - nameOrID: A string representing the name or ID of the application port profile to search for. This parameter is required.
// - vdcOrVDCGroup: An interface representing the edge gateway or VDC group within which to search for the application port profile.
//
// Returns:
// - A pointer to a FirewallGroupAppPortProfiles struct containing the found application port profiles.
// - An error if no application port profile is found or if multiple profiles with the same name are found.
func findFirewallAppPortProfile(nameOrID string, vdcOrVDCGroup idOrNameInterface) (*FirewallGroupAppPortProfiles, error) {
	if nameOrID == "" {
		return nil, fmt.Errorf("the name or ID must be provided")
	}

	appProfiles := make([]*FirewallGroupAppPortProfileModelResponse, 0)

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	if urn.IsAppPortProfile(nameOrID) {
		app, err := c.Org.GetNsxtAppPortProfileById(nameOrID)
		if err != nil {
			return nil, err
		}

		x := &FirewallGroupAppPortProfileModelResponse{}
		x.fromGovcdtypesNsxtAppPortProfile(app.NsxtAppPortProfile)
		x.Scope = FirewallGroupAppPortProfileModelScope(app.NsxtAppPortProfile.Scope)

		appProfiles = append(appProfiles, x)
	} else {
		scopes := []FirewallGroupAppPortProfileModelScope{FirewallGroupAppPortProfileModelScopeTenant, FirewallGroupAppPortProfileModelScopeProvider, FirewallGroupAppPortProfileModelScopeSystem}

		for _, scope := range scopes {
			queryParams := url.Values{}
			queryParams.Add("filter", fmt.Sprintf("name==%s;scope==%s;_context==%s", nameOrID, string(scope), vdcOrVDCGroup.GetID()))
			appPortProfiles, _ := c.Org.GetAllNsxtAppPortProfiles(queryParams, "")
			// Error is ignored because we want to continue searching in other scopes if not found
			if len(appPortProfiles) > 0 {
				if len(appPortProfiles) > 1 {
					return nil, fmt.Errorf("found multiple application port profiles with the same name %s", nameOrID)
				}

				x := &FirewallGroupAppPortProfileModelResponse{}
				x.fromGovcdtypesNsxtAppPortProfile(appPortProfiles[0].NsxtAppPortProfile)
				x.Scope = FirewallGroupAppPortProfileModelScope(appPortProfiles[0].NsxtAppPortProfile.Scope)

				appProfiles = append(appProfiles, x)
			}
		}
	}

	// If no app port profile is found, error not found
	if len(appProfiles) == 0 {
		return nil, fmt.Errorf("application port profile %w with the name or ID %s", errors.ErrNotFound, nameOrID)
	}

	return &FirewallGroupAppPortProfiles{
		vdcOrVDCGroup:   vdcOrVDCGroup,
		org:             c.Org,
		AppPortProfiles: appProfiles,
	}, nil
}

// getFirewallAppPortProfile retrieves a FirewallGroupAppPortProfile by its name or ID.
// It takes a string `nameOrID` which can be either the name or ID of the profile,
// and an `edgeGatewayOrVDCGroup` which is an interface representing the edge gateway or VDC group.
//
// If the `nameOrID` is empty, it returns an error indicating that the name or ID must be provided.
// It initializes a new client and attempts to retrieve the NsxtAppPortProfile either by ID or by name.
// If the retrieval is successful, it converts the profile to a FirewallGroupAppPortProfileModelResponse
// and returns a populated FirewallGroupAppPortProfile struct.
//
// Parameters:
// - nameOrID: string representing the name or ID of the application port profile.
// - vdcOrVDCGroup: idOrNameInterface representing the edge gateway or VDC group.
//
// Returns:
// - *FirewallGroupAppPortProfile: a pointer to the retrieved FirewallGroupAppPortProfile.
// - error: an error if the profile could not be retrieved or if the nameOrID is empty.
func getFirewallAppPortProfile(nameOrID string, vdcOrVDCGroup idOrNameInterface) (*FirewallGroupAppPortProfile, error) {
	if nameOrID == "" {
		return nil, fmt.Errorf("the name or ID must be provided")
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	var app *govcd.NsxtAppPortProfile

	if urn.IsAppPortProfile(nameOrID) {
		app, err = c.Org.GetNsxtAppPortProfileById(nameOrID)
	} else {
		app, err = c.Org.GetNsxtAppPortProfileByName(nameOrID, govcdtypes.ApplicationPortProfileScopeTenant)
	}

	if err != nil {
		return nil, err
	}

	x := &FirewallGroupAppPortProfileModelResponse{}
	x.fromGovcdtypesNsxtAppPortProfile(app.NsxtAppPortProfile)
	x.Scope = FirewallGroupAppPortProfileModelScope(app.NsxtAppPortProfile.Scope)

	return &FirewallGroupAppPortProfile{
		appProfile:                               app,
		FirewallGroupAppPortProfileModelResponse: x,
		vdcOrVDCGroup:                            vdcOrVDCGroup,
		org:                                      c.Org,
	}, nil
}

// createFirewallAppPortProfile creates a new firewall application port profile based on the provided configuration.
// It validates the configuration, initializes a new client, and creates the NSX-T application port profile.
//
// Parameters:
// - appPortProfileConfig: A pointer to FirewallGroupAppPortProfileModel containing the configuration for the application port profile.
// - vdcOrVDCGroup: An interface representing either an edge gateway or a VDC group.
//
// Returns:
// - A pointer to FirewallGroupAppPortProfile containing the created application port profile and associated metadata.
// - An error if the configuration is nil, validation fails, or the creation process encounters an issue.
func createFirewallAppPortProfile(appPortProfileConfig *FirewallGroupAppPortProfileModel, vdcOrVDCGroup idOrNameInterface) (*FirewallGroupAppPortProfile, error) {
	if appPortProfileConfig == nil {
		return nil, fmt.Errorf("appPortProfileConfig is nil")
	}

	if err := appPortProfileConfig.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	v, err := c.Org.CreateNsxtAppPortProfile(appPortProfileConfig.toGovcdtypesNsxtAppPortProfile(c.Org.Org.ID, vdcOrVDCGroup.GetID()))
	if err != nil {
		return nil, err
	}

	appPortProfileResponse := &FirewallGroupAppPortProfileModelResponse{}
	appPortProfileResponse.fromGovcdtypesNsxtAppPortProfile(v.NsxtAppPortProfile)
	appPortProfileResponse.Scope = FirewallGroupAppPortProfileModelScope(v.NsxtAppPortProfile.Scope)

	return &FirewallGroupAppPortProfile{
		appProfile:                               v,
		FirewallGroupAppPortProfileModelResponse: appPortProfileResponse,
		vdcOrVDCGroup:                            vdcOrVDCGroup,
		org:                                      c.Org,
	}, nil
}
