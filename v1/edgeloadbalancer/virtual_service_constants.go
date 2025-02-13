/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

const (
	// Application Profile Types.
	VirtualServiceApplicationProfileHTTP  VirtualServiceModelApplicationProfile = "HTTP"
	VirtualServiceApplicationProfileHTTPS VirtualServiceModelApplicationProfile = "HTTPS"
	VirtualServiceApplicationProfileL4TCP VirtualServiceModelApplicationProfile = "L4_TCP"
	VirtualServiceApplicationProfileL4UDP VirtualServiceModelApplicationProfile = "L4_UDP"
	VirtualServiceApplicationProfileL4TLS VirtualServiceModelApplicationProfile = "L4_TLS"

	// Service Port Types. (unexported).
	virtualServiceServicePortTypeTCPProxy    VirtualServiceModelServicePortType = "TCP_PROXY"
	virtualServiceServicePortTypeTCPFastPath VirtualServiceModelServicePortType = "TCP_FAST_PATH"
	virtualServiceServicePortTypeUDPFastPath VirtualServiceModelServicePortType = "UDP_FAST_PATH"

	// Health Status Types.
	VirtualServiceHealthStatusUP          VirtualServiceModelHealthStatus = "UP"
	VirtualServiceHealthStatusDOWN        VirtualServiceModelHealthStatus = "DOWN"
	VirtualServiceHealthStatusRUNNING     VirtualServiceModelHealthStatus = "RUNNING"
	VirtualServiceHealthStatusUNAVAILABLE VirtualServiceModelHealthStatus = "UNAVAILABLE"
	VirtualServiceHealthStatusUNKNOWN     VirtualServiceModelHealthStatus = "UNKNOWN"
)
