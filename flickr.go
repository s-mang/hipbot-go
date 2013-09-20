package main

// @botling image me <query>
// Search Flickr's photos for <query> via their REST API
// return an HTML image tag or else a text no-results response

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const FLICKR_ENDPOINT = "http://api.flickr.com/services/rest"

var flickrApiKey = os.Getenv("FLICKR_API_KEY")

type PhotoResponse struct {
	Page Page `json:"photos"`
}

type Page struct {
	PhotoList []Photo `json:"photo"`
}

type Photo struct {
	Id     string      `json:"id"`
	Secret string      `json:"secret"`
	Server string      `json:"server"`
	Farm   json.Number `json:"farm"`
}

func flickrSearch(query string) string {
	// Set request args for flickr.photos.search
	requestArgs := "?method=flickr.photos.search"
	requestArgs += "&api_key=" + flickrApiKey
	requestArgs += "&format=json"
	requestArgs += "&tags=" + url.QueryEscape(query)
	// Only get one image, sort by relevance
	requestArgs += "&page=1&per_page=1&sort=relevance&media=photo"

	// Send GET request, collect response
	res, err := http.Get(FLICKR_ENDPOINT + requestArgs)
	if err != nil {
		log.Println("Error in HTTP GET:", err)
		return "error"
	}

	defer res.Body.Close()

	// Flickr wraps its response json in a function
	body, err := unwrappedJSON(res.Body)
	if err != nil {
		log.Println("Error parsing Flickr JSON response:", err)
		return "error"
	}

	// Decode JSON body
	decoder := json.NewDecoder(body)
	photos := new(PhotoResponse)
	decoder.Decode(photos)

	imageCount := len(photos.Page.PhotoList)

	// Check if we got any photos back
	if imageCount == 0 {
		return "I found nothing! So sorry."
	} else {
		// Of the photos we got, return one at random
		// so that if bot is asked to run the same query,
		// the user is likely to get a different image.
		r := randNum(imageCount)
		src := photoUrl(photos.Page.PhotoList[r])

		return "<img src='" + src + "'>"
	}

}

// Flickr API responds with wrapped JSON:
// jsonFlickrApi({"photos":{"page":1, "pages":8291105...}, "stat":"ok"})
// This method returns the JSON response without the wrapper
func unwrappedJSON(body io.Reader) (io.Reader, error) {
	bodyBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	stringBody := strings.TrimPrefix(string(bodyBytes), "jsonFlickrApi(")
	stringBody = strings.TrimSuffix(stringBody, ")")

	return strings.NewReader(stringBody), nil

}

// Flickr's photo URL construction is complicated. The JSON we get back from a photo
// search does not include the direct URL. It does include the parameters needed to construct
// one though. So grudgingly, we construct it ourselves -
func photoUrl(photo Photo) string {
	src := "http://farm"
	src += string(photo.Farm)
	src += ".staticflickr.com/"
	src += photo.Server + "/"
	src += photo.Id + "_"
	src += photo.Secret + ".jpg"

	return src
}

// Returns a random number between 0 and max
func randNum(max int) int {
	// Seed with current time
	source := rand.NewSource(time.Now().UnixNano())
	rander := rand.New(source)
	rInt := rander.Int()

	// make the int smaller (between 0 and max) with mod (%)
	smallerRInt := rInt % max
	return smallerRInt
}
