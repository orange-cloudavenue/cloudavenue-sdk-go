package v1

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

type (
	CAVAdminVDC struct{}
	AdminVDC    struct {
		*govcd.AdminVdc
	}
)

// GetAdminVdc return the admin vdc using the name provided in the provider.
func (v CAVAdminVDC) Get(vdcName string) (*AdminVDC, error) {
	if vdcName == "" {
		return nil, fmt.Errorf("%w", ErrEmptyVDCNameProvided)
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	av, err := c.AdminOrg.GetAdminVDCByName(vdcName, true)
	if err != nil {
		return nil, fmt.Errorf("%w: %s %w", ErrRetrievingVDC, vdcName, err)
	}

	return &AdminVDC{av}, nil
}
