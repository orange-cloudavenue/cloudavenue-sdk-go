package edgeloadbalancer

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func (c *client) ListPools(ctx context.Context, edgeGatewayID string) ([]*PoolModel, error) {
	if edgeGatewayID == "" {
		return nil, fmt.Errorf("edgeGatewayID is %w. Please provide a valid edgeGatewayID", errors.ErrEmpty)
	}

	if !urn.IsEdgeGateway(edgeGatewayID) {
		return nil, fmt.Errorf("edgeGatewayID has %w. Please provide a valid edgeGatewayID", errors.ErrInvalidFormat)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	allAlbPoolSummaries, err := c.clientGoVCD.GetAllAlbPoolSummaries(edgeGatewayID, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving all ALB Pool summaries: %w", err)
	}

	// Loop over all Summaries and retrieve complete information
	allAlbPools := make([]*PoolModel, len(allAlbPoolSummaries))
	for index := range allAlbPoolSummaries {
		allAlbPools[index], err = c.GetPool(ctx, allAlbPoolSummaries[index].NsxtAlbPool.GatewayRef.ID, allAlbPoolSummaries[index].NsxtAlbPool.ID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving complete ALB Pool: %w", err)
		}
	}

	return allAlbPools, nil
}

// GetPool retrieves a pool by name or ID.
func (c *client) GetPool(ctx context.Context, edgeGatewayID, poolNameOrID string) (*PoolModel, error) {
	if edgeGatewayID == "" {
		return nil, fmt.Errorf("edgeGatewayID is %w. Please provide a valid edgeGatewayID", errors.ErrEmpty)
	}

	if !urn.IsEdgeGateway(edgeGatewayID) {
		return nil, fmt.Errorf("edgeGatewayID has %w. Please provide a valid edgeGatewayID", errors.ErrInvalidFormat)
	}

	if poolNameOrID == "" {
		return nil, fmt.Errorf("poolNameOrID is %w. Please provide a valid poolNameOrID", errors.ErrEmpty)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	var (
		albPool *govcd.NsxtAlbPool
		err     error
	)

	if urn.IsLoadBalancerPool(poolNameOrID) {
		albPool, err = c.clientGoVCD.GetAlbPoolById(poolNameOrID)
	} else {
		albPool, err = c.clientGoVCD.GetAlbPoolByName(edgeGatewayID, poolNameOrID)
	}
	if err != nil {
		return nil, err
	}

	return &PoolModel{
		ID:                       albPool.NsxtAlbPool.ID,
		Name:                     albPool.NsxtAlbPool.Name,
		Description:              albPool.NsxtAlbPool.Description,
		GatewayRef:               albPool.NsxtAlbPool.GatewayRef,
		Enabled:                  albPool.NsxtAlbPool.Enabled,
		Algorithm:                PoolAlgorithm(albPool.NsxtAlbPool.Algorithm),
		DefaultPort:              albPool.NsxtAlbPool.DefaultPort,
		GracefulTimeoutPeriod:    albPool.NsxtAlbPool.GracefulTimeoutPeriod,
		PassiveMonitoringEnabled: albPool.NsxtAlbPool.PassiveMonitoringEnabled,
		HealthMonitors: func() []PoolModelHealthMonitor {
			monitors := make([]PoolModelHealthMonitor, len(albPool.NsxtAlbPool.HealthMonitors))
			for i, monitor := range albPool.NsxtAlbPool.HealthMonitors {
				monitors[i] = PoolModelHealthMonitor{
					Name: monitor.Name,
					Type: PoolHealthMonitorType(monitor.Type),
				}
			}
			return monitors
		}(),
		Members: func() []PoolModelMember {
			members := make([]PoolModelMember, len(albPool.NsxtAlbPool.Members))
			for i, member := range albPool.NsxtAlbPool.Members {
				members[i] = PoolModelMember{
					Enabled:               member.Enabled,
					IPAddress:             member.IpAddress,
					Port:                  member.Port,
					Ratio:                 member.Ratio,
					MarkedDownBy:          member.MarkedDownBy,
					HealthStatus:          member.HealthStatus,
					DetailedHealthMessage: member.DetailedHealthMessage,
				}
			}
			return members
		}(),
		MemberGroupRef:         albPool.NsxtAlbPool.MemberGroupRef,
		CaCertificateRefs:      albPool.NsxtAlbPool.CaCertificateRefs,
		CommonNameCheckEnabled: albPool.NsxtAlbPool.CommonNameCheckEnabled,
		DomainNames:            albPool.NsxtAlbPool.DomainNames,
		PersistenceProfile: func() *PoolModelPersistenceProfile {
			if albPool.NsxtAlbPool.PersistenceProfile == nil {
				return nil
			}
			return &PoolModelPersistenceProfile{
				Name:  albPool.NsxtAlbPool.PersistenceProfile.Name,
				Type:  PoolPersistenceProfileType(albPool.NsxtAlbPool.PersistenceProfile.Type),
				Value: albPool.NsxtAlbPool.PersistenceProfile.Value,
			}
		}(),
		MemberCount:        albPool.NsxtAlbPool.MemberCount,
		EnabledMemberCount: albPool.NsxtAlbPool.EnabledMemberCount,
		UpMemberCount:      albPool.NsxtAlbPool.UpMemberCount,
		HealthMessage:      albPool.NsxtAlbPool.HealthMessage,
		VirtualServiceRefs: albPool.NsxtAlbPool.VirtualServiceRefs,
		SSLEnabled:         albPool.NsxtAlbPool.SslEnabled,
	}, nil
}
