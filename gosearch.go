package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
)

const GO_DOC_URL = "http://api.godoc.org/search?q="

type Result struct {
	Path     string `json:"path"`
	Synopsis string `json:"synopsis"`
}

type Response struct {
	Results []*Result `json:"results"`
}

func goSearch(query string) string {
	res, err := http.Get(GO_DOC_URL + url.QueryEscape(query))
	
	if err != nil {
		fmt.Printf("Error occurred in HTTP GET: %s", err)
		return "error"
	}
	
	defer res.Body.Close()
	
	decoder := json.NewDecoder(res.Body)
	response := new(Response)
	
	decoder.Decode(response)
	
	return (*(response.Results[0])).Synopsis
}