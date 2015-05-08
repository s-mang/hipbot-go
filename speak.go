package main

// This file is mainly comprised of a switch statment that maps a hipchat.Message
// to the appropriate function, and returns the result along with a string reprenting
// the format of the response ("html" or "text")

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Returns the appropriate reply message for a given ping
func replyMessage(message string) (reply, kind string) {
	switch {

	// // @hipbot register fork jhonnas/goblob
	// case strings.Contains(message, "register fork"):
	// 	fork := strings.Split(message, "register fork ")[1]
	// 	return registerFork(fork), "html"

	// // @hipbot list forks
	// case strings.Contains(message, "list forks"):
	// 	return listWatchingForks(), "html"

	// // @hipbot forks
	// case strings.Contains(message, "forks"):
	// 	return behindForksHTML(), "html"

	// @hipbot search me HMAC
	case strings.Contains(message, "search me"):
		query := strings.Split(message, "search me ")[1]
		return webSearch(query), "html"

		// @hipbot thesaurus me challenge
	case strings.Contains(message, "thesaurus me"):
		query := strings.Split(message, "thesaurus me ")[1]
		return synonyms(query), "html"

		// @hipbot nearby sushi
	case strings.Contains(message, "nearby"):
		query := strings.Split(message, "nearby ")[1]
		return places(query), "html"

		// @hipbot nytimes technology
	case strings.Contains(message, "nytimes"):
		query := strings.Split(message, "nytimes ")[1]
		return nytimes(query), "html"

		// @hipbot image me sunset
	case strings.Contains(message, "image me"):
		query := strings.Split(message, "image me ")[1]
		return bingImageSearch(query), "html"

		// @hipbot weather me today
	case strings.Contains(message, "weather me"):
		query := strings.Split(message, "weather me ")[1]
		return weather(query), "html"

		// @hipbot trivia me today
	case strings.Contains(message, "trivia me today"):
		return numberTrivia("today"), "text"

		// @hipbot trivia me 123
	case strings.Contains(message, "trivia me"):
		query := strings.Split(message, "trivia me ")[1]
		return numberTrivia(query), "text"

		// @hipbot wolfram me pi
	case strings.Contains(message, "wolfram me"):
		query := strings.Split(message, "wolfram me ")[1]
		return wolframSearch(query), "html"

		// @hipbot gopkg math
	case strings.Contains(message, "gopkg"):
		query := strings.Split(message, "gopkg ")[1]
		return goSearch(query), "text"

	// 	// @hipbot logo
	// case strings.Contains(message, "logo"):
	// 	return "<img src='" + logoUrl + "'/>", "html"

	// @hipbot goodnight
	case strings.Contains(message, "goodnight"):
		return "Goodnight. You're awesome.", "text"

		// @hipbot foo
	default:
		return "Hello!", "text"
	}
}

// Post Hipbot's reply via Hipchat's API (for html messages)
// Note that text responses will be submitted to Hipchat via XMPP (see hipbot.go)
func speakInHTML(message string, notify bool) {
	var ioReader io.Reader
	messageURL := htmlPostUrl + "&message=" + url.QueryEscape(message)
	if notify {
		messageURL += "&notify=1"
	}

	log.Println(messageURL)

	resp, err := http.Post(messageURL, "application/x-www-form-urlencoded", ioReader)

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
