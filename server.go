package nexmo

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// MessageType can be one of the following:
//  - TextMessage
//	- UnicodeMessage
//	- BinaryMessage
type MessageType int

// Message types
const (
	TextMessage = iota + 1
	UnicodeMessage
	BinaryMessage
)

var messageTypeMap = map[string]MessageType{
	"text":    TextMessage,
	"unicode": UnicodeMessage,
	"binary":  BinaryMessage,
}

var messageTypeIntMap = map[MessageType]string{
	TextMessage:    "text",
	UnicodeMessage: "unicode",
	BinaryMessage:  "binary",
}

func (m MessageType) String() string {
	if m < 1 || m > 3 {
		return "undefined"
	}

	return messageTypeIntMap[m]
}

// ReceivedMessage represents a message that was received from the Nexmo API.
type ReceivedMessage struct {
	// Expected values are "text" or "binary".
	Type MessageType

	// Recipient number (your long virtual number).
	To string

	// Sender ID.
	MSISDN string

	// Optional unique identifier of a mobile network MCCMNC.
	NetworkCode string

	// Nexmo message ID.
	ID string

	// Time when Nexmo started to push the message to you.
	Timestamp time.Time

	// Parameters for conactenated messages:
	Concatenated bool // Set to true if a MO concatenated message is detected.
	Concat       struct {

		// Transaction reference. All message parts will share the same
		//transaction reference.
		Reference string

		// Total number of parts in this concatenated message set.
		Total int

		// The part number of this message withing the set (starts at 1).
		Part int
	}

	// When Type == text:
	Text string // Content of the message

	Keyword string // First word in the message body, typically used with short codes
	// When type == binary:

	// Content of the message.
	Data []byte

	// User Data Header.
	UDH []byte
}

// DeliveryReceipt is a delivery receipt for a single SMS sent via the Nexmo API
type DeliveryReceipt struct {
	To              string    `json:"to"`
	NetworkCode     string    `json:"network-code"`
	MessageID       string    `json:"messageId"`
	MSISDN          string    `json:"msisdn"`
	Status          string    `json:"status"`
	ErrorCode       string    `json:"err-code"`
	Price           string    `json:"price"`
	SCTS            time.Time `json:"scts"`
	Timestamp       time.Time `json:"message-timestamp"`
	ClientReference string    `json:"client-ref"`
}

// ParseReceivedMessage unmarshals and processes the form data in a Nexmo request
// and returns a DeliveryReceipt struct.
func ParseDeliveryReceipt(req *http.Request) (*DeliveryReceipt, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse form data: %v", err)
	}

	// Decode the form data
	m := new(DeliveryReceipt)

	m.To = req.FormValue("to")
	m.NetworkCode = req.FormValue("network-code")
	m.MessageID = req.FormValue("messageId")
	m.MSISDN = req.FormValue("msisdn")
	m.Status = req.FormValue("status")
	m.ErrorCode = req.FormValue("err-code")
	m.Price = req.FormValue("price")
	m.ClientReference = req.FormValue("client-ref")

	{
		t, err := url.QueryUnescape(req.FormValue("scts"))
		if err != nil {
			return nil, fmt.Errorf("failed to unescape field 'scts': %v", err)
		}

		m.SCTS, err = parseSCTS(t)
		if err != nil {
			return nil, err
		}
	}

	{
		t, err := url.QueryUnescape(req.FormValue("message-timestamp"))
		if err != nil {
			return nil, fmt.Errorf("failed to unescape field 'message-timestamp': %v", err)
		}

		m.Timestamp, err = parseMessageTimestamp(t)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func parseSCTS(t string) (time.Time, error) {
	if t == "" {
		return time.Time{}, nil
	}

	timestamp, err := time.Parse("0601021504", t)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse timestamp for field 'scts': %v", err)
	}

	return timestamp, nil
}

