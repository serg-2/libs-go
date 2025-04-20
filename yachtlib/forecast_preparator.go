package yachtlib

import (
	"fmt"
	"time"

	. "github.com/serg-2/libs-go/marinelib"
)

func PrepareForecast(fore ForecastData, location *time.Location) string {
	var currentLocation *time.Location
	if location == nil {
		currentLocation, _ = time.LoadLocation("Etc/UTC")
	} else {
		currentLocation = location
	}

	var result string = ""
	var header string = "Прогноз погоды:\n"
	timeNow := int(time.Now().Unix())
	for index, w := range fore.WeatherTimeStamps {
		// Save time
		if index == 0 {
			header += fmt.Sprintf(
				"Актуальность на %s\n==================\n",
				time.Unix(int64(w.Dt-3600*3), 0).In(currentLocation).Format(TimeForCoordinateLayout))
		}

		if timeNow > w.Dt {
			continue
		}
		result += fmt.Sprintf("In %6s Wind %.0f (%.0f) from %s (%d\u00B0) P: %d\n",
			prepareSecs(w.Dt-timeNow),
			w.Wind.Speed*1.94384449244061,
			w.Wind.Gust*1.94384449244061,
			GetRhumb(float64(w.Wind.Deg)),
			w.Wind.Deg,
			w.Main.Pressure,
		)
	}
	if result == "" {
		return "Прогноз погоды очень старый."
	}
	return header + result
}

func prepareSecs(diff int) string {
	dDays := diff / (3600 * 24)
	dHours := (diff - dDays*(3600*24)) / 3600

	if dDays != 0 {
		return fmt.Sprintf("%dd %dh",
			dDays,
			dHours,
		)
	} else {
		return fmt.Sprintf("%d h",
			dHours,
		)
	}
}
