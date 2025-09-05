package yachtlib

type ForecastData struct {
	// Internal
	Cod string `json:"cod"`
	// internal
	Message int `json:"message"`
	// A number of timestamps returned in the API response
	Cnt int `json:"cnt"`
	// List of forecasts
	WeatherTimeStamps []struct {
		// Timestamp
		Dt int `json:"dt"`
		// TempKf - Internal parameter
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Visibility int `json:"visibility"`
		// Probability of precipitation. The values of the parameter vary between 0 and 1, where 0 is equal to 0%, 1 is equal to 100%
		Pop  float64 `json:"pop"`
		Rain struct {
			H3 float64 `json:"3h"`
		} `json:"rain"`
		Snow struct {
			H3 float64 `json:"3h"`
		} `json:"snow"`
		// Part of the day (n - night, d - day)
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		// Time of data forecasted, ISO, UTC - SAME AS DT Timestamp!
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	// Id of city
	// Name - name of City
	// Country code
	// Population. Not working?
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}
