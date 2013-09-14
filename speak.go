package main

import (
	"fmt"
	"io"
	"net/http"
)

// Post Botling's reply either via Hipchat's API (for html) or XMPP (for text)
func replyWithHtml(url string) {
	var ioReader io.Reader	
	resp, err := http.Post(url, "html", ioReader)
	
	if err != nil {
		fmt.Printf("Error occurred in HTTP POST to Hipchat API: %s\n", err)
		return
	}
	
	resp.Body.Close()
}