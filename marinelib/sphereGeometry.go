package marinelib

import "math"

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func radToDeg(rad float64) float64 {
	return rad * (180 / math.Pi)
}

func getReducedLat(rad float64) float64 {
	return math.Atan((1 - f) * math.Tan(rad))
}

func getCentralAngleRad(lat1Rad float64, lon1Rad float64, lat2Rad float64, lon2Rad float64) float64 {
	_diffLong := math.Abs(lon2Rad - lon1Rad)
	_diffLat := math.Abs(lat2Rad - lat1Rad)
	_a1 := math.Pow(math.Sin(_diffLat/2), 2)
	_a2 := math.Cos(lat2Rad)
	_a3 := math.Cos(lat1Rad)
	_a4 := math.Pow(math.Sin(_diffLong/2), 2)
	_a := _a1 + _a2*_a3*_a4
	return 2 * math.Atan2(math.Sqrt(_a), math.Sqrt(1-_a))
}
