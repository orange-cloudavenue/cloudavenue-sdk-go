package edgeloadbalancer

// PoolAlgorithm is the algorithm for choosing a member within the pools list of available members for each new connection.
const (
	// Algorithm for choosing a member within the pools list of available members for each new connection.
	// Default value is LEAST_CONNECTIONS
	// Supported algorithms are:
	// * LEAST_CONNECTIONS
	// * ROUND_ROBIN
	// * CONSISTENT_HASH (uses Source IP Address hash)
	// * FASTEST_RESPONSE
	// * LEAST_LOAD
	// * FEWEST_SERVERS
	// * RANDOM
	// * FEWEST_TASKS
	// * CORE_AFFINITY.

	// PoolAlgorithmLeastConnections is the least connections algorithm.
	PoolAlgorithmLeastConnections PoolAlgorithm = "LEAST_CONNECTIONS"
	// PoolAlgorithmRoundRobin is the round robin algorithm.
	PoolAlgorithmRoundRobin PoolAlgorithm = "ROUND_ROBIN"
	// PoolAlgorithmConsistentHash is the consistent hash algorithm.
	PoolAlgorithmConsistentHash PoolAlgorithm = "CONSISTENT_HASH"
	// PoolAlgorithmFastestResponse is the fastest response algorithm.
	PoolAlgorithmFastestResponse PoolAlgorithm = "FASTEST_RESPONSE"
	// PoolAlgorithmLeastLoad is the least load algorithm.
	PoolAlgorithmLeastLoad PoolAlgorithm = "LEAST_LOAD"
	// PoolAlgorithmFewestServers is the fewest servers algorithm.
	PoolAlgorithmFewestServers PoolAlgorithm = "FEWEST_SERVERS"
	// PoolAlgorithmRandom is the random algorithm.
	PoolAlgorithmRandom PoolAlgorithm = "RANDOM"
	// PoolAlgorithmFewestTasks is the fewest tasks algorithm.
	PoolAlgorithmFewestTasks PoolAlgorithm = "FEWEST_TASKS"
	// PoolAlgorithmCoreAffinity is the core affinity algorithm.
	PoolAlgorithmCoreAffinity PoolAlgorithm = "CORE_AFFINITY"
)

// PoolHealthMonitorType is the type of health monitor.
const (
	// PoolHealthMonitorTypeHTTP is the http health monitor.
	// HTTP request/response is used to validate health.
	PoolHealthMonitorTypeHTTP PoolHealthMonitorType = "HTTP"

	// PoolHealthMonitorTypeHTTPS is the https health monitor.
	// Used against HTTPS encrypted web servers to validate health.
	PoolHealthMonitorTypeHTTPS PoolHealthMonitorType = "HTTPS"

	// PoolHealthMonitorTypeTCP is the tcp health monitor.
	// TCP connection is used to validate health.
	PoolHealthMonitorTypeTCP PoolHealthMonitorType = "TCP"

	// PoolHealthMonitorTypeUDP is the udp health monitor.
	// UDP connection is used to validate health.
	PoolHealthMonitorTypeUDP PoolHealthMonitorType = "UDP"

	// PoolHealthMonitorTypePING is the ping health monitor.
	// ICMP ping is used to validate health.
	PoolHealthMonitorTypePING PoolHealthMonitorType = "PING"
)

// PoolPersistenceProfileType is the type of persistence profile.
const (
	// PoolPersistenceProfileTypeClientIP is the client IP persistence profile.
	PoolPersistenceProfileTypeClientIP PoolPersistenceProfileType = "CLIENT_IP"

	// PoolPersistenceProfileTypeHTTPCookie is the HTTP cookie persistence profile.
	PoolPersistenceProfileTypeHTTPCookie PoolPersistenceProfileType = "HTTP_COOKIE"

	// PoolPersistenceProfileTypeCustomHTTPHeader is the custom HTTP header persistence profile.
	PoolPersistenceProfileTypeCustomHTTPHeader PoolPersistenceProfileType = "CUSTOM_HTTP_HEADER"

	// PoolPersistenceProfileTypeAPPCookie is the APP cookie persistence profile.
	PoolPersistenceProfileTypeAPPCookie PoolPersistenceProfileType = "APP_COOKIE"

	// PoolPersistenceProfileTypeTLS is the TLS persistence profile.
	PoolPersistenceProfileTypeTLS PoolPersistenceProfileType = "TLS"
)
