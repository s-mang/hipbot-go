package main

// @botling nytimes <section-name-query>
// Get 3 most recent NY Times article in the section <section-name-query>
// return a html list of articles or else text no-results response

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const NYTIMES_ENDPOINT = "http://api.nytimes.com/svc/search/v2/articlesearch.json"

var nytimesApiKey = os.Getenv("NYTIMES_API_KEY")

type NytimesResponse struct {
	ResponseData ResponseData `json:"response"`
}

type ResponseData struct {
	Docs []Doc `json:"docs"`
}

type Doc struct {
	Url         string   `json:"web_url"`
	Snippet     string   `json:"snippet"`
	Headline    Headline `json:"headline"`
	PublishDate string   `json:"pub_date"`
}

type Headline struct {
	Main string `json:"main"`
}

func nytimes(subject string) string {
	// Set request args for nytimes search
	// Sort - newest
	// query by section name (ie. Technology)
	additionalParams := "?sort=newest&fq=section_name:" + url.QueryEscape(subject) + "&api-key=" + nytimesApiKey

	// Send GET request, collect response
	res, err := http.Get(NYTIMES_ENDPOINT + additionalParams)

	if err != nil {
		log.Println("Error occurred in HTTP GET:", err)
		return "error"
	}

	defer res.Body.Close()

	// Decode JSON body
	decoder := json.NewDecoder(res.Body)
	response := new(NytimesResponse)
	decoder.Decode(response)

	return htmlArticleList(response.ResponseData.Docs, subject)
}

// HTML formatter for a []Doc instance
func htmlArticleList(docs []Doc, querySubject string) string {
	// Check that ther actually are any results
	if len(docs) == 0 {
		return "I found nothing! So sorry."
	}

	// Title
	html := "<strong>NYTIMES ON " + strings.ToUpper(querySubject) + "</strong><br>"

	// Unordered list of first 3 articles
	html += "<ul>"
	for i := range docs {
		if i > 2 {
			break
		}

		html += "<li>"
		html += "<a href='" + docs[i].Url + "'>" + docs[i].Headline.Main + "</a>: <br>"
		html += "<ul style='list-style:none'><li>" + docs[i].Snippet + "</li></ul>"
		html += "</li><br><br>"
	}

	return html

}
