package payrail

import (
	"errors"

	"payrail/core"
	"payrail/provider/lenco"
)

// Client is a thin wrapper around a concrete payment provider.  The
// provider is chosen during construction and implements the core.Provider
// interface.
type Client struct {
	provider core.Provider
}

// NewClient creates a Payrail client bound to the named provider.  Currently
// only "lenco" is supported; passing any other string returns an error.  The
// apiKey argument will be forwarded to the provider's configuration.
func NewClient(apiKey string, provider string) (*Client, error) {
	switch provider {
	case "lenco":
		p := lenco.NewProvider(lenco.Config{APIKey: apiKey})
		return &Client{provider: p}, nil
	default:
		return nil, errors.New("unsupported provider: " + provider)
	}
}
