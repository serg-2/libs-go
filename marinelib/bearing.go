package marinelib

import "math"

func CalculateBearing(a [2]float64, b [2]float64) float64 {
	_lat1 := degToRad(a[0])
	_lat2 := degToRad(b[0])
	_diffLong := degToRad(b[1] - a[1])

	_x := math.Sin(_diffLong) * math.Cos(_lat2)
	_y := math.Cos(_lat1)*math.Sin(_lat2) - (math.Sin(_lat1) * math.Cos(_lat2) * math.Cos(_diffLong))

	_initial_bearing_rad := math.Atan2(_x, _y)
	_initial_bearing_deg := radToDeg(_initial_bearing_rad)
	compass_bearing := math.Mod(_initial_bearing_deg+360, 360)
	return compass_bearing
}
