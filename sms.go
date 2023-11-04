package nexmo

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
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

// MessageClass will be one of the following:
//   - Flash
//   - Standard
//   - SIMData
//   - Forward
type MessageClass int

// SMS message classes.
const (
	// This type of SMS message is displayed on the mobile screen without being
	// saved in the message store or on the SIM card; unless explicitly saved
	// by the mobile user.
	Flash MessageClass = iota

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

var messageClassMap = map[MessageClass]string{
	Flash:    "flash",
	Standard: "standard",
	SIMData:  "SIM data",
	Forward:  "forward",
}

func (m MessageClass) String() string {
	return messageClassMap[m]
}

// MarshalJSON implements the json.Marshaller interface
func (m *SMSMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		APIKey    string `json:"api_key"`
		APISecret string `json:"api_secret"`
		SMSMessage
	}{
		APIKey:     m.apiKey,
		APISecret:  m.apiSecret,
		SMSMessage: *m,
	})
}

// SMSMessage defines a single SMS message.
type SMSMessage struct {
	apiKey               string
	apiSecret            string
	From                 string       `json:"from"`
	To                   string       `json:"to"`
	Type                 string       `json:"type"`
	Text                 string       `json:"text,omitempty"`              // Optional.
	StatusReportRequired int          `json:"status-report-req,omitempty"` // Optional.
	ClientReference      string       `json:"client-ref,omitempty"`        // Optional.
	NetworkCode          string       `json:"network-code,omitempty"`      // Optional.
	VCard                string       `json:"vcrad,omitempty"`             // Optional.
	VCal                 string       `json:"vcal,omitempty"`              // Optional.
	TTL                  int          `json:"ttl,omitempty"`               // Optional.
	Class                MessageClass `json:"message-class,omitempty"`     // Optional.
	Callback             string       `json:"callback,omitempty"`          // Optional.
	Body                 []byte       `json:"body,omitempty"`              // Required for Binary message.
	UDH                  []byte       `json:"udh,omitempty"`               // Required for Binary message.

	// The following is only for type=wappush

	Title    string `json:"title,omitempty"`    // Title shown to recipient
	URL      string `json:"url,omitempty"`      // WAP Push URL
	Validity int    `json:"validity,omitempty"` // Duration WAP Push is available in milliseconds
}

// A ResponseCode will be returned
// whenever an SMSMessage is sent.
type ResponseCode int

// String implements the fmt.Stringer interface
func (c ResponseCode) String() string {
	return responseCodeMap[c]
}

// Possible response codes
const (
	ResponseSuccess ResponseCode = iota
	ResponseThrottled
	ResponseMissingParams
	ResponseInvalidParams
	ResponseInvalidCredentials
	ResponseInternalError
	ResponseInvalidMessage
	ResponseNumberBarred
	ResponsePartnerAcctBarred
	ResponsePartnerQuotaExceeded
	ResponseUnused //This is not used  yet.Left blank by  Nexmo for the time being.
	ResponseRESTNotEnabled
	ResponseMessageTooLong
	ResponseCommunicationFailed
	ResponseInvalidSignature
	ResponseInvalidSenderAddress
	ResponseInvalidTTL
	ResponseFacilityNotAllowed
	ResponseInvalidMessageClass
)

var responseCodeMap = map[ResponseCode]string{
	ResponseSuccess:              "Success",
	ResponseThrottled:            "Throttled",
	ResponseMissingParams:        "Missing params",
	ResponseInvalidParams:        "Invalid params",
	ResponseInvalidCredentials:   "Invalid credentials",
	ResponseInternalError:        "Internal error",
	ResponseInvalidMessage:       "Invalid message",
	ResponseNumberBarred:         "Number barred",
	ResponsePartnerAcctBarred:    "Partner account barred",
	ResponsePartnerQuotaExceeded: "Partner quota exceeded",
	ResponseRESTNotEnabled:       "Account not enabled for REST",
	ResponseMessageTooLong:       "Message too long",
	ResponseCommunicationFailed:  "Communication failed",
	ResponseInvalidSignature:     "Invalid signature",
	ResponseInvalidSenderAddress: "Invalid sender address",
	ResponseInvalidTTL:           "Invalid TTL",
	ResponseFacilityNotAllowed:   "Facility not allowed",
	ResponseInvalidMessageClass:  "Invalid message class",
}

// MessageReport is the "status report" for a single SMS sent via the Nexmo API
type MessageReport struct {
	Status           ResponseCode `json:"status,string"`
	MessageID        string       `json:"message-id"`
	To               string       `json:"to"`
	ClientReference  string       `json:"client-ref"`
	RemainingBalance string       `json:"remaining-balance"`
	MessagePrice     string       `json:"message-price"`
	Network          string       `json:"network"`
	ErrorText        string       `json:"error-text"`
}

