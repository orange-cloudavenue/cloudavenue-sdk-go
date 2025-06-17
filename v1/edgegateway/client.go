/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
)

//go:generate mockgen -source=client.go -destination=zz_generated_client_test.go -self_package github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/edgegateway -package edgegateway -copyright_file "../../mock_header.txt"

type (
	// Exposed client interface.
	Client interface {
		ListEdgeGateway(ctx context.Context) ([]*EdgeGatewayModel, error)
		GetEdgeGateway(ctx context.Context, edgeGatewayNameOrID string) (*EdgeGateway, error)
		CreateEdgeGateway(ctx context.Context, edgeGateway *EdgeGatewayModelRequest) (*EdgeGatewayModel, error)
		UpdateEdgeGateway(ctx context.Context, edgeGateway *EdgeGatewayModelUpdate) error
		DeleteEdgeGateway(ctx context.Context, edgeGatewayNameOrID string) error
	}

	// Internal client interfaces.
	clientInterface interface {
		clientGoVCDOrg
		clientCloudavenue
	}

	client struct {
		clientGoVCDOrg
		clientCloudavenue
	}

	clientGoVCDOrg interface {
		GetAllNsxtEdgeGateways(queryParameters url.Values) ([]*govcd.NsxtEdgeGateway, error)
		GetNsxtEdgeGatewayById(id string) (*govcd.NsxtEdgeGateway, error)
		GetNsxtEdgeGatewayByName(name string) (*govcd.NsxtEdgeGateway, error)

		GetVDCById(vdcID string, refresh bool) (*govcd.Vdc, error)
		GetVdcGroupById(id string) (*govcd.VdcGroup, error)
	}

	clientCloudavenue interface {
		Refresh() error
		R() *resty.Request
		GetClient() *http.Client
	}
)

// NewClient creates a new edgegateway load balancer client.
func NewClient() (Client, error) {
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	return &client{
		clientCloudavenue: c,
		clientGoVCDOrg:    c.Org,
	}, nil
}

// NewFakeClient creates a new fake Org client used for testing.
func NewFakeClient(i clientInterface) (Client, error) {
	return &client{
		clientCloudavenue: i,
		clientGoVCDOrg:    i,
	}, nil
}

func newFakeEdgeGatewayClient(i clientInterface) *EdgeGateway {
	return &EdgeGateway{
		internalClient: i,
	}
}
