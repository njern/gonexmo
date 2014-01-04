// +build messages

// Tests in this file will only be run if the build tag messages is set:
// go test -tag messages
package nexmo

import (
	"strconv"
	"testing"
	"time"
)

func TestUssdPushMessage(t *testing.T) {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client object with error:", err)
	}

	messageResponse, err := nexmo.SendUssdPush(TEST_FROM, TEST_PHONE_NUMBER,
		"Gonexmo test USSD push message, sent at "+time.Now().String(),
		"gonexmo-test "+strconv.FormatInt(time.Now().Unix(), 10), false)
	if err != nil {
		t.Error("Failed to send USSD push message with error:", err)
	}

	t.Logf("Sent USSD push, response was: %#v\n", messageResponse)
}

func TestUssdPromptMessage(t *testing.T) {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client object with error:", err)
	}

	messageResponse, err := nexmo.SendUssdPrompt(TEST_FROM, TEST_PHONE_NUMBER,
		"Gonexmo test USSD prompt message, sent at "+time.Now().String(),
		"gonexmo-test "+strconv.FormatInt(time.Now().Unix(), 10), false)
	if err != nil {
		t.Error("Failed to send USSD prompt message with error:", err)
	}

	t.Logf("Sent USSD prompt, response was: %#v\n", messageResponse)
}

func TestSendTextMessage(t *testing.T) {
	// TODO(inhies): Create an internal rate limiting system and do away with
	// this hacky 1 second delay.
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client object with error:", err)
	}

	messageResponse, err := nexmo.SendTextMessage(TEST_FROM, TEST_PHONE_NUMBER,
		"Gonexmo test text message, sent at "+time.Now().String(),
		"gonexmo-test "+strconv.FormatInt(time.Now().Unix(), 10), false)
	if err != nil {
		t.Error("Failed to send text message with error:", err)
	}

	t.Logf("Sent SMS, response was: %#v\n", messageResponse)
}

func TestFlashMessage(t *testing.T) {
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client object with error:", err)
	}

	messageResponse, err := nexmo.SendFlashMessage(TEST_FROM, TEST_PHONE_NUMBER,
		"Gonexmo test flash  message, sent at "+time.Now().String(),
		"gonexmo-test "+strconv.FormatInt(time.Now().Unix(), 10), false)
	if err != nil {
		t.Error("Failed to send flash message (class 0 SMS) with error:", err)
	}

	t.Logf("Sent Flash SMS, response was: %#v\n", messageResponse)
}
