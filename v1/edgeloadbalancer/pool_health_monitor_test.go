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
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

const monitorURNPrefix = "urn:vcloud:loadBalancerHealthMonitor:"

var (
	testMonitorIDHTTP = monitorURNPrefix + uuid.New().String()
	testMonitorIDTCP  = monitorURNPrefix + uuid.New().String()
	testMonitorIDPING = monitorURNPrefix + uuid.New().String()

	testAssignmentIDHTTP = monitorURNPrefix + uuid.New().String()
	testAssignmentIDPING = monitorURNPrefix + uuid.New().String()

	testHealthMonitors = []poolHealthMonitor{
		{ID: testMonitorIDHTTP, Name: "System-HTTP", Type: string(PoolHealthMonitorTypeHTTP)},
		{ID: testMonitorIDTCP, Name: "System-TCP", Type: string(PoolHealthMonitorTypeTCP)},
		{ID: testMonitorIDPING, Name: "System-Ping", Type: string(PoolHealthMonitorTypePING)},
	}

	healthMonitorAssignmentsURL = &url.URL{Path: "/" + healthMonitorAssignmentsEndpoint}
)

func stubWaitTask(t *testing.T) {
	t.Helper()
	previous := waitTask
	waitTask = func(govcd.Task) error { return nil }
	t.Cleanup(func() { waitTask = previous })
}

// stubHealthMonitorAPI answers any health monitor API exchange permissively:
// the gateway serves testHealthMonitors and the pool has no assignment.
// Used by the pool lifecycle tests that do not target the health monitors.
func stubHealthMonitorAPI(t *testing.T, clientCAV *MockclientFake) {
	t.Helper()
	stubWaitTask(t)

	urlRef := &url.URL{Path: "/stub"}
	clientCAV.EXPECT().OpenApiBuildEndpoint(gomock.Any()).Return(urlRef, nil).AnyTimes()
	clientCAV.EXPECT().
		OpenApiGetAllItems(clientcloudavenue.DefaultVCDAPIVersion, gomock.Eq(urlRef), gomock.Nil(), gomock.AssignableToTypeOf(&[]poolHealthMonitor{}), gomock.Nil()).
		DoAndReturn(func(_ string, _ *url.URL, _ url.Values, outType interface{}, _ map[string]string) error {
			*outType.(*[]poolHealthMonitor) = testHealthMonitors
			return nil
		}).AnyTimes()
	clientCAV.EXPECT().
		OpenApiGetAllItems(clientcloudavenue.DefaultVCDAPIVersion, gomock.Eq(urlRef), gomock.Not(gomock.Nil()), gomock.AssignableToTypeOf(&[]poolHealthMonitorAssignment{}), gomock.Nil()).
		DoAndReturn(func(_ string, _ *url.URL, _ url.Values, outType interface{}, _ map[string]string) error {
			*outType.(*[]poolHealthMonitorAssignment) = nil
			return nil
		}).AnyTimes()
	clientCAV.EXPECT().
		OpenApiPostItemAsync(clientcloudavenue.DefaultVCDAPIVersion, gomock.Eq(urlRef), gomock.Nil(), gomock.Any()).
		Return(govcd.Task{}, nil).AnyTimes()
}

func expectHealthMonitorSummariesList(clientCAV *MockclientFake, edgeGatewayID string, monitors []poolHealthMonitor) {
	urlRef := &url.URL{Path: "/" + fmt.Sprintf(healthMonitorSummariesEndpointFormat, edgeGatewayID)}
	clientCAV.EXPECT().
		OpenApiBuildEndpoint(gomock.Eq(fmt.Sprintf(healthMonitorSummariesEndpointFormat, edgeGatewayID))).
		Return(urlRef, nil)
	clientCAV.EXPECT().
		OpenApiGetAllItems(clientcloudavenue.DefaultVCDAPIVersion, gomock.Eq(urlRef), gomock.Nil(), gomock.AssignableToTypeOf(&[]poolHealthMonitor{}), gomock.Nil()).
		DoAndReturn(func(_ string, _ *url.URL, _ url.Values, outType interface{}, _ map[string]string) error {
			*outType.(*[]poolHealthMonitor) = monitors
			return nil
		})
}

