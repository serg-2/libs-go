package yachtlib

type WeatherData struct {
	// Coordinates
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	// ID - weather condition Id
	// Main - Group of weather parameters (Rain, Snow, Clouds etc.)
	// Description - description
	// Icon - icon id
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	// Base - Internal parameter
	Base string `json:"base"`
	// Temp      - temperature
	// FeelsLike - temperature feels like
	// TempMin   - Minimum temperature at the moment. This is minimal currently observed temperature (within large megalopolises and urban areas).
	// TempMax   - Maximum temperature at the moment. This is maximal currently observed temperature (within large megalopolises and urban areas).
	// Pressure  - Atmospheric pressure on the sea level. hPa.
	// Humidity  - Humidity. %
	// SeaLevel  - Atmospheric pressure on the sea level, hPa
	// GrndLevel - Atmospheric pressure on the ground level, hPa
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	// Visibility, meter. The maximum value of the visibility is 10000 m
	Visibility int `json:"visibility"`
	// Speed - m/s
	// Deg - Wind direction, degrees (meteorological)
	// Gust - Wind gust. m/s
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	// OPTIONAL!!!
	// H1 - Rain volume for the last 1 hour, mm. 
	// H3 - Rain volume for the last 3 hour, mm. 
	Rain struct {
		H1 float64 `json:"1h"`
		H3 float64 `json:"3h"`
	} `json:"rain"`
	// OPTIONAL!!!
	// H1 - Snow volume for the last 1 hour, mm. 
	// H3 - Snow volume for the last 3 hour, mm. 
	Snow struct {
		H1 float64 `json:"1h"`
		H3 float64 `json:"3h"`
	} `json:"snow"`
	// All - Cloudiness, %
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	// Time of data calculation, unix, UTC
	Dt int `json:"dt"`
	// Country - Country code (GB, JP etc.) 
	// Sunrise - Sunrise time, unix, UTC
	// Sunset -  Sunrise time, unix, UTC
	Sys struct {
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	// Timezone - Shift in seconds from UTC
	Timezone int    `json:"timezone"`
	// City ID. Deprecated
	ID       int    `json:"id"`
	// City Name. Deprecated
	Name     string `json:"name"`
	// Internal parameter
	Cod      int    `json:"cod"`
}
