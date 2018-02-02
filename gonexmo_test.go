package nexmo

import (
	"fmt"
	"os"
	"testing"
)

var (
	testAPIKey      string
	testAPISecret   string
	testPhoneNumber string
	testFrom        string
)

func init() {
	testAPIKey = os.Getenv("NEXMO_KEY")
	if testAPIKey == "" {
		fmt.Println("No API key specified. Please set NEXMO_KEY")
		os.Exit(1)
	}

	testAPISecret = os.Getenv("NEXMO_SECRET")
	if testAPISecret == "" {
		fmt.Println("No API secret specified. Please set NEXMO_SECRET")
		os.Exit(1)
	}

	testPhoneNumber = os.Getenv("NEXMO_NUM")

	// Set a custom from value, or use the default. If you get error 15 when
	// sending a message ("Illegal Sender Address - rejected") try setting this
	// to your nexmo phone number.
	testFrom = os.Getenv("NEXMO_FROM")
	if testFrom == "" {
		testFrom = "gonexmo/test"
	}
}

func TestNexmoCreation(t *testing.T) {
	_, err := NewClient(testAPIKey, testAPISecret)
	if err != nil {
		t.Error("failed to create Client with error:", err)
	}
}
