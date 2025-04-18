package endpoints

// List of API endpoints.
const (
	EdgeGatewayCreateFromVDC      = "/api/customers/v2.0/vdcs/{vdc-name}/edges"
	EdgeGatewayCreateFromVDCGroup = "/api/customers/v2.0/vdc-groups/{vdc-group-name}/edges"
	EdgeGatewayGet                = "/api/customers/v2.0/edges/{edge-id}"
	EdgeGatewayList               = "/api/customers/v2.0/edges"
	EdgeGatewayDelete             = EdgeGatewayGet
	EdgeGatewayUpdate             = EdgeGatewayGet

	NetworkServiceGet    = "/api/customers/v2.0/network"
	NetworkServiceCreate = "/api/customers/v2.0/services"
	NetworkServiceDelete = "/api/customers/v2.0/services/{service-id}"
)
