package main

import (
	"fmt"
	"os"
	"strings"
	"io"
	"net/http"
	"net/url"
	"github.com/daneharrigan/hipchat"
	"./handler"
)

const (
	POST_URL = "https://api.hipchat.com/v1/rooms/message"
	POST_COLOR = "gray"
)

var (
	username = os.Getenv("BOT_USERNAME")
	mentionname = os.Getenv("BOT_MENTIONNAME")
	fullname = os.Getenv("BOT_FULLNAME")
	password = os.Getenv("BOT_PASSWORD")
	roomJid = os.Getenv("ROOM_JID")
	roomId = os.Getenv("ROOM_ID")
	roomApiId = os.Getenv("ROOM_APIID")
	resource = "bot"
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
			reply, kind := handler.Reply(*message)
			if kind == "html" {
				url := POST_URL +
					"?room_id=" + url.QueryEscape(roomId) + 
					"&auth_token=" + url.QueryEscape(roomApiId) +
					"&from=" + url.QueryEscape(fullname) +
					"&message_format=html" +
					"&color=" + POST_COLOR +
					"&message=" + url.QueryEscape(reply)
				
				fmt.Println(url)
				
				var ioReader io.Reader	
				resp, err := http.Post(url, "html", ioReader)
				if err != nil {
					fmt.Printf("Error occurred in HTTP POST to Hipchat API: %s\n", err)
					return
				}
				resp.Body.Close()
			} else {
				botling.Say(roomJid, mentionname, reply)
			}
		}
	}
}



