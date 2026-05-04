/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

const (
	testSuccess      = "success"
	testCertName     = "test"
	testCertName2    = "test2"
	testRefreshError = "refresh-error"
)

func TestClient_ListCertificatesInLibrary(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockinternalClient(ctrl)

	c, _ := NewFakeClient(clientCAV)

	testCases := []struct {
		name              string
		mockFunc          func()
		expectedCertValue CertificatesModel
		expectedErr       bool
		err               error
	}{
		{
			name: testSuccess,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllCertificatesFromLibrary(nil).Return([]*govcd.Certificate{
					{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:          testCertName,
							Alias:       testCertName,
							Description: testCertName,
							Certificate: testCertName,
						},
					},
					{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:          testCertName2,
							Alias:       testCertName2,
							Description: testCertName2,
							Certificate: testCertName2,
						},
					},
				}, nil)
			},
			expectedCertValue: CertificatesModel{
				{
					ID:          testCertName,
					Name:        testCertName,
					Description: testCertName,
					Certificate: testCertName,
				},
				{
					ID:          testCertName2,
					Name:        testCertName2,
					Description: testCertName2,
					Certificate: testCertName2,
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: testRefreshError,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name: "error-get-all-certificates",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllCertificatesFromLibrary(nil).Return(nil, errors.New("error"))
			},
			expectedCertValue: nil,
			expectedErr:       true,
		},
		{
			name: "error-get-all-certificates-nil",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllCertificatesFromLibrary(nil).Return(nil, nil)
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               errors.New("no certificates found in the library"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			certificates, err := c.ListCertificatesInLibrary(context.Background())
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificates)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificates)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, certificates)
		})
	}
}

func TestClient_GetCertificateFromLibrary(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockinternalClient(ctrl)

	c, _ := NewFakeClient(clientCAV)

	generatedValidID := urn.CertificateLibraryItem.String() + uuid.New().String()

	testCases := []struct {
		name              string
		mockFunc          func()
		nameOrID          string
		expectedCertValue CertificateModel
		expectedErr       bool
		err               error
	}{
		{
			name:     testSuccess,
			nameOrID: testCertName,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryByName(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          testCertName,
						Alias:       testCertName,
						Description: testCertName,
						Certificate: testCertName,
					},
				}, nil)
			},
			expectedCertValue: CertificateModel{
				ID:          testCertName,
				Name:        testCertName,
				Description: testCertName,
				Certificate: testCertName,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:     "success-id",
			nameOrID: generatedValidID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          generatedValidID,
						Alias:       testCertName,
						Description: testCertName,
						Certificate: testCertName,
					},
				}, nil)
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        testCertName,
				Description: testCertName,
				Certificate: testCertName,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:     testRefreshError,
			nameOrID: testCertName,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name:     "error-get-cert-by-name",
			nameOrID: testCertName,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryByName(gomock.Any()).Return(nil, errors.New("error"))
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
		{
			name:     "error-get-cert-by-id",
			nameOrID: generatedValidID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(nil, errors.New("error"))
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			certificate, err := c.GetCertificateFromLibrary(context.Background(), tc.nameOrID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificate)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificate)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, *certificate)
		})
	}
}

func TestClient_CreateCertificateInLibrary(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockinternalClient(ctrl)

	c, _ := NewFakeClient(clientCAV)

	generatedValidID := urn.CertificateLibraryItem.String() + uuid.New().String()

	testCases := []struct {
		name              string
		mockFunc          func()
		certificate       CertificateCreateRequest
		expectedCertValue CertificateModel
		expectedErr       bool
		err               error
	}{
		{
			name: testSuccess,
			certificate: CertificateCreateRequest{
				Name:        testCertName,
				Description: testCertName,
				Certificate: testCertName,
				PrivateKey:  testCertName,
				Passphrase:  testCertName,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().AddCertificateToLibrary(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:                   generatedValidID,
						Alias:                testCertName,
						Description:          testCertName,
						Certificate:          testCertName,
						PrivateKey:           testCertName,
						PrivateKeyPassphrase: testCertName,
					},
				}, nil)
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        testCertName,
				Description: testCertName,
				Certificate: testCertName,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-no-description",
			certificate: CertificateCreateRequest{
				Name:        testCertName,
				Certificate: testCertName,
				PrivateKey:  testCertName,
				Passphrase:  testCertName,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().AddCertificateToLibrary(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          generatedValidID,
						Alias:       testCertName,
						Description: "",
						Certificate: testCertName,
					},
				}, nil)
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        testCertName,
				Description: "",
				Certificate: testCertName,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-validation",
			certificate: CertificateCreateRequest{
				Name:        testCertName,
				Certificate: "",
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
		{
			name: testRefreshError,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name: "error-add-cert",
			certificate: CertificateCreateRequest{
				Name:        testCertName,
				Certificate: testCertName,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().AddCertificateToLibrary(gomock.Any()).Return(nil, errors.New("error"))
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			certificate, err := c.CreateCertificateInLibrary(context.Background(), &tc.certificate)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificate)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificate)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, *certificate)
		})
	}
}

