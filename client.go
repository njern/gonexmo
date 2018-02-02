package nexmo

import (
	"errors"
	"net/http"
)

// Client encapsulates the Nexmo functions.
// Should be created with NewClient()
type Client struct {
	Account    *Account
	SMS        *SMS
	USSD       *USSD
	Verify     *Verification
	HTTPClient *http.Client
	apiKey     string
	apiSecret  string
	useOauth   bool
}

// NewClient creates a new Client type with the
// provided API key / API secret.
func NewClient(apiKey, apiSecret string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey can not be empty")
	} else if apiSecret == "" {
		return nil, errors.New("apiSecret can not be empty")
	}

	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		useOauth:  false,
	}

	c.Account = &Account{c}
	c.SMS = &SMS{c}
	c.USSD = &USSD{c}
	c.Verify = &Verification{c}
	c.HTTPClient = http.DefaultClient
	return c, nil
}
