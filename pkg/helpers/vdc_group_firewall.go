/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package helpers

import (
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

// ParseVDCGroupFirewallTypeRuleDirection parses the string and returns the corresponding firewall rule direction.
// If the string is not valid, it returns an error.
func ParseVDCGroupFirewallRuleDirection(direction string) (v1.VDCGroupFirewallTypeRuleDirection, error) {
	switch strings.ToLower(direction) {
	case "in":
		return v1.VDCGroupFirewallTypeRuleDirectionIn, nil
	case "out":
		return v1.VDCGroupFirewallTypeRuleDirectionOut, nil
	case "inout", "in_out", "in-out":
		return v1.VDCGroupFirewallTypeRuleDirectionInOut, nil
	default:
		return "", fmt.Errorf("%w. Use one of %v", errors.ErrInvalidFirewallRuleDirection, v1.VDCGroupFirewallTypeRuleDirections)
	}
}

// ParseVDCGroupFirewallRuleIPProtocol parses the string and returns the corresponding firewall rule IP protocol.
// If the string is not valid, it returns an error.
func ParseVDCGroupFirewallRuleIPProtocol(protocol string) (v1.VDCGroupFirewallTypeRuleIPProtocol, error) {
	switch strings.ToLower(protocol) {
	case "ipv4":
		return v1.VDCGroupFirewallTypeRuleIPProtocolIPv4, nil
	case "ipv6":
		return v1.VDCGroupFirewallTypeRuleIPProtocolIPv6, nil
	case "ipv4-ipv6", "ipv4_ipv6", "ipv4ipv6":
		return v1.VDCGroupFirewallTypeRuleIPProtocolIPv4IPv6, nil
	default:
		return "", fmt.Errorf("%w. Use one of %v", errors.ErrInvalidFirewallRuleIPProtocol, v1.VDCGroupFirewallTypeRuleIPProtocols)
	}
}

// ParseVDCGroupFirewallRuleAction parses the string and returns the corresponding firewall rule action.
// If the string is not valid, it returns an error.
func ParseVDCGroupFirewallRuleAction(action string) (v1.VDCGroupFirewallTypeRuleAction, error) {
	switch strings.ToLower(action) {
	case "allow":
		return v1.VDCGroupFirewallTypeRuleActionAllow, nil
	case "drop":
		return v1.VDCGroupFirewallTypeRuleActionDrop, nil
	case "reject":
		return v1.VDCGroupFirewallTypeRuleActionReject, nil
	default:
		return "", fmt.Errorf("%w. Use one of %v", errors.ErrInvalidFirewallRuleAction, v1.VDCGroupFirewallTypeRuleActions)
	}
}
