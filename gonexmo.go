/*
Package gonexmo implements a simple client library for accessing the Nexmo API.

Usage is simple. Create a Nexmo instance with NexmoWithKeyAndSecret(), providing
your API key and API secret. Then send messages with SendTextMessage() or
SendFlashMessage(). The API will return a MessageResponse which you can
use to see if your message went through, how much it cost, etc.
*/
package nexmo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// MessageReport is the "status report" for a single SMS sent via the Nexmo API
type MessageReport struct {
	Status           string `json:"status"`
	MessageID        string `json:"message-id"`
	To               string `json:"to"`
	ClientReference  string `json:"client-ref"`
	RemainingBalance string `json:"remaining-balance"`
	MessagePrice     string `json:"message-price"`
	Network          string `json:"network"`
	ErrorText        string `json:"error-text"`
}

// MessageResponse contains the response from Nexmo's API after we attempt to send any kind of message.
// It will contain one MessageReport for every 160 chars sent.
type MessageResponse struct {
	MessageCount string          `json:"message-count"`
	Messages     []MessageReport `json:"messages"`
}

// AccountBalance represents the "balance" object we get back when calling GET /account/get-balance
type AccountBalance struct {
	Value float64 `json:"value"`
}

// Nexmo encapsulates the Nexmo functions - must be created with NexmoWithKeyAndSecret()
type Conn struct {
	apiKey    string
	apiSecret string
}

// NewConn creates a Nexmo object with the provided API key / API secret.
func NewConn(apiKey, apiSecret string) (*Conn, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey can not be empty!")
	} else if apiSecret == "" {
		return nil, errors.New("apiSecret can not be empty!")
	}

	nexmo := &Conn{apiKey, apiSecret}
	return nexmo, nil
}

func (nexmo *Conn) sendMessage(from string, to string, text string, clientReference string, statusReportRequired bool, is_flash_message bool) (*MessageResponse, error) {
	var messageResponse *MessageResponse

	values := make(url.Values)
	values.Set("api_key", nexmo.apiKey)
	values.Set("api_secret", nexmo.apiSecret)
	values.Set("type", "text")

	values.Set("to", to)
	values.Set("from", from)
	values.Set("text", text)
	values.Set("client_ref", clientReference)

	if statusReportRequired {
		values.Set("status_report_req", string(1))

	}
	if is_flash_message {
		values.Set("message_class", "0")
	}

	client := &http.Client{}
	r, _ := http.NewRequest("POST", "https://rest.nexmo.com/sms/json", bytes.NewBufferString(values.Encode())) // <-- URL-encoded payload
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &messageResponse)
	if err != nil {
		return nil, err
	} else {
		return messageResponse, nil
	}
}

// SendTextMessage() sends a normal SMS
func (nexmo *Conn) SendTextMessage(from, to, text, clientReference string, statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference, statusReportRequired, false)
}

// SendFlashMessage() sends a class 0 SMS (Flash message).
func (nexmo *Conn) SendFlashMessage(from, to, text, clientReference string, statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference, statusReportRequired, true)
}

// GetBalance() retrieves the current balance of your Nexmo account in Euros (€)
func (nexmo *Conn) GetBalance() (float64, error) {
	var accBalance *AccountBalance

	client := &http.Client{}
	r, _ := http.NewRequest("GET", "https://rest.nexmo.com/account/get-balance/"+nexmo.apiKey+"/"+nexmo.apiSecret, nil)
	r.Header.Add("Accept", "application/json")

	resp, err := client.Do(r)
	defer resp.Body.Close()

	if err != nil {
		return 0.0, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &accBalance)
	if err != nil {
		return 0.0, err
	} else {
		return accBalance.Value, nil
	}
}
