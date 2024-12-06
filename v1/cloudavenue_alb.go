package v1

import "fmt"

type (
	ALB struct {
		VirtualServices []VirtualService `json:"VirtualService"`
		Pools           []Pool           `json:"Pool"`
	}

	VirtualService struct {
		ID                     string                      `json:"id"`
		Name                   string                      `json:"name"`
		EdgeGatewayName        string                      `json:"edge_gateway_name"`
		EdgeGatewayID          string                      `json:"edge_gateway_id"`
		Description            string                      `json:"description"`
		Enabled                bool                        `json:"enabled"`
		PoolName               string                      `json:"pool_name"`
		PoolID                 string                      `json:"pool_id"`
		ServiceEngineGroupName string                      `json:"service_engine_group_name"`
		VirtualIP              string                      `json:"virtual_ip"`
		ServiceType            string                      `json:"service_type"`
		CertificateID          string                      `json:"certificate_id"`
		ServicePorts           []VirtualServiceServicePort `json:"service_ports"`
		PreserveClientIP       bool                        `json:"preserve_client_ip"`
	}

	VirtualServiceServicePort struct {
		PortStart int64  `json:"port_start"`
		PortEnd   int64  `json:"port_end"`
		PortType  string `json:"port_type"`
		PortSSL   bool   `json:"port_ssl"`
	}

	Pool struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

// ! ALB

// Get All Virtual Services
func (v *ALB) GetALLVirtualService() []VirtualService {
	return v.VirtualServices
}

// Get Virtual Service by name
func (v *ALB) GetVirtualServiceByName(name string) (response *VirtualService, err error) {
	// For each BMS find the one with the same hostname
	for _, albvs := range v.VirtualServices {
		if albvs.Name == name {
			return &albvs, nil
		}
	}
	return nil, fmt.Errorf("ALB Virtual Service with name %s not found", name)
}

// Get Virtual Service by ID
func (v *ALB) GetVirtualServiceByID(id string) (response *VirtualService, err error) {
	// For each BMS find the one with the same hostname
	for _, albvs := range v.VirtualServices {
		if albvs.ID == id {
			return &albvs, nil
		}
	}
	return nil, fmt.Errorf("ALB Virtual Service with ID %s not found", id)
}

// Get Virtual Service by Edge Gateway Name
func (v *ALB) GetVirtualServiceByEdgeGatewayName(name string) (response []VirtualService, err error) {
	// For each BMS find the one with the same hostname
	for _, albvs := range v.VirtualServices {
		if albvs.EdgeGatewayName == name {
			response = append(response, albvs)
		}
	}
	return response, nil
}

// Get ALL Service Ports
func (v *VirtualService) GetServicePorts() []VirtualServiceServicePort {
	return v.ServicePorts
}

// TODO - Add more methods to the struct : NEW METHOD

// New - Creates a new Virtual Service
func (v *ALB) NewVirtualService() (response *VirtualService, err error) {
	return &VirtualService{}, nil
}
