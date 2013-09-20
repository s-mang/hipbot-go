package main

// @botling gopkg <package-name>
// Search Godoc.org's docs for <package-name> via their API
// return a text 1-sentence explanation of package or else a text no-results response

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

const GO_DOC_URL = "http://api.godoc.org/search"

type Result struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis"`
}

type Response struct {
	Results []*Result `json:"results"`
}

// Search Godoc.org's docs via their API
func goSearch(query string) string {
	// Send GET request, collect response
	res, err := http.Get(GO_DOC_URL + "?q=" + url.QueryEscape(query))

	if err != nil {
		log.Println("Error in HTTP GET:", err)
		return "error"
	}

	defer res.Body.Close()

	// Decode JSON body
	decoder := json.NewDecoder(res.Body)
	response := new(Response)
	decoder.Decode(response)

	if len(response.Results) == 0 {
		return "I found nothing! So sorry."
	} else {
		firstResult := (*(response.Results[0]))
		textResponse := "Synopsis: " + firstResult.Synopsis
		textResponse += "\nPath: \"" + firstResult.Path + "\""
		return textResponse
	}
}
