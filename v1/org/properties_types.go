/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package org

type (
	// Properties represents a properties in the organization settings.
	PropertiesModel struct {
		Email        string `json:"customerMail,omitempty" validate:"omitempty,email"`
		Description  string `json:"description,omitempty" validate:"omitempty"`
		FullName     string `json:"fullName,omitempty" validate:"omitempty"`
		BillingModel string `json:"internetBillingMode" validate:"required,oneof=PAYG TRAFFIC_VOLUME"`
	}

	PropertiesRequest PropertiesModel

	propertiesResponse struct {
		Email               string `json:"customerMail,omitempty"`
		Description         string `json:"description,omitempty"`
		FullName            string `json:"fullName,omitempty"`
		InternetBillingMode string `json:"internetBillingMode"`
		IsEnabled           bool   `json:"isEnabled"`
		IsSuspended         bool   `json:"isSuspended"`
		Name                string `json:"name"`
	}
)
