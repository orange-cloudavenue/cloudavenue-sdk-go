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
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	mock "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/org/mock"
)

func TestClient_ListCertificatesInLibrary(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for admin org.
	clientAdminOrg := mock.NewMockclientGoVCDAdminOrg(ctrl)

	// Mock client for cloudavenue.
	clientCAV := mock.NewMockclientCloudavenue(ctrl)

	c := &Client{
		clientGoVCDAdminOrg: clientAdminOrg,
		clientCloudavenue:   clientCAV,
	}

	testCases := []struct {
		name              string
		mockFunc          []func()
		expectedCertValue CertificatesModel
		expectedErr       bool
		err               error
	}{
		{
			name: "success",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetAllCertificatesFromLibrary(nil).Return([]*govcd.Certificate{
						{
							CertificateLibrary: &govcdtypes.CertificateLibraryItem{
								Id:          "test",
								Alias:       "test",
								Description: "test",
								Certificate: "test",
							},
						},
						{
							CertificateLibrary: &govcdtypes.CertificateLibraryItem{
								Id:          "test2",
								Alias:       "test2",
								Description: "test2",
								Certificate: "test2",
							},
						},
					}, nil)
				},
			},
			expectedCertValue: CertificatesModel{
				{
					ID:          "test",
					Name:        "test",
					Description: "test",
					Certificate: "test",
					PrivateKey:  "",
					Passphrase:  "",
				},
				{
					ID:          "test2",
					Name:        "test2",
					Description: "test2",
					Certificate: "test2",
					PrivateKey:  "",
					Passphrase:  "",
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "refresh-error",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(errors.New("error"))
				},
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name: "error-get-all-certificates",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetAllCertificatesFromLibrary(nil).Return(nil, errors.New("error"))
				},
			},
			expectedCertValue: nil,
			expectedErr:       true,
		},
		{
			name: "error-get-all-certificates-nil",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetAllCertificatesFromLibrary(nil).Return(nil, nil)
				},
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               errors.New("no certificates found in the library"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, f := range tc.mockFunc {
				f()
			}

			certificates, err := c.ListCertificatesInLibrary()
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

	// Mock client for admin org.
	clientAdminOrg := mock.NewMockclientGoVCDAdminOrg(ctrl)

	// Mock client for cloudavenue.
	clientCAV := mock.NewMockclientCloudavenue(ctrl)

	c := &Client{
		clientGoVCDAdminOrg: clientAdminOrg,
		clientCloudavenue:   clientCAV,
	}

	generatedValidID := urn.CertificateLibraryItem.String() + uuid.New().String()

	testCases := []struct {
		name              string
		mockFunc          []func()
		nameOrID          string
		expectedCertValue CertificateModel
		expectedErr       bool
		err               error
	}{
		{
			name:     "success",
			nameOrID: "test",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetCertificateFromLibraryByName(gomock.Any()).Return(&govcd.Certificate{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:          "test",
							Alias:       "test",
							Description: "test",
							Certificate: "test",
						},
					}, nil)
				},
			},
			expectedCertValue: CertificateModel{
				ID:          "test",
				Name:        "test",
				Description: "test",
				Certificate: "test",
				PrivateKey:  "",
				Passphrase:  "",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:     "success-id",
			nameOrID: generatedValidID,
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(&govcd.Certificate{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:          generatedValidID,
							Alias:       "test",
							Description: "test",
							Certificate: "test",
						},
					}, nil)
				},
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        "test",
				Description: "test",
				Certificate: "test",
				PrivateKey:  "",
				Passphrase:  "",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:     "refresh-error",
			nameOrID: "test",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(errors.New("error"))
				},
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name:     "error-get-cert-by-name",
			nameOrID: "test",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetCertificateFromLibraryByName(gomock.Any()).Return(nil, errors.New("error"))
				},
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
		{
			name:     "error-get-cert-by-id",
			nameOrID: generatedValidID,
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().GetCertificateFromLibraryById(gomock.Any()).Return(nil, errors.New("error"))
				},
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, f := range tc.mockFunc {
				f()
			}

			certificate, err := c.GetCertificateFromLibrary(tc.nameOrID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificate)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificate)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, certificate.Certificate)
		})
	}
}

func TestClient_CreateCertificateInLibrary(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for admin org.
	clientAdminOrg := mock.NewMockclientGoVCDAdminOrg(ctrl)

	// Mock client for cloudavenue.
	clientCAV := mock.NewMockclientCloudavenue(ctrl)

	c := &Client{
		clientGoVCDAdminOrg: clientAdminOrg,
		clientCloudavenue:   clientCAV,
	}

	generatedValidID := urn.CertificateLibraryItem.String() + uuid.New().String()

	testCases := []struct {
		name              string
		mockFunc          []func()
		certificate       CertificateModel
		expectedCertValue CertificateModel
		expectedErr       bool
		err               error
	}{
		{
			name: "success",
			certificate: CertificateModel{
				Name:        "test",
				Description: "test",
				Certificate: "test",
				PrivateKey:  "test",
				Passphrase:  "test",
			},
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().AddCertificateToLibrary(gomock.Any()).Return(&govcd.Certificate{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:                   generatedValidID,
							Alias:                "test",
							Description:          "test",
							Certificate:          "test",
							PrivateKey:           "test",
							PrivateKeyPassphrase: "test",
						},
					}, nil)
				},
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        "test",
				Description: "test",
				Certificate: "test",
				PrivateKey:  "",
				Passphrase:  "",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-no-description",
			certificate: CertificateModel{
				Name:        "test",
				Certificate: "test",
				PrivateKey:  "test",
				Passphrase:  "test",
			},
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().AddCertificateToLibrary(gomock.Any()).Return(&govcd.Certificate{
						CertificateLibrary: &govcdtypes.CertificateLibraryItem{
							Id:          generatedValidID,
							Alias:       "test",
							Description: "",
							Certificate: "test",
						},
					}, nil)
				},
			},
			expectedCertValue: CertificateModel{
				ID:          generatedValidID,
				Name:        "test",
				Description: "",
				Certificate: "test",
				PrivateKey:  "",
				Passphrase:  "",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-validation",
			certificate: CertificateModel{
				Name:        "test",
				Certificate: "",
			},
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
		{
			name: "refresh-error",
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(errors.New("error"))
				},
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name: "error-add-cert",
			certificate: CertificateModel{
				Name:        "test",
				Certificate: "test",
			},
			mockFunc: []func(){
				func() {
					clientCAV.EXPECT().Refresh().Return(nil)
				},
				func() {
					clientAdminOrg.EXPECT().AddCertificateToLibrary(gomock.Any()).Return(nil, errors.New("error"))
				},
			},
			expectedCertValue: CertificateModel{},
			expectedErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, f := range tc.mockFunc {
				f()
			}

			certificate, err := c.CreateCertificateInLibrary(tc.certificate)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificate)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificate)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, certificate.Certificate)
		})
	}
}
