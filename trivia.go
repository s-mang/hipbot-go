package main

// @hipbot trivia me <numerical-query>
// @hipbot trivia me today
// Get number/date trivia about the query (either a number or else "today")
// return a text 1-sentence factoid

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const NUMBERS_API_ENDPOINT = "http://numbersapi.com/"

// Get number/date trivia from numbersapi.com (free number trivia API. Wheee!)
// Return text string for Hipchat response
func numberTrivia(query string) string {
	// Compiler barfs if we define these with `:=` in the if/else block
	var triviaResp *http.Response
	var triviaErr error

	if query == "today" {
		// Get today's date, convert month and day to int -> string, and send GET to API
		today := time.Now()
		_, month, day := today.Date()

		triviaResp, triviaErr = http.Get(NUMBERS_API_ENDPOINT + strconv.Itoa(int(month)) + "/" + strconv.Itoa(int(day)) + "/date")
	} else {
		// Straight numbers are simpler
		triviaResp, triviaErr = http.Get(NUMBERS_API_ENDPOINT + query + "/math")
	}

	// Check for error
	if triviaErr != nil {
		log.Println("Error in HTTP GET:", triviaErr)
		return "error"
	}

	// Response is text, not JSON - so we can just read it and voiala!
	byteBody, stringErr := ioutil.ReadAll(triviaResp.Body)

	if stringErr != nil {
		log.Println("Error reading response body:", stringErr)
		return "error"
	}

	return string(byteBody)
}
