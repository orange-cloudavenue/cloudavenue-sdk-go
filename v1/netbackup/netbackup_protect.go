/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package netbackup

type protectBody struct {
	ProtectionLevelID int      `json:"ProtectionLevelId"`
	Paths             []string `json:"Paths,omitempty"`
}

// ProtectUnprotectRequest - Is the request structure for the Protects and Unprotects APIs.
type ProtectUnprotectRequest struct {
	// One of the following protection level settings must be specified
	// The ID of the protection level in the netbackup system
	ProtectionLevelID *int // Optional if ProtectionLevelName is specified
	// The name of the protection level in the netbackup system
	ProtectionLevelName string // Optional if ProtectionLevelID is specified
}

type protectionLevelAppliedResponse struct {
	Data struct {
		ProtectedLevels []struct {
			EntityType      string          `json:"EntityType"`
			ProtectionLevel ProtectionLevel `json:"ProtectionLevel"`
		} `json:"ProtectedLevels"`
	} `json:"data"`
}
