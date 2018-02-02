package nexmo

import (
	"testing"
)

func TestGetAccountBalance(t *testing.T) {
	client, err := NewClient(testAPIKey, testAPISecret)
	if err != nil {
		t.Error("Failed to create Client with error:", err)
	}

	balance, err := client.Account.GetBalance()
	if err != nil {
		t.Error("Failed to get account balance with error:", err)
	}

	t.Log("Got account balance: ", balance, "€")
}
