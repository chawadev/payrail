package payrail

import "payrail/core"

func (c *Client) Charge(req core.ChargeRequest) (*core.ChargeResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	resp := &core.ChargeResponse{
		Status:    "pending",
		Reference: req.Reference,
	}
	return resp, nil
}
