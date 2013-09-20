package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const WOLFRAM_URL = "http://api.wolframalpha.com/v2/query"

var (
	appid   = os.Getenv("WOLFRAM_API_ID")
	postUrl = WOLFRAM_URL +
		"?format=html" +
		"&appid=" + appid
)

type QueryResult struct {
	XMLName     xml.Name     `xml:"queryresult"`
	Pods        []Pod        `xml:"pod"`
	Didyoumeans []Didyoumean `xml:"didyoumeans"`
}

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
	resp, err := http.Get(postUrl + "&input=" + url.QueryEscape(query))

	if err != nil {
		fmt.Printf("Error in HTTP GET: %s", err)
		return "error"
	}

	decoder := xml.NewDecoder(resp.Body)
	results := new(QueryResult)

	decoder.Decode(results)

	return fullResponse(*results)
}

func fullResponse(results QueryResult) string {
	output := ""

	pods := (results).Pods
	didyoumeans := (results).Didyoumeans

	output += searchResults(pods)

	if len(didyoumeans) > 0 {
		output += didYouMeanText(didyoumeans)
	}

	return output
}

func searchResults(pods []Pod) string {
	output := ""
	replacer := strings.NewReplacer("<![CDATA[", "", "]]>", "", "<img", "<br>&nbsp;&nbsp;&nbsp;<img")

	for i := range pods {
		if i > 5 {
			break
		}

		output += replacer.Replace(pods[i].Markup)
		output += "<br>"
	}

	if output == "" {
		return "I found nothing! So sorry. Your query may be too general. <br>"
	} else {
		return output
	}

}

func didYouMeanText(didyoumeans []Didyoumean) string {
	output := ""
	for i := range didyoumeans {
		if i == 0 {
			output += "Why don't you try one of the following: <br>"
		}

		output += "&nbsp;&nbsp;* " + didyoumeans[i].Didyoumean
		output += "<br>"
	}

	return output
}
