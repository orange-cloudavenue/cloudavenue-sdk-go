/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

import govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

type (
	// PoolModel represents an ALB Pool to an Edge Gateway.
	PoolModel struct {
		ID          string
		Name        string
		Description string

		GatewayRef govcdtypes.OpenApiReference
		Enabled    *bool

		// The heart of a load balancer is its ability to effectively distribute traffic across healthy servers. If persistence is enabled, only the first connection from a client is load balanced. While the persistence remains in effect, subsequent connections or requests from a client are directed to the same server.
		// Default value is PoolAlgorithmLeastConnections.
		// Supported algorithms are:
		// * PoolAlgorithmLeastConnections
		// * PoolAlgorithmRoundRobin
		// * PoolAlgorithmConsistentHash
		// * PoolAlgorithmFastestResponse
		// * PoolAlgorithmLeastLoad
		// * PoolAlgorithmFewestServers
		// * PoolAlgorithmRandom
		// * PoolAlgorithmFewestTasks
		// * PoolAlgorithmCoreAffinity
		Algorithm PoolAlgorithm

		// DefaultPort defines destination server port used by the traffic sent to the member.
		DefaultPort *int

		// GracefulTimeoutPeriod sets maximum time (in minutes) to gracefully disable a member. Virtual service waits for the
		// specified time before terminating the existing connections to the pool members that are disabled.
		//
		// Special values: 0 represents Immediate, -1 represents Infinite.
		GracefulTimeoutPeriod *int

		// PassiveMonitoringEnabled sets if client traffic should be used to check if pool member is up or down.
		PassiveMonitoringEnabled *bool

		// HealthMonitors check member servers health. It can be monitored by using one or more health monitors. Active
		// monitors generate synthetic traffic and mark a server up or down based on the response.
		HealthMonitors []PoolModelHealthMonitor

		// Members field defines list of destination servers which are used by the Load Balancer Pool to direct load balanced
		// traffic.
		//
		// Note. Only one of Members or MemberGroupRef can be specified
		Members []PoolModelMember

		// MemberGroupRef contains reference to the Edge Firewall Group (`types.NsxtFirewallGroup`)
		// representing destination servers which are used by the Load Balancer Pool to direct load
		// balanced traffic.
		//
		// Note. Only one of Members or MemberGroupRef can be specified
		MemberGroupRef *govcdtypes.OpenApiReference

		// CaCertificateRefs point to root certificates to use when validating certificates presented by the pool members.
		CaCertificateRefs []govcdtypes.OpenApiReference

		// CommonNameCheckEnabled specifies whether to check the common name of the certificate presented by the pool member.
		// This cannot be enabled if no caCertificateRefs are specified.
		CommonNameCheckEnabled *bool

		// DomainNames holds a list of domain names which will be used to verify the common names or subject alternative
		// names presented by the pool member certificates. It is performed only when common name check
		// (CommonNameCheckEnabled) is enabled. If common name check is enabled, but domain names are not specified then the
		// incoming host header will be used to check the certificate.
		DomainNames []string

		// PersistenceProfile of a Load Balancer Pool. Persistence profile will ensure that the same user sticks to the same
		// server for a desired duration of time. If the persistence profile is unmanaged by Cloud Director, updates that
		// leave the values unchanged will continue to use the same unmanaged profile. Any changes made to the persistence
		// profile will cause Cloud Director to switch the pool to a profile managed by Cloud Director.
		PersistenceProfile *PoolModelPersistenceProfile

		// MemberCount is a read only value that reports number of members added
		MemberCount int

		// EnabledMemberCount is a read only value that reports number of enabled members
		EnabledMemberCount int

		// UpMemberCount is a read only value that reports number of members that are serving traffic
		UpMemberCount int

		// HealthMessage shows a pool health status (e.g. "The pool is unassigned.")
		HealthMessage string

		// VirtualServiceRefs holds list of Load Balancer Virtual Services associated with this Load balancer Pool.
		VirtualServiceRefs []govcdtypes.OpenApiReference

		// SslEnabled is required when CA Certificates are used starting with API V37.0
		SSLEnabled *bool
	}

	// PoolModelHealthMonitor checks member servers health.
	// Active monitor generates synthetic traffic and mark a server up or down based on the response.
	PoolModelHealthMonitor struct {
		Name string
		// Type
		// * PoolHealthMonitorTypeHTTP - HTTP request/response is used to validate health.
		// * PoolHealthMonitorTypeHTTPS - Used against HTTPS encrypted web servers to validate health.
		// * PoolHealthMonitorTypeTCP - TCP connection is used to validate health.
		// * PoolHealthMonitorTypeUDP - A UDP datagram is used to validate health.
		// * PoolHealthMonitorTypePING - An ICMP ping is used to validate health.
		Type PoolHealthMonitorType
	}

	// PoolModelMember defines a single destination server which is used by the Load Balancer
	// Pool to direct load balanced traffic.
	PoolModelMember struct {
		// Enabled defines if member is enabled (will receive incoming requests) or not
		Enabled bool
		// IPAddress of the Load Balancer Pool member.
		IPAddress string

		// Port number of the Load Balancer Pool member.
		// If unset, the port that the client used to connect will be used.
		Port int

		// Ratio of selecting eligible servers in the pool.
		Ratio *int

		// MarkedDownBy gives the names of the health monitors that marked the member as down when it is DOWN.
		// If a monitor cannot be determined, the value will be UNKNOWN.
		MarkedDownBy []string

		// HealthStatus of the pool member. Possible values are:
		// * UP - The member is operational
		// * DOWN - The member is down
		// * DISABLED - The member is disabled
		// * UNKNOWN - The state is unknown
		HealthStatus string

		// DetailedHealthMessage contains non-localized detailed message on the health of the pool member.
		DetailedHealthMessage string
	}

	// PoolModelPersistenceProfile olds Persistence Profile of a Load Balancer Pool. Persistence profile will ensure that
	// the same user sticks to the same server for a desired duration of time. If the persistence profile is unmanaged by
	// Cloud Director, updates that leave the values unchanged will continue to use the same unmanaged profile. Any changes
	// made to the persistence profile will cause Cloud Director to switch the pool to a profile managed by Cloud Director.
	PoolModelPersistenceProfile struct {
		// Name field is tricky. It remains empty in some case, but if it is sent it can become computed.
		// (e.g. setting 'CUSTOM_HTTP_HEADER' results in value being
		// 'VCD-LoadBalancer-3510eae9-53bb-49f1-b7aa-7aedf5ce3a77-CUSTOM_HTTP_HEADER')
		Name string

		// Type of persistence strategy to use. Supported values are:
		// * PoolPersistenceProfileTypeClientIP - The clients IP is used as the identifier and mapped to the server.
		// * PoolPersistenceProfileTypeHTTPCookie - Load Balancer inserts a cookie into HTTP responses. Cookie name must be provided as value.
		// * PoolPersistenceProfileTypeCustomHTTPHeader - Custom, static mappings of header values to specific servers are used. Header name must be provided as value.
		// * PoolPersistenceProfileTypeAPPCookie - Load Balancer reads existing server cookies or URI embedded data such as JSessionID. Cookie name must be provided as value.
		// * PoolPersistenceProfileTypeTLS - Information is embedded in the client's SSL/TLS ticket ID. This will use default system profile System-Persistence-TLS.
		Type PoolPersistenceProfileType

		// Value of attribute based on selected persistence type.
		// This is required for PoolPersistenceProfileTypeHTTPCookie, PoolPersistenceProfileTypeCustomHTTPHeader and PoolPersistenceProfileTypeAPPCookie persistence types.
		//
		// PoolPersistenceProfileTypeHTTPCookie, PoolPersistenceProfileTypeAPPCookie must have cookie name set as the value and PoolPersistenceProfileTypeCustomHTTPHeader must have header name set as
		// the value.
		Value string
	}

	PoolAlgorithm              string
	PoolHealthMonitorType      string
	PoolPersistenceProfileType string
)

