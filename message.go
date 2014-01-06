package nexmo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Message identifiers to be used with sendMessage()
const (
	sms = iota
	flash
	ussdPush
	ussdPrompt
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
	statusReportRequired bool, class int) (*MessageResponse, error) {
	if len(clientReference) > 40 {
		return nil, errors.New("client reference too long")
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

	client := &http.Client{}
	valuesReader := bytes.NewReader([]byte(values.Encode()))
	var r *http.Request
	switch class {
	case sms:
		r, _ = http.NewRequest("POST", apiRoot+"/sms/json", valuesReader)
	case flash:
		values.Set("message_class", "0")
		r, _ = http.NewRequest("POST", apiRoot+"/sms/json", valuesReader)
	case ussdPush:
		r, _ = http.NewRequest("POST", apiRoot+"/ussd/json", valuesReader)
	case ussdPrompt:
		r, _ = http.NewRequest("POST", apiRoot+"/ussd-prompt/json", valuesReader)
	}

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
	}

	return messageResponse, nil
}

// SendUssdPush sends a USSD push message.
func (nexmo *Client) SendUssdPush(from, to, text, clientReference string,
	statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference,
		statusReportRequired, ussdPush)
}

// SendUssdPrompt sends a USSD prompt message. You must have a callback URL
//  set up and the 'from' field must be a long virtual number.
func (nexmo *Client) SendUssdPrompt(from, to, text, clientReference string,
	statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference,
		statusReportRequired, ussdPrompt)
}

// SendTextMessage sends a normal SMS.
func (nexmo *Client) SendTextMessage(from, to, text, clientReference string,
	statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference,
		statusReportRequired, sms)
}

// SendFlashMessage sends a class 0 SMS (Flash message).
func (nexmo *Client) SendFlashMessage(from, to, text, clientReference string,
	statusReportRequired bool) (*MessageResponse, error) {
	return nexmo.sendMessage(from, to, text, clientReference,
		statusReportRequired, flash)
}
