package main

import "strings"

// Grabs a user's full name from a hipchat.Message.From string
func name(from string) (nick string) {
	names := strings.Split(from, "/")
	return names[1]
}