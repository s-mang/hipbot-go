package main

import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
)

const (
	WEATHER_ENDPOINT = "https://api.forecast.io/forecast/"
	SPROUTLING_LAT = "37.776266"
	SPROUTLING_LNG = "-122.397550"
	WEATHER_ICON_ENDPOINT = "https://cdn1.iconfinder.com/data/icons/sketchy-weather-icons-by-azuresol/64/"
)

var weatherApiKey = os.Getenv("WEATHER_API_KEY")


type WeatherResults struct {
	Current Current `json:"currently"`
	Day Day `json:"daily"`
}

type Current struct {
	Temperature json.Number `json:"temperature"`
	Icon string `json:"icon"`
}

type Day struct {
	DailyData []*DailyData `json:"data"`
}

type DailyData struct {
	Summary string `json:"summary"`
	PrecipProbability json.Number `json:"precipProbability"`
	TempMin json.Number `json:"temperatureMin"`
	TempMax json.Number `json:"temperatureMax"`
}

func weather(query string) string {
	wResp, wReqErr := http.Get(WEATHER_ENDPOINT + weatherApiKey +"/"+SPROUTLING_LAT+","+SPROUTLING_LNG)
	if wReqErr != nil {
		fmt.Println(wReqErr)
		return "error"
	}
	
	defer wResp.Body.Close()
	
	weatherDecoder := json.NewDecoder(wResp.Body)
	weatherResults := new(WeatherResults)
	
	weatherDecoder.Decode(weatherResults)
	
	fmt.Println(weatherResults)
	
	return formattedWeather(*weatherResults)
	
}

func formattedWeather(weather WeatherResults) string {
	currentTemp := string(weather.Current.Temperature)
	currentIcon := weather.Current.Icon
	daySummary := weather.Day.DailyData[0].Summary
	precipProb := string(weather.Day.DailyData[0].PrecipProbability)
	tempMin := string(weather.Day.DailyData[0].TempMin)
	tempMax := string(weather.Day.DailyData[0].TempMax)

	weatherHtml := "&nbsp;&nbsp;<strong>Today's weather</strong>: "+daySummary+".<br>"
	weatherHtml += "<ul><li>High: "+tempMax+"&deg;, Low: "+tempMin+"&deg;</li>"
	weatherHtml += "<li>Precipitation: "+precipProb+"&#37; chance</li>"
	weatherHtml += "<li>Currently "+currentTemp+"&deg;</li>"
	weatherHtml += "<li><img src='"+WEATHER_ICON_ENDPOINT+weatherIcon(currentIcon)+"'></li></ul>"

	return weatherHtml
}

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


