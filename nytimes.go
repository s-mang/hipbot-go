package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"net/url"
)

const NYTIMES_QUERY_URL = "http://api.nytimes.com/svc/search/v2/articlesearch.json?sort=newest"

var nytimesKey = os.Getenv("NYTIMES_KEY")

type NytimesResponse struct {
	ResponseData ResponseData `json:"response"`
}

type ResponseData struct {
	Docs []Doc `json:"docs"`
}

type Doc struct {
	Url string `json:"web_url"`
	Snippet string `json:"snippet"`
	Headline Headline `json:"headline"`
	PublishDate string `json:"pub_date"`
}

type Headline struct {
	Main string `json:"main"`
}

func nytimes(subject string) string {
	additionalParams := "&fq=section_name:"+url.QueryEscape(subject)+"&api-key="+nytimesKey

	res, err := http.Get(NYTIMES_QUERY_URL + additionalParams)
	
	if err != nil {
		fmt.Printf("Error occurred in HTTP GET: %s", err)
		return "error"
	}
	
	defer res.Body.Close()
	
	bd, _ := ioutil.ReadAll(res.Body)
	iface := new(NytimesResponse)
	
	_ = json.Unmarshal(bd, &iface)
	
	return htmlArticleList(iface.ResponseData.Docs, subject)
}

func htmlArticleList(docs []Doc, querySubject string) string {
	html := "<strong>NYTIMES ON "+strings.ToUpper(querySubject)+" "+prettyDate(docs[0].PublishDate)+"</strong><br>"
	html += "<ul>"
	for i := range docs {
		if i > 2 {
			break
		}
		
		html += "<li>"
		html += "<a href='"+docs[i].Url+"'>"+docs[i].Headline.Main+"</a>: <br>"
		html += "<ul style='list-style:none'><li>"+docs[i].Snippet+"</li></ul>" 
		html += "</li><br><br>"
	}
	
	return html
	
}

func prettyDate(date string) string {
	splitDate := strings.Split(date, "-")
	year := splitDate[0]
	month := splitDate[1]
	day := strings.Split(splitDate[2], "T")[0]
	
	return (month + "/" + day + "/" + year)
}






