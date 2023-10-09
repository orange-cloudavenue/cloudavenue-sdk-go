package commoncloudavenue

import "fmt"

type APIErrorResponse struct {
	Code    string `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// FormatError - Formats the error
func (e *APIErrorResponse) FormatError() string {
	return fmt.Sprintf("ErrorCode:%s - ErrorReason:%s - ErrorMessage:%s", e.Code, e.Reason, e.Message)
}

// ToError - Converts an APIErrorResponse to an error
func ToError(e *APIErrorResponse) error {
	return fmt.Errorf("error on API call: %s", e.FormatError())
}
