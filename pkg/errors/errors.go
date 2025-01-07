/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

import (
	"errors"
	"fmt"
)

var (

	// * Generic.
	ErrNotFound      = errors.New("not found")
	ErrEmpty         = errors.New("empty")
	ErrInvalidFormat = errors.New("invalid format")

	// * Client.
	ErrConfigureVmwareClient       = errors.New("unable to configure vmware client")
	ErrOrganizationFormatIsInvalid = fmt.Errorf("organization has an %w", ErrInvalidFormat)

	// * VDCGroup
	// * VDCGroupFirewall.
	ErrInvalidFirewallRuleDirection  = fmt.Errorf("firewall rule direction has an %w", ErrInvalidFormat)
	ErrInvalidFirewallRuleIPProtocol = fmt.Errorf("firewall rule ipProtocol has an %w", ErrInvalidFormat)
	ErrInvalidFirewallRuleAction     = fmt.Errorf("firewall rule action has an %w", ErrInvalidFormat)

	// * FirewallAppPortProfile.
	ErrInvalidFirewallAppPortProfileProtocol = fmt.Errorf("firewall app port profile protocol has an %w", ErrInvalidFormat)
)