func parseMessageTimestamp(t string) (time.Time, error) {
	if t == "" {
		return time.Time{}, nil
	}

	// nexmo is just doing some crazy stuff lately
	formats := []string{
		"2006-01-02 15:04:05 -0700", // actually valid
		"2006-01-02 15:04:05 0000",  // where did the plus go
		"2006-01-02 15:04:05  0000", // oh you forgot to URL encode it? very cool
		"2006-01-02 15:04:05",       // you know what, forget timezones
	}

	for _, f := range formats {
		timestamp, err := time.Parse(f, t)
		if err == nil {
			return timestamp, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse timestamp '%s' for field 'message-timestamp'", t)
}

// NewDeliveryHandler creates a new http.HandlerFunc that can be used to listen
// for delivery receipts from the Nexmo server. Any receipts received will be
// decoded and passed to the out chan.
func NewDeliveryHandler(out chan *DeliveryReceipt, verifyIPs bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if verifyIPs {
			// Check if the request came from Nexmo
			host, _, err := net.SplitHostPort(req.RemoteAddr)
			if !IsTrustedIP(host) || err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		// Check if the query is empty. If it is, it's just Nexmo
		// making sure our service is up, so we don't want to return
		// an error.
		if req.URL.RawQuery == "" {
			return
		}

		receipt, err := ParseDeliveryReceipt(req)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Pass it out on the chan
		out <- receipt
	}
}

// ParseReceivedMessage unmarshals and processes the form data in a Nexmo request
// and returns a ReceivedMessage struct.
func ParseReceivedMessage(req *http.Request) (*ReceivedMessage, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("failed to parse form data: %v", err)
	}

	// Decode the form data
	m := new(ReceivedMessage)
	switch t := req.FormValue("type"); t {
	case "text":
		m.Text, err = url.QueryUnescape(req.FormValue("text"))
		if err != nil {
			return nil, fmt.Errorf("failed to unescape field 'text': %v", err)
		}
		m.Type = TextMessage

	case "unicode":
		m.Text, err = url.QueryUnescape(req.FormValue("text"))
		if err != nil {
			return nil, fmt.Errorf("failed to unescape field 'text': %v", err)
		}
		m.Type = UnicodeMessage

	case "binary":
		// TODO: I have no idea if this data stuff works, as I'm unable to
		// send data SMS messages.
		data, err := url.QueryUnescape(req.FormValue("data"))
		if err != nil {
			return nil, fmt.Errorf("failed to unescape field 'data': %v", err)
		}
		m.Data = []byte(data)

		udh, err := url.QueryUnescape(req.FormValue("udh"))
		if err != nil {
			return nil, fmt.Errorf("failed to unescape field 'udh': %v", err)
		}
		m.UDH = []byte(udh)
		m.Type = BinaryMessage

	default:
		//error
		return nil, fmt.Errorf("unrecognized message type %s", t)
	}

	m.To = req.FormValue("to")
	m.MSISDN = req.FormValue("msisdn")
	m.NetworkCode = req.FormValue("network-code")
	m.ID = req.FormValue("messageId")

	m.Keyword = req.FormValue("keyword")
	t, err := url.QueryUnescape(req.FormValue("message-timestamp"))
	if err != nil {
		return nil, fmt.Errorf("failed to unescape field 'message-timestamp': %v", err)
	}

	// Convert the timestamp to a time.Time.
	timestamp, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp for field 'message-timestamp': %v", err)
	}

	m.Timestamp = timestamp

	// TODO: I don't know if this works as I've been unable to send an SMS
	// message longer than 160 characters that doesn't get concatenated
	// automatically.
	if req.FormValue("concat") == "true" {
		m.Concatenated = true
		m.Concat.Reference = req.FormValue("concat-ref")
		m.Concat.Total, err = strconv.Atoi(req.FormValue("concat-total"))
		if err != nil {
			return nil, fmt.Errorf("failed to convert field 'concat-total' to int: %v", err)
		}
		m.Concat.Part, err = strconv.Atoi(req.FormValue("concat-part"))
		if err != nil {
			return nil, fmt.Errorf("failed to convert field 'concat-part' to int: %v", err)
		}
	}

	return m, nil
}

// NewMessageHandler creates a new http.HandlerFunc that can be used to listen
// for new messages from the Nexmo server. Any new messages received will be
// decoded and passed to the out chan.
func NewMessageHandler(out chan *ReceivedMessage, verifyIPs bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if verifyIPs {
			// Check if the request came from Nexmo
			host, _, err := net.SplitHostPort(req.RemoteAddr)
			if !IsTrustedIP(host) || err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		// Check if the query is empty. If it is, it's just Nexmo
		// making sure our service is up, so we don't want to return
		// an error.
		if req.URL.RawQuery == "" {
			return
		}

		message, err := ParseReceivedMessage(req)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Pass it out on the chan
		out <- message
	}
}
