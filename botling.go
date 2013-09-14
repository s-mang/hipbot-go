package main

import (
	"fmt"
	"strings"
	"github.com/daneharrigan/hipchat"
)

const (
	USERNAME = "65829_460351"
	MENTIONNAME = "botling"
	FULLNAME = "Botling Sprout"
	PASSWORD = "password"
	RESOURCE = "bot"
	ROOMJID = "65829_s._adams_codes@conf.hipchat.com"
)

func main() {
	botling, err := hipchat.NewClient(USERNAME, PASSWORD, RESOURCE)
	
	if err != nil {
		fmt.Printf("Client error occurred: %s\n", err)
		return
	}
	
	botling.Status("chat")
	botling.Join(ROOMJID, FULLNAME)
	
	botling.Say(ROOMJID, MENTIONNAME, "I'm Here!")
	
	for message := range botling.Messages() {
		if strings.HasPrefix(message.Body, "@"+MENTIONNAME) {
			botling.Say(ROOMJID, MENTIONNAME, "Hello, "+nickname(message.From))
		}
	}
}

func nickname(from string) (nick string) {
	names := strings.Split(from, "/")
	return names[1]
}



