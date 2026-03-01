package lenco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chawadev/payrail/payrailCore"
)

// LencoProvider implements payrailCore.Provider and is configured with an API key
// and optional base URL.  It encapsulates all the details required to talk
// to Lenco's mobile-money collection endpoint.

type LencoProvider struct {
	config Config
}

// NewProvider returns a provider instance. If BaseURL is empty it defaults
// to the public Lenco API host.
func NewProvider(cfg Config) *LencoProvider {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.lenco.co"
	}
	return &LencoProvider{config: cfg}
}

// internal payload used by the mobile-money endpoint.  The fields mirror the
// JSON body that the example in the user's request shows.
type chargePayload struct {
	Operator  string  `json:"operator"`
	Bearer    string  `json:"bearer"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
	Country   string  `json:"country"`
	Phone     string  `json:"phone"`
}

// minimal response structure used to decode the portions of Lenco's JSON
// response that we care about.  We ignore most fields because the framework
// only needs status, ids and a copy of the raw bytes for debugging.
type lencoResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID             string `json:"id"`
		Reference      string `json:"reference"`
		LencoReference string `json:"lencoReference"`
		Status         string `json:"status"`
	} `json:"data"`
}

// Charge constructs a JSON request from the generic ChargeRequest, sends it to
// the Lenco mobile‑money endpoint and translates the provider response into a
// payrailCore.ChargeResponse.
func (l *LencoProvider) Charge(req payrailCore.ChargeRequest) (*payrailCore.ChargeResponse, error) {
	// build request body
	payload := chargePayload{
		Operator:  string(req.Operator),
		Bearer:    string(req.Bearer),
		Amount:    req.Amount,
		Reference: req.Reference,
		Country:   string(req.Country),
		Phone:     req.Phone,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := l.config.BaseURL + "/access/v2/collections/mobile-money"
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+l.config.APIKey)

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("lenco api returned status %d: %s", res.StatusCode, string(respBytes))
	}

	var lr lencoResponse
	if err := json.Unmarshal(respBytes, &lr); err != nil {
		return nil, err
	}

	return &payrailCore.ChargeResponse{
		Status:        lr.Data.Status,
		TransactionID: lr.Data.ID,
		Reference:     lr.Data.LencoReference,
		RawResponse:   respBytes,
	}, nil
}

// Veryfi checks the status of a payment using the reference.  It queries
// the Lenco status endpoint and returns the current payment state.
func (l *LencoProvider) Veryfi(reference string) (*payrailCore.StatusResponse, error) {
	if reference == "" {
		return nil, fmt.Errorf("reference cannot be empty")
	}

	url := l.config.BaseURL + "/access/v2/collections/status/" + reference
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+l.config.APIKey)

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("lenco api returned status %d: %s", res.StatusCode, string(respBytes))
	}

	var lr lencoResponse
	if err := json.Unmarshal(respBytes, &lr); err != nil {
		return nil, err
	}

	return &payrailCore.StatusResponse{
		Status:        lr.Data.Status,
		TransactionID: lr.Data.ID,
		Reference:     lr.Data.LencoReference,
		RawResponse:   respBytes,
	}, nil
}
