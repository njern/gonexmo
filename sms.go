package nexmo

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// SMS represents the SMS API functions for sending text messages.
type SMS struct {
	client *Client
}

// SMS message types.
const (
	Text = iota + 1
	Binary
	WAPPush
	Unicode
	VCal
	VCard
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
	From                 string
	To                   string
	Type                 int    // Optional: default to Text.
	Text                 string // Optional.
	StatusReportRequired bool   // Optional.
	ClientReference      string // Optional.
	NetworkCode          string // Optional.
	VCard                string // Optional.
	VCal                 string // Optional.
	TTL                  int    // Optional.
	Class                int    // Optional.
	Body                 []byte // Required for Binary message.
	UDH                  []byte // Required for Binary message.
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

	values := make(url.Values)

	switch msg.Type {
	// 0 would be default
	case 0, Text:
		if len(msg.Text) <= 0 {
			return nil, errors.New("Invalid message text")
		} else {
			// TODO(inhies): UTF8 and URL encode before setting
			values.Set("type", "text")
			values.Set("text", msg.Text)
		}
	case Binary:
		if len(msg.UDH) == 0 || len(msg.Body) == 0 {
			return nil, errors.New("Invalid binary message")
		}
		values.Set("type", "binary")
	case WAPPush:
		values.Set("type", "wappush")
	case Unicode:
		values.Set("type", "unicode")
	case VCal:
		values.Set("type", "vcal")
	case VCard:
		values.Set("type", "vcard")
	}

	if !c.client.useOauth {
		values.Set("api_key", c.client.apiKey)
		values.Set("api_secret", c.client.apiSecret)
	}

	if msg.StatusReportRequired {
		values.Set("status_report_req", "1")
	}

	if msg.ClientReference != "" {
		values.Set("client_ref", msg.ClientReference)
	}

	if msg.NetworkCode != "" {
		values.Set("network-code", msg.NetworkCode)
	}

	if msg.VCard != "" {
		values.Set("vcard", msg.VCard)
	}

	if msg.VCal != "" {
		values.Set("vcal", msg.VCal)
	}

	if msg.TTL != 0 {
		values.Set("ttl", strconv.Itoa(msg.TTL))
	}

	if msg.Class != 0 {
		values.Set("message-class", strconv.Itoa(msg.Class))
	}

	values.Set("to", msg.To)
	values.Set("from", msg.From)

	client := &http.Client{}
	valuesReader := bytes.NewReader([]byte(values.Encode()))
	var r *http.Request
	r, _ = http.NewRequest("POST", apiRoot+"/sms/json", valuesReader)

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
