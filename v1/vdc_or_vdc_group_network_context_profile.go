package v1

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/uuid"
)

// getNetworkContextProfile retrieves a network context profile
func getNetworkContextProfile(networkContextProfileIDOrName, vdcOrVDCGroupID string, scope VDCOrVDCGroupNetworkContextProfileScope) (*VDCOrVDCGroupNetworkContextProfile, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	urlRef, err := c.Vmware.Client.OpenApiBuildEndpoint(govcdtypes.OpenApiPathVersion1_0_0 + govcdtypes.OpenApiEndpointNetworkContextProfiles)
	if err != nil {
		return nil, err
	}

	parameters := map[string]string{
		"filter": func() string {
			if uuid.IsNetworkContextProfile(networkContextProfileIDOrName) {
				return "id==" + networkContextProfileIDOrName
			}
			return "name==" + networkContextProfileIDOrName
		}(),
		"_context": vdcOrVDCGroupID,
		"scope":    string(scope),
	}

	queryParameters := url.Values{}
	for key, value := range parameters {
		queryParameters.Add(key, value)
	}

	typeResponses := []*govcdtypes.NsxtNetworkContextProfile{}
	err = c.Vmware.Client.OpenApiGetAllItems(c.Vmware.Client.APIVersion, urlRef, queryParameters, &typeResponses, nil)
	if err != nil {
		return nil, err
	}

	if len(typeResponses) == 0 {
		return nil, fmt.Errorf("%w : Network Context profile %s is not found in scope %s", govcd.ErrorEntityNotFound, networkContextProfileIDOrName, scope)
	}

	if len(typeResponses) > 1 {
		return nil, fmt.Errorf("multiple network context profiles found with name %s in scope %s", networkContextProfileIDOrName, scope)
	}

	return &VDCOrVDCGroupNetworkContextProfile{
		NsxtNetworkContextProfile: typeResponses[0],
	}, nil
}

func createNetworkContextProfile(networkContextProfile *govcdtypes.NsxtNetworkContextProfile) (*VDCOrVDCGroupNetworkContextProfile, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	urlRef, err := c.Vmware.Client.OpenApiBuildEndpoint(govcdtypes.OpenApiPathVersion1_0_0 + govcdtypes.OpenApiEndpointNetworkContextProfiles)
	if err != nil {
		return nil, err
	}

	output := &govcdtypes.NsxtNetworkContextProfile{}

	if err := c.Vmware.Client.OpenApiPostItem(c.Vmware.Client.APIVersion, urlRef, nil, networkContextProfile, output, nil); err != nil {
		return nil, err
	}

	return &VDCOrVDCGroupNetworkContextProfile{
		NsxtNetworkContextProfile: output,
	}, nil
}

func listNetworkContextProfileAttributes() any {
	c, err := clientcloudavenue.New()
	if err != nil {
		return err
	}

	urlRef, err := c.Vmware.Client.OpenApiBuildEndpoint(govcdtypes.OpenApiPathVersion1_0_0 + govcdtypes.OpenApiEndpointNetworkContextProfiles + "/attributes")
	if err != nil {
		return err
	}

	output := new(interface{})

	if err := c.Vmware.Client.OpenApiGetItem(c.Vmware.Client.APIVersion, urlRef, nil, output, nil); err != nil {
		return err
	}

	return output
}
