package v1

import (
	"github.com/avast/retry-go/v4"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// GetNetworkIsolated returns the isolated network by its name or ID.
func (g *VDCGroup) GetNetworkIsolated(nameOrID string) (*VDCNetworkIsolated, error) {
	net, err := g.genericGetNetwork(nameOrID)
	if err != nil {
		return nil, err
	}

	x := &VDCNetworkIsolated{
		VDCNetwork: VDCNetwork[*VDCNetworkIsolatedModel]{
			v:    g,
			net:  net,
			data: &VDCNetworkIsolatedModel{},
		},
	}

	x.data.fromVDCNetworkModel(net.OpenApiOrgVdcNetwork)
	x.VDCNetworkIsolatedModel = x.data
	return x, nil
}

func (g VDCGroup) genericGetNetwork(nameOrID string) (*govcd.OpenApiOrgVdcNetwork, error) {
	var values *govcd.OpenApiOrgVdcNetwork

	err := retry.Do(
		func() error {
			var err error
			if urn.IsNetwork(nameOrID) {
				values, err = g.getVDCNetworkByID(nameOrID)
			} else {
				values, err = g.getVDCNetworkByName(nameOrID)
			}

			return err
		},
		retry.RetryIf(govcd.ContainsNotFound),
		retry.Attempts(5),
	)

	return values, err
}

// CreateNetworkIsolated creates an isolated network.
func (g *VDCGroup) CreateNetworkIsolated(model *VDCNetworkIsolatedModel) (*VDCNetworkIsolated, error) {
	net, err := g.createVDCNetwork(model.toVDCNetworkModel(g))
	if err != nil {
		return nil, err
	}

	x := &VDCNetworkIsolated{
		VDCNetwork: VDCNetwork[*VDCNetworkIsolatedModel]{
			v:    g,
			net:  net,
			data: &VDCNetworkIsolatedModel{},
		},
	}

	x.data.fromVDCNetworkModel(net.OpenApiOrgVdcNetwork)
	x.VDCNetworkIsolatedModel = x.data
	return x, nil
}
