package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/url"
)

type WebResults struct {
	Definition string `json:"Definition"`
	Heading string `json:"Heading"`
	ImageUrl string `json:"Image"`
	AbstractText string `json:"AbstractText"`
	MoreInfoSource string `json:"AbstractSource"`
	MoreInfoUrl string `json:"AbstractURL"`
}

func webSearch(query string) string {
	queryParams := "?q="+url.QueryEscape(query)+"&format=json"
	res, err := http.Get(DUCK_DUCK_GO_ENDPOINT + queryParams)
	
	if err != nil {
		fmt.Printf("Error occurred in HTTP GET: %s", err)
		return "error"
	}
	
	defer res.Body.Close()
	
	decoder := json.NewDecoder(res.Body)
	response := new(WebResults)
	
	decoder.Decode(response)
	
	html := htmlWebResults(*response)
	
	if html == "<ul></ul>" {
		return "I found nothing! So sorry."
	} else {
		return html
	}
}

func htmlWebResults(response WebResults) string {
	html := ""
	
	if len(response.ImageUrl) != 0 {
		html += "<img src='"+response.ImageUrl+"'>&nbsp;&nbsp;&nbsp;"
	}
	
	if len(response.Heading) != 0 {
		html += "<strong>"+response.Heading+"</strong><br>"
	}
	
	html += "<ul>"
	if len(response.Definition) != 0 {
		html += "<li><strong>Definition</strong>: "+response.Definition+"</li>"
	}
	
	if len(response.AbstractText) != 0 {
		html += "<li><strong>Abstract</strong>: "+response.AbstractText+"</li>"
		html += "<li>More info at <a href='"+response.MoreInfoUrl+"'>"+response.MoreInfoSource+"</li>"
	}
	
	html += "</ul>"
	return html
}

