package marinelib

import "math"

const R = 6373000

func CalculateDistance(a [2]float64, b [2]float64) float64 {
	_diffLong := (math.Pi / 180) * math.Abs(b[1]-a[1])
	_diffLat := (math.Pi / 180) * math.Abs(b[0]-a[0])

	_a1 := math.Pow(math.Sin(_diffLat/2), 2)
	_a2 := math.Cos((math.Pi / 180) * b[0])
	_a3 := math.Cos((math.Pi / 180) * a[0])
	_a4 := math.Pow(math.Sin(_diffLong/2), 2)
	_a := _a1 + _a2*_a3*_a4
	_c := 2 * math.Atan2(math.Sqrt(_a), math.Sqrt(1-_a))
	return R * _c
}

func CalculateBearing(a [2]float64, b [2]float64) float64 {
	_lat1 := (math.Pi / 180) * a[0]
	_lat2 := (math.Pi / 180) * b[0]
	_diffLong := (math.Pi / 180) * (b[1] - a[1])

	_x := math.Sin(_diffLong) * math.Cos(_lat2)
	_y := math.Cos(_lat1)*math.Sin(_lat2) - (math.Sin(_lat1) * math.Cos(_lat2) * math.Cos(_diffLong))

	_initial_bearing := math.Atan2(_x, _y)
	_initial_bearing = _initial_bearing * (180 / math.Pi)
	compass_bearing := math.Mod(_initial_bearing+360, 360)
	return compass_bearing
}
