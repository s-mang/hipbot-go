package main

// @botling Wolfram me pi
// Get first 6 Wolfram search results
// return an HTML list of these results (includes embedded images)
// Wolfram (Alpha) is a search engine for computational information
// For information on Wolfram alpha search API, visit http://products.wolframalpha.com/api/‎

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const WOLFRAM_ENDPOINT = "http://api.wolframalpha.com/v2/query"

var (
	appid          = os.Getenv("WOLFRAM_API_ID")
	wolframPostUrl = WOLFRAM_ENDPOINT +
		"?format=html" +
		"&appid=" + appid
)

type WolframResult struct {
	XMLName     xml.Name     `xml:"queryresult"`
	Pods        []Pod        `xml:"pod"`
	Didyoumeans []Didyoumean `xml:"didyoumeans"`
}

// A pod is just a wrapper for one kind of response information
// ie. "Definition", "Synonyms", "Graph", "Usage", etc.
type Pod struct {
	Markup string `xml:"markup"`
}

type Value struct {
	Name string `xml:"name,attr"`
	Desc string `xml:"desc,attr"`
}

type Didyoumean struct {
	Didyoumean string `xml:"didyoumean"`
}

func wolframSearch(query string) string {
	// Set request input (query) for Wolfram alpha search
	// Send GET request, collect response
	resp, err := http.Get(wolframPostUrl + "&input=" + url.QueryEscape(query))

	if err != nil {
		fmt.Printf("Error in HTTP GET: %s", err)
		return "error"
	}

	defer resp.Body.Close()

	// Decode XML response
	decoder := xml.NewDecoder(resp.Body)
	results := new(WolframResult)
	decoder.Decode(results)

	return fullResponse(*results)
}

// Returns a formatted HTML response with embedded images IF results exists
// Also, returns some "Did you mean.." text IF there are alternative suggestions
func fullResponse(results WolframResult) string {
	output := ""

	pods := (results).Pods
	didyoumeans := (results).Didyoumeans

	output += htmlWolframRespose(pods)

	if len(didyoumeans) > 0 {
		output += didYouMeanText(didyoumeans)
	}

	return output
}

// Formats the response pods as HTML in a nice, easy to read way
func htmlWolframRespose(pods []Pod) string {
	// Check for no pod results
	if len(pods) == 0 {
		return "I found nothing! So sorry. Your query may be too general. <br>"
	}

	// Format pods in an unordered list
	output := "<ul>"
	// Remove the xml wrapper ("<![CDATA[" and "]]")
	// Put images on a new line with a 3-space indentation
	replacer := strings.NewReplacer("<![CDATA[", "", "]]>", "", "<img", "<br>&nbsp;&nbsp;&nbsp;<img")

	// If there is too much text, Hipchat will reject the POST
	// Since Wolfram can return a huge number of result pods, we cut them off at 6
	for i := range pods[:6] {
		output += "<li>" + replacer.Replace(pods[i].Markup) + "</li>"
	}
	output += "</ul>"

	return output
}

// didyoumeans are suggestions that relate to the <query>
// Usually didyoumeans correspond with no pod results, but NOT ALWAYS!
func didYouMeanText(didyoumeans []Didyoumean) string {
	if len(didyoumeans) == 0 {
		return ""
	}

	// Header text for didyoumeans
	output := "Why don't you try one of the following: <br>"

	// Format didyoumeans in an unordered list
	output += "<ul>"
	for i := range didyoumeans {
		output += "<li>" + didyoumeans[i].Didyoumean + "</li>"
	}
	output += "</ul>"

	return output
}
