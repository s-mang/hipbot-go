package main

// @botling weather me today
// @botling weather me tomorrow
// Get all weather info for either today or tomorrow
// return an HTML list of weather attributes along with a pretty icon representation
// Many thanks to azuresol for the AWESOME & FREE ICONS! (http://azuresol.deviantart.com/)

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	WEATHER_ICON_ENDPOINT     = "https://cdn1.iconfinder.com/data/icons/sketchy-weather-icons-by-azuresol/64/"
	WEATHER_FORECAST_ENDPOINT = "https://api.forecast.io/forecast/"
)

var weatherApiKey = os.Getenv("FORECAST_IO_API_KEY")

type WeatherResults struct {
	Current Current `json:"currently"`
	Day     Day     `json:"daily"`
}

type Current struct {
	Temperature json.Number `json:"temperature"`
	Icon        string      `json:"icon"`
}

type Day struct {
	DailyData []*DailyData `json:"data"`
}

type DailyData struct {
	Summary           string      `json:"summary"`
	PrecipProbability json.Number `json:"precipProbability"`
	TempMin           json.Number `json:"temperatureMin"`
	TempMax           json.Number `json:"temperatureMax"`
}

// Get the weather forecast for <query> (today or tomorrow)
func weather(query string) string {
	// Set query URL - uses latLngPair for location
	queryUrl := WEATHER_FORECAST_ENDPOINT + weatherApiKey + "/" + latLngPair

	// By default, the API will retrieve today's weather.
	// To change that (give a date), we append a formatted date to the query URL
	if query == "tomorrow" {
		tomorrow := time.Now().AddDate(0, 0, 1)
		queryUrl += "," + formattedTime(tomorrow)
	}

	// Send GET request, collect response
	res, err := http.Get(queryUrl)
	if err != nil {
		log.Println(err)
		return "error"
	}

	defer res.Body.Close()

	// Decode JSON response
	decoder := json.NewDecoder(res.Body)
	results := new(WeatherResults)
	decoder.Decode(results)

	// Return a nicely formatted HTML version of results, complete with awesome icon
	return formattedWeather(*weatherResults, query)

}

// Format the time.Time instance so that forecast.io API can read it
// return something like: "2013-09-15T16:37:00"
func formattedTime(t time.Time) string {
	// t.String() will look something like c 21:59:55.459808262 -0700 PDT
	// Split time string on spaces (splits between: date/time, time/zone_diff, zone_diff/zone)
	timeParts := strings.Split(t.String(), " ")
	// Get date in format 2009-11-10
	stringDate := timeParts[0]
	// Get time without part-seconds (21:59:55)
	stringTime := strings.Split(timeParts[1], ".")[0]

	// Final output will look like 2009-11-10T21:59:55
	return (stringDate + "T" + stringTime)
}

// Format the weather forecast struct as pretty HTML for the Hipchat POST
func formattedWeather(weather WeatherResults, query string) string {
	dailyData := weather.Day.DailyData[0]
	formattedTitle := strings.Title(query) + "'s Weather"

	// Slightly-indented title (ex: Today's Weather: partly cloudy.)
	weatherHtml := "&nbsp;&nbsp;<strong>" + formattedTitle + "</strong>: " + dailyData.Summary + ".<br>"

	// High & Low temperatures for <query> in Farenheit
	weatherHtml += "<ul><li>High: " + dailyData.TempMax + "&deg;, Low: " + dailyData.TempMin + "&deg;</li>"

	// Percent chance of precipitation (ex: Precipitation: 20% chance)
	weatherHtml += "<li>Precipitation: " + dailyData.PrecipProbability + "&#37; chance</li>"

	// Show current weather if <query> == "today"
	if query == "today" {
		weatherHtml += "<li>Currently " + string(weather.Current.Temperature) + "&deg;</li>"
	}

	// Map the response weather icon string to an actual icon
	// icon string ex: "clear-night"
	weatherHtml += "<li><img src='" + WEATHER_ICON_ENDPOINT + weatherIcon(weather.Current.Icon) + "'></li></ul>"

	return weatherHtml
}

// Maps a weather type to its corresponding icon path
// Note that these icons are free for non-commercial use.
// Info on the author at http://azuresol.deviantart.com/
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
