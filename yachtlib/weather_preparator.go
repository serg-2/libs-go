package yachtlib

import (
	"fmt"
	"time"
)

func PrepareWeather(w WeatherData, location *time.Location) string {
	var currentLocation *time.Location
	if location == nil {
		currentLocation, _ = time.LoadLocation("Etc/UTC")
	} else {
		currentLocation = location
	}

	// Прогноз на какое время. Time zone not using +w.Timezone
	tPrognoz := time.Unix(int64(w.Dt), 0)

	var res string = fmt.Sprintf("Текущая погода (%s)\n====================\n", tPrognoz.In(currentLocation).Format(TimeForCoordinateLayout))

	// Описание
	if len(w.Weather) > 0 {
		res += fmt.Sprintf("Описание: %s\n", w.Weather[0].Description)
	}
	res += fmt.Sprintf("Температура %.1f, чувствуется, как %.1f\n", w.Main.Temp, w.Main.FeelsLike)
	res += fmt.Sprintf("Давление %d гПа, Влажность %d %%\n", w.Main.Pressure, w.Main.Humidity)
	// Видимость
	if w.Visibility == 10000 {
		res += "Видимость отличная.\n"
	} else {
		res += fmt.Sprintf("Видимость %d метров\n", w.Visibility)
	}
	// Ветер
	if w.Wind.Gust == 0 {
		res += fmt.Sprintf("Ветер %.1f ровный, направление %d\n", w.Wind.Speed*1.94384449244061, w.Wind.Deg)
	} else {
		res += fmt.Sprintf("Ветер %.1f, порывы %.1f, направление %d\n", w.Wind.Speed*1.94384449244061, w.Wind.Gust*1.94384449244061, w.Wind.Deg)
	}

	// Дождь
	if w.Rain.H1 != 0 {
		res += fmt.Sprintf("Осадки за последний час: %.1f\n", w.Rain.H1)
	}
	if w.Rain.H3 != 0 {
		res += fmt.Sprintf("Осадки за 3 часа: %.1f\n", w.Rain.H3)
	}
	// Снег
	if w.Snow.H1 != 0 {
		res += fmt.Sprintf("Осадки за последний час: %.1f\n", w.Snow.H1)
	}
	if w.Snow.H3 != 0 {
		res += fmt.Sprintf("Осадки за 3 часа: %.1f\n", w.Snow.H3)
	}
	// Облачность
	if w.Clouds.All != 0 {
		res += fmt.Sprintf("Облачность %d %%\n", w.Clouds.All)
	}
	// Timezone not using +w.Timezone
	tSunrise := time.Unix(int64(w.Sys.Sunrise), 0)
	tSunset := time.Unix(int64(w.Sys.Sunset), 0)
	//debug
	// log.Printf("Время восхода полученное: %d изменение %d", w.Sys.Sunrise, w.Timezone)

	// Восход
	res += fmt.Sprintf("Восход: %s\n", tSunrise.In(currentLocation).Format(TimeForCoordinateLayout))
	// Закат
	res += fmt.Sprintf("Закат: %s\n", tSunset.In(currentLocation).Format(TimeForCoordinateLayout))

	return res
}
