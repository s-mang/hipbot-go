package main

const (
	// API endpoints
	QUERY_URL             = "https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" + "--LTLNGPR--" + "&sensor=false&rankby=distance"
	MAPS_ENDPOINT         = "https://maps.googleapis.com/maps/api/staticmap?center=" + "--LTLNGPR--" + "&zoom=15&size=600x200&sensor=false"
	WEATHER_ICON_ENDPOINT = "https://cdn1.iconfinder.com/data/icons/sketchy-weather-icons-by-azuresol/64/"
	WEATHER_ENDPOINT      = "https://api.forecast.io/forecast/"
	WOLFRAM_URL           = "http://api.wolframalpha.com/v2/query"
	THESAURUS_URL         = "http://words.bighugelabs.com/api/2/"
	DUCK_DUCK_GO_ENDPOINT = "https://api.duckduckgo.com/"

	// Misc
	LOGO_URL = "https://1.gravatar.com/avatar/fcc942ea417a208d0f5d835b8427fcc4"
)
