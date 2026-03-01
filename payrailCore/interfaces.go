package payrailCore

type Provider interface {
	Charge(req ChargeRequest) (*ChargeResponse, error)
	Veryfi(reference string) (*StatusResponse, error)
}
