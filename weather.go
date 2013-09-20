package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var weatherApiKey = os.Getenv("WEATHER_API_KEY")

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

func weather(query string) string {
	queryUrl := WEATHER_ENDPOINT + weatherApiKey + "/" + latLngPair
	if query == "tomorrow" {
		tomorrow := time.Now().AddDate(0, 0, 1)
		queryUrl += "," + formattedTime(tomorrow)
	}

	wResp, wReqErr := http.Get(queryUrl)
	if wReqErr != nil {
		fmt.Println(wReqErr)
		return "error"
	}

	defer wResp.Body.Close()

	weatherDecoder := json.NewDecoder(wResp.Body)
	weatherResults := new(WeatherResults)

	weatherDecoder.Decode(weatherResults)

	return formattedWeather(*weatherResults, query)

}

func formattedTime(t time.Time) string {
	// format: "2013-09-15T16:37:00"
	timeParts := strings.Split(t.String(), " ")
	stringDate := timeParts[0]
	stringTime := strings.Split(timeParts[1], ".")[0]

	return (stringDate + "T" + stringTime)
}

func formattedWeather(weather WeatherResults, query string) string {
	currentTemp := string(weather.Current.Temperature)
	currentIcon := weather.Current.Icon
	daySummary := weather.Day.DailyData[0].Summary
	precipProb := string(weather.Day.DailyData[0].PrecipProbability)
	tempMin := string(weather.Day.DailyData[0].TempMin)
	tempMax := string(weather.Day.DailyData[0].TempMax)

	formattedTitle := strings.Title(query) + "'s Weather"

	weatherHtml := "&nbsp;&nbsp;<strong>" + formattedTitle + "</strong>: " + daySummary + ".<br>"
	weatherHtml += "<ul><li>High: " + tempMax + "&deg;, Low: " + tempMin + "&deg;</li>"
	weatherHtml += "<li>Precipitation: " + precipProb + "&#37; chance</li>"
	if query == "today" {
		weatherHtml += "<li>Currently " + currentTemp + "&deg;</li>"
	}
	weatherHtml += "<li><img src='" + WEATHER_ICON_ENDPOINT + weatherIcon(currentIcon) + "'></li></ul>"

	return weatherHtml
}
