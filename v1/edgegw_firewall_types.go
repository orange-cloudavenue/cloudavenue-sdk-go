package v1

import "github.com/vmware/go-vcloud-director/v2/govcd"

type (
	EdgeGatewayFirewall struct {
		client *EdgeClient
		*govcd.NsxtFirewall
	}
)
