package main

const (
	// General
	POST_URL = "https://api.hipchat.com/v1/rooms/message"
	POST_COLOR = "gray"
	LAT_LNG = "37.776266,-122.397550" // Sproutling headquarters lat/lng
	
	// API endpoints
	GO_DOC_URL = "http://api.godoc.org/search?q="
	FLICKR_ENDPOINT = "http://api.flickr.com/services/rest"
	NUMBERS_API_URL = "http://numbersapi.com/"
	NYTIMES_QUERY_URL = "http://api.nytimes.com/svc/search/v2/articlesearch.json?sort=newest"
	QUERY_URL = "https://maps.googleapis.com/maps/api/place/nearbysearch/json?location="+LAT_LNG+"&sensor=false&rankby=distance"
	MAPS_ENDPOINT = "https://maps.googleapis.com/maps/api/staticmap?center="+LAT_LNG+"&zoom=15&size=600x200&sensor=false"
	WEATHER_ICON_ENDPOINT = "https://cdn1.iconfinder.com/data/icons/sketchy-weather-icons-by-azuresol/64/"
	WEATHER_ENDPOINT = "https://api.forecast.io/forecast/"
	WOLFRAM_URL = "http://api.wolframalpha.com/v2/query"
	
	// Misc
	LOGO_URL = "https://1.gravatar.com/avatar/fcc942ea417a208d0f5d835b8427fcc4"
)