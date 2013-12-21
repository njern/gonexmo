# go-nexmo

go-nexmo is a [Go](http://golang.org/) library specifically for sending SMS's with [Nexmo](https://www.nexmo.com/). You can also check your balance and in the future I may add support for the rest of their API (check prices before sending messages, etc)


## Installation

Assuming you have a working Go environment, installation is simply:

    go get github.com/njern/gonexmo

You can take a look at the documentation locally with:

	godoc github.com/njern/gonexmo

The included tests in `gonexmo_test.go` also illustrate usage of the package.

**Note:** You must enter valid API credentials and a valid phone number in `gonexmo_test.go` or the tests will fail!


## Usage
    import "github.com/njern/gonexmo"

    nexmo, _ := gonexmo.NexmoWithKeyAndSecret("API_KEY_GOES_HERE", "API_SECRET_GOES_HERE")

    // Test if it works by retrieving your account balance
    balance, err := nexmo.GetBalance()

    // Send an SMS (from, to, text, reference_id, status_report_required)
    // See https://docs.nexmo.com/index.php/sms-api/send-message for details.
    messageResponse, err := nexmo.SendTextMessage("go-nexmo", "00358123412345", "Looks like go-nexmo works great, we should definitely buy that njern guy a beer!", "001", false)

