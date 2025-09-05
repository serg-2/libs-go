package yachtlib

import (
	"fmt"
)

// Иконка. 4x - size (10d.png, 10d@2x.png, 10d@4x.png)
func GetPictureUrl(weather WeatherData) string {
	return fmt.Sprintf("https://openweathermap.org/img/wn/%s@4x.png", weather.Weather[0].Icon)
}
