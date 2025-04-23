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
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

func (*client) edgeGatewayNameOrIDValidation(edgeGatewayNameOrID string) error {
	if edgeGatewayNameOrID == "" {
		return fmt.Errorf("edgeGatewayNameOrID is %w. Please provide a valid edgeGatewayNameOrID", errors.ErrEmpty)
	}

	return nil
}
