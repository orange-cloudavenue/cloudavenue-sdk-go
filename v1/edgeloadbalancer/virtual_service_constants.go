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
