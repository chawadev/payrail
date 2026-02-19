package payrail

import "net/http"

type Client struct {
	APIKey     string
	Provider   string
	BaseURL    string
	httpClient *http.Client
}

func NewClient(apiKey string, provider string) *Client {
	return &Client{
		APIKey:     apiKey,
		Provider:   provider,
		httpClient: &http.Client{},
	}
}
