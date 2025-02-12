package edgeloadbalancer

const (
	// Application Profile Types.
	VirtualServiceApplicationProfileHTTP  VirtualServiceModelApplicationProfile = "HTTP"
	VirtualServiceApplicationProfileHTTPS VirtualServiceModelApplicationProfile = "HTTPS"
	VirtualServiceApplicationProfileL4    VirtualServiceModelApplicationProfile = "L4"
	VirtualServiceApplicationProfileL4TLS VirtualServiceModelApplicationProfile = "L4_TLS"

	// Service Port Types.
	VirtualServiceServicePortTypeTCPProxy    VirtualServiceModelServicePortType = "TCP_PROXY"
	VirtualServiceServicePortTypeTCPFastPath VirtualServiceModelServicePortType = "TCP_FAST_PATH"
	VirtualServiceServicePortTypeUDPFastPath VirtualServiceModelServicePortType = "UDP_FAST_PATH"

	// Health Status Types.
	VirtualServiceHealthStatusUP          VirtualServiceModelHealthStatus = "UP"
	VirtualServiceHealthStatusDOWN        VirtualServiceModelHealthStatus = "DOWN"
	VirtualServiceHealthStatusRUNNING     VirtualServiceModelHealthStatus = "RUNNING"
	VirtualServiceHealthStatusUNAVAILABLE VirtualServiceModelHealthStatus = "UNAVAILABLE"
	VirtualServiceHealthStatusUNKNOWN     VirtualServiceModelHealthStatus = "UNKNOWN"
)
