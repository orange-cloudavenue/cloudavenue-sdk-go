/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

import (
	"context"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

//go:generate mockgen -source=client.go -destination=zz_generated_client_test.go -self_package github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgeloadbalancer -package edgeloadbalancer -copyright_file "../../mock_header.txt"

type (
	// Exposed client interface.
	Client interface {
		// * Service Engine Groups
		ListServiceEngineGroups(ctx context.Context, edgeGatewayID string) ([]*ServiceEngineGroupModel, error)
		GetServiceEngineGroup(ctx context.Context, edgeGatewayID, nameOrID string) (*ServiceEngineGroupModel, error)
		GetFirstServiceEngineGroup(ctx context.Context, edgeGatewayID string) (*ServiceEngineGroupModel, error)

		// * Pools
		CreatePool(ctx context.Context, pool PoolModelRequest) (*PoolModel, error)
		ListPools(ctx context.Context, edgeGatewayID string) ([]*PoolModel, error)
		GetPool(ctx context.Context, edgeGatewayID, poolNameOrID string) (*PoolModel, error)
		UpdatePool(ctx context.Context, poolID string, pool PoolModelRequest) (*PoolModel, error)
		DeletePool(ctx context.Context, poolID string) error

		// * Pools Health Monitors
		// Since the VCD API 39.0 the health monitors are assigned with a dedicated API.
		// CreatePool and UpdatePool reconcile them when the request defines some, the
		// methods below allow to manage them independently from the pool lifecycle.
		ListPoolHealthMonitors(ctx context.Context, edgeGatewayID, poolID string) ([]PoolModelHealthMonitor, error)
		AssignPoolHealthMonitor(ctx context.Context, edgeGatewayID, poolID string, monitorType PoolHealthMonitorType) error
		UnassignPoolHealthMonitor(ctx context.Context, edgeGatewayID, poolID string, monitorType PoolHealthMonitorType) error
		ReconcilePoolHealthMonitors(ctx context.Context, edgeGatewayID, poolID string, wanted []PoolHealthMonitorType) error

		// * Virtual Services
		ListVirtualServices(ctx context.Context, edgeGatewayID string) ([]*VirtualServiceModel, error)
		GetVirtualService(ctx context.Context, edgeGatewayID, virtualServiceNameOrID string) (*VirtualServiceModel, error)
		CreateVirtualService(ctx context.Context, vsr VirtualServiceModelRequest) (*VirtualServiceModel, error)
		UpdateVirtualService(ctx context.Context, virtualServiceID string, vsr VirtualServiceModelRequest) (*VirtualServiceModel, error)
		DeleteVirtualService(ctx context.Context, virtualServiceID string) error

		// * Policies
		// ? Request
		GetPoliciesHTTPRequest(ctx context.Context, virtualServiceID string) (*PoliciesHTTPRequestModel, error)
		UpdatePoliciesHTTPRequest(ctx context.Context, policies *PoliciesHTTPRequestModel) (*PoliciesHTTPRequestModel, error)
		DeletePoliciesHTTPRequest(ctx context.Context, virtualServiceID string) error
		// ? Response
		GetPoliciesHTTPResponse(ctx context.Context, virtualServiceID string) (*PoliciesHTTPResponseModel, error)
		UpdatePoliciesHTTPResponse(ctx context.Context, policies *PoliciesHTTPResponseModel) (*PoliciesHTTPResponseModel, error)
		DeletePoliciesHTTPResponse(ctx context.Context, virtualServiceID string) error
		// ? Security
		GetPoliciesHTTPSecurity(ctx context.Context, virtualServiceID string) (*PoliciesHTTPSecurityModel, error)
		UpdatePoliciesHTTPSecurity(ctx context.Context, policies *PoliciesHTTPSecurityModel) (*PoliciesHTTPSecurityModel, error)
		DeletePoliciesHTTPSecurity(ctx context.Context, virtualServiceID string) error
	}

	// Internal client interfaces.
	clientFake interface {
		clientGoVCD
		clientGoVCDOpenAPI
		clientCloudavenue
	}

	client struct {
		clientGoVCD        clientGoVCD
		clientGoVCDOpenAPI clientGoVCDOpenAPI
		clientCloudavenue  clientCloudavenue

		// apiVersion is the OpenAPI version of the VCD instance, required by the
		// generic OpenAPI calls which take it as a parameter.
		apiVersion string
	}

	clientGoVCD interface {
		// Service Engine Groups
		GetAllAlbServiceEngineGroupAssignments(queryParameters url.Values) ([]*govcd.NsxtAlbServiceEngineGroupAssignment, error)

		// Pools
		GetAlbPoolById(id string) (*govcd.NsxtAlbPool, error)
		GetAlbPoolByName(edgeGatewayID, name string) (*govcd.NsxtAlbPool, error)
		GetAllAlbPoolSummaries(edgeGatewayID string, queryParameters url.Values) ([]*govcd.NsxtAlbPool, error)
		CreateNsxtAlbPool(albPoolConfig *govcdtypes.NsxtAlbPool) (*govcd.NsxtAlbPool, error)

		// Virtual Services
		GetAlbVirtualServiceByName(edgeGatewayID, name string) (*govcd.NsxtAlbVirtualService, error)
		GetAlbVirtualServiceById(id string) (*govcd.NsxtAlbVirtualService, error)
		GetAllAlbVirtualServiceSummaries(edgeGatewayID string, queryParameters url.Values) ([]*govcd.NsxtAlbVirtualService, error)
		CreateNsxtAlbVirtualService(albVirtualServiceConfig *govcdtypes.NsxtAlbVirtualService) (*govcd.NsxtAlbVirtualService, error)
	}

	// Generic OpenAPI calls of go-vcloud-director. They live on the underlying
	// govcd.Client (a named field of VCDClient, not an embedded one), so they
	// cannot join the clientGoVCD interface, which is satisfied by *govcd.VCDClient.
	// Used by the entities go-vcloud-director does not expose, like the Load
	// Balancer Pool health monitor assignments.
	clientGoVCDOpenAPI interface {
		OpenApiBuildEndpoint(endpoint ...string) (*url.URL, error)
		OpenApiGetItem(apiVersion string, urlRef *url.URL, params url.Values, outType interface{}, additionalHeader map[string]string) error
		OpenApiGetAllItems(apiVersion string, urlRef *url.URL, queryParams url.Values, outType interface{}, additionalHeader map[string]string) error
		OpenApiPostItemAsync(apiVersion string, urlRef *url.URL, params url.Values, payload interface{}) (govcd.Task, error)
		OpenApiDeleteItem(apiVersion string, urlRef *url.URL, params url.Values, additionalHeader map[string]string) error
	}

	clientCloudavenue interface {
		Refresh() error
		R() *resty.Request
	}
)

// NewClient creates a new edgegateway load balancer client.
func NewClient() (Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return &client{
		clientCloudavenue:  c,
		clientGoVCD:        c.Vmware,
		clientGoVCDOpenAPI: &c.Vmware.Client,
		apiVersion:         c.Vmware.Client.APIVersion,
	}, nil
}

// NewFakeClient creates a new fake Org client used for testing.
func NewFakeClient(i clientFake) (Client, error) {
	return &client{
		clientCloudavenue:  i,
		clientGoVCD:        i,
		clientGoVCDOpenAPI: i,
		apiVersion:         clientcloudavenue.DefaultVCDAPIVersion,
	}, nil
}
