package nexmo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Account represents the user's account. Used when retrieving e.g current
// balance.
type Account struct {
	client *Client
}

// GetBalance retrieves the current balance of your Nexmo account in Euros (€)
func (nexmo *Account) GetBalance() (float64, error) {
	// Declare this locally, since we are only going to return a float64.
	type AccountBalance struct {
		Value float64 `json:"value"`
	}

	var accBalance *AccountBalance

	client := &http.Client{}
	r, _ := http.NewRequest("GET", apiRoot+"/account/get-balance/"+
		nexmo.client.apiKey+"/"+nexmo.client.apiSecret, nil)
	r.Header.Add("Accept", "application/json")

	resp, err := client.Do(r)
	defer resp.Body.Close()

	if err != nil {
		return 0.0, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &accBalance)
	if err != nil {
		return 0.0, err
	}
	return accBalance.Value, nil
}
