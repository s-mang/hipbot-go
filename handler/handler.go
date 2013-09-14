package handler

import (
	"github.com/daneharrigan/hipchat"
	"strings"
)

func Reply(message hipchat.Message) (reply, kind string) {
	if strings.Contains(message.Body, "logo") {
		return "<img src='https://1.gravatar.com/avatar/fcc942ea417a208d0f5d835b8427fcc4'/>", "html"
	} else {
		return "Hello, "+nickname(message.From), "text"
	}
}

func nickname(from string) (nick string) {
	names := strings.Split(from, "/")
	return names[1]
}