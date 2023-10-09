package commonnetbackup

import "fmt"

type APIError []struct {
	PropertyName string `json:"PropertyName"`
	Message      string `json:"Message"`
}

// Convert to Error
func ToError(err *APIError) error {
	return fmt.Errorf("API Error: %v", err)
}
