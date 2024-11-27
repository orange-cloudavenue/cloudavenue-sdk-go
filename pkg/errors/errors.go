package errors

import (
	"errors"
	"fmt"
)

var (

	// * Generic
	ErrNotFound      = errors.New("not found")
	ErrEmpty         = errors.New("empty")
	ErrInvalidFormat = errors.New("invalid format")

	// * Client
	ErrConfigureVmwareClient       = errors.New("unable to configure vmware client")
	ErrOrganizationFormatIsInvalid = fmt.Errorf("organization has an %w", ErrInvalidFormat)

	// * VDCGroup
	// * VDCGroupFirewall
	ErrInvalidFirewallRuleDirection  = fmt.Errorf("firewall rule direction has an %w", ErrInvalidFormat)
	ErrInvalidFirewallRuleIPProtocol = fmt.Errorf("firewall rule ipProtocol has an %w", ErrInvalidFormat)
	ErrInvalidFirewallRuleAction     = fmt.Errorf("firewall rule action has an %w", ErrInvalidFormat)
)
