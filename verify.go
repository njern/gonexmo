package nexmo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Verification wraps a client to be able to use local verify methods.
type Verification struct {
	client *Client
}

// MarshalJSON returns a byte slice with the serialized JSON of the
// VerifyMessageRequest struct.
func (m *VerifyMessageRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		APIKey    string `json:"api_key"`
		APISecret string `json:"api_secret"`
		VerifyMessageRequest
	}{
		APIKey:               m.apiKey,
		APISecret:            m.apiSecret,
		VerifyMessageRequest: *m,
	})
}

// VerifyMessageRequest is the request struct for initiating the verification process
// for a phone number.
type VerifyMessageRequest struct {
	apiKey    string
	apiSecret string

	Number        string `json:"number"`
	Brand         string `json:"brand"`
	SenderID      string `json:"sender_id,omitempty"`
	Country       string `json:"country,omitempty"`
	Language      string `json:"lg,omitempty"`
	CodeLength    int    `json:"code_length,omitempty"`
	PINExpiry     int    `json:"pin_expiry,omitempty"`
	NextEventWait int    `json:"next_event_wait,omitempty"`
}

// VerifyMessageResponse is the struct for the response from the verify
// endpoint.
type VerifyMessageResponse struct {
	Status    ResponseCode `json:"status,string"`
	RequestID string       `json:"request_id"`
	ErrorText string       `json:"error_text"`
}

// Send makes the actual HTTP request to the endpoint and returns the
// response.
func (c *Verification) Send(m *VerifyMessageRequest) (*VerifyMessageResponse, error) {
	if len(m.Number) == 0 {
		return nil, errors.New("Invalid Number field specified")
	}

	if len(m.Brand) == 0 {
		return nil, errors.New("Invalid Brand field specified")
	}

	var verifyMessageResponse *VerifyMessageResponse

	if !c.client.useOauth {
		m.apiKey = c.client.apiKey
		m.apiSecret = c.client.apiSecret
	}

	var r *http.Request
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("invalid message struct - can not convert to JSON")
	}

	b := bytes.NewBuffer(buf)
	r, err = http.NewRequest("POST", apiRootv2+"/verify/json", b)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &verifyMessageResponse)
	if err != nil {
		return nil, err
	}
	return verifyMessageResponse, nil
}

// MarshalJSON implements the json.Marshaler interface
func (m *VerifyCheckRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		APIKey    string `json:"api_key"`
		APISecret string `json:"api_secret"`
		VerifyCheckRequest
	}{
		APIKey:             m.apiKey,
		APISecret:          m.apiSecret,
		VerifyCheckRequest: *m,
	})
}

// A VerifyCheckRequest is sent to Nexmo
// when we want to verify a user has the
// phone number he says he does.
type VerifyCheckRequest struct {
	apiKey    string
	apiSecret string

	RequestID string `json:"request_id"`
	Code      string `json:"code"`
	IPAddress string `json:"ip_address,omitempty"`
}

// A VerifyCheckResponse is received from Nexmo
// after verifying a user has the
// phone number he says he does.
type VerifyCheckResponse struct {
	Status    ResponseCode `json:"status,string"`
	EventID   string       `json:"event_id"`
	Price     string       `json:"price"`
	Currency  string       `json:"currency"`
	ErrorText string       `json:"error_text"`
}

// Check (by sending a PIN to a user) whether a user can be contacted at his given phone number.
// https://developer.nexmo.com/api/verify#verify-check
func (c *Verification) Check(m *VerifyCheckRequest) (*VerifyCheckResponse, error) {
	if len(m.RequestID) == 0 {
		return nil, errors.New("Invalid RequestID field specified")
	}

	if len(m.Code) == 0 {
		return nil, errors.New("Invalid Code field specified")
	}

	var verifyCheckResponse *VerifyCheckResponse

	if !c.client.useOauth {
		m.apiKey = c.client.apiKey
		m.apiSecret = c.client.apiSecret
	}

	var r *http.Request
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("invalid message struct - unable to convert to JSON")
	}
	b := bytes.NewBuffer(buf)
	r, err = http.NewRequest("POST", apiRootv2+"/verify/check/json", b)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &verifyCheckResponse)
	if err != nil {
		return nil, err
	}
	return verifyCheckResponse, nil
}

// MarshalJSON implements the json.Marshaler interface
func (m *VerifySearchRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		APIKey    string `json:"api_key"`
		APISecret string `json:"api_secret"`
		VerifySearchRequest
	}{
		APIKey:              m.apiKey,
		APISecret:           m.apiSecret,
		VerifySearchRequest: *m,
	})
}

// A VerifySearchRequest is sent to Nexmo
// when searching for the status of a Verify
// request.
type VerifySearchRequest struct {
	apiKey    string
	apiSecret string

	RequestID string `json:"request_id,omitempty"`
}

// A VerifySearchResponse is received from Nexmo in
// response to a VerifySearchRequest
type VerifySearchResponse struct {
	RequestID      string `json:"request_id"`
	AccountID      string `json:"account_id"`
	Number         string `json:"number"`
	SenderID       string `json:"sender_id"`
	DateSubmitted  string `json:"date_submitted"`
	DateFinalized  string `json:"date_finalized"`
	FirstEventDate string `json:"first_event_date"`
	LastEventDate  string `json:"last_event_date"`
	Status         string `json:"status"`
	Checks         []struct {
		DateReceived string `json:"date_received"`
		Code         string `json:"code"`
		Status       string `json:"status"`
		IPAddress    string `json:"ip_address,omitempty"`
	} `json:"checks"`
	Price     string `json:"price"`
	Currency  string `json:"currency"`
	ErrorText string `json:"error_text"`
}

// Search sends the verify search request to Nexmo.
// https://developer.nexmo.com/api/verify#verify-search
func (c *Verification) Search(m *VerifySearchRequest) (*VerifySearchResponse, error) {
	var verifySearchResponse *VerifySearchResponse

	if !c.client.useOauth {
		m.apiKey = c.client.apiKey
		m.apiSecret = c.client.apiSecret
	}

	var r *http.Request
	buf, err := json.Marshal(m)
	if err != nil {
		return nil, errors.New("invalid message struct - unable to convert to JSON")
	}

	b := bytes.NewBuffer(buf)
	r, err = http.NewRequest("POST", apiRootv2+"/verify/search/json", b)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := c.client.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &verifySearchResponse)
	if err != nil {
		return nil, err
	}

	return verifySearchResponse, nil
}
