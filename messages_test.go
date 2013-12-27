// +build messages

// Tests in this file will only be run if the build tag messages is set:
// go test -tag messages
package nexmo

import (
	"testing"
)

func TestSendTextMessage(t *testing.T) {
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := New(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}

	messageResponse, err := nexmo.SendTextMessage("go-nexmo", TEST_PHONE_NUMBER,
		"Looks like go-nexmo works great,"+
			" we should definitely buy that njern guy a beer!", "001", false)
	if err != nil {
		t.Error("Failed to send text message with error:", err)
	}

	t.Logf("Sent SMS, response was: %#v\n", messageResponse)
}

func TestFlashMessage(t *testing.T) {
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := New(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}

	messageResponse, err := nexmo.SendFlashMessage("gonexmo", TEST_PHONE_NUMBER,
		"Looks like go-nexmo works great,"+
			" we should definitely buy that njern guy a beer!", "001", false)
	if err != nil {
		t.Error("Failed to send flash message (class 0 SMS) with error:", err)
	}

	t.Logf("Sent Flash SMS, response was: %#v\n", messageResponse)
}
