package edgeloadbalancer

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func TestVirtualServiceIDValidator(t *testing.T) {
	client := &client{}

	tests := []struct {
		name             string
		virtualServiceID string
		expectedError    error
	}{
		{
			name:             "Empty virtualServiceID",
			virtualServiceID: "",
			expectedError:    errors.ErrEmpty,
		},
		{
			name:             "Invalid format virtualServiceID",
			virtualServiceID: "invalid-id",
			expectedError:    errors.ErrInvalidFormat,
		},
		{
			name:             "Valid virtualServiceID",
			virtualServiceID: urn.LoadBalancerVirtualService.String() + uuid.New().String(),
			expectedError:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.virtualServiceIDValidator(tt.virtualServiceID)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
