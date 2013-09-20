package main

import (
	"fmt"
	"github.com/daneharrigan/hipchat"
	"net/url"
	"os"
	"strings"
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

	htmlPostUrl = POST_URL +
		"?room_id=" + url.QueryEscape(roomId) +
		"&auth_token=" + url.QueryEscape(roomApiId) +
		"&from=" + url.QueryEscape(fullname) +
		"&message_format=html" +
		"&color=" + POST_COLOR +
		"&message="
)

func main() {
	botling, err := hipchat.NewClient(username, password, resource)

	if err != nil {
		fmt.Printf("Client error occurred: %s\n", err)
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
				postUrl := htmlPostUrl + url.QueryEscape(reply)
				replyWithHtml(postUrl)
			} else {
				// Plain text messages sent to Hipchat via XMPP
				botling.Say(roomJid, mentionname, reply)
			}
		}
	}
}

// Help Botling join the hipchat room with a status set to 'chat',
// make Botling say hello, and run Botling as a goroutine
func welcomeBotling(botling hipchat.Client) {
	botling.Status("chat")
	botling.Join(roomJid, fullname)

	go botling.KeepAlive()
}
