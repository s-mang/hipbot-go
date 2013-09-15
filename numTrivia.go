package main

import (
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"strconv"
)

const NUMBERS_API_URL = "http://numbersapi.com/"


func numTrivia(query string) string {
	var triviaResp *http.Response
	var triviaErr error
	
	if query == "today" {
		today := time.Now()
		day, month, _ := today.Date()
	
		triviaResp, triviaErr = http.Get(NUMBERS_API_URL + "/"+strconv.Itoa(int(month))+"/"+strconv.Itoa(int(day))+"/date")
	} else {
		triviaResp, triviaErr = http.Get(NUMBERS_API_URL + "/"+query+"/math")
	}
	if triviaErr != nil {
		fmt.Printf("Error occurred in HTTP GET: %s", triviaErr)
		return "error"
	}
	
	stringBody, stringErr := ioutil.ReadAll(triviaResp.Body)
	
	if stringErr != nil {
		fmt.Printf("Error reading response body: %s", stringErr)
		return "error"
	}
	
	return string(stringBody)
	
}