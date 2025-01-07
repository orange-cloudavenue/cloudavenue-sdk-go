/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package clientnetbackup

import (
	"errors"
	"testing"

	cavErrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

func TestOpts_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *Opts
		wantErr error
	}{
		{
			name: "should return an error if the organization is empty",
			opts: &Opts{
				org: "",
				URL: "",
			},
			wantErr: cavErrors.ErrEmpty,
		},
		{
			name: "should return an error if the endpoint and the url are empty and the organization is invalid",
			opts: &Opts{
				org: "tel01ev01ocb0001234",
			},
			wantErr: cavErrors.ErrOrganizationFormatIsInvalid,
		},
		{
			name: "should not return an error if the endpoint or the url are not empty and the organization is provided",
			opts: &Opts{
				org:      "cav01ev01ocb0001234",
				Endpoint: "https://backup4.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
			wantErr: nil,
		},
		{
			name: "should not return an error if the endpoint or the url are not empty and the organization is empty",
			opts: &Opts{
				URL: "https://backup4.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
			wantErr: nil,
		},
		{
			name: "Validate console1",
			opts: &Opts{
				org: "cav01ev01ocb0001234",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()

			if tt.wantErr == nil && err != nil {
				t.Errorf("Opts.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr != nil && err == nil {
				t.Errorf("Opts.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Opts.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && err == nil {
				// Check if URL and Endpoint are not empty
				if tt.opts.URL == "" && tt.opts.Endpoint == "" {
					t.Errorf("Opts.Validate() URL and Endpoint are empty")
				}
			}
		})
	}
}

func TestInit(t *testing.T) {
	tests := []struct {
		name         string
		org          string
		opts         *Opts
		expectedErr  error
		expectedUser string
		expectedPass string
		expectedURL  string
	}{
		{
			name: "should return an error if options validation fails",
			org:  "",
			opts: &Opts{
				URL: "",
			},
			expectedErr: cavErrors.ErrEmpty,
		},
		{
			name: "should not set username, password, endpoint, and debug if not provided",
			org:  "cav01ev01ocb0001234",
			opts: &Opts{
				Endpoint: "https://backup4.cloudavenue.orange-business.com/NetBackupSelfService/Api",
			},
			expectedUser: "",
			expectedPass: "",
			expectedURL:  "",
		},
		{
			name: "should set username, password, and debug if provided",
			org:  "cav01ev01ocb0001234",
			opts: &Opts{
				Username: "testuser",
				Password: "testpass",
				Debug:    true,
			},
			expectedUser: "testuser",
			expectedPass: "testpass",
			expectedURL:  "https://backup1.cloudavenue.orange-business.com/NetBackupSelfService/Api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(tt.opts, tt.org)

			if tt.expectedErr != nil && err == nil {
				t.Errorf("Init() error = %v, expected error %v", err, tt.expectedErr)
				return
			}

			if tt.expectedErr == nil && err != nil {
				t.Errorf("Init() error = %v, expected no error", err)
				return
			}

			if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("Init() error = %v, expected error %v", err, tt.expectedErr)
				return
			}

			if tt.expectedUser != c.token.username {
				t.Errorf("Init() username = %v, expected %v", c.token.username, tt.expectedUser)
				return
			}

			if tt.expectedPass != c.token.password {
				t.Errorf("Init() password = %v, expected %v", c.token.password, tt.expectedPass)
				return
			}

			if tt.expectedURL != c.token.endpoint {
				t.Errorf("Init() endpoint = %v, expected %v", c.token.endpoint, tt.expectedURL)
				return
			}

			if tt.opts.Debug != c.token.debug {
				t.Errorf("Init() debug = %v, expected %v", c.token.debug, tt.opts.Debug)
				return
			}
		})
	}
}
