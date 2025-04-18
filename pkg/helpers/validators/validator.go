/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package validators

import "github.com/go-playground/validator/v10"

// New creates a new validator.
func New() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	_ = v.RegisterValidation(DisallowUpper.Key, DisallowUpper.Func)
	_ = v.RegisterValidation(DisallowSpace.Key, DisallowSpace.Func)
	_ = v.RegisterValidation(KeyValue.Key, KeyValue.Func)
	_ = v.RegisterValidation(URN.Key, URN.Func)

	// * Network
	_ = v.RegisterValidation(IPV4Range.Key, IPV4Range.Func)
	_ = v.RegisterValidation(TCPUDPPort.Key, TCPUDPPort.Func)
	_ = v.RegisterValidation(TCPUDPPortRange.Key, TCPUDPPortRange.Func)

	// * HTTP
	_ = v.RegisterValidation(HTTPStatusCode.Key, HTTPStatusCode.Func)
	_ = v.RegisterValidation(HTTPStatusCodeRange.Key, HTTPStatusCodeRange.Func)

	return v
}
