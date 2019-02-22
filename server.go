package nexmo

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

// ToString provides string represenstation of struct
func (dlr *DeliveryReceipt) ToString() string {
	s, err := json.Marshal(dlr)
	if err != nil {
		log.Printf("Failed to Marshal dlr with error %v, returning empty string", %v)
		return ""
	}
	return string(s)
}

type rawDeliveryReceipt struct {
	To              string `json:"to"`
	NetworkCode     string `json:"network-code"`
	MessageID       string `json:"messageId"`
	MSISDN          string `json:"msisdn"`
	Status          string `json:"status"`
	ErrorCode       string `json:"err-code"`
	Price           string `json:"price"`
	SCTS            string `json:"scts"`
	Timestamp       string `json:"message-timestamp"`
	ClientReference string `json:"client-ref"`
}

// NewDeliveryHandler creates a new http.HandlerFunc that can be used to listen
// for delivery receipts from the Nexmo server. Any receipts received will be
// decoded nad passed to the out chan.
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
		log.Println("merf")

		var err error
		// Check if the query and body are empty. If it is, it's just Nexmo
		// making sure our service is up, so we don't want to return
		// an error.
		if req.URL.RawQuery == "" && req.ContentLength <= 0 {
			w.WriteHeader(http.StatusOK)
			return
		}

		contentType, ok := req.Header["Content-Type"]
		if !ok {
			log.Println("foo")
			http.Error(w, "Content-Type not set", http.StatusBadRequest)
			return
		}

		//  nexmo claims the response is going to be of type application/json
		//  add support for application/json and maintain www-form-urlencoded support
		var rm *rawDeliveryReceipt
		if contentType[0] == "application/json" {
			rm, err = parseJSON(w, req)
			if err != nil {
				return
			}
		} else if contentType[0] == "application/x-www-form-urlencoded" {
			rm, err = parseForm(w, req)
			if err != nil {
				return
			}
		} else {
			http.Error(w, "Content-Type "+contentType[0]+" not supported.", http.StatusBadRequest)
			return
		}

		m, err := convertTimestamps(rm, w)
		if err != nil {
			return
		}

		// Pass it out on the chan
		out <- m
	}

}

//  convert data supplied in form to all string representation of values
func parseForm(w http.ResponseWriter, req *http.Request) (*rawDeliveryReceipt, error) {
	err := req.ParseForm()
	if err != nil {
		log.Printf("form parse error %v", err)
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return nil, err
	}

	// Decode the form data
	m := &rawDeliveryReceipt{
		To:              req.FormValue("to"),
		NetworkCode:     req.FormValue("network-code"),
		MessageID:       req.FormValue("messageId"),
		MSISDN:          req.FormValue("msisdn"),
		Status:          req.FormValue("status"),
		ErrorCode:       req.FormValue("err-code"),
		Price:           req.FormValue("price"),
		ClientReference: req.FormValue("client-ref"),
		SCTS:            req.FormValue("client-ref"),
		Timestamp:       req.FormValue("message-timestamp"),
	}

	return m, nil
}

// extract json blob into structure
func parseJSON(w http.ResponseWriter, req *http.Request) (*rawDeliveryReceipt, error) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Failed to read request body %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return nil, err
	}

	rawDLR := &rawDeliveryReceipt{}
	err = json.Unmarshal(body, rawDLR)
	if err != nil {
		log.Printf("Failed to parse request body %s with error %v", string(body), err)
		http.Error(w, "", http.StatusInternalServerError)
		return nil, err
	}

	return rawDLR, nil

}

// parse timestamps, throw failures if timestamps not present or not formatted correctly
func convertTimestamps(rawDLR *rawDeliveryReceipt, w http.ResponseWriter) (*DeliveryReceipt, error) {
	dlr := &DeliveryReceipt{
		To:              rawDLR.To,
		NetworkCode:     rawDLR.NetworkCode,
		MessageID:       rawDLR.MessageID,
		MSISDN:          rawDLR.MSISDN,
		Status:          rawDLR.Status,
		ErrorCode:       rawDLR.ErrorCode,
		Price:           rawDLR.Price,
		ClientReference: rawDLR.ClientReference,
	}

	t, err := url.QueryUnescape(rawDLR.SCTS)
	if err != nil {
		log.Printf("Unable to unescape SCTS from %v with error %v\n", rawDLR, err)
		http.Error(w, "unable to return formvalue scts", http.StatusInternalServerError)
		return nil, err
	}

	// Convert the timestamp to a time.Time.
	timestamp, err := time.Parse("0601021504", t)
	if err != nil {
		log.Printf("Unable to time.Parse SCTS from %v with error %v\n", rawDLR, err)
		http.Error(w, "unable to parse time value from scts", http.StatusInternalServerError)
		return nil, err
	}

	dlr.SCTS = timestamp

	t, err = url.QueryUnescape(rawDLR.Timestamp)
	if err != nil {
		log.Printf("Unable to unescape message-timestamp from %v with error %v\n", rawDLR, err)
		http.Error(w, "unable to return formvalue message-timestamp", http.StatusInternalServerError)
		return nil, err
	}

	// Convert the timestamp to a time.Time.
	timestamp, err = time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		log.Printf("Unable to time.Parse message-timestamp from %v with error %v\n", rawDLR, err)
		http.Error(w, "unable to parse time value from message-timestamp", http.StatusInternalServerError)
		return nil, err
	}

	dlr.Timestamp = timestamp

	return dlr, nil
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

		var err error

		// Check if the query is empty. If it is, it's just Nexmo
		// making sure our service is up, so we don't want to return
		// an error.
		if req.URL.RawQuery == "" {
			return
		}

		req.ParseForm()
		// Decode the form data
		m := new(ReceivedMessage)
		switch req.FormValue("type") {
		case "text":
			m.Text, err = url.QueryUnescape(req.FormValue("text"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Type = TextMessage
		case "unicode":
			m.Text, err = url.QueryUnescape(req.FormValue("text"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Type = UnicodeMessage

			// TODO: I have no idea if this data stuff works, as I'm unable to
			// send data SMS messages.
		case "binary":
			data, err := url.QueryUnescape(req.FormValue("data"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Data = []byte(data)

			udh, err := url.QueryUnescape(req.FormValue("udh"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.UDH = []byte(udh)
			m.Type = BinaryMessage

		default:
			//error
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.To = req.FormValue("to")
		m.MSISDN = req.FormValue("msisdn")
		m.NetworkCode = req.FormValue("network-code")
		m.ID = req.FormValue("messageId")

		m.Keyword = req.FormValue("keyword")
		t, err := url.QueryUnescape(req.FormValue("message-timestamp"))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Convert the timestamp to a time.Time.
		timestamp, err := time.Parse("2006-01-02 15:04:05", t)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
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
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Concat.Part, err = strconv.Atoi(req.FormValue("concat-part"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		// Pass it out on the chan
		out <- m
	}

}
