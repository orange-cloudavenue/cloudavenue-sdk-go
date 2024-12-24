package helpers

import (
	"fmt"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

// ParseFirewallAppPortProfileProtocol parses the string and returns the corresponding firewall application port profile protocol.
// If the string is not valid, it returns an error.
func ParseFirewallAppPortProfileProtocol(protocol string) (v1.FirewallGroupAppPortProfileModelPortProtocol, error) {
	switch strings.ToLower(protocol) {
	case "icmpv4":
		return v1.FirewallGroupAppPortProfileModelPortProtocolICMPv4, nil
	case "icmpv6":
		return v1.FirewallGroupAppPortProfileModelPortProtocolICMPv6, nil
	case "tcp":
		return v1.FirewallGroupAppPortProfileModelPortProtocolTCP, nil
	case "udp":
		return v1.FirewallGroupAppPortProfileModelPortProtocolUDP, nil
	default:
		return "", fmt.Errorf("%w. Use one of %v", errors.ErrInvalidFirewallAppPortProfileProtocol, v1.FirewallGroupAppPortProfileModelPortProtocols)
	}
}
