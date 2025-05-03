package edgegateway

import (
	"crypto/sha256"
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	FirewallModel struct{}

	FirewallModelRules struct {
		Rules []*FirewallModelRule `validate:"required,dive"`
	}

	FirewallModelRule struct {
		// ID contains UUID (e.g. d0bf5d51-f83a-489a-9323-1661024874b8)
		ID string `validate:"omitempty" json:"id,omitempty"`
		// Name - API does not enforce uniqueness
		Name string `validate:"required" json:"name"`
		// Comment
		Comment string `validate:"omitempty,len=2048" json:"comment,omitempty"`

		// Priority defines the order in which the rules are applied. The lower the number, the higher the priority.
		Priority *int `validate:"required,min=1,max=1000" json:"-"`

		// ActionValue replaces deprecated field Action and defines action to be applied to all the
		// traffic that meets the firewall rule criteria. It determines if the rule permits or blocks
		// traffic. Property is required if action is not set. Below are valid values:
		// * ALLOW permits traffic to go through the firewall.
		// * DROP blocks the traffic at the firewall. No response is sent back to the source.
		// * REJECT blocks the traffic at the firewall. A response is sent back to the source.
		Action string `validate:"required,oneof=ALLOW DROP REJECT" json:"actionValue"`

		// Enabled allows to enable or disable the rule
		Enabled bool `validate:"omitempty" default:"true" json:"active"`

		// * Sourcess

		// SourceIPAddresses contains a list of IP addresses. Empty list means 'Any'
		SourceIPAddresses []string `validate:"omitempty,dive,ipv4|cidr|ipv4_range" json:"sourceFirewallIpAddresses,omitempty"`
		// SourceFirewallGroups contains a list of references to Firewall Groups. Empty list means 'Any'
		SourceFirewallGroups []govcdtypes.OpenApiReference `validate:"omitempty" json:"sourceFirewallGroups,omitempty"`

		// * Destinations

		// DestinationIPAddresses contains a list of IP addresses. Empty list means 'Any'
		DestinationIPAddresses []string `validate:"omitempty,dive,ipv4|cidr|ipv4_range" json:"destinationFirewallIpAddresses,omitempty"`
		// DestinationFirewallGroups contains a list of references to Firewall Groups. Empty list means 'Any'
		DestinationFirewallGroups []govcdtypes.OpenApiReference `validate:"omitempty" json:"destinationFirewallGroups,omitempty"`

		// * Properties

		// NetworkContextProfiles contains a list of references to Network Context Profiles. Empty list means 'Any'
		NetworkContextProfiles []govcdtypes.OpenApiReference `validate:"omitempty" json:"networkContextProfiles,omitempty"`
		// ApplicationPortProfiles contains a list of references to Application Port Profiles. Empty list means 'Any'
		ApplicationPortProfiles []govcdtypes.OpenApiReference `validate:"omitempty" json:"applicationPortProfiles,omitempty"`

		// TODO rawPortProtocols (https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/2.0.0/edgeGateways/gatewayId/firewall/rules/post/)

		// IpProtocol 'IPV4', 'IPV6', 'IPV4_IPV6'
		IPProtocol string `validate:"omitempty,oneof=IPV4 IPV6 IPV4_IPV6" default:"IPV4" json:"ipProtocol"`

		Logging bool `validate:"omitempty" json:"logging"`

		// Direction 'IN_OUT', 'OUT', 'IN'
		Direction string `validate:"required,oneof=IN OUT IN_OUT" json:"direction"`

		// Hash is a sha256 hash of the rule's Name, Action, and Direction.
		// It is used to identify the rule in the API in the create operation.
		Hash string `validate:"omitempty" json:"-"`
	}

	firewallModelAPIRequest struct {
		SystemRules      []*FirewallModelRule `json:"systemRules"`
		UserDefinedRules []*FirewallModelRule `json:"userDefinedRules"`
		DefaultRules     []*FirewallModelRule `json:"defaultRules"`
	}
)

// Hash is used to identify the rule in the API.
// It is a combination of Name, Action, Direction
func (r FirewallModelRule) ComputeHash() string {
	return FirewallHashRule(r.Name, r.Action, r.Direction)
}

func FirewallHashRule(name, action, direction string) string {
	// Hash is used to identify the rule in the API.
	// It is a combination of Name, Action, Direction
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s/%s/%s",
		name,
		action,
		direction,
	)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
