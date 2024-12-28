package v1

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

type (
	EdgeGateway struct{}
	Bandwidth   int
	OwnerType   string
)

// IsOwnerVDC - Returns true if the owner type is VDC.
func (o OwnerType) IsVDC() bool {
	return o == OwnerVDC
}

// IsVDCGROUP - Returns true if the owner type is VDCGROUP.
func (o OwnerType) IsVDCGROUP() bool {
	return o == ownerVDCGROUP
}

// GetVmwareEdgeGateway - Returns the VMware Edge Gateway.
func (e *EdgeClient) GetVmwareEdgeGateway() (*govcd.NsxtEdgeGateway, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return c.Org.GetNsxtEdgeGatewayById(urn.Normalize(urn.Gateway, e.GetID()).String())
}

// List - Returns the list of edge gateways.
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

var (
	allowedRateLimitVRFStandard        = []int{5, 25, 50, 75, 100, 150, 200, 250, 300}                                     // 5, 25, 50, 75, 100, 150, 200, 250, 300
	allowedRateLimitVRFPremium         = append(allowedRateLimitVRFStandard, []int{400, 500, 600, 700, 800, 900, 1000}...) // 5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000
	allowedRateLimitVRFDedicatedMedium = append(allowedRateLimitVRFPremium, []int{2000}...)                                // 5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000
	allowedRateLimitVRFDedicatedLarge  = append(allowedRateLimitVRFDedicatedMedium, []int{3000, 4000, 5000, 6000}...)      // 5, 25, 50, 75, 100, 150, 200, 250, 300, 400, 500, 600, 700, 800, 900, 1000, 2000, 3000, 4000, 5000, 6000
)

// GetAllowedBandwidthValues - Returns the allowed rate limit value.
func (v *EdgeGateway) GetAllowedBandwidthValues(t0VrfName string) (allowedValues []int, err error) {
	t0, err := (&Tier0{}).GetT0(t0VrfName)
	if err != nil {
		return
	}

	switch t0.GetClassService() {
	case ClassServiceVRFPremium:
		allowedValues = allowedRateLimitVRFPremium
	case ClassServiceVRFStandard:
		allowedValues = allowedRateLimitVRFStandard
	case ClassServiceVRFDedicatedMedium:
		allowedValues = allowedRateLimitVRFDedicatedMedium
	case ClassServiceVRFDedicatedLarge:
		allowedValues = allowedRateLimitVRFDedicatedLarge
	}

	return
}

// GetBandwidthCapacityRemaining - Returns the bandwidth capacity remaining in Mbps.
func (e *EdgeGateways) GetBandwidthCapacityRemaining(t0VrfName string) (response int, err error) {
	t0, err := (&Tier0{}).GetT0(t0VrfName)
	if err != nil {
		return
	}

	t0BandwidthCapacity, err := t0.GetBandwidthCapacity()
	if err != nil {
		return
	}

	for _, edgeGateway := range *e {
		if edgeGateway.GetT0() == t0VrfName {
			t0BandwidthCapacity -= int(edgeGateway.GetBandwidth())
		}
	}

	// 5 Mbps is the minimum bandwidth capacity
	if t0BandwidthCapacity < 5 {
		return 0, fmt.Errorf("no bandwidth capacity remaining")
	}

	return t0BandwidthCapacity, nil
}

// * New

// New - Creates a new edge gateway.
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

// NewFromVDCGroup - Creates a new edge gateway from a VDC Group.
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

func (v *EdgeGateway) Get(edgeGatewayNameOrID string) (edgeClient *EdgeClient, err error) {
	if edgeGatewayNameOrID == "" {
		return nil, fmt.Errorf("edge gateway name or ID is empty")
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	edgeClient = new(EdgeClient)

	if urn.IsUUIDV4(edgeGatewayNameOrID) || urn.IsEdgeGateway(edgeGatewayNameOrID) {
		nameOrID := urn.Normalize(urn.Gateway, edgeGatewayNameOrID)

		// wait group to wait for all goroutines to finish
		var wg errgroup.Group

		wg.Go(func() error {
			vmwareEdgeClient, err := c.Org.GetNsxtEdgeGatewayById(nameOrID.String())
			if err != nil {
				return err
			}
			edgeClient.EdgeVCDInterface = vmwareEdgeClient
			edgeClient.vcdEdge = vmwareEdgeClient
			return nil
		})

		wg.Go(func() error {
			r, err := c.R().
				SetResult(&EdgeGatewayType{}).
				SetError(&commoncloudavenue.APIErrorResponse{}).
				SetPathParams(map[string]string{
					"EdgeID": strings.TrimPrefix(nameOrID.String(), urn.Gateway.String()),
				}).
				Get("/api/customers/v2.0/edges/{EdgeID}")
			if err != nil {
				return err
			}

			if r.IsError() {
				return fmt.Errorf("error on get edge gateway: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
			}

			edgeClient.EdgeGatewayType = r.Result().(*EdgeGatewayType)

			return nil
		})
		return edgeClient, wg.Wait()
	}

	// * GetByName
	edgeGateways, err := v.List()
	if err != nil {
		return nil, err
	}

	for _, edgeGateway := range *edgeGateways {
		if edgeGateway.EdgeName == edgeGatewayNameOrID {
			return v.Get(edgeGateway.EdgeID)
		}
	}

	return nil, fmt.Errorf("%w: not found the edgegateway with name %s", govcd.ErrorEntityNotFound, edgeGatewayNameOrID)
}

// Get - Returns the edge gateway
//
// Deprecated: Use Get instead.
func (v *EdgeGateway) GetByName(edgeGatewayName string) (edgeClient *EdgeClient, err error) {
	return v.Get(edgeGatewayName)
}

// GetByID - Returns the edge gateway ID
// ID format is UUID
//
// Deprecated: Use Get instead.
func (v *EdgeGateway) GetByID(edgeGatewayID string) (edgeClient *EdgeClient, err error) {
	return v.Get(edgeGatewayID)
}

// * Delete

// Delete - Deletes the edge gateway.
func (e *EdgeGatewayType) Delete() (job *commoncloudavenue.JobStatus, err error) {
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

// * Bandwidth

// GetBandwidth - Returns the rate limit.
func (e *EdgeGatewayType) GetBandwidth() Bandwidth {
	return e.Bandwidth
}

// UpdateBandwidth - Updates the bandwidth.
func (e *EdgeGatewayType) UpdateBandwidth(rateLimit int) (job *commoncloudavenue.JobStatus, err error) {
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
		return job, fmt.Errorf("error on set bandwidth: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
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

// IsServiceZone - Returns true if the network type is ServiceZone.
func (n NetworkType) IsServiceZone() bool {
	return n.Type == NetworkTypeServiceZone
}

// GetStartAddress - Returns the StartAddress.
func (n NetworkType) GetStartAddress() string {
	return n.StartAddress
}

// GetPrefixLength - Returns the PrefixLength.
func (n NetworkType) GetPrefixLength() int {
	return n.PrefixLength
}

// GetEndAddress - Returns the EndAddress.
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

// ListNetworksType - Returns the list of networks by type configured on the edge gateway.
func (e *EdgeGatewayType) ListNetworksType() (response *NetworkTypes, err error) {
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
