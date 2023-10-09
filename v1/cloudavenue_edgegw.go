package v1

import (
	"fmt"
	"net"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type (
	EdgeGateway struct{}
	RateLimit   int
	OwnerType   string
)

// IsOwnerVDC - Returns true if the owner type is VDC
func (o OwnerType) IsVDC() bool {
	return o == OwnerVDC
}

// IsVDCGROUP - Returns true if the owner type is VDCGROUP
func (o OwnerType) IsVDCGROUP() bool {
	return o == ownerVDCGROUP
}

const (
	OwnerVDC      OwnerType = "vdc"
	ownerVDCGROUP OwnerType = "vdc-group"
)

type (
	EdgeGateways []EdgeGw
	EdgeGw       struct {
		Tier0VrfName string    `json:"tier0VrfId"`
		EdgeID       string    `json:"edgeId"`
		EdgeName     string    `json:"edgeName"`
		OwnerType    OwnerType `json:"ownerType"`
		OwnerName    string    `json:"ownerName"`
		Description  string    `json:"description"`
		RateLimit    RateLimit `json:"rateLimit"`
	}
)

// GetTier0VrfID - Returns the Tier0VrfID
func (e *EdgeGw) GetTier0VrfID() string {
	return e.Tier0VrfName
}

// GetT0 - Returns the Tier0VrfID (alias)
func (e *EdgeGw) GetT0() string {
	return e.Tier0VrfName
}

// GetID - Returns the EdgeID
func (e *EdgeGw) GetID() string {
	return e.EdgeID
}

// GetName - Returns the EdgeName
func (e *EdgeGw) GetName() string {
	return e.EdgeName
}

// GetOwnerType - Returns the OwnerType
func (e *EdgeGw) GetOwnerType() OwnerType {
	return e.OwnerType
}

// GetOwnerName - Returns the OwnerName
func (e *EdgeGw) GetOwnerName() string {
	return e.OwnerName
}

// GetDescription - Returns the Description
func (e *EdgeGw) GetDescription() string {
	return e.Description
}

// List - Returns the list of edge gateways
func (v *EdgeGateway) List() (response *EdgeGateways, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&EdgeGateways{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get("/api/customers/v2.0/edges")
	if err != nil {
		return
	}

	if r.IsError() {
		return response, fmt.Errorf("error on list edge gateways: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*EdgeGateways), nil
}

// * New

// New - Creates a new edge gateway
func (v *EdgeGateway) New(vdcName, tier0VrfName string) (job *commoncloudavenue.JobStatus, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetBody(map[string]interface{}{
			"tier0VrfId": tier0VrfName,
		}).
		SetPathParam("VdcName", vdcName).
		Post("/api/customers/v2.0/vdcs/{VdcName}/edges")
	if err != nil {
		return
	}

	if r.IsError() {
		return job, fmt.Errorf("error on create edge gateway: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus), nil
}

// NewFromVDCGroup - Creates a new edge gateway from a VDC Group
func (v *EdgeGateway) NewFromVDCGroup(vdcGroupName, tier0VrfName string) (job *commoncloudavenue.JobStatus, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetBody(map[string]interface{}{
			"tier0VrfId": tier0VrfName,
		}).
		SetPathParam("VdcGroupName", vdcGroupName).
		Post("/api/customers/v2.0/vdc-groups/{VdcGroupName}/edges")
	if err != nil {
		return
	}

	if r.IsError() {
		return job, fmt.Errorf("error on create edge gateway from VDC Group: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus), nil
}

// * Get

// Get - Returns the edge gateway
func (v *EdgeGateway) GetByName(edgeGatewayName string) (response *EdgeGw, err error) {
	edgeGateways, err := v.List()
	if err != nil {
		return
	}

	for _, edgeGateway := range *edgeGateways {
		if edgeGateway.EdgeName == edgeGatewayName {
			return &edgeGateway, nil
		}
	}

	return response, fmt.Errorf("edge gateway %s not found", edgeGatewayName)
}

// GetByID - Returns the edge gateway ID
// ID format is UUID
func (v *EdgeGateway) GetByID(edgeGatewayID string) (response *EdgeGw, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&EdgeGw{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParams(map[string]string{
			"EdgeID": edgeGatewayID,
		}).
		Get("/api/customers/v2.0/edges/{EdgeID}")
	if err != nil {
		return
	}

	if r.IsError() {
		return response, fmt.Errorf("error on get edge gateway: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*EdgeGw), nil
}

// * Delete

// Delete - Deletes the edge gateway
func (e *EdgeGw) Delete() (job *commoncloudavenue.JobStatus, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParams(map[string]string{
			"EdgeID": e.EdgeID,
		}).
		Delete("/api/customers/v2.0/edges/{EdgeID}")
	if err != nil {
		return
	}

	if r.IsError() {
		return job, fmt.Errorf("error on delete edge gateway: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus), nil
}

// * RateLimit

// CalculateRateLimitAvailable - Returns the rate limit available

// GetRateLimit - Returns the RateLimit
func (e *EdgeGw) GetRateLimit() RateLimit {
	return e.RateLimit
}

// SetRateLimit - Sets the rate limit
func (e *EdgeGw) SetRateLimit(rateLimit int) (job *commoncloudavenue.JobStatus, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetBody(map[string]interface{}{
			"rateLimit": rateLimit,
		}).
		SetPathParams(map[string]string{
			"EdgeID": e.EdgeID,
		}).
		Put("/api/customers/v2.0/edges/{EdgeID}")
	if err != nil {
		return
	}

	if r.IsError() {
		return job, fmt.Errorf("error on set rate limit: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus), nil
}

// * NetworkTypes

type TypeOfNetwork string

const (
	NetworkTypeServiceZone TypeOfNetwork = "sv"
)

type (
	NetworkTypes []NetworkType
	NetworkType  struct {
		Type         TypeOfNetwork `json:"networkType"`
		PrefixLength int           `json:"prefixLength"`
		StartAddress string        `json:"startAddress"`
	}
)

// IsServiceZone - Returns true if the network type is ServiceZone
func (n NetworkType) IsServiceZone() bool {
	return n.Type == NetworkTypeServiceZone
}

// GetStartAddress - Returns the StartAddress
func (n NetworkType) GetStartAddress() string {
	return n.StartAddress
}

// GetPrefixLength - Returns the PrefixLength
func (n NetworkType) GetPrefixLength() int {
	return n.PrefixLength
}

// GetEndAddress - Returns the EndAddress
func (n NetworkType) GetEndAddress() string {
	// determine the end address from the start address and the prefix length
	// StartAddress is a IPv4 address
	// PrefixLength is an integer
	// EndAddress is a IPv4 address
	start := net.ParseIP(n.StartAddress).To4()
	if start == nil {
		return ""
	}

	prefixLength := n.PrefixLength
	if prefixLength < 0 || prefixLength > 32 {
		return ""
	}

	numIPs := uint32(1) << (32 - prefixLength)
	broadcast := numIPs - 1
	end := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		end[i] = start[i] + byte((broadcast>>uint(8*i))&0xff)
	}

	return end.String()
}

// ListAllAddressesAvailable - Returns the list of all addresses available from start address and end address
// func (e *EdgeGw) ListAllAddressesAvailable() (response []string) {
// 	// Returns the list of all addresses available from start address and end address
// 	// StartAddress is a IPv4 address
// 	// EndAddress is a IPv4 address
// 	// Returns a list of IPv4 addresses

// 	networkTypes, err := e.GetNetworkTypes()
// 	if err != nil {
// 		return
// 	}

// 	for _, networkType := range *networkTypes {
// 		startAddress := networkType.GetStartAddress()
// 		endAddress := networkType.GetEndAddress()

// 		if startAddress == "" || endAddress == "" {
// 			continue
// 		}

// 		start := net.ParseIP(startAddress).To4()
// 		end := net.ParseIP(endAddress).To4()

// 		if start == nil || end == nil {
// 			continue
// 		}

// 		for ip := start; bytes.Compare(ip, end) <= 0; incrementIP(ip) {
// 			response = append(response, ip.String())
// 		}
// 	}

// 	return
// }

// func incrementIP(ip net.IP) {
// 	for j := len(ip) - 1; j >= 0; j-- {
// 		ip[j]++
// 		if ip[j] > 0 {
// 			break
// 		}
// 	}
// }

// ListNetworksType - Returns the list of networks by type configured on the edge gateway
func (e *EdgeGw) ListNetworksType() (response *NetworkTypes, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return
	}

	r, err := c.R().
		SetResult(&NetworkTypes{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParams(map[string]string{
			"EdgeID": e.EdgeID,
		}).
		Get("/api/customers/v2.0/edges/{EdgeID}/networks")
	if err != nil {
		return
	}

	if r.IsError() {
		return response, fmt.Errorf("error on list networks type: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*NetworkTypes), nil
}
