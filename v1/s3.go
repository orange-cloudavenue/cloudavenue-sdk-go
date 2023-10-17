package v1

import (
	"github.com/aws/aws-sdk-go/service/s3"
	clientS3 "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/s3"
)

func (v *V1) S3() *s3.S3 {
	c, _ := clientS3.New()
	return c.S3
}
