package nexmo

import (
	"errors"
	"net/http"
)

// Client encapsulates the Nexmo functions - must be created with
// NewClientFromAPI()
type Client struct {
	Account    *Account
	SMS        *SMS
	USSD       *USSD
	HTTPClient *http.Client
	apiKey     string
	apiSecret  string
	useOauth   bool
}

// NewClientFromAPI creates a new Client type with the
// provided API key / API secret.
func NewClientFromAPI(apiKey, apiSecret string) (*Client, error) {
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
	c.HTTPClient = http.DefaultClient
	return c, nil
}
