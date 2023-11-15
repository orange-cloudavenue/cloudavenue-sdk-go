package v1

import (
	"fmt"

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

type OSEError struct {
	Status  OSEErrorStatus `json:"status"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
}

type OSEErrorStatus int

func (e *OSEError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *OSEError) GetStatus() OSEErrorStatus {
	return e.Status
}

func (e *OSEError) GetCode() string {
	return e.Code
}

func (e *OSEError) GetMessage() string {
	return e.Message
}

// IsNotFountError returns true if the error is a 404 error
func (e *OSEError) IsNotFountError() bool {
	return e.Status == 404
}
