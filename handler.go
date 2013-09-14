package main

import (
	"github.com/daneharrigan/hipchat"
	"strings"
)

// Returns the appropriate reply message for a given ping
func replyMessage(message hipchat.Message) (reply, kind string) {
	if strings.Contains(message.Body, "gopkg") { 
		query := strings.Split(message.Body, "gopkg")[1]
		return goSearch(query), "text"
	} else if strings.Contains(message.Body, "logo") {
		return "<img src='" + LOGO_URL + "'/>", "html"
	} else if strings.Contains(message.Body, "goodnight") {
		return "Goodnight, " + name(message.From) + ". You're awesome.", "text"
	} else {
		return "Hello, " + name(message.From), "text"
	}
}