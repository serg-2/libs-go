package yachtlib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	. "github.com/serg-2/libs-go/httplib"
	. "github.com/serg-2/libs-go/marinelib"
)

func GetForecast(coord Point, apikey string) ForecastData {
	req, _ := http.NewRequest("GET", "https://api.openweathermap.org/data/2.5/forecast", nil)

	vars := map[string]string{
		"lat":   fmt.Sprintf("%9f", coord.Lat),
		"lon":   fmt.Sprintf("%9f", coord.Long),
		"appid": apikey,
		// Optional
		"units": "metric",
	}
	AddRequestVars(req, vars)

	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:129.0) Gecko/20100101 Firefox/129.0",
		"Accept":          "application/json, text/javascript, */*; q=0.01",
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate, br, zstd",
		"Content-Type":    "application/json",
		"Connection":      "keep-alive",
	}
	AddHeaders(req, headers)

	resp := GetResponse(req, 3, 10, "Forecast", false)
	// upon fail
	if resp == nil {
		log.Println("Empty forecast data")
		return ForecastData{}
	}

	defer resp.Body.Close()

	var doc ForecastData

	err := json.NewDecoder(resp.Body).Decode(&doc)
	if err != nil {
		log.Println("Can't unmarshall forecast JSON. JSON:")
		log.Fatalf("error: %v", err)
	}
	return doc
}
