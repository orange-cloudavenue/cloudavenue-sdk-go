/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// The health monitor assignment is managed by a dedicated API since the VCD API 39.0
// (VCD 10.6). The property `healthMonitors` is not part of the Load Balancer Pool
// payload anymore and the platform rejects it with:
// "Load Balancer Health Monitor Assignment is not supported with legacy API".
//
// go-vcloud-director does not expose these endpoints (as of v2.26.2), so they are
// called through the generic OpenAPI helpers, the same way the network context
// profiles are handled within this SDK.
const (
	// healthMonitorAssignmentsEndpoint lists and creates the health monitor assignments
	// of the Load Balancer Pools.
	healthMonitorAssignmentsEndpoint = govcdtypes.OpenApiPathVersion1_0_0 + "loadBalancer/healthMonitors/assignments"

	// healthMonitorAssignmentEndpoint reads or deletes one health monitor assignment.
	healthMonitorAssignmentEndpoint = govcdtypes.OpenApiPathVersion1_0_0 + "loadBalancer/healthMonitors/assignment/"

	// healthMonitorSummariesEndpointFormat lists the health monitors available for
	// a given Edge Gateway.
	healthMonitorSummariesEndpointFormat = govcdtypes.OpenApiPathVersion1_0_0 + "edgeGateways/%s/loadBalancer/healthMonitorSummaries"
)

type (
	// poolHealthMonitorAssignment is the API representation of a health monitor
	// assigned to a Load Balancer Pool.
	poolHealthMonitorAssignment struct {
		ID                  string                      `json:"id,omitempty"`
		HealthMonitorRef    govcdtypes.OpenApiReference `json:"healthMonitorRef"`
		LoadBalancerPoolRef govcdtypes.OpenApiReference `json:"loadBalancerPoolRef"`
	}

	// poolHealthMonitor is a health monitor available for a given Edge Gateway.
	poolHealthMonitor struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
)

// waitTask waits the completion of a task returned by an asynchronous endpoint.
// It is a seam used to keep the client testable.
var waitTask = func(task govcd.Task) error {
	return task.WaitTaskCompletion()
}

func healthMonitorTypes(monitors []PoolModelHealthMonitor) []PoolHealthMonitorType {
	types := make([]PoolHealthMonitorType, 0, len(monitors))
	for _, monitor := range monitors {
		types = append(types, monitor.Type)
	}
	return types
}

// ListPoolHealthMonitors retrieves the health monitors assigned to a Load Balancer Pool.
func (c *client) ListPoolHealthMonitors(ctx context.Context, edgeGatewayID, poolID string) ([]PoolModelHealthMonitor, error) {
	if err := c.checkHealthMonitorArguments(edgeGatewayID, poolID); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	assignments, err := c.listHealthMonitorAssignments(poolID)
	if err != nil {
		return nil, err
	}

	monitors, err := c.listHealthMonitors(edgeGatewayID)
	if err != nil {
		return nil, err
	}

	healthMonitors := make([]PoolModelHealthMonitor, 0, len(assignments))
	for _, assignment := range assignments {
		healthMonitors = append(healthMonitors, PoolModelHealthMonitor{
			Name: assignment.HealthMonitorRef.Name,
			Type: healthMonitorType(assignment.HealthMonitorRef, monitors),
		})
	}

	return healthMonitors, nil
}

