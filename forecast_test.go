package main

import (
	"fmt"
	"testing"
	"time"

	js "github.com/serg-2/libs-go/jsonlib"
	. "github.com/serg-2/libs-go/marinelib"
	. "github.com/serg-2/libs-go/yachtlib"
)

func TestForecast(t *testing.T) {
	// t.Skip("Skipping")
	// Point 1
	apiKey := ""

	coord := &Point{
		Lat:  40,
		Long: 14,
	}

	// Main Call
	forecast := GetForecast(*coord, apiKey)
	fmt.Println(
		js.JsonAsString(forecast),
	)
	// Preparator
	currentLocation, _ := time.LoadLocation("Europe/Moscow")

	fmt.Println(
		PrepareForecast(forecast, currentLocation),
	)

}
