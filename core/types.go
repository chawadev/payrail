package core

import "errors"

// ChargeResponse represents what comes back from provider
type ChargeResponse struct {
	Status        string
	TransactionID string
	Reference     string
	RawResponse   []byte
}

type Operator string

const (
	OperatorAirtel Operator = "airtel"
	OperatorMTN    Operator = "mtn"
)

type Country string

const (
	CountryZM Country = "zm"
	CountryMW Country = "mw"
)

type Bearer string

const (
	BearerMerchant Bearer = "merchant"
	BearerCustomer Bearer = "customer"
)

type ChargeRequest struct {
	Amount     float64
	Reference  string
	Phone      string
	Operator   Operator
	Country    Country
	Bearer     Bearer
	CustomerID string
}

// Validate ensures the request is safe before sending to provider
func (r ChargeRequest) Validate() error {
	if r.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if r.CustomerID == "" {
		return errors.New("customer ID is required")
	}

	if r.Reference == "" {
		return errors.New("reference is required")
	}

	if len(r.Reference) < 6 {
		return errors.New("reference too short")
	}

	return nil
}
