package edgeloadbalancer

import (
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// sliceAnyToSliceString converts a slice of any type to a slice of strings.
// It takes a generic slice as input and returns a slice of strings where each element
// is the string representation of the corresponding element in the input slice.
//
// T: any type
// input: a slice of any type T
// returns: a slice of strings
func sliceAnyToSliceString[T any](input []T) []string {
	var result []string
	for _, value := range input {
		result = append(result, fmt.Sprint(value))
	}
	return result
}

// virtualServiceIDValidator validates the given virtualServiceID.
// It returns an error if the virtualServiceID is empty or not in the correct format.
//
// Parameters:
//   - virtualServiceID: The ID of the virtual service to be validated.
//
// Returns:
//   - error: An error indicating whether the virtualServiceID is empty or has an invalid format.
func (*client) virtualServiceIDValidator(virtualServiceID string) error {
	if virtualServiceID == "" {
		return fmt.Errorf("virtualServiceID is %w. Please provide a valid virtualServiceID", errors.ErrEmpty)
	}

	if !urn.IsLoadBalancerVirtualService(virtualServiceID) {
		return fmt.Errorf("virtualServiceID has %w. Please provide a valid virtualServiceID", errors.ErrInvalidFormat)
	}

	return nil
}
