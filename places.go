package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var googleApiKey = os.Getenv("GOOGLE_API_KEY")

type Place struct {
	Icon      string      `json:"icon"`
	Name      string      `json:"name"`
	OpenHours OpenHours   `json:"opening_hours"`
	Rating    json.Number `json:"rating"`
	Address   string      `json:"vicinity"`
	Geometry  Geometry    `json:"geometry"`
}

type Geometry struct {
	PlaceLocation PlaceLocation `json:"location"`
}

type PlaceLocation struct {
	Lat json.Number `json:"lat"`
	Lng json.Number `json:"lng"`
}

type OpenHours struct {
	OpenNow bool `json:"open_now"`
}

type PlacesResponse struct {
	Places []Place `json:"results"`
}

func places(query string) string {
	additionalParams := "key=" + googleApiKey + "&keyword=" + url.QueryEscape(query)
	fullQueryUrl := QUERY_URL + "&" + additionalParams

	res, err := http.Get(fullQueryUrl)

	if err != nil {
		fmt.Printf("Error occurred in HTTP GET: %s", err)
		return "error"
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	response := new(PlacesResponse)

	decoder.Decode(response)

	return htmlPlaces(response.Places, query)

}

func htmlPlaces(places []Place, query string) string {
	html := "<strong>Results for Nearby " + strings.Title(query) + "</strong><br><ul>"
	markers := ""
	for i := range places {
		if i > 2 {
			break
		}

		html += "<li>" + places[i].Name + "<br>"
		html += places[i].Address + "<br>"
		html += "<em>Rating: " + stringRating(places[i].Rating) + "</em> | "
		html += openNowHtml(places[i].OpenHours.OpenNow) + "<br></li>"

		markers += "&markers=color:blue|label:" + alphabet(i) + "|" + NewLatLngPair(places[i].Geometry.PlaceLocation)
	}

	html += "</ul><br>"

	html += "<img src='" + MAPS_ENDPOINT + markers + "'>"

	return html
}

func stringRating(rating json.Number) string {
	if len(rating) != 0 {
		return string(rating)
	} else {
		return "N/A"
	}
}

func openNowHtml(isOpen bool) string {
	if isOpen {
		return "<strong>Open Now</strong>"
	} else {
		return "<strong>Closed</strong>"
	}
}