var (
	PoolAlgorithms = map[PoolAlgorithm]string{
		PoolAlgorithmLeastConnections: "New connections are sent to the server that currently has the least number of outstanding concurrent connections.",
		PoolAlgorithmRoundRobin:       "New connections are sent to the next eligible server in the pool in sequential order.",
		PoolAlgorithmConsistentHash:   "New connections are distributed across the servers by using the IP address of the client to generate an IP hash.",
		PoolAlgorithmFastestResponse:  "New connections are sent to the server that is currently providing the fastest response to new connections or requests.",
		PoolAlgorithmLeastLoad:        "New connections are sent to the server with the lightest load, regardless of the number of connections that server has.",
		PoolAlgorithmFewestServers:    "Instead of attempting to distribute all connections or requests across all servers, the fewest number of servers which are required to satisfy the current client load will be determined.",
		PoolAlgorithmRandom:           "Picks servers at random",
		PoolAlgorithmFewestTasks:      "Load is adaptively balanced, based on server feedback.",
		PoolAlgorithmCoreAffinity:     "Each CPU core uses a subset of servers, and each server is used by a subset of cores. Essentially it provides a many-to-many mapping between servers and cores.",
	}

	PoolHealthMonitorTypes = []PoolHealthMonitorType{
		PoolHealthMonitorTypeHTTP,
		PoolHealthMonitorTypeHTTPS,
		PoolHealthMonitorTypeTCP,
		PoolHealthMonitorTypeUDP,
		PoolHealthMonitorTypePING,
	}

	PoolPersistenceProfileTypes = map[PoolPersistenceProfileType]string{
		PoolPersistenceProfileTypeClientIP:         "The clients IP is used as the identifier and mapped to the server.",
		PoolPersistenceProfileTypeHTTPCookie:       "Load Balancer inserts a cookie into HTTP responses. Cookie name must be provided as value.",
		PoolPersistenceProfileTypeCustomHTTPHeader: "Custom, static mappings of header values to specific servers are used. Header name must be provided as value.",
		PoolPersistenceProfileTypeAPPCookie:        "Load Balancer reads existing server cookies or URI embedded data such as JSessionID. Cookie name must be provided as value.",
		PoolPersistenceProfileTypeTLS:              "Information is embedded in the client's SSL/TLS ticket ID. This will use default system profile System-Persistence-TLS.",
	}
)
