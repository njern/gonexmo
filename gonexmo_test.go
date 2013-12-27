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
}

func TestNexmoCreation(t *testing.T) {
	_, err := New(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}
}

func TestGetAccountBalance(t *testing.T) {
	nexmo, err := New(API_KEY, API_SECRET)
	if err != nil {
		t.Error("Failed to create Nexmo object with error:", err)
	}

	balance, err := nexmo.GetBalance()
	if err != nil {
		t.Error("Failed to get account balance with error:", err)
	}

	t.Log("Got account balance: ", balance, "€")
}