// AssignPoolHealthMonitor assigns a health monitor to a Load Balancer Pool.
// It does nothing when the health monitor is already assigned.
func (c *client) AssignPoolHealthMonitor(ctx context.Context, edgeGatewayID, poolID string, monitorType PoolHealthMonitorType) error {
	if err := c.checkHealthMonitorArguments(edgeGatewayID, poolID); err != nil {
		return err
	}

	if monitorType == "" {
		return fmt.Errorf("monitorType is %w. Please provide a valid monitorType", errors.ErrEmpty)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	monitors, err := c.listHealthMonitors(edgeGatewayID)
	if err != nil {
		return err
	}

	monitor, err := findHealthMonitor(monitors, monitorType)
	if err != nil {
		return err
	}

	assignments, err := c.listHealthMonitorAssignments(poolID)
	if err != nil {
		return err
	}

	if slices.ContainsFunc(assignments, healthMonitorIsAssigned(monitor.ID)) {
		return nil
	}

	return c.assignHealthMonitor(monitor, poolID)
}

// UnassignPoolHealthMonitor removes a health monitor assignment from a Load Balancer Pool.
// It does nothing when the health monitor is not assigned.
func (c *client) UnassignPoolHealthMonitor(ctx context.Context, edgeGatewayID, poolID string, monitorType PoolHealthMonitorType) error {
	if err := c.checkHealthMonitorArguments(edgeGatewayID, poolID); err != nil {
		return err
	}

	if monitorType == "" {
		return fmt.Errorf("monitorType is %w. Please provide a valid monitorType", errors.ErrEmpty)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	monitors, err := c.listHealthMonitors(edgeGatewayID)
	if err != nil {
		return err
	}

	monitor, err := findHealthMonitor(monitors, monitorType)
	if err != nil {
		return err
	}

	assignments, err := c.listHealthMonitorAssignments(poolID)
	if err != nil {
		return err
	}

	for _, assignment := range assignments {
		if assignment.HealthMonitorRef.ID != monitor.ID {
			continue
		}
		if err := c.deleteHealthMonitorAssignment(assignment.ID); err != nil {
			return err
		}
	}

	return nil
}

// ReconcilePoolHealthMonitors makes the health monitors assigned to a Load Balancer
// Pool match the wanted list: the missing ones are assigned and the others are
// unassigned. Health monitors unknown to the Edge Gateway are left untouched because
// they are not managed by this SDK.
func (c *client) ReconcilePoolHealthMonitors(ctx context.Context, edgeGatewayID, poolID string, wanted []PoolHealthMonitorType) error {
	if err := c.checkHealthMonitorArguments(edgeGatewayID, poolID); err != nil {
		return err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	return c.reconcileHealthMonitors(edgeGatewayID, poolID, wanted)
}

// reconcileHealthMonitors applies the reconcile without refreshing the session,
// CreatePool and UpdatePool having refreshed it already.
func (c *client) reconcileHealthMonitors(edgeGatewayID, poolID string, wanted []PoolHealthMonitorType) error {
	monitors, err := c.listHealthMonitors(edgeGatewayID)
	if err != nil {
		return err
	}

	assignments, err := c.listHealthMonitorAssignments(poolID)
	if err != nil {
		return err
	}

	// Resolve everything before any change to avoid a partial update.
	wantedIDs := make(map[string]bool, len(wanted))
	monitorsToAssign := make([]*poolHealthMonitor, 0, len(wanted))
	assigned := make(map[string]bool, len(assignments))
	for _, assignment := range assignments {
		assigned[assignment.HealthMonitorRef.ID] = true
	}
	for _, monitorType := range wanted {
		monitor, err := findHealthMonitor(monitors, monitorType)
		if err != nil {
			return err
		}
		if wantedIDs[monitor.ID] {
			continue
		}
		wantedIDs[monitor.ID] = true

		if !assigned[monitor.ID] {
			monitorsToAssign = append(monitorsToAssign, monitor)
		}
	}

	for _, monitor := range monitorsToAssign {
		if err := c.assignHealthMonitor(monitor, poolID); err != nil {
			return err
		}
	}

	for _, assignment := range assignments {
		if wantedIDs[assignment.HealthMonitorRef.ID] {
			continue
		}
		// A health monitor unknown to this Edge Gateway is not managed by this SDK.
		if _, found := findHealthMonitorByID(monitors, assignment.HealthMonitorRef.ID); !found {
			continue
		}
		if err := c.deleteHealthMonitorAssignment(assignment.ID); err != nil {
			return err
		}
	}

	return nil
}

// assignHealthMonitor assigns an already resolved health monitor to a Load Balancer Pool.
func (c *client) assignHealthMonitor(monitor *poolHealthMonitor, poolID string) error {
	urlRef, err := c.clientGoVCDOpenAPI.OpenApiBuildEndpoint(healthMonitorAssignmentsEndpoint)
	if err != nil {
		return fmt.Errorf("error building Load Balancer Pool health monitor assignment endpoint: %w", err)
	}

	payload := &poolHealthMonitorAssignment{
		HealthMonitorRef:    govcdtypes.OpenApiReference{ID: monitor.ID, Name: monitor.Name},
		LoadBalancerPoolRef: govcdtypes.OpenApiReference{ID: poolID},
	}

	task, err := c.clientGoVCDOpenAPI.OpenApiPostItemAsync(c.apiVersion, urlRef, nil, payload)
	if err != nil {
		return fmt.Errorf("error assigning the %q Load Balancer Pool health monitor: %w", monitor.Type, err)
	}

	if err := waitTask(task); err != nil {
		return fmt.Errorf("error waiting the %q Load Balancer Pool health monitor assignment: %w", monitor.Type, err)
	}

	return nil
}

// deleteHealthMonitorAssignment removes one health monitor assignment.
func (c *client) deleteHealthMonitorAssignment(assignmentID string) error {
	urlRef, err := c.clientGoVCDOpenAPI.OpenApiBuildEndpoint(healthMonitorAssignmentEndpoint + assignmentID)
	if err != nil {
		return fmt.Errorf("error building Load Balancer Pool health monitor assignment endpoint: %w", err)
	}

	if err := c.clientGoVCDOpenAPI.OpenApiDeleteItem(c.apiVersion, urlRef, nil, nil); err != nil {
		return fmt.Errorf("error unassigning the Load Balancer Pool health monitor: %w", err)
	}

	return nil
}

// listHealthMonitorAssignments retrieves the health monitor assignments of a Load Balancer Pool.
func (c *client) listHealthMonitorAssignments(poolID string) ([]poolHealthMonitorAssignment, error) {
	urlRef, err := c.clientGoVCDOpenAPI.OpenApiBuildEndpoint(healthMonitorAssignmentsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("error building Load Balancer Pool health monitor assignments endpoint: %w", err)
	}

	queryParameters := url.Values{
		"filter": []string{fmt.Sprintf("(loadBalancerPoolRef.id==%s)", poolID)},
	}

	var assignments []poolHealthMonitorAssignment
	if err := c.clientGoVCDOpenAPI.OpenApiGetAllItems(c.apiVersion, urlRef, queryParameters, &assignments, nil); err != nil {
		return nil, fmt.Errorf("error retrieving Load Balancer Pool health monitor assignments: %w", err)
	}

	return assignments, nil
}

// listHealthMonitors retrieves the health monitors available for an Edge Gateway.
func (c *client) listHealthMonitors(edgeGatewayID string) ([]poolHealthMonitor, error) {
	urlRef, err := c.clientGoVCDOpenAPI.OpenApiBuildEndpoint(fmt.Sprintf(healthMonitorSummariesEndpointFormat, edgeGatewayID))
	if err != nil {
		return nil, fmt.Errorf("error building Load Balancer health monitors endpoint: %w", err)
	}

	var monitors []poolHealthMonitor
	if err := c.clientGoVCDOpenAPI.OpenApiGetAllItems(c.apiVersion, urlRef, nil, &monitors, nil); err != nil {
		return nil, fmt.Errorf("error retrieving Load Balancer health monitors: %w", err)
	}

	return monitors, nil
}

// checkHealthMonitorArguments validates the arguments shared by all the health
// monitor methods.
func (c *client) checkHealthMonitorArguments(edgeGatewayID, poolID string) error {
	if edgeGatewayID == "" {
		return fmt.Errorf("edgeGatewayID is %w. Please provide a valid edgeGatewayID", errors.ErrEmpty)
	}
	if !urn.IsEdgeGateway(edgeGatewayID) {
		return fmt.Errorf("edgeGatewayID has %w. Please provide a valid edgeGatewayID", errors.ErrInvalidFormat)
	}
	if poolID == "" {
		return fmt.Errorf("poolID is %w. Please provide a valid poolID", errors.ErrEmpty)
	}
	if !urn.IsLoadBalancerPool(poolID) {
		return fmt.Errorf("poolID has %w. Please provide a valid poolID", errors.ErrInvalidFormat)
	}

	return nil
}

// healthMonitorIsAssigned returns a predicate telling whether a health monitor is
// already assigned to the pool.
func healthMonitorIsAssigned(monitorID string) func(poolHealthMonitorAssignment) bool {
	return func(assignment poolHealthMonitorAssignment) bool {
		return assignment.HealthMonitorRef.ID == monitorID
	}
}

// findHealthMonitor returns the health monitor matching the given type.
func findHealthMonitor(monitors []poolHealthMonitor, monitorType PoolHealthMonitorType) (*poolHealthMonitor, error) {
	for i := range monitors {
		if strings.EqualFold(monitors[i].Type, string(monitorType)) {
			return &monitors[i], nil
		}
	}

	supported := make([]string, 0, len(PoolHealthMonitorTypes))
	for _, t := range PoolHealthMonitorTypes {
		supported = append(supported, string(t))
	}

	return nil, fmt.Errorf("monitorType %q has %w. This health monitor is not available on the Edge Gateway. Supported values are: %s",
		monitorType, errors.ErrInvalidFormat, strings.Join(supported, ", "))
}

// findHealthMonitorByID returns the health monitor matching the given identifier.
func findHealthMonitorByID(monitors []poolHealthMonitor, monitorID string) (*poolHealthMonitor, bool) {
	for i := range monitors {
		if monitors[i].ID == monitorID {
			return &monitors[i], true
		}
	}

	return nil, false
}

// healthMonitorType resolves the type of an assigned health monitor.
func healthMonitorType(ref govcdtypes.OpenApiReference, monitors []poolHealthMonitor) PoolHealthMonitorType {
	if monitor, found := findHealthMonitorByID(monitors, ref.ID); found {
		return PoolHealthMonitorType(monitor.Type)
	}

	// The health monitor is unknown to the Edge Gateway. Its name is the only
	// information provided by the platform (e.g. `System-TCP`).
	name := ref.Name
	if index := strings.LastIndex(name, "-"); index >= 0 {
		name = name[index+1:]
	}

	return PoolHealthMonitorType(strings.ToUpper(name))
}
