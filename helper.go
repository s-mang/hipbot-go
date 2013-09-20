package main

import (
	"math/rand"
	"strings"
	"time"
)

// Grabs a user's full name from a hipchat.Message.From string
func name(from string) (nick string) {
	names := strings.Split(from, "/")
	return names[1]
}

// Returns a random number between 0 and max
func randNum(max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	rander := rand.New(source)
	rInt := rander.Int()

	smallerRInt := rInt % max
	return smallerRInt
}

// Maps an integer (0 - 6) to an upper-case letter
func alphabet(i int) string {
	alphab := [7]string{"A", "B", "C", "D", "E", "F", "G"}
	return alphab[i]
}

// Stringifies lat & lng and concatenates them together with a comma
func NewLatLngPair(location PlaceLocation) string {
	return (string(location.Lat) + "," + string(location.Lng))
}

// Maps a weather type to its corresponding icon
func weatherIcon(weatherType string) string {
	switch weatherType {
	case "clear-day":
		return "32_cloud_weather.png"
	case "clear-night":
		return "31_cloud_weather.png"
	case "rain":
		return "11_cloud_weather.png"
	case "snow":
		return "41_cloud_weather.png"
	case "sleet":
		return "40_cloud_weather.png"
	case "wind":
		return "24_cloud_weather.png"
	case "fog":
		return "20_cloud_weather.png"
	case "cloudy":
		return "26_cloud_weather.png"
	case "partly-cloudy-day":
		return "30_cloud_weather.png"
	case "partly-cloudy-night":
		return "29_cloud_weather.png"
	default:
		return "na_cloud_weather.png"
	}
}
