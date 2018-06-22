// +build messages

// Tests in this file will only be run if the build tag messages is set:
// `go test -tag messages`
// Test with only sending one message using:
// `go test -test.run UssdPush -tags messages`
package nexmo

import (
	"strconv"
	"testing"
	"time"
)

func TestUssdPushMessage(t *testing.T) {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if testPhoneNumber == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClient(testAPIKey, testAPISecret)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}
	message := &USSDMessage{
		From:            testFrom,
		To:              testPhoneNumber,
		Text:            "Gonexmo test USSD push message, sent at " + time.Now().String(),
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
	}

	messageResponse, err := nexmo.USSD.Send(message)
	if err != nil {
		t.Error("Failed to send USSD push message with error:", err)
	}

	t.Logf("Sent USSD push, response was: %#v\n", messageResponse)
}

func TestUssdPromptMessage(t *testing.T) {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if testPhoneNumber == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClient(testAPIKey, testAPISecret)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &USSDMessage{
		From:            testFrom,
		To:              testPhoneNumber,
		Text:            "Gonexmo test USSD prompt message, sent at " + time.Now().String(),
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
		Prompt:          true,
	}

	messageResponse, err := nexmo.USSD.Send(message)
	if err != nil {
		t.Error("Failed to send USSD prompt message with error:", err)
	}

	t.Logf("Sent USSD prompt, response was: %#v\n", messageResponse)
}
