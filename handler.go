package main

import (
	"github.com/daneharrigan/hipchat"
	"strings"
)

// Returns the appropriate reply message for a given ping
func replyMessage(message hipchat.Message) (reply, kind string) {
	switch {
		// @botling weather me today
	case strings.Contains(message.Body, "weather me"):
		query := strings.Split(message.Body, "weather me ")[1]
		return weather(query), "html"
		// @botling trivia me today
 	case strings.Contains(message.Body, "trivia me today"):
		return numTrivia("today"), "text"
		// @botling trivia me number 123
	case strings.Contains(message.Body, "trivia me number"):
		query := strings.Split(message.Body, "trivia me number ")[1]
		return numTrivia(query), "text"
		// @botling wolfram me pi
	case strings.Contains(message.Body, "wolfram me"):
		query := strings.Split(message.Body, "wolfram me ")[1]
		return wolframSearch(query), "html"
		// @botling gopkg math
	case strings.Contains(message.Body, "gopkg"):
		query := strings.Split(message.Body, "gopkg ")[1]
		return goSearch(query), "text"
		// @botling logo
	case strings.Contains(message.Body, "logo"):
		return "<img src='" + LOGO_URL + "'/>", "html"
		// @botling goodnight
	case strings.Contains(message.Body, "goodnight"): 
		return "Goodnight, " + name(message.From) + ". You're awesome.", "text"
	default:
		return "Hello, " + name(message.From), "text"
	}
}