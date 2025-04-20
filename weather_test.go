package main

import (
	"fmt"
	"testing"
	"time"

	js "github.com/serg-2/libs-go/jsonlib"
	. "github.com/serg-2/libs-go/marinelib"
	. "github.com/serg-2/libs-go/yachtlib"
)

func TestWeather(t *testing.T) {
	//t.Skip("Skipping")
	// Point 1
	apiKey := ""

	coord := &Point{
		Lat:  55,
		Long: 37,
	}

	// Main Call
	weather := GetWeather(*coord, apiKey)
	fmt.Println(
		js.JsonAsString(weather),
	)
	
	// Preparator
	currentLocation, _ := time.LoadLocation("Europe/Moscow")
	fmt.Println(
		PrepareWeather(weather, currentLocation),
	)

	// Picture Url
	fmt.Println(
		GetPictureUrl(weather),
	)
}
