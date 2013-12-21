package gonexmo

import (
	"testing"
)

const (
	API_KEY           = "YOUR NEXMO API KEY GOES HERE"
	API_SECRET        = "YOUR NEXMO API SECRET GOES HERE"
	TEST_PHONE_NUMBER = "YOUR PHONE NUMBER GOES HERE"
)

func TestNexmoCreation(t *testing.T) {
	_, err := NexmoWithKeyAndSecret(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}
}

func TestGetAccountBalance(t *testing.T) {
	nexmo, err := NexmoWithKeyAndSecret(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}

	balance, err := nexmo.GetBalance()
	if err != nil {
		t.Error("Failed to get account balance with error:", err)
	}

	t.Log("Got account balance: ", balance, "€")
}

func TestSendTextMessage(t *testing.T) {
	nexmo, err := NexmoWithKeyAndSecret(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}

	messageResponse, err := nexmo.SendTextMessage("go-nexmo", "00358123412345", "Looks like go-nexmo works great, we should definitely buy that njern guy a beer!", "001", false)
	if err != nil {
		t.Error("Failed to send text message with error:", err)
	}

	t.Log("Sent SMS, response was: ", messageResponse)
}

func TestFlashMessage(t *testing.T) {
	nexmo, err := NexmoWithKeyAndSecret(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}

	messageResponse, err := nexmo.SendFlashMessage("go-nexmo", "00358123412345", "Looks like go-nexmo works great, we should definitely buy that njern guy a beer!", "001", false)
	if err != nil {
		t.Error("Failed to send flash message (class 0 SMS) with error:", err)
	}

	t.Log("Sent Flash SMS, response was: ", messageResponse)
}