func expectHealthMonitorAssignmentsList(clientCAV *MockclientFake, poolID string, assignments []poolHealthMonitorAssignment) {
	clientCAV.EXPECT().
		OpenApiBuildEndpoint(gomock.Eq(healthMonitorAssignmentsEndpoint)).
		Return(healthMonitorAssignmentsURL, nil).AnyTimes()
	clientCAV.EXPECT().
		OpenApiGetAllItems(clientcloudavenue.DefaultVCDAPIVersion, gomock.Eq(healthMonitorAssignmentsURL),
			gomock.Eq(url.Values{"filter": []string{"(loadBalancerPoolRef.id==" + poolID + ")"}}),
			gomock.AssignableToTypeOf(&[]poolHealthMonitorAssignment{}), gomock.Nil()).
		DoAndReturn(func(_ string, _ *url.URL, _ url.Values, outType interface{}, _ map[string]string) error {
			*outType.(*[]poolHealthMonitorAssignment) = assignments
			return nil
		})
}

func expectHealthMonitorAssign(clientCAV *MockclientFake, monitorID, monitorName, poolID string) {
	clientCAV.EXPECT().
		OpenApiPostItemAsync(clientcloudavenue.DefaultVCDAPIVersion, gomock.Eq(healthMonitorAssignmentsURL), gomock.Nil(),
			gomock.Eq(&poolHealthMonitorAssignment{
				HealthMonitorRef:    govcdtypes.OpenApiReference{ID: monitorID, Name: monitorName},
				LoadBalancerPoolRef: govcdtypes.OpenApiReference{ID: poolID},
			})).
		Return(govcd.Task{}, nil)
}

func expectHealthMonitorUnassign(clientCAV *MockclientFake, assignmentID string) {
	urlRef := &url.URL{Path: "/" + healthMonitorAssignmentEndpoint + assignmentID}
	clientCAV.EXPECT().
		OpenApiBuildEndpoint(gomock.Eq(healthMonitorAssignmentEndpoint+assignmentID)).
		Return(urlRef, nil)
	clientCAV.EXPECT().
		OpenApiDeleteItem(clientcloudavenue.DefaultVCDAPIVersion, gomock.Eq(urlRef), gomock.Nil(), gomock.Nil()).
		Return(nil)
}

func TestClient_ListPoolHealthMonitors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientCAV := NewMockclientFake(ctrl)
	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()

	t.Run("success", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, []poolHealthMonitorAssignment{
			{ID: testAssignmentIDHTTP, HealthMonitorRef: govcdtypes.OpenApiReference{ID: testMonitorIDHTTP, Name: "System-HTTP"}},
		})
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)

		monitors, err := c.ListPoolHealthMonitors(context.Background(), edgeGatewayID, poolID)
		assert.NoError(t, err)
		assert.Equal(t, []PoolModelHealthMonitor{{Name: "System-HTTP", Type: PoolHealthMonitorTypeHTTP}}, monitors)
	})

	t.Run("unknown-monitor-falls-back-on-name", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, []poolHealthMonitorAssignment{
			{ID: testAssignmentIDPING, HealthMonitorRef: govcdtypes.OpenApiReference{ID: monitorURNPrefix + "unknown", Name: "Custom-TCP"}},
		})
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)

		monitors, err := c.ListPoolHealthMonitors(context.Background(), edgeGatewayID, poolID)
		assert.NoError(t, err)
		assert.Equal(t, []PoolModelHealthMonitor{{Name: "Custom-TCP", Type: PoolHealthMonitorTypeTCP}}, monitors)
	})

	t.Run("error-invalid-arguments", func(t *testing.T) {
		_, err := c.ListPoolHealthMonitors(context.Background(), "", poolID)
		assert.Error(t, err)

		_, err = c.ListPoolHealthMonitors(context.Background(), edgeGatewayID, "not-an-urn")
		assert.Error(t, err)
	})
}