func TestClient_UpdateCertificateInLibrary(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockinternalClient(ctrl)

	c, _ := NewFakeClient(clientCAV)

	generatedValidID := urn.CertificateLibraryItem.String() + uuid.New().String()

	testCases := []struct {
		name              string
		mockFunc          func()
		certificateID     string
		certificate       CertificateUpdateRequest
		expectedCertValue CertificateModel
		expectedErr       bool
		err               error
	}{
		{
			name:          testSuccess,
			certificateID: generatedValidID,
			certificate: CertificateUpdateRequest{
				Name:        testCertName,
				Description: testCertName,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          generatedValidID,
						Alias:       testCertName,
						Description: testCertName,
						Certificate: testCertName,
					},
				}, nil)
				updateCertificateInLibrary = func(_ internalCertificateClient) (*govcd.Certificate, error) {
					return &govcd.Certificate{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:          generatedValidID,
							Alias:       testCertName,
							Description: testCertName,
							Certificate: testCertName,
						},
					}, nil
				}
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        testCertName,
				Description: testCertName,
				Certificate: testCertName,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "success-no-description",
			certificateID: generatedValidID,
			certificate: CertificateUpdateRequest{
				Name: testCertName,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          generatedValidID,
						Alias:       testCertName,
						Description: testCertName,
						Certificate: testCertName,
					},
				}, nil)
				updateCertificateInLibrary = func(_ internalCertificateClient) (*govcd.Certificate, error) {
					return &govcd.Certificate{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:          generatedValidID,
							Alias:       testCertName,
							Description: "",
							Certificate: testCertName,
						},
					}, nil
				}
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        testCertName,
				Description: "",
				Certificate: testCertName,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "error-validation",
			certificateID: generatedValidID,
			certificate: CertificateUpdateRequest{
				Name: "",
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
		{
			name:          testRefreshError,
			certificateID: testCertName,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name:          "error-get-cert",
			certificateID: generatedValidID,
			certificate: CertificateUpdateRequest{
				Name: testCertName,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(nil, errors.New("error"))
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
		{
			name:          "error-update-cert",
			certificateID: generatedValidID,
			certificate: CertificateUpdateRequest{
				Name: testCertName,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          generatedValidID,
						Alias:       testCertName,
						Description: testCertName,
						Certificate: testCertName,
					},
				}, nil)
				updateCertificateInLibrary = func(_ internalCertificateClient) (*govcd.Certificate, error) {
					return nil, errors.New("error")
				}
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			certificate, err := c.UpdateCertificateInLibrary(context.Background(), tc.certificateID, &tc.certificate)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificate)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificate)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, *certificate)
		})
	}
}

func TestClient_DeleteCertificateFromLibrary(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockinternalClient(ctrl)

	c, _ := NewFakeClient(clientCAV)

	generatedValidID := urn.CertificateLibraryItem.String() + uuid.New().String()

	testCases := []struct {
		name          string
		mockFunc      func()
		certificateID string
		expectedErr   bool
		err           error
	}{
		{
			name:          testSuccess,
			certificateID: generatedValidID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          generatedValidID,
						Alias:       testCertName,
						Description: testCertName,
						Certificate: testCertName,
					},
				}, nil)
				deleteCertificateFromLibrary = func(_ internalCertificateClient) error {
					return nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "error-get-cert",
			certificateID: generatedValidID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(nil, errors.New("error"))
			},
			expectedErr: true,
		},
		{
			name:          "error-delete-cert",
			certificateID: generatedValidID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(&govcd.Certificate{
					CertificateLibrary: &govcdtypes.CertificateLibraryItem{
						Id:          generatedValidID,
						Alias:       testCertName,
						Description: testCertName,
						Certificate: testCertName,
					},
				}, nil)
				deleteCertificateFromLibrary = func(_ internalCertificateClient) error {
					return errors.New("error")
				}
			},
			expectedErr: true,
		},
		{
			name:          testRefreshError,
			certificateID: testCertName,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := c.DeleteCertificateFromLibrary(context.Background(), tc.certificateID)
			if !tc.expectedErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
