/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commonnetbackup

import "fmt"

type APIError []struct {
	PropertyName string `json:"PropertyName"`
	Message      string `json:"Message"`
}

// Convert to Error.
func ToError(err *APIError) error {
	return fmt.Errorf("API Error: %v", err)
}
