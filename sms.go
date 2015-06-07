package nexmo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// SMS represents the SMS API functions for sending text messages.
type SMS struct {
	client *Client
}

// SMS message types.
const (
	Text    = "text"
	Binary  = "binary"
	WAPPush = "wappush"
	Unicode = "unicode"
	VCal    = "vcal"
	VCard   = "vcard"
)

// SMS message classes.
const (
	// This type of SMS message is displayed on the mobile screen without being
	// saved in the message store or on the SIM card; unless explicitly saved
	// by the mobile user.
	Flash = iota

	// This message is to be stored in the device memory or the SIM card
	// (depending on memory availability).
	Standard

	// This message class carries SIM card data. The SIM card data must be
	// successfully transferred prior to sending acknowledgment to the service
	// center. An error message is sent to the service center if this
	// transmission is not possible.
	SIMData

	// This message is forwarded from the receiving entity to an external
	// device. The delivery acknowledgment is sent to the service center
	// regardless of whether or not the message was forwarded to the external
	// device.
	Forward
)

// Type SMSMessage defines a single SMS message.
type SMSMessage struct {
	ApiKey               string `json:"api_key"`
	ApiSecret            string `json:"api_secret"`
	From                 string `json:"from"`
	To                   string `json:"to"`
	Type                 string `json:"type"`
	Text                 string `json:"text,omitempty"`              // Optional.
	StatusReportRequired int    `json:"status-report-req,omitempty"` // Optional.
	ClientReference      string `json:"client-ref,omitempty"`        // Optional.
	NetworkCode          string `json:"network-code,omitempty"`      // Optional.
	VCard                string `json:"vcrad,omitempty"`             // Optional.
	VCal                 string `json:"vcal,omitempty"`              // Optional.
	TTL                  int    `json:"ttl,omitempty"`               // Optional.
	Class                int    `json:"message-class,omitempty"`     // Optional.
	Body                 []byte `json:"body,omitempty"`              // Required for Binary message.
	UDH                  []byte `json:"udh,omitempty"`               // Required for Binary message.
}

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

// Send the message using the specified SMS client.
func (c *SMS) Send(msg *SMSMessage) (*MessageResponse, error) {
	if len(msg.From) <= 0 {
		return nil, errors.New("Invalid From field specified")
	}

	if len(msg.To) <= 0 {
		return nil, errors.New("Invalid To field specified")
	}

	if len(msg.ClientReference) > 40 {
		return nil, errors.New("Client reference too long")
	}

	var messageResponse *MessageResponse

	switch msg.Type {
	// 0 would be default
	case Text:
		if len(msg.Text) <= 0 {
			return nil, errors.New("Invalid message text")
		}
	case Binary:
		if len(msg.UDH) == 0 || len(msg.Body) == 0 {
			return nil, errors.New("Invalid binary message")
		}
	}

	if !c.client.useOauth {
		msg.ApiKey = c.client.apiKey
		msg.ApiSecret = c.client.apiSecret
	}

	client := &http.Client{}

	var r *http.Request
	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.New("Invalid message struct. Cannot convert to json.")
	}
	b := bytes.NewBuffer(buf)
	r, _ = http.NewRequest("POST", apiRoot+"/sms/json", b)

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &messageResponse)
	if err != nil {
		return nil, err
	}
	return messageResponse, nil
}
