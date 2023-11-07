package v1

import (
	"github.com/aws/aws-sdk-go/service/s3"
	clientS3 "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/s3"
)

type S3Client struct {
	*s3.S3
}

func (v *V1) S3() S3Client {
	c, _ := clientS3.New()
	return S3Client{c.S3}
}
