package edgegateway

import (
	"fmt"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

func (*client) edgeGatewayNameOrIDValidation(edgeGatewayNameOrID string) error {
	if edgeGatewayNameOrID == "" {
		return fmt.Errorf("edgeGatewayNameOrID is %w. Please provide a valid edgeGatewayNameOrID", errors.ErrEmpty)
	}

	return nil
}
