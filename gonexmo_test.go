package nexmo

import (
	"fmt"
	"os"
	"testing"
)

var (
	API_KEY, API_SECRET, TEST_PHONE_NUMBER string
)

func init() {
	API_KEY = os.Getenv("NEXMO_KEY")
	if API_KEY == "" {
		fmt.Println("No API key specified. Please set NEXMO_KEY")
		os.Exit(1)
	}

	API_SECRET = os.Getenv("NEXMO_SECRET")
	if API_SECRET == "" {
		fmt.Println("No API secret specified. Please set NEXMO_SECRET")
		os.Exit(1)
	}

	TEST_PHONE_NUMBER = os.Getenv("NEXMO_NUM")
	if TEST_PHONE_NUMBER == "" {
		fmt.Println("No test phone number specified. Please set NEXMO_NUM")
		os.Exit(1)
	}
}

func TestNexmoCreation(t *testing.T) {
	_, err := NewConn(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}
}

func TestGetAccountBalance(t *testing.T) {
	nexmo, err := NewConn(API_KEY, API_SECRET)
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
	nexmo, err := NewConn(API_KEY, API_SECRET)
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
	nexmo, err := NewConn(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}

	messageResponse, err := nexmo.SendFlashMessage("go-nexmo", "00358123412345", "Looks like go-nexmo works great, we should definitely buy that njern guy a beer!", "001", false)
	if err != nil {
		t.Error("Failed to send flash message (class 0 SMS) with error:", err)
	}

	t.Log("Sent Flash SMS, response was: ", messageResponse)
}
