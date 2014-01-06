/*
Package nexmo implements a simple client library for accessing the Nexmo API.

Usage is simple. Create a nexmo.Client instance with NewClientFromAPI(), providing
your API key and API secret. Then send messages with SendTextMessage() or
SendFlashMessage(). The API will return a MessageResponse which you can
use to see if your message went through, how much it cost, etc.
*/
package nexmo

const (
	apiRoot = "https://rest.nexmo.com"
)
