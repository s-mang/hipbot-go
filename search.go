package main

// @botling search me <query>
// Get relevant search results from duckduckgo.com's API
// return an HTML list of results or else text no-results response

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const DUCK_DUCK_GO_ENDPOINT = "https://api.duckduckgo.com/"

type WebResults struct {
	Definition     string `json:"Definition"`
	Heading        string `json:"Heading"`
	ImageUrl       string `json:"Image"`
	AbstractText   string `json:"AbstractText"`
	MoreInfoSource string `json:"AbstractSource"`
	MoreInfoUrl    string `json:"AbstractURL"`
}

// Search duckduckgo.com for <query> via their API
func webSearch(query string) string {
	// Set the query params
	queryParams := "?q=" + url.QueryEscape(query) + "&format=json"

	// Send GET request, collect response
	res, err := http.Get(DUCK_DUCK_GO_ENDPOINT + queryParams)

	if err != nil {
		log.Printf("Error in HTTP GET:", err)
		return "error"
	}

	defer res.Body.Close()

	// Decode JSON Response
	decoder := json.NewDecoder(res.Body)
	response := new(WebResults)
	decoder.Decode(response)

	// Get HTML formatted version of response
	html := htmlWebResults(*response)

	if html == "<ul></ul>" { // If there were no results
		return "I found nothing! So sorry."
	} else {
		return html
	}
}

// Format the JSON search results and return a nice, pretty, HTML version
func htmlWebResults(response WebResults) string {
	html := ""

	// Search results from duckduckgo.com's API often include an image
	// For funzies, we use the image in the hipchat ping
	if len(response.ImageUrl) != 0 {
		html += "<img src='" + response.ImageUrl + "'>&nbsp;&nbsp;&nbsp;"
	}

	// Heading is what duckduckgo assumed was meant by <query>
	// (Often response.Heading just equals <query>)
	if len(response.Heading) != 0 {
		html += "<strong>" + response.Heading + "</strong><br>"
	}

	// Start list of results
	html += "<ul>"

	// Duckduckgo's definition of <query> (often blank)
	if len(response.Definition) != 0 {
		html += "<li><strong>Definition</strong>: " + response.Definition + "</li>"
	}

	// "Abstrac Text" (the first few lines of the wikipedia page for <query>)
	if len(response.AbstractText) != 0 {
		html += "<li><strong>Abstract</strong>: " + response.AbstractText + "</li>"
		html += "<li>More info at <a href='" + response.MoreInfoUrl + "'>" + response.MoreInfoSource + "</li>"
	}

	// End list
	html += "</ul>"
	return html
}
