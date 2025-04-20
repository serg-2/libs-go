package yachtlib

import (
	"fmt"
)

func GetPictureUrl(weather WeatherData) string {
	return fmt.Sprintf("https://openweathermap.org/img/wn/%s@4x.png", weather.Weather[0].Icon)
}
