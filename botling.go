package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/daneharrigan/hipchat"
)

var (
	username = os.Getenv("BOT_USERNAME")
	mentionname = os.Getenv("BOT_MENTIONNAME")
	fullname = os.Getenv("BOT_FULLNAME")
	password = os.Getenv("BOT_PASSWORD")
	resource = "bot"
	roomJid = os.Getenv("ROOM_JID")
)

func main() {
	botling, err := hipchat.NewClient(username, password, resource)
	
	if err != nil {
		fmt.Printf("Client error occurred: %s\n", err)
		return
	}
	
	botling.Status("chat")
	botling.Join(roomJid, fullname)
	
	botling.Say(roomJid, mentionname, "Hello all, Botling here.")
	go botling.KeepAlive()
	
	for message := range botling.Messages() {
		if strings.HasPrefix(message.Body, "@"+mentionname) {
			botling.Say(roomJid, mentionname, "Hello, "+nickname(message.From))
		}
	}
}

func nickname(from string) (nick string) {
	names := strings.Split(from, "/")
	return names[1]
}



