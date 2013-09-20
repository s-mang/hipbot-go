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
	HIPCHAT_HTML_POST_COLOR    = "gray"
)

var (
	resource    = "bot"
	username    = os.Getenv("BOT_USERNAME")
	mentionname = os.Getenv("BOT_MENTIONNAME")
	fullname    = os.Getenv("BOT_FULLNAME")
	password    = os.Getenv("BOT_PASSWORD")
	roomJid     = os.Getenv("ROOM_JID")
	roomId      = os.Getenv("ROOM_ID")
	roomApiId   = os.Getenv("ROOM_APIID")
	latLngPair  = os.Getenv("LAT_LNG_PAIR")

	htmlPostUrl = HIPCHAT_HTML_POST_ENDPOINT +
		"?room_id=" + url.QueryEscape(roomId) +
		"&auth_token=" + url.QueryEscape(roomApiId) +
		"&from=" + url.QueryEscape(fullname) +
		"&color=" + HIPCHAT_HTML_POST_COLOR +
		"&message_format=html"
)

func main() {
	botling, err := hipchat.NewClient(username, password, resource)

	if err != nil {
		log.Println("Client error:", err)
		return
	}

	welcomeBotling(*botling)

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

// Help Botling join the hipchat room with a status set to 'chat',
// Run Botling as a goroutine
func welcomeBotling(botling hipchat.Client) {
	botling.Status("chat")
	botling.Join(roomJid, fullname)

	go botling.KeepAlive()
}
