package v1

import (
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/v1/netbackup"
)

type V1 struct {
	Netbackup   netbackup.Netbackup
	PublicIP    PublicIP
	EdgeGateway EdgeGateway
	T0          Tier0
	VCDA        VCDA
	// VDC         VDC is a method of the V1 struct that returns a pointer to the CAVVdc struct
	// S3          *s3.S3 - S3 is a method of the V1 struct that returns a pointer to the AWS S3 client preconfigured
}

func (v *V1) VDC() *CAVVdc {
	return &CAVVdc{}
}
