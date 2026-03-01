package payrailCore

import "errors"

// ChargeStatus represents the state of a payment after being sent to a provider.
type ChargeStatus string

const (
	StatusPending    ChargeStatus = "pending"
	StatusSuccessful ChargeStatus = "successful"
	StatusFailed     ChargeStatus = "failed"
	StatusPayOffline ChargeStatus = "pay-offline"
)

// ChargeResponse represents what comes back from provider
type ChargeResponse struct {
	Status        string // one of the ChargeStatus constants above
	TransactionID string
	Reference     string
	RawResponse   []byte
}

type Operator string

type Country string

const (
	Zambia Country = "zm"
	Malawi Country = "mw"
)

type Bearer string

const (
	BearerMerchant Bearer = "merchant"
	BearerCustomer Bearer = "customer"
)

type ChargeRequest struct {
	Amount    float64
	Reference string
	Phone     string
	Operator  string
	Country   Country
	Bearer    Bearer
}

// Validate ensures the request is safe before sending to provider
func (r ChargeRequest) Validate() error {
	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if r.Reference == "" {
		return errors.New("reference is required")
	}

	if len(r.Reference) < 6 {
		return errors.New("reference too short")
	}

	return nil
}

// StatusResponse represents the response from a status check query
type StatusResponse struct {
	Status           string // payment status: pending, successful, failed, pay-offline, 3ds-auth-required
	TransactionID    string // provider's transaction ID
	Reference        string // merchant reference
	Amount           string
	Currency         string
	Fee              string
	SettlementStatus string
	ReasonForFailure string
	RawResponse      []byte
}

// VeryfiRequest is used to query payment status
type VeryfiRequest struct {
	Reference string
}
