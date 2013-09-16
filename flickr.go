// Not currently in use - photos are not very relevant
package main

import (
	"fmt"
	"os"
	"net/http"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"strings"
	"math/rand"
)

const FLICKR_ENDPOINT = "http://api.flickr.com/services/rest"

var flickrApiKey = os.Getenv("FLICKER_API_KEY")

type PhotoResponse struct {
	Page Page `json:"photos"` 
}

type Page struct {
	PhotoList []Photo `json:"photo"` 
}

type Photo struct {
	Id string `json:"id"`
	Secret string `json:"secret"`
	Server string `json:"server"`
	Farm json.Number `json:"farm"`
}

func flickrSearch(query string) string {
	requestArgs := "?method=flickr.photos.search"
	requestArgs += "&api_key="+flickrApiKey
	requestArgs += "&format=json"
	requestArgs += "&page=1&per_page=1&sort=relevance&is_getty=true&media=photo"
	requestArgs += "&tags="+url.QueryEscape(query)
	
	res, err := http.Get(FLICKR_ENDPOINT + requestArgs)
	
	if err != nil {
		fmt.Printf("Error occurred in HTTP GET: %s", err)
		return "error"
	}
	
	defer res.Body.Close()
	
	bdy, err := ioutil.ReadAll(res.Body)
	
	if err != nil {
		fmt.Println(err)
		return "error"
	}
	
	stringBdy := withoutFlickrWrapper(string(bdy))
	reader := strings.NewReader(stringBdy)	
	decoder := json.NewDecoder(reader)
	response := new(PhotoResponse)
	
	decoder.Decode(response)
	
	r := randNum(len(response.Page.PhotoList))
	src := photoUrl(response.Page.PhotoList[r])
	
	return "<img src='"+src+"'>"
}

func withoutFlickrWrapper(body string) string {
	stringBdy := strings.TrimPrefix(string(body), "jsonFlickrApi(")
	return strings.TrimSuffix(stringBdy, ")")
}

func photoUrl(photo Photo) string {
	src := "http://farm"
	src += string(photo.Farm)
	src += ".staticflickr.com/"
	src += photo.Server + "/"
	src += photo.Id + "_"
	src += photo.Secret + ".jpg"
	
	return src
}

func randNum(max int) int {
	seed := int64(2)
	source := rand.NewSource(seed)
	rander := rand.New(source)
	rInt := rander.Int()
	
	smallerRInt := rInt%max
	return smallerRInt
}


