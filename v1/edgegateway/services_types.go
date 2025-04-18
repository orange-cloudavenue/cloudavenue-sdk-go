package edgegateway

type (
	NetworkServicesModelSvcs struct {
		LoadBalancer *NetworkServicesModelSvcLoadBalancer
		PublicIP     []*NetworkServicesModelSvcPublicIP
		Service      *NetworkServicesModelSvcService
	}

	NetworkServicesModelSvc struct {
		// ID is the identifier of the network service
		ID string
		// Name is the name of the network service
		Name string
	}

	NetworkServicesModelSvcLoadBalancer struct {
		NetworkServicesModelSvc

		ClassOfService     string
		MaxVirtualServices int
	}

	NetworkServicesModelSvcPublicIP struct {
		NetworkServicesModelSvc

		// IP is the public IP address
		IP string
		// Announced represents if the public IP address is announced
		Announced bool
	}

	NetworkServicesModelSvcService struct {
		NetworkServicesModelSvc

		// Network is the network of the service ip/cidr
		Network string
		// DedicatedIPForService is the dedicated IP for the service
		// Used for the NAT to connect to the service
		DedicatedIPForService string
		// Services is the list of services
		ServiceDetails []ServiceModelDetails
	}

	// NetworkServicesModelSvcService.

	NetworkServicesModelSvcServiceDetailsPorts struct {
		// Port is the port of the service
		Port int
		// Protocol is the protocol of the service
		Protocol string
	}

	networkServicesResponse []struct {
		Type     string                            `json:"type"`
		Name     string                            `json:"name"`
		Children []networkServicesResponseChildren `json:"children,omitempty"`
	}

	networkServicesResponseChildren struct {
		Type        string `json:"type"`
		Name        string `json:"name,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		Properties  struct {
			// EdgeGateway
			RateLimit int    `json:"rateLimit,omitempty"`
			EdgeUUID  string `json:"edgeUUID,omitempty"`

			// Load Balancer
			ClassOfService     string `json:"classOfService,omitempty"`
			MaxVirtualServices int    `json:"maxVirtualServices,omitempty"`

			// Public IP
			IP        string `json:"ip,omitempty"`
			Announced bool   `json:"announced,omitempty"`

			// Service
			Ranges []string `json:"ranges,omitempty"`
		} `json:"properties,omitempty"`
		Children  []networkServicesResponseChildren `json:"children,omitempty"`
		ServiceID string                            `json:"serviceId,omitempty"`
	}
)
