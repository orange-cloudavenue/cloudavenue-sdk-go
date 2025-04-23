package edgegateway

import (
	"context"
	"fmt"
	"net"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/endpoints"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func (e *EdgeGateway) getNetworkServices(ctx context.Context) error {
	// Clean up the previous services
	e.EdgeGatewayModel.Services = NetworkServicesModelSvcs{
		Service:      nil,
		LoadBalancer: nil,
		PublicIP:     nil,
	}

	// Get network services
	resp, err := e.R().
		SetContext(ctx).
		SetResult(&networkServicesResponse{}).
		Get(endpoints.NetworkServiceGet)
	if err != nil {
		return fmt.Errorf("error getting network services: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("error getting network services: %s", resp.String())
	}

	originalResponse := resp.Result().(*networkServicesResponse)
	if originalResponse == nil || len(*originalResponse) == 0 {
		return fmt.Errorf("no network services found")
	}

	// Parse the original response and populate the NetworkServicesModel
	for _, ns := range *originalResponse {
		if ns.Type == "tier-0-vrf" && ns.Name != e.UplinkT0 {
			continue
		}
		for _, child := range ns.Children {
			if child.Type == "edge-gateway" && child.Properties.EdgeUUID == urn.ExtractUUID(e.ID) {
				// Found the edge gateway
				// iterate over the children to find the services
				for _, service := range child.Children {
					switch service.Type {
					case "load-balancer":
						// Found load balancer service
						e.EdgeGatewayModel.Services.LoadBalancer = &NetworkServicesModelSvcLoadBalancer{
							NetworkServicesModelSvc: NetworkServicesModelSvc{
								ID:   service.Name,        // The name is the ID
								Name: service.DisplayName, // The display name is the name
							},
							ClassOfService:     service.Properties.ClassOfService,
							MaxVirtualServices: service.Properties.MaxVirtualServices,
						}
					case "service":
						// service is a generic service
						// the name of the service define the type of service
						switch service.Name {
						case "cav-services", "cav_services": // Match both cav-services and cav_services
							// Found cav-services
							e.EdgeGatewayModel.Services.Service = &NetworkServicesModelSvcService{
								NetworkServicesModelSvc: NetworkServicesModelSvc{
									ID:   service.ServiceID,   // The ServiceID is the ID
									Name: service.DisplayName, // The display name is the name
								},
								Network: service.Properties.Ranges[0], // The first range is the network
								DedicatedIPForService: func() string {
									// Parse Network (ip/cidr) to get the first IP of the network
									// and use it as the dedicated IP for the service

									ip, _, err := net.ParseCIDR(service.Properties.Ranges[0])
									if err != nil {
										return ""
									}
									return ip.String()
								}(),
								ServiceDetails: ListOfServices,
							}

						case "internet":
							// Found internet service
							publicIP := &NetworkServicesModelSvcPublicIP{
								NetworkServicesModelSvc: NetworkServicesModelSvc{
									ID:   service.ServiceID,     // The ServiceID is the ID
									Name: service.Properties.IP, // The IP don't have a name use IP instead
								},
								IP:        service.Properties.IP,
								Announced: service.Properties.Announced,
							}

							// Prevent nil pointer dereference
							if e.EdgeGatewayModel.Services.PublicIP == nil {
								e.EdgeGatewayModel.Services.PublicIP = make([]*NetworkServicesModelSvcPublicIP, 0)
							}

							// Append the public IP to the list
							e.EdgeGatewayModel.Services.PublicIP = append(e.EdgeGatewayModel.Services.PublicIP, publicIP)
						}
					}
				}
			}
		}
	}

	return nil
}

func (e *EdgeGateway) EnableNetworkService(ctx context.Context) error {
	type networkServicesModelCreation struct {
		NetworkType string `json:"networkType"`
		EdgeGateway string `json:"edgeGateway"`
		Properties  struct {
			PrefixLength int `json:"prefixLength"`
		}
	}

	nsmc := networkServicesModelCreation{
		NetworkType: "cav-services",
		EdgeGateway: urn.ExtractUUID(e.ID),
		Properties: struct {
			PrefixLength int `json:"prefixLength"`
		}{
			PrefixLength: 27,
		},
	}

	if e.NetworkServiceIsEnabled() {
		return fmt.Errorf("the service is already enabled")
	}

	// Enable network services
	resp, err := e.R().
		SetContext(ctx).
		SetBody(nsmc).
		SetResult(&commoncloudavenue.JobStatus{}).
		Post(endpoints.NetworkServiceCreate)
	if err != nil {
		return fmt.Errorf("error enabling network services: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("error enabling network services: %s", resp.String())
	}

	job := resp.Result().(*commoncloudavenue.JobStatus)

	// Wait for the job to finish
	if err := job.WaitWithContext(ctx, 2); err != nil {
		return err
	}

	return e.getNetworkServices(ctx)
}

func (e *EdgeGateway) DisableNetworkService(ctx context.Context) error {
	if e.EdgeGatewayModel.Services.Service == nil {
		return fmt.Errorf("network service is not enabled")
	}

	// Disable network services
	resp, err := e.R().
		SetContext(ctx).
		SetResult(&commoncloudavenue.JobStatus{}).
		Delete(endpoints.InlineTemplate(endpoints.NetworkServiceDelete, map[string]string{
			"service-id": e.EdgeGatewayModel.Services.Service.ID,
		}))
	if err != nil {
		return fmt.Errorf("error disabling network services: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("error disabling network services: %s", resp.String())
	}

	job := resp.Result().(*commoncloudavenue.JobStatus)

	// Wait for the job to finish
	if err := job.WaitWithContext(ctx, 2); err != nil {
		return err
	}

	return e.getNetworkServices(ctx)
}

func (e *EdgeGateway) NetworkServiceIsEnabled() bool {
	return e.EdgeGatewayModel.Services.Service != nil
}
