package org

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
)

func TestClient_GetProperties(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer httpmock.DeactivateAndReset()

	// Mock client for cloudavenue.
	clientCAV := NewMockinternalClient(ctrl)

	c, _ := NewFakeClient(clientCAV)

	testCases := []struct {
		name           string
		mockFunc       func()
		expectedValues PropertiesModel
		expectedErr    bool
	}{
		{
			name: "GetProperties success",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"fullName":"John Doe","description":"John Doe description","customerMail":"","internetBillingMode":"PAYG","isEnabled":true,"isSuspended":false}`))
					if err != nil {
						t.Fatal(err)
					}
					httpmock.RegisterResponder("GET", "/api/customers/v2.0/configurations", responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedValues: PropertiesModel{
				FullName:     "John Doe",
				Description:  "John Doe description",
				Email:        "",
				BillingModel: "PAYG",
			},
			expectedErr: false,
		},
		{
			name: "GetProperties error on refresh",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(fmt.Errorf("error on refresh"))
			},
			expectedValues: PropertiesModel{},
			expectedErr:    true,
		},
		{
			name: "GetProperties error on get properties",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					errResponse := commoncloudavenue.APIErrorResponse{
						Code:    "err-0002",
						Reason:  "fail-001",
						Message: "error on get properties",
					}

					b, err := json.Marshal(errResponse)
					if err != nil {
						t.Fatal(err)
					}

					responder, err := httpmock.NewJsonResponder(500, b)
					if err != nil {
						t.Fatal(err)
					}
					httpmock.RegisterResponder("GET", "/api/customers/v2.0/configurations", responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedValues: PropertiesModel{},
			expectedErr:    true,
		},
		{
			name: "GetProperties error on get properties with invalid response",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder := httpmock.NewStringResponder(200, ``)

					httpmock.RegisterResponder("GET", "/api/customers/v2.0/bad/path", responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedValues: PropertiesModel{},
			expectedErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Reset()
			httpmock.DeactivateAndReset()
			// Run the mock functions.
			tc.mockFunc()

			properties, err := c.GetProperties(context.Background())
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValues, *properties)
		})
	}
}

func TestClient_UpdateProperties(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer httpmock.DeactivateAndReset()

	// Mock client for cloudavenue.
	clientCAV := NewMockinternalClient(ctrl)

	c, _ := NewFakeClient(clientCAV)

	testCases := []struct {
		name        string
		mockFunc    func()
		properties  PropertiesRequest
		expectedErr bool
	}{
		{
			name: "UpdateProperties success",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"16f61fe1-481c-4ec1-b62d-046af22b51ad"}`))
					if err != nil {
						t.Fatal(err)
					}
					httpmock.BodyContainsBytes(json.RawMessage(`{"fullName":"John Doe","description":"John Doe description","customerMail":"","internetBillingMode":"PAYG"}`))
					httpmock.RegisterResponder("PUT", "/api/customers/v2.0/configurations", responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			properties: PropertiesRequest{
				FullName:     "John Doe",
				Description:  "John Doe description",
				Email:        "",
				BillingModel: "PAYG",
			},
			expectedErr: false,
		},
		{
			name: "UpdateProperties error on refresh",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(fmt.Errorf("error on refresh"))
			},
			properties:  PropertiesRequest{},
			expectedErr: true,
		},

		{
			name: "UpdateProperties error validate",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
			},
			properties:  PropertiesRequest{},
			expectedErr: true,
		},
		{
			name: "UpdateProperties error on update properties",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().DoAndReturn(func() error {
					clientcloudavenue.MockClient()
					return nil
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					errResponse := commoncloudavenue.APIErrorResponse{
						Code:    "err-0002",
						Reason:  "fail-001",
						Message: "error on update properties",
					}

					b, err := json.Marshal(errResponse)
					if err != nil {
						t.Fatal(err)
					}

					responder, err := httpmock.NewJsonResponder(500, b)
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("PUT", "/api/customers/v2.0/configurations", responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			properties: PropertiesRequest{
				FullName:     "John Doe",
				Description:  "John Doe description",
				Email:        "",
				BillingModel: "PAYG",
			},
			expectedErr: true,
		},
		{
			name: "UpdateProperties error on update properties with invalid response",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().DoAndReturn(func() error {
					clientcloudavenue.MockClient()
					return nil
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder := httpmock.NewStringResponder(200, ``)

					httpmock.RegisterResponder("PUT", "/api/customers/v2.0/bad/path", responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			properties: PropertiesRequest{
				FullName:     "John Doe",
				Description:  "John Doe description",
				Email:        "",
				BillingModel: "PAYG",
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.Reset()
			httpmock.DeactivateAndReset()
			// Run the mock functions.
			tc.mockFunc()

			_, err := c.UpdateProperties(context.Background(), &tc.properties)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
