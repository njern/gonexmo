/*
Package nexmo implements a simple client library for accessing the Nexmo API.

<<<<<<< HEAD
Usage is simple. Create a nexmo.Client instance with NewClientFromAPI(), providing
your API key and API secret. Then send messages with SendTextMessage() or
SendFlashMessage(). The API will return a MessageResponse which you can
=======
Usage is simple. Create a nexmo.Client instance with NewClientFromAPI(),
providing your API key and API secret. Compose a new Message and then call
Client.SMS.Send(Message). The API will return a MessageResponse which you can
>>>>>>> 17fbd3f4f8eef122e74920e99ad6f75630e5bed2
use to see if your message went through, how much it cost, etc.
*/
package nexmo

const (
	apiRoot = "https://rest.nexmo.com"
)