// MessageResponse contains the response from Nexmo's API after we attempt to
// send any kind of message.
// It will contain one MessageReport for every 160 chars sent.
type MessageResponse struct {
	MessageCount int             `json:"message-count,string"`
	Messages     []MessageReport `json:"messages"`
}

type InvalidResponseError struct {
	Message string
	Err     error
	Body    []byte
}

type SendConnectionError struct {
	Message string
	Err     error
	Body    []byte
	Debug   []string
}

func (e SendConnectionError) Error() string {
	return e.Message
}

func (e InvalidResponseError) Error() string {
	return e.Message
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
	case Text:
	case Unicode:
		if len(msg.Text) <= 0 {
			return nil, errors.New("Invalid message text")
		}
	case Binary:
		if len(msg.UDH) == 0 || len(msg.Body) == 0 {
			return nil, errors.New("Invalid binary message")
		}

	case WAPPush:
		if len(msg.URL) == 0 || len(msg.Title) == 0 {
			return nil, errors.New("Invalid WAP Push parameters")
		}
	}
	if !c.client.useOauth {
		msg.apiKey = c.client.apiKey
		msg.apiSecret = c.client.apiSecret
	}

	var r *http.Request

	buf, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.New("invalid message struct - unable to convert to JSON")
	}
	b := bytes.NewBuffer(buf)
	r, _ = http.NewRequest("POST", apiRoot+"/sms/json", b)

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	debug, trace := getRequestTrace()
	r = r.WithContext(httptrace.WithClientTrace(r.Context(), trace))

	resp, err := c.client.HTTPClient.Do(r)

	if err != nil {
		sendErr := SendConnectionError{
			Message: "nexmo http send failed",
			Err:     err,
			Debug:   *debug,
		}

		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
			body, bodyErr := ioutil.ReadAll(resp.Body)
			if bodyErr != nil {
				sendErr.Body = body
			}
		}

		return nil, sendErr
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, InvalidResponseError{
			Message: "failed to read response from Nexmo",
			Body:    body,
			Err:     err,
		}
	}

	err = json.Unmarshal(body, &messageResponse)
	if err != nil {
		return nil, InvalidResponseError{
			Message: "failed to unmarshal response from Nexmo",
			Body:    body,
			Err:     err,
		}
	}

	return messageResponse, nil
}

func getRequestTrace() (*[]string, *httptrace.ClientTrace) {

	debugTrace := &[]string{}

	return debugTrace, &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			*debugTrace = append(*debugTrace, fmt.Sprintf("Initiating connecting to %s", hostPort))
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			if connInfo.Reused {
				*debugTrace = append(*debugTrace, "Re-using existing connection")
			} else {
				*debugTrace = append(*debugTrace, "New connection successfully established")
			}
		},
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			*debugTrace = append(*debugTrace, fmt.Sprintf("Resolving DNS for %s", dnsInfo.Host))
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			if dnsInfo.Err != nil {
				*debugTrace = append(*debugTrace, fmt.Sprintf("Error resolving DNS (%s)", dnsInfo.Err.Error()))
			} else {
				*debugTrace = append(*debugTrace, "DNS resolved successfully")
			}
		},
		ConnectStart: func(network string, addr string) {
			*debugTrace = append(*debugTrace, fmt.Sprintf("Initiating connecting to %s %s", network, addr))
		},
		ConnectDone: func(network string, addr string, err error) {
			if err != nil {
				*debugTrace = append(*debugTrace, fmt.Sprintf("Error connecting to %s %s (%s)", network, addr, err.Error()))
			} else {
				*debugTrace = append(*debugTrace, fmt.Sprintf("Connection complete to %s %s", network, addr))
			}
		},
		GotFirstResponseByte: func() {
			*debugTrace = append(*debugTrace, "Read first byte of response headers")
		},
		TLSHandshakeStart: func() {
			*debugTrace = append(*debugTrace, "TLS handshake started")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			if err != nil {
				*debugTrace = append(*debugTrace, fmt.Sprintf("TLS handshake error (%s)", err.Error()))
			} else {
				*debugTrace = append(*debugTrace, "TLS handshake complete")
			}
		},
		WroteHeaders: func() {
			*debugTrace = append(*debugTrace, "Request headers successfully written")
		},
		WroteRequest: func(requestInfo httptrace.WroteRequestInfo) {
			if requestInfo.Err != nil {
				*debugTrace = append(*debugTrace, fmt.Sprintf("Error while writing http request (%s)", requestInfo.Err.Error()))
			} else {
				*debugTrace = append(*debugTrace, "Full request successfully written")
			}
		},
	}
}
