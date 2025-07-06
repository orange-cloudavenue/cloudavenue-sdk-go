package edgegateway

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/endpoints"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/utils"
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	commonvmware "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/vmware"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers/validators"
)

// CreateFirewallRules validates and creates new firewall rules on the EdgeGateway.
// It refreshes the gateway state, merges user-defined rules, sets rule names, and updates the rules via API.
func (e *EdgeGateway) CreateFirewallRules(ctx context.Context, rules FirewallModelRules) error {
	if err := validators.New().Struct(rules); err != nil {
		return err
	}

	if err := e.internalClient.Refresh(); err != nil {
		return err
	}

	fwR, err := e.getFirewallRules(ctx)
	if err != nil {
		return err
	}

	prioritizeRules(append(fwR.UserDefinedRules, rules.Rules...))

	// Set the name of the rule
	setRulesName(fwR.UserDefinedRules)

	// https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/2.0.0/edgeGateways/gatewayId/firewall/rules/put/
	resp, err := e.internalClient.R().
		SetContext(ctx).
		SetBody(fwR).
		SetPathParams(map[string]string{
			"edge-id": e.ID,
		}).
		Put(endpoints.FirewallRules)
	if err != nil {
		return err
	}

	return commonvmware.NewAndWait(&clientcloudavenue.GetClient().Vmware.Client, resp)
}

// UpdateFirewallRules updates existing firewall rules on the EdgeGateway. If the rule does not exist, a new rule will be created.
// Returns an error if validation fails, the rule does not exist, or the update request fails.
func (e *EdgeGateway) UpdateFirewallRules(ctx context.Context, rules FirewallModelRules) error {
	if err := validators.New().Struct(rules); err != nil {
		return err
	}

	if err := e.internalClient.Refresh(); err != nil {
		return err
	}

	fwR, err := e.getFirewallRules(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules.Rules {
		added := false

		for i, existingRule := range fwR.UserDefinedRules {
			if ((existingRule.ID != "") && (existingRule.ID == rule.ID)) || (existingRule.ComputeHash() == rule.ComputeHash()) {
				added = true
				fwR.UserDefinedRules[i] = rule
				break
			}
		}

		if !added {
			// If the rule does not exist, add it to the list
			fwR.UserDefinedRules = append(fwR.UserDefinedRules, rule)
		}
	}

	prioritizeRules(fwR.UserDefinedRules)

	// Set the name of the rule
	setRulesName(fwR.UserDefinedRules)

	// https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/2.0.0/edgeGateways/gatewayId/firewall/rules/put/
	resp, err := e.internalClient.R().
		SetContext(ctx).
		SetBody(fwR).
		SetPathParams(map[string]string{
			"edge-id": e.ID,
		}).
		Put(endpoints.FirewallRules)
	if err != nil {
		return err
	}

	return commonvmware.NewAndWait(&clientcloudavenue.GetClient().Vmware.Client, resp)
}

// GetFirewallRules retrieves the current firewall rules for the EdgeGateway.
func (e *EdgeGateway) GetFirewallRules(ctx context.Context) (FirewallModelRules, error) {
	if err := e.internalClient.Refresh(); err != nil {
		return FirewallModelRules{}, err
	}

	fwR, err := e.getFirewallRules(ctx)
	if err != nil {
		return FirewallModelRules{}, err
	}

	return FirewallModelRules{
		Rules: fwR.UserDefinedRules,
	}, nil
}

// DeleteFirewallRules deletes multiple firewall rules identified by their IDs from the EdgeGateway.
func (e *EdgeGateway) DeleteFirewallRules(ctx context.Context, rulesID []string) error {
	if err := e.internalClient.Refresh(); err != nil {
		return err
	}

	// https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/2.0.0/edgeGateways/gatewayId/firewall/rules/ruleId/delete/
	for _, rule := range rulesID {
		resp, err := e.internalClient.R().
			SetContext(ctx).
			SetPathParams(map[string]string{
				"edge-id": e.ID,
				"rule-id": rule,
			}).
			Delete(endpoints.FirewallRule)
		if err != nil {
			return err
		}

		if err := commonvmware.NewAndWait(&clientcloudavenue.GetClient().Vmware.Client, resp); err != nil {
			return err
		}
	}

	return nil
}

// DeleteAllFirewallRules deletes all firewall rules from the EdgeGateway.
func (e *EdgeGateway) DeleteAllFirewallRules(ctx context.Context) error {
	if err := e.internalClient.Refresh(); err != nil {
		return err
	}

	// https://developer.broadcom.com/xapis/vmware-cloud-director-openapi/v38.1/cloudapi/2.0.0/edgeGateways/gatewayId/firewall/rules/delete/
	resp, err := e.internalClient.R().
		SetContext(ctx).
		SetPathParams(map[string]string{
			"edge-id": e.ID,
		}).
		Delete(endpoints.FirewallRules)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return resp.Error().(error)
	}

	// Wait for the job to finish
	return commonvmware.NewAndWait(&clientcloudavenue.GetClient().Vmware.Client, resp)
}

func (e *EdgeGateway) getFirewallRules(ctx context.Context) (*firewallModelAPIRequest, error) {
	responseRules := &firewallModelAPIRequest{}

	// Get network services
	resp, err := e.internalClient.R().
		SetContext(ctx).
		SetResult(responseRules).
		SetPathParams(map[string]string{
			"edge-id": e.ID,
		}).
		Get(endpoints.FirewallRules)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return nil, govcd.ErrorEntityNotFound
		}
		return nil, resp.Error().(error)
	}

	extractRulesPriority(responseRules.UserDefinedRules)

	prioritizeRules(responseRules.UserDefinedRules)

	for _, rule := range responseRules.UserDefinedRules {
		// Compute the hash of the rule
		rule.Hash = rule.ComputeHash()
	}

	return responseRules, nil
}

func prioritizeRules(rules []*FirewallModelRule) {
	// calculate the order of rules based on the priority
	// 1 is the highest priority
	// 1000 is the lowest priority

	// priority is not a valid field in the API to overcome that we need to prefix the name with the priority
	// e.g 1_rule_name
	// e.g 100_rule_name

	// But a lots of rules have the same priority
	// so we need to sort them by the name
	// e.g 1_rule_name
	// e.g 1_rule_name_2
	// e.g 1_rule_name_3

	// Sort rules by priority and name
	sort.SliceStable(rules, func(i, j int) bool {
		if rules[i].Priority == rules[j].Priority {
			return rules[i].Name < rules[j].Name
		}
		return *rules[i].Priority < *rules[j].Priority
	})
}

func setRulesName(rules []*FirewallModelRule) {
	// Set the name of the rule
	for index, rule := range rules {
		rule.Name = fmt.Sprintf("%d_%s", *rule.Priority, rule.Name)
		rules[index] = rule
	}
}

func extractRulesPriority(rules []*FirewallModelRule) {
	// Extract the priority of the rules
	for _, rule := range rules {
		// Parse the name and set the priority
		// e.g 1_rule_name = priority 1 & name rule_name
		// e.g 100_rule_name = priority 100 & name rule_name

		before, after, found := strings.Cut(rule.Name, "_")
		if found {
			priorityInt, err := strconv.Atoi(before)
			if err != nil {
				rule.Priority = utils.ToPTR(0)
			} else {
				rule.Priority = utils.ToPTR(priorityInt)
			}
			rule.Name = after
		} else {
			rule.Priority = utils.ToPTR(0)
		}
	}
}
