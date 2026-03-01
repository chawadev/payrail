package payrail

import (
	"errors"

	"github.com/chawadev/payrail/payrailCore"
)

// Charge validates the request and forwards it to the configured provider.
func (c *Client) Charge(req payrailCore.ChargeRequest) (*payrailCore.ChargeResponse, error) {
	if c == nil || c.provider == nil {
		return nil, errors.New("client not initialized with provider")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return c.provider.Charge(req)
}

// Veryfi checks the status of a payment using its reference.  It delegates
// to the provider to fetch the latest payment state.
func (c *Client) Veryfi(reference string) (*payrailCore.StatusResponse, error) {
	if c == nil || c.provider == nil {
		return nil, errors.New("client not initialized with provider")
	}

	if reference == "" {
		return nil, errors.New("reference is required")
	}

	return c.provider.Veryfi(reference)
}
