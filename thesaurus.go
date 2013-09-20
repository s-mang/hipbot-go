package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var thesaurusApiKey = os.Getenv("THESAURUS_API_KEY")

type SynonymResult struct {
	NounList      WordList `json:"noun"`
	VerbList      WordList `json:"verb"`
	AdjectiveList WordList `json:"adjective"`
	AdverbList    WordList `json:"adverb"`
}

type WordList struct {
	Synonyms []string `json:"syn"`
}

func synonyms(query string) string {
	queryUrl := THESAURUS_URL + thesaurusApiKey + "/" + url.QueryEscape(query) + "/json"

	resp, err := http.Get(queryUrl)

	if err != nil {
		fmt.Println(err)
		return "error"
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	results := new(SynonymResult)

	decoder.Decode(results)

	if err != nil {
		fmt.Println(err)
		return "error"
	}

	return formattedSynonyms(*results)
}

func formattedSynonyms(results SynonymResult) string {
	html := ""
	nouns := results.NounList.Synonyms
	verbs := results.VerbList.Synonyms
	adverbs := results.AdverbList.Synonyms
	adjectives := results.AdjectiveList.Synonyms

	html += words(nouns, "Noun")
	html += words(verbs, "Verb")
	html += words(adverbs, "Adverb")
	html += words(adjectives, "Adjective")

	if html == "" {
		html += "I found nothing! So sorry."
	}

	return strings.TrimSuffix(html, "<br><br>")
}

func words(synonyms []string, kind string) string {
	html := ""
	if len(synonyms) != 0 {
		html += "<strong>" + kind + ":</strong><br>"
		html += strings.Join(synonyms, ", ")
		html += "<br><br>"
	}

	return html
}
