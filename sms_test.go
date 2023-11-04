//go:build messages
// +build messages

// Tests in this file will only be run if the build tag messages is set:
// `go test -tag messages`
// Test with only sending one message using:
// `go test -test.run SendText -tags messages`
package nexmo

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TODO(inhies): Only create a Client once in an init() function.

func TestSendTextMessage(t *testing.T) {
	InitEnv()

	// TODO(inhies): Create an internal rate limiting system and do away with
	// this hacky 1 second delay.
	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &SMSMessage{
		From:            TEST_FROM,
		To:              TEST_PHONE_NUMBER,
		Type:            Text,
		Text:            "Gonexmo test SMS message, sent at " + time.Now().String(),
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
		Class:           Standard,
	}

	messageResponse, err := nexmo.SMS.Send(message)
	if err != nil {
		t.Error("Failed to send text message with error:", err)
	}

	t.Logf("Sent SMS, response was: %#v\n", messageResponse)
}

func TestFlashMessage(t *testing.T) {
	InitEnv()

	time.Sleep(1 * time.Second) // Sleep 1 second due to API limitation
	if TEST_PHONE_NUMBER == "" {
		t.Fatal("No test phone number specified. Please set NEXMO_NUM")
	}
	nexmo, err := NewClientFromAPI(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	message := &SMSMessage{
		From:            TEST_FROM,
		To:              TEST_PHONE_NUMBER,
		Type:            Text,
		Text:            "Gonexmo test flash SMS message, sent at " + time.Now().String(),
		ClientReference: "gonexmo-test " + strconv.FormatInt(time.Now().Unix(), 10),
		Class:           Flash,
	}

	messageResponse, err := nexmo.SMS.Send(message)
	if err != nil {
		t.Error("Failed to send flash message (class 0 SMS) with error:", err)
	}

	t.Logf("Sent Flash SMS, response was: %#v\n", messageResponse)
}

func TestCallbackAttributeShouldBeFilled(t *testing.T) {
	smsMessageWithCallbackString := `{"to": "5534999998888", "callback": "https://mycustomcallback.com.br"}`
	smsMessageWithoutCallbackString := `{"to": "5534988887777"}`

	smsMessageWithCallback := &SMSMessage{}
	smsMessageWithoutCallback := &SMSMessage{}

	errWithCallback := json.Unmarshal(smsMessageWithCallbackString.([]byte), smsMessageWithCallback)
	errWithoutCallback := json.Unmarshal(smsMessageWithoutCallbackString.([]byte), smsMessageWithoutCallback)

	if errWithCallback != nil || errWithoutCallback != nil {
		t.Error("Failed to unmarshal Json string.")
	}

	if smsMessageWithCallback.Callback != "https://mycustomcallback.com.br" {
		t.Error("Callback attribute wasn't filled as expected.")
	}

	if smsMessageWithoutCallback.Callback != "" {
		t.Error("Callback attribute was filled when it shouldn't be.")
	}

	t.Log("Callback attribute works as it should be.")
}

func TestCallbackAttributeShouldBeOmited(t *testing.T) {
	to := "5534999998888"
	callback := "https://mycustomcallback.com.br"

	smsMessageWithCallback := &SMSMessage{}
	smsMessageWithCallback.To = to
	smsMessageWithCallback.Callback = callback

	smsMessageWithoutCallback := &SMSMessage{}
	smsMessageWithoutCallback.To = to

	smsMessageWithCallbackByte, errWithCallback := json.Marshal(smsMessageWithCallback)
	smsMessageWithoutCallbackByte, errWithoutCallback := json.Marshal(smsMessageWithoutCallback)

	if errWithCallback != nil || errWithoutCallback != nil {
		t.Error("Failed to marshal SMSMessage.")
	}

	if !strings.Contains(str(smsMessageWithCallbackByte), callback) {
		t.Error("Callback attribute was omited.")
	}

	if strings.Contains(str(smsMessageWithoutCallbackByte), "callback") {
		t.Error("Callback attribute wasn't omited.")
	}

	t.Log("Callback attribute works as it should be.")
}
