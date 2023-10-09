package v1

import (
	"fmt"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type PublicIP struct{}

type IPs struct {
	InternalIP    string `json:"internalIp"`
	NetworkConfig []IP   `json:"networkConfig"`
}

type IP struct {
	UplinkIP        string `json:"uplinkIp"`
	TranslatedIP    string `json:"translatedIp"`
	EdgeGatewayName string `json:"edgeGatewayName"`
}

// GetIP - Returns the public IP
func (i *IP) GetIP() string {
	return i.UplinkIP
}

// GetIPs - Returns the list of public IPs
func (v *PublicIP) GetIPs() (response *IPs, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&IPs{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/ip")
	if err != nil {
		return
	}

	if r.IsError() {
		return response, fmt.Errorf("error on get public IPs: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*IPs), nil
}

// GetIPsByEdgeGateway - Returns the list of public IPs by edge gateway
func (v *PublicIP) GetIPsByEdgeGateway(edgeGatewayName string) (response *IPs, err error) {
	ipS, err := v.GetIPs()
	if err != nil {
		return
	}

	var ipsByEdgeGateway IPs
	for _, ip := range ipS.NetworkConfig {
		if ip.EdgeGatewayName == edgeGatewayName {
			ipsByEdgeGateway.NetworkConfig = append(ipsByEdgeGateway.NetworkConfig, ip)
		}
	}

	return &ipsByEdgeGateway, nil
}

// GetIP - Returns the public IP by edge gateway
func (v *PublicIP) GetIP(publicIP string) (response *IP, err error) {
	ipS, err := v.GetIPs()
	if err != nil {
		return
	}

	for _, ip := range ipS.NetworkConfig {
		if ip.TranslatedIP == publicIP {
			return &ip, nil
		}
	}

	return nil, fmt.Errorf("public IP %s not found", publicIP)
}

// New - Returns a new PublicIP
func New(edgeGatewayName string) (job *commoncloudavenue.JobStatus, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		panic(err)
	}

	if edgeGatewayName == "" {
		return nil, fmt.Errorf("edgeGatewayName is empty")
	}

	r, err := c.R().
		SetHeader("X-VDC-Edge-Name", edgeGatewayName).
		SetResult(&commoncloudavenue.JobCreatedAPIResponse{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Post("/api/customers/v1.0/ip")
	if err != nil {
		return
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on create public IP: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobCreatedAPIResponse).GetJobStatus()
}

// Delete - Deletes a public IP
func (i *IP) Delete() (job *commoncloudavenue.JobStatus, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		panic(err)
	}

	r, err := c.R().
		SetPathParams(map[string]string{
			"PublicIP": i.UplinkIP,
		}).
		SetResult(&commoncloudavenue.JobCreatedAPIResponse{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Delete("/api/customers/v1.0/ip/{PublicIP}/")
	if err != nil {
		return
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on delete public IP: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobCreatedAPIResponse).GetJobStatus()
}
