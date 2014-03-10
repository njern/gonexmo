package nexmo

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// RecvdMessage represents a message that was received from the Nexmo API.
type RecvdMessage struct {
	// Expected values are "text" or "binary".
	Type int

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

	// When type == binary:

	// Content of the message.
	Data []byte

	// User Data Header.
	UDH []byte
}

// NewMessageHandler creates a new http.HandlerFunc that can be used to listen
// for new messages from the Nexmo server. Any new messages received will be
// decoded and passed to the out chan.
func NewMessageHandler(out chan *RecvdMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		// Decode the form data
		m := new(RecvdMessage)
		var err error
		switch req.FormValue("type") {
		case "text":
			m.Text, err = url.QueryUnescape(req.FormValue("text"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			// TODO: I have no idea if this data stuff works, as I'm unable to
			// send data SMS messages.
		case "binary":
			data, err := url.QueryUnescape(req.FormValue("text"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.Data = []byte(data)

			udh, err := url.QueryUnescape(req.FormValue("text"))
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			m.UDH = []byte(udh)

		default:
			//error
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.To = req.FormValue("to")
		m.MSISDN = req.FormValue("msisdn")
		m.NetworkCode = req.FormValue("network-code")
		m.ID = req.FormValue("messageId")

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
