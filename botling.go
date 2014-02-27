// Botling is a neat little bot with some awesome functionality.
// He sits in your Hipchat room and obeys your every request (well, the ones he's familiar with anyway).
// At the end of the day, botling likes to remind you of how awesome you are.
// He knows how to search for nearby restaurants, get an image given a tag, search the New York Times,
// get a weather forecast, and much more.
//
// For full details on setup, implementation, and usage, see Readme.md

package main

import (
	"fmt"
	"github.com/Sproutling/go-xmpp"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	// HipChat jabber info
	HIPCHAT_JABBER_CONNECT_URL  = "chat.hipchat.com"
	HIPCHAT_JABBER_CONNECT_PORT = "5223"

	HIPCHAT_HTML_POST_ENDPOINT = "https://api.hipchat.com/v1/rooms/message"

	// Color is for HTML responses ONLY! (roughly 3/4 of commands respond in HTML)
	// Available colors are  "yellow", "red", "green", "purple", "gray", or "random"
	HIPCHAT_HTML_POST_COLOR = "gray"
)

var (
	resource = "bot" // Kind of Hipchat user (probably shouldn't change this)

	// Database connection string
	dburi = os.Getenv("DATABASE_URL")

	// Company github organization name, used for checking for newly-updated forks
	forkOwner = os.Getenv("GITHUB_FORK_OWNER")

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

var DB gorm.DB

// Init a Hipchat client
// Set up Botling in your Hipchat room
// Parse incoming messages & determine if Botling needs to respond
// Get response from replyMessage(*message) (defined in speak.go)
// Speak the response via HTTP POST (HTML) or XMPP (plain text)
func main() {
	var err error
	DB, err = gorm.Open("postgres", dburi)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database. Error: '%v'", err))
	}

	var botling *xmpp.Client
	fullConnectURL := HIPCHAT_JABBER_CONNECT_URL + ":" + HIPCHAT_JABBER_CONNECT_PORT
	jabberId := username + "@" + HIPCHAT_JABBER_CONNECT_URL

	opts := xmpp.Options{
		Host:     fullConnectURL,
		User:     jabberId,
		Password: password,
		Debug:    false,
		Resource: resource,
	}

	// Initialize client
	botling, err = opts.NewClient()
	if err != nil {
		log.Println("Client error:", err)
		return
	}

	defer botling.Close()

	// Join main room
	botling.JoinMUC(roomJid, fullname)

	// Set up fork notifications
	go scheduleForkUpdates(24*time.Hour, "12:40")

	// Check for @botling in messages & respond accordingly
	for {
		message, err := botling.Recv()
		if err != nil {
			log.Fatal(err)
		}
		if chatMsg, ok := message.(xmpp.Chat); ok {
			if strings.HasPrefix(chatMsg.Text, "@"+mentionname) {
				// Get appropriate reply message
				reply, kind := replyMessage(chatMsg.Text)

				if kind == "html" {
					// HTML messages sent via POST to Hipchat API
					speakInHTML(reply, false)
				} else {
					// Plain text messages sent to Hipchat via XMPP
					botling.Send(xmpp.Chat{Remote: roomJid, Type: "chat", Text: reply})
				}
			}
		}
	}
}