func TestClient_AssignPoolHealthMonitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientCAV := NewMockclientFake(ctrl)
	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()

	t.Run("success", func(t *testing.T) {
		stubWaitTask(t)
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, nil)
		expectHealthMonitorAssign(clientCAV, testMonitorIDTCP, "System-TCP", poolID)

		err := c.AssignPoolHealthMonitor(context.Background(), edgeGatewayID, poolID, PoolHealthMonitorTypeTCP)
		assert.NoError(t, err)
	})

	t.Run("already-assigned-is-a-noop", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, []poolHealthMonitorAssignment{
			{ID: testAssignmentIDHTTP, HealthMonitorRef: govcdtypes.OpenApiReference{ID: testMonitorIDHTTP, Name: "System-HTTP"}},
		})

		err := c.AssignPoolHealthMonitor(context.Background(), edgeGatewayID, poolID, PoolHealthMonitorTypeHTTP)
		assert.NoError(t, err)
	})

	t.Run("error-unknown-monitor-type", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)

		err := c.AssignPoolHealthMonitor(context.Background(), edgeGatewayID, poolID, PoolHealthMonitorTypeUDP)
		assert.Error(t, err)
	})

	t.Run("error-empty-monitor-type", func(t *testing.T) {
		err := c.AssignPoolHealthMonitor(context.Background(), edgeGatewayID, poolID, "")
		assert.Error(t, err)
	})
}

func TestClient_UnassignPoolHealthMonitor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientCAV := NewMockclientFake(ctrl)
	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()

	t.Run("success", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, []poolHealthMonitorAssignment{
			{ID: testAssignmentIDHTTP, HealthMonitorRef: govcdtypes.OpenApiReference{ID: testMonitorIDHTTP, Name: "System-HTTP"}},
		})
		expectHealthMonitorUnassign(clientCAV, testAssignmentIDHTTP)

		err := c.UnassignPoolHealthMonitor(context.Background(), edgeGatewayID, poolID, PoolHealthMonitorTypeHTTP)
		assert.NoError(t, err)
	})

	t.Run("not-assigned-is-a-noop", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, nil)

		err := c.UnassignPoolHealthMonitor(context.Background(), edgeGatewayID, poolID, PoolHealthMonitorTypeHTTP)
		assert.NoError(t, err)
	})
}

func TestClient_ReconcilePoolHealthMonitors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clientCAV := NewMockclientFake(ctrl)
	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()

	t.Run("assign-missing-and-remove-extra", func(t *testing.T) {
		stubWaitTask(t)
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, []poolHealthMonitorAssignment{
			{ID: testAssignmentIDHTTP, HealthMonitorRef: govcdtypes.OpenApiReference{ID: testMonitorIDHTTP, Name: "System-HTTP"}},
			{ID: testAssignmentIDPING, HealthMonitorRef: govcdtypes.OpenApiReference{ID: testMonitorIDPING, Name: "System-Ping"}},
		})
		expectHealthMonitorAssign(clientCAV, testMonitorIDTCP, "System-TCP", poolID)
		expectHealthMonitorUnassign(clientCAV, testAssignmentIDHTTP)

		err := c.ReconcilePoolHealthMonitors(context.Background(), edgeGatewayID, poolID,
			[]PoolHealthMonitorType{PoolHealthMonitorTypeTCP, PoolHealthMonitorTypePING})
		assert.NoError(t, err)
	})

	t.Run("no-change-is-a-noop", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, []poolHealthMonitorAssignment{
			{ID: testAssignmentIDHTTP, HealthMonitorRef: govcdtypes.OpenApiReference{ID: testMonitorIDHTTP, Name: "System-HTTP"}},
		})

		err := c.ReconcilePoolHealthMonitors(context.Background(), edgeGatewayID, poolID,
			[]PoolHealthMonitorType{PoolHealthMonitorTypeHTTP})
		assert.NoError(t, err)
	})

	t.Run("empty-wanted-list-removes-all-managed", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, []poolHealthMonitorAssignment{
			{ID: testAssignmentIDHTTP, HealthMonitorRef: govcdtypes.OpenApiReference{ID: testMonitorIDHTTP, Name: "System-HTTP"}},
		})
		expectHealthMonitorUnassign(clientCAV, testAssignmentIDHTTP)

		err := c.ReconcilePoolHealthMonitors(context.Background(), edgeGatewayID, poolID, nil)
		assert.NoError(t, err)
	})

	t.Run("unknown-monitor-type-fails-before-any-change", func(t *testing.T) {
		clientCAV.EXPECT().Refresh().Return(nil)
		expectHealthMonitorSummariesList(clientCAV, edgeGatewayID, testHealthMonitors)
		expectHealthMonitorAssignmentsList(clientCAV, poolID, nil)

		err := c.ReconcilePoolHealthMonitors(context.Background(), edgeGatewayID, poolID,
			[]PoolHealthMonitorType{PoolHealthMonitorTypeHTTPS})
		assert.Error(t, err)
	})
}
