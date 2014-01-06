package nexmo

import (
	"errors"
)

// Client encapsulates the Nexmo functions - must be created with
// NewClientFromAPI()
type Client struct {
	apiKey    string
	apiSecret string
	useOauth  bool
}

// NewClientFromAPI creates a new Client type with the
// provided API key / API secret.
func NewClientFromAPI(apiKey, apiSecret string) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey can not be empty")
	} else if apiSecret == "" {
		return nil, errors.New("apiSecret can not be empty")
	}

	return &Client{apiKey, apiSecret, false}, nil
}
