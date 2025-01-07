/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	VDCGroupFirewall struct {
		vdcGroup *govcd.VdcGroup
		vgf      *govcd.DistributedFirewall
	}

	VDCGroupFirewallType struct {
		Enabled bool

		// Rules contains a list of firewall rules.
		Rules VDCGroupFirewallTypeRules
	}

	VDCGroupFirewallTypeRules          []*VDCGroupFirewallTypeRule
	VDCGroupFirewallTypeRuleDirection  string
	VDCGroupFirewallTypeRuleIPProtocol string
	VDCGroupFirewallTypeRuleAction     string

	VDCGroupFirewallTypeRule struct {
		// Required
		Name    string
		Enabled bool

		// Direction 'IN_OUT', 'OUT', 'IN'.
		//
		// Use helpers.ParseVDCGroupFirewallTypeRuleDirection to parse string to VDCGroupFirewallTypeRuleDirection
		Direction VDCGroupFirewallTypeRuleDirection

		// IpProtocol 'IPV4', 'IPV6', 'IPV4_IPV6'.
		//
		// Use helpers.ParseVDCGroupFirewallTypeRuleIPProtocol to parse string to VDCGroupFirewallTypeRuleIPProtocol
		IPProtocol VDCGroupFirewallTypeRuleIPProtocol

		// Action defines action to be applied to all the traffic that meets the firewall rule criteria.
		// It determines if the rule permits or blocks traffic.
		// Use helpers.ParseVDCGroupFirewallTypeRuleAction to parse string to VDCGroupFirewallTypeRuleAction
		//
		// - ALLOW permits traffic to go through the firewall.
		//
		// - DROP blocks the traffic at the firewall. No response is sent back to the source.
		//
		// - REJECT blocks the traffic at the firewall. A response is sent back to the source.
		Action VDCGroupFirewallTypeRuleAction

		// Optional
		ID          string
		Description string // Length cannot exceed 2048 characters.
		Logging     bool

		// ApplicationPortProfiles contains a list of references to Application Port Profiles. Empty
		// list means 'Any'
		ApplicationPortProfiles []govcdtypes.OpenApiReference

		// SourceFirewallGroups contains a list of references to Firewall Groups. Empty list means 'Any'
		SourceFirewallGroups []govcdtypes.OpenApiReference
		// DestinationFirewallGroups contains a list of references to Firewall Groups. Empty list means
		// 'Any'
		DestinationFirewallGroups []govcdtypes.OpenApiReference

		// SourceGroupsExcluded reverses the list specified in SourceFirewallGroups and the rule gets
		// applied on all the groups that are NOT part of the SourceFirewallGroups. If false, the rule
		// applies to the all the groups including the source groups.
		SourceGroupsExcluded *bool

		// DestinationGroupsExcluded reverses the list specified in DestinationFirewallGroups and the
		// rule gets applied on all the groups that are NOT part of the DestinationFirewallGroups. If
		// false, the rule applies to the all the groups in DestinationFirewallGroups.
		DestinationGroupsExcluded *bool
	}
)

const (
	VDCGroupFirewallTypeRuleDirectionInOut VDCGroupFirewallTypeRuleDirection = "IN_OUT"
	VDCGroupFirewallTypeRuleDirectionOut   VDCGroupFirewallTypeRuleDirection = "OUT"
	VDCGroupFirewallTypeRuleDirectionIn    VDCGroupFirewallTypeRuleDirection = "IN"

	VDCGroupFirewallTypeRuleIPProtocolIPv4     VDCGroupFirewallTypeRuleIPProtocol = "IPV4"
	VDCGroupFirewallTypeRuleIPProtocolIPv6     VDCGroupFirewallTypeRuleIPProtocol = "IPV6"
	VDCGroupFirewallTypeRuleIPProtocolIPv4IPv6 VDCGroupFirewallTypeRuleIPProtocol = "IPV4_IPV6"

	// Permits traffic to go through the firewall.
	VDCGroupFirewallTypeRuleActionAllow VDCGroupFirewallTypeRuleAction = "ALLOW"
	// Blocks the traffic at the firewall. No response is sent back to the source.
	VDCGroupFirewallTypeRuleActionDrop VDCGroupFirewallTypeRuleAction = "DROP"
	// Blocks the traffic at the firewall. A response is sent back to the source.
	VDCGroupFirewallTypeRuleActionReject VDCGroupFirewallTypeRuleAction = "REJECT"
)

var (
	VDCGroupFirewallTypeRuleDirections = []VDCGroupFirewallTypeRuleDirection{
		VDCGroupFirewallTypeRuleDirectionInOut,
		VDCGroupFirewallTypeRuleDirectionOut,
		VDCGroupFirewallTypeRuleDirectionIn,
	}

	VDCGroupFirewallTypeRuleIPProtocols = []VDCGroupFirewallTypeRuleIPProtocol{
		VDCGroupFirewallTypeRuleIPProtocolIPv4,
		VDCGroupFirewallTypeRuleIPProtocolIPv6,
		VDCGroupFirewallTypeRuleIPProtocolIPv4IPv6,
	}

	VDCGroupFirewallTypeRuleActions = []VDCGroupFirewallTypeRuleAction{
		VDCGroupFirewallTypeRuleActionAllow,
		VDCGroupFirewallTypeRuleActionDrop,
		VDCGroupFirewallTypeRuleActionReject,
	}
)
