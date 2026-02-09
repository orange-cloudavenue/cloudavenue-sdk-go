/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

import (
	"fmt"
	"regexp"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/endpoints"
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

type PublicIP struct{}

// IPs represents the list of public IPs (legacy format for backward compatibility).
type IPs struct {
	NetworkConfig []IP `json:"networkConfig"`
}

// IP represents a public IP address.
type IP struct {
	UplinkIP        string `json:"uplinkIp"`
	EdgeGatewayName string `json:"edgeGatewayName"`
	Announced       bool   `json:"announced"`
	// ServiceID is the service identifier used for delete operations.
	ServiceID string `json:"serviceId"`
}

// GetIP - Returns the public IP address.
func (i *IP) GetIP() string {
	return i.UplinkIP
}

// networkHierarchyResponse represents the response from /network endpoint.
type networkHierarchyResponse []networkHierarchyItem

type networkHierarchyItem struct {
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	DisplayName string                 `json:"displayName,omitempty"`
	Properties  networkItemProperties  `json:"properties,omitempty"`
	Children    []networkHierarchyItem `json:"children,omitempty"`
	ServiceID   string                 `json:"serviceId,omitempty"`
}

type networkItemProperties struct {
	// Edge Gateway properties
	RateLimit int    `json:"rateLimit,omitempty"`
	EdgeUUID  string `json:"edgeUUID,omitempty"`
	// Internet service properties
	IP        string `json:"ip,omitempty"`
	Announced bool   `json:"announced,omitempty"`
}

// GetIPs - Returns the list of public IPs from the network hierarchy.
func (v *PublicIP) GetIPs() (response *IPs, err error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	r, err := c.R().
		SetResult(&networkHierarchyResponse{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Get(endpoints.NetworkServiceGet)
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on get public IPs: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	hierarchy := r.Result().(*networkHierarchyResponse)
	if hierarchy == nil {
		return &IPs{NetworkConfig: []IP{}}, nil
	}

	// Extract public IPs from the hierarchy
	ips := &IPs{NetworkConfig: []IP{}}
	for _, vrf := range *hierarchy {
		if vrf.Type != "tier-0-vrf" {
			continue
		}
		for _, child := range vrf.Children {
			if child.Type == "edge-gateway" {
				edgeGatewayName := child.Name
				for _, service := range child.Children {
					if service.Type == "service" && service.Name == "internet" {
						ips.NetworkConfig = append(ips.NetworkConfig, IP{
							UplinkIP:        service.Properties.IP,
							EdgeGatewayName: edgeGatewayName,
							Announced:       service.Properties.Announced,
							ServiceID:       service.ServiceID,
						})
					}
				}
			}
		}
	}

	return ips, nil
}

// GetIPsByEdgeGateway - Returns the list of public IPs by edge gateway name.
func (v *PublicIP) GetIPsByEdgeGateway(edgeGatewayName string) (response *IPs, err error) {
	ipS, err := v.GetIPs()
	if err != nil {
		return nil, err
	}

	var ipsByEdgeGateway IPs
	for _, ip := range ipS.NetworkConfig {
		if ip.EdgeGatewayName == edgeGatewayName {
			ipsByEdgeGateway.NetworkConfig = append(ipsByEdgeGateway.NetworkConfig, ip)
		}
	}

	return &ipsByEdgeGateway, nil
}

// GetIP - Returns the public IP by IP address.
func (v *PublicIP) GetIP(publicIP string) (response *IP, err error) {
	ipS, err := v.GetIPs()
	if err != nil {
		return nil, err
	}

	for _, ip := range ipS.NetworkConfig {
		if ip.UplinkIP == publicIP {
			return &ip, nil
		}
	}

	return nil, fmt.Errorf("public IP %s not found", publicIP)
}

// GetIPByJob - Returns the public IP by job.
func (v *PublicIP) GetIPByJob(job *commoncloudavenue.JobStatus) (response *IP, err error) {
	if job == nil {
		return nil, fmt.Errorf("job is nil")
	}

	for _, action := range job.Actions {
		// if actions details have a public IP then return it
		// regex IPV4
		reg := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
		if reg.MatchString(action.Details) {
			return v.GetIP(reg.FindString(action.Details))
		}
	}

	return nil, fmt.Errorf("public IP not found")
}

// internetServiceRequest represents the request body for creating an internet service.
type internetServiceRequest struct {
	NetworkType string `json:"networkType"`
	EdgeGateway string `json:"edgeGateway"`
}

// New - Creates a new public IP on the specified edge gateway.
// The edgeGatewayID must be the UUID of the edge gateway (not the name).
func (v *PublicIP) New(edgeGatewayID string) (job *commoncloudavenue.JobStatus, err error) {
	if edgeGatewayID == "" {
		return nil, fmt.Errorf("edgeGatewayID is empty")
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	req := internetServiceRequest{
		NetworkType: "internet",
		EdgeGateway: edgeGatewayID,
	}

	r, err := c.R().
		SetBody(req).
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Post(endpoints.NetworkServiceCreate)
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on create public IP: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus), nil
}

// Delete - Deletes a public IP.
func (i *IP) Delete() (job *commoncloudavenue.JobStatus, err error) {
	if i.ServiceID == "" {
		return nil, fmt.Errorf("serviceID is empty, cannot delete public IP %s", i.UplinkIP)
	}

	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	r, err := c.R().
		SetPathParams(map[string]string{
			"service-id": i.ServiceID,
		}).
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		Delete(endpoints.NetworkServiceDelete)
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on delete public IP: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	return r.Result().(*commoncloudavenue.JobStatus), nil
}
