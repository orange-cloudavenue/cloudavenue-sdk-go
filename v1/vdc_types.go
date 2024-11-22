package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/infrapi"
)

type (
	VDC struct {
		*govcd.Vdc
		infrapi *infrapi.CAVVirtualDataCenter
	}
)
