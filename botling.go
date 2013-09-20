// Botling is a neat little bot with some awesome functionality.
// He sits in your Hipchat room and obeys your every request (well, the ones he's familiar with anyway).
// At the end of the day, botling likes to remind you of how awesome you are.
// He knows how to search for nearby restaurants, get an image given a tag, search the New York Times,
// get a weather forecast, and much more.
//
// For full details on setup, implementation, and usage, see Readme.md

package main

import (
	"github.com/daneharrigan/hipchat"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	HIPCHAT_HTML_POST_ENDPOINT = "https://api.hipchat.com/v1/rooms/message"

	// Color is for HTML responses ONLY! (roughly 3/4 of commands respond in HTML)
	// Available colors are  "yellow", "red", "green", "purple", "gray", or "random"
	HIPCHAT_HTML_POST_COLOR = "gray"
)

var (
	resource = "bot" // Kind of Hipchat user (probably shouldn't change this)

	// Vars needed for Botling to ping Hipchat:
	username     = os.Getenv("BOT_USERNAME")
	mentionname  = os.Getenv("BOT_MENTIONNAME")
	fullname     = os.Getenv("BOT_FULLNAME")
	password     = os.Getenv("BOT_PASSWORD")
	roomJid      = os.Getenv("ROOM_JID")
	roomId       = os.Getenv("ROOM_ID")
	roomApiToken = os.Getenv("ROOM_API_TOKEN")

	// Var needed for location-based commands (ie. weather, nearby)
	latLngPair = os.Getenv("LAT_LNG_PAIR")

	// Var needed for Botling to respond to a request for the company logos
	logoUrl = os.Getenv("COMPANY_LOGO_URL")

	// URL used to post HTML to your Hipchat room, complete with query params
	htmlPostUrl = HIPCHAT_HTML_POST_ENDPOINT +
		"?room_id=" + url.QueryEscape(roomId) +
		"&auth_token=" + url.QueryEscape(roomApiToken) +
		"&from=" + url.QueryEscape(fullname) +
		"&color=" + HIPCHAT_HTML_POST_COLOR +
		"&message_format=html"
)

// Init a Hipchat client
// Set up Botling in your Hipchat room
// Parse incoming messages & determine if Botling needs to respond
// Get response from replyMessage(*message) (defined in speak.go)
// Speak the response via HTTP POST (HTML) or XMPP (plain text)
func main() {
	botling, err := hipchat.NewClient(username, password, resource)

	if err != nil {
		log.Println("Client error:", err)
		return
	}

	// Get Botling all set up in your Hipchat room
	botling.Status("chat")
	botling.Join(roomJid, fullname)

	// Run botling as a goroutine
	go botling.KeepAlive()

	// Check for @botling in messages & respond accordingly
	for message := range botling.Messages() {
		if strings.HasPrefix(message.Body, "@"+mentionname) {

			// Get appropriate reply message
			reply, kind := replyMessage(*message)

			if kind == "html" {
				// HTML messages sent via POST to Hipchat API
				fullPostUrl := htmlPostUrl + "&message=" + url.QueryEscape(reply)
				replyWithHtml(fullPostUrl)
			} else {
				// Plain text messages sent to Hipchat via XMPP
				botling.Say(roomJid, mentionname, reply)
			}
		}
	}
}
