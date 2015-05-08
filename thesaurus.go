package main

// @hipbot thesaurus me <word-query>
// Get all synonyms from bighugelabs.com's API
// return an HTML list of synonyms or else text no-results response

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const BIG_HUGE_LABS_ENDPOINT = "http://words.bighugelabs.com/api/2/"

var bigHugeLabsApiKey = os.Getenv("BIG_HUGE_LABS_API_KEY")

type SynonymResult struct {
	NounList      WordList `json:"noun"`
	VerbList      WordList `json:"verb"`
	AdjectiveList WordList `json:"adjective"`
	AdverbList    WordList `json:"adverb"`
}

type WordList struct {
	Synonyms []string `json:"syn"`
}

// Get synonyms for <query> via bighugelabs's API
// Return pretty HTML formatted response of synonyms for all
// versions of the <query> word (ie. noun, adjective, adverb, verb, etc)
func synonyms(query string) string {
	// Set query url with params - ask for JSON response
	queryUrl := BIG_HUGE_LABS_ENDPOINT + bigHugeLabsApiKey + "/" + url.QueryEscape(query) + "/json"

	// Send GET request, collect response
	resp, err := http.Get(queryUrl)

	if err != nil {
		log.Println(err)
		return "error"
	}

	defer resp.Body.Close()

	// Decode JSON response
	decoder := json.NewDecoder(resp.Body)
	results := new(SynonymResult)
	decoder.Decode(results)

	if err != nil {
		log.Println(err)
		return "error"
	}

	// Return nicely formatted HTML synonym results
	return formattedSynonyms(*results)
}

func formattedSynonyms(results SynonymResult) string {
	html := ""

	// We only use the noun, verb, adverb, and adjective synonyms
	html += words(results.NounList.Synonyms, "Noun")
	html += words(results.VerbList.Synonyms, "Verb")
	html += words(results.AdverbList.Synonyms, "Adverb")
	html += words(results.AdjectiveList.Synonyms, "Adjective")

	if html == "" {
		html += "I found nothing! So sorry."
	}

	// Get rid of the line breaks below the last entry
	return strings.TrimSuffix(html, "<br><br>")
}

// Return the word kind in bold, followed by a comma separated
// list of the corresponding synonyms
// Insert two line breaks at the bottom of each section for aesthetics
func words(synonyms []string, kind string) string {
	html := ""
	if len(synonyms) != 0 {
		// Bold title (kind of word)
		html += "<strong>" + kind + ":</strong><br>"
		// Comma separated list of synonyms
		html += strings.Join(synonyms, ", ")
		// Double line break after list
		html += "<br><br>"
	}

	return html
}
