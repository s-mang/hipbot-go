package main

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"os"
	"strings"
	"fmt"
)

var (
	appid = os.Getenv("WOLFRAM_APIID")
	postUrl = WOLFRAM_URL + 
		"?format=html"+
		"&appid=" + appid
)

type QueryResult struct {
    XMLName xml.Name `xml:"queryresult"`
    Pods   []Pod  `xml:"pod"`
	Assumptions []Assumption `xml:"assumptions"`
	Didyoumeans []Didyoumean `xml:"didyoumeans"`
}

type Pod struct {
	Markup string `xml:"markup"`
}

type Assumption struct {
	Type string `xml:"type,attr"`
	Values []Value `xml:"value"`
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
	assumptions := (results).Assumptions
	didyoumeans := (results).Didyoumeans
	
	output += searchResults(pods)
	
	if len(assumptions) > 0 && assumptions[0].Type == "Clash" {
		output += assumptionCategories(assumptions[0])
	}
	
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

func assumptionCategories(clashes Assumption) string {
	output := ""
	for i := range clashes.Values {
		if i == 0 {
			output += "For better results, prepend your query with one of the following categories+':'."+
			" For example, 'acronym:RAID'<br>Categories: "
		}
		
		output += "" + clashes.Values[i].Name
		output += "<br>"
	}
	
	
	return output
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




