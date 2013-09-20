package main

// This file is mainly comprised of a switch statment that maps a hipchat.Message
// to the appropriate function, and returns the result along with a string reprenting
// the format of the response ("html" or "text")

import (
	"github.com/daneharrigan/hipchat"
	"io"
	"log"
	"net/http"
	"strings"
)

// Returns the appropriate reply message for a given ping
func replyMessage(message hipchat.Message) (reply, kind string) {
	switch {

	// @botling search me HMAC
	case strings.Contains(message.Body, "search me"):
		query := strings.Split(message.Body, "search me ")[1]
		return webSearch(query), "html"

		// @botling thesaurus me challenge
	case strings.Contains(message.Body, "thesaurus me"):
		query := strings.Split(message.Body, "thesaurus me ")[1]
		return synonyms(query), "html"

		// @botling nearby sushi
	case strings.Contains(message.Body, "nearby"):
		query := strings.Split(message.Body, "nearby ")[1]
		return places(query), "html"

		// @botling nytimes technology
	case strings.Contains(message.Body, "nytimes"):
		query := strings.Split(message.Body, "nytimes ")[1]
		return nytimes(query), "html"

		// @botling image me sunset
	case strings.Contains(message.Body, "image me"):
		query := strings.Split(message.Body, "image me ")[1]
		return flickrSearch(query), "html"

		// @botling weather me today
	case strings.Contains(message.Body, "weather me"):
		query := strings.Split(message.Body, "weather me ")[1]
		return weather(query), "html"

	// @botling trivia me today
	case strings.Contains(message.Body, "trivia me today"):
		return numberTrivia("today"), "text"

		// @botling trivia me 123
	case strings.Contains(message.Body, "trivia me"):
		query := strings.Split(message.Body, "trivia me ")[1]
		return numberTrivia(query), "text"

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
		return "<img src='" + logoUrl + "'/>", "html"

		// @botling goodnight
	case strings.Contains(message.Body, "goodnight"):
		return "Goodnight, " + name(message.From) + ". You're awesome.", "text"

		// @botling foo
	default:
		return "Hello, " + name(message.From), "text"
	}
}

// Post Botling's reply via Hipchat's API (for html responses)
// Note that text responses will be submitted to Hipchat via XMPP (see botling.go)
func replyWithHtml(url string) {
	var ioReader io.Reader
	resp, err := http.Post(url, "html", ioReader)

	if err != nil {
		log.Println("Error in POST to Hipchat API:", err)
		return
	}

	resp.Body.Close()
}

// Grabs a user's full name from a hipchat.Message.From string
func name(from string) (nick string) {
	names := strings.Split(from, "/")
	return names[1]
}
