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

const (
	apiRoot = "https://rest.nexmo.com"
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

// MessageResponse contains the response from Nexmo's API after we attempt to
// send any kind of message.
// It will contain one MessageReport for every 160 chars sent.
type MessageResponse struct {
	MessageCount string          `json:"message-count"`
	Messages     []MessageReport `json:"messages"`
}

func (nexmo *Client) sendMessage(from, to, text, clientReference string,
	statusReportRequired bool, isFlashMessage bool) (*MessageResponse, error) {
	if len(clientReference) > 40 {
		return nil, errors.New("Client reference too long")
	}
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
	if isFlashMessage {
		values.Set("message_class", "0")
	}

	client := &http.Client{}
	var r *http.Request
	r, _ = http.NewRequest("POST", apiRoot+"/sms/json", bytes.NewReader([]byte(values.Encode())))
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
func (nexmo *Client) SendTextMessage(from, to, text, clientReference string,
	statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference,
		statusReportRequired, false)
}

// SendFlashMessage() sends a class 0 SMS (Flash message).
func (nexmo *Client) SendFlashMessage(from, to, text, clientReference string,
	statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference,
		statusReportRequired, true)
}
