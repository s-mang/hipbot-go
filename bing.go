package main

// NOT USED FOR `image me` QUERIES -- See bing.go

// @botling image me <query>
// Search Bing's photos for <query> via their API
// return an HTML image tag or else a text no-results response

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const BING_ENDPOINT = "https://api.datamarket.azure.com/Bing/Search/Image"

// Encode your application key at www.base64encode.org
// You must encode a string consisting of your key twice, separated by a colon.
//   ie. MY_KEY:MY_KEY
var bingAuth = os.Getenv("BING_BASE64_API_KEY")

type BingResponse struct {
	ResultList BingPhotoResultList `json:"d"`
}

type BingPhotoResultList struct {
	BingPhotoResults []BingPhotoResult `json:"results"`
}

type BingPhotoResult struct {
	Thumbnail Thumbnail
}

type Thumbnail struct {
	MediaUrl string
}

// Search Azure/Bing's database of photos for ones with tags matching <query>
// Return an HTML image tag to be POSTed to Hipchat room
func bingImageSearch(query string) string {
	requestArgs := "?$format=json"
	requestArgs += "&Query=%27" + url.QueryEscape(query) + "%27"
	requestArgs += "&$top=1"
	// Get a random one of the first 20 results
	requestArgs += "&$skip=" + randNumParam(10)

	req, err := http.NewRequest("GET", BING_ENDPOINT+requestArgs, nil)
	if err != nil {
		log.Println("Error forming request", err.Error())
		return "error"
	}

	req.Header.Add("Authorization", "Basic "+bingAuth)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error in HTTP GET:", err)
		return "error"
	}

	defer resp.Body.Close()

	// Decode JSON body
	decoder := json.NewDecoder(resp.Body)
	photos := new(BingResponse)
	decoder.Decode(photos)

	imageCount := len(photos.ResultList.BingPhotoResults)

	// Check if we got any photos back
	if imageCount == 0 {
		return "I found nothing! So sorry."
	} else {
		src := photos.ResultList.BingPhotoResults[0].Thumbnail.MediaUrl
		return "<img src='" + src + "'>"
	}

}

// Returns a random number between 0 and max in string format
func randNumParam(max int) string {
	// Seed with current time
	source := rand.NewSource(time.Now().UnixNano())
	rander := rand.New(source)
	rInt := rander.Int()

	// make the int smaller (between 0 and max) with mod (%)
	smallerRInt := rInt % max
	return strconv.Itoa(smallerRInt)
}
