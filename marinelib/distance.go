package marinelib

import (
	"math"
)

// Point type
type Point struct {
	Lat  float64
	Long float64
}

func CalculateDistanceLambert(a [2]float64, b [2]float64) float64 {
	_diffLat := (math.Pi / 180) * math.Abs(b[0]-a[0])
	_diffLong := (math.Pi / 180) * math.Abs(b[1]-a[1])

	_a1 := math.Pow(math.Sin(_diffLat/2), 2)
	_a2 := math.Cos((math.Pi / 180) * b[0])
	_a3 := math.Cos((math.Pi / 180) * a[0])
	_a4 := math.Pow(math.Sin(_diffLong/2), 2)
	_a := _a1 + _a2*_a3*_a4
	_c := 2 * math.Atan2(math.Sqrt(_a), math.Sqrt(1-_a))

	// _c - central angle in radians

	// f - flattening
	f := float64(1) / float64(298.25642)

	// eqR - equatorial radius of the chosen spheroid
	eqR := float64(6378136.6)

	// Reduced Latitudes
	rLat0 := math.Atan((1 - f) * (math.Pi / 180) * math.Tan(a[0]))
	rLat1 := math.Atan((1 - f) * (math.Pi / 180) * math.Tan(b[0]))

	P := (rLat0 + rLat1) / 2
	Q := (rLat1 - rLat0) / 2
	X := (_c - math.Sin(_c)) * ((math.Sin(P) * math.Sin(P) * math.Cos(Q) * math.Cos(Q)) / (math.Cos(_c/float64(2)) * math.Cos(_c/float64(2))))
	Y := (_c + math.Sin(_c)) * ((math.Sin(Q) * math.Sin(Q) * math.Cos(P) * math.Cos(P)) / (math.Sin(_c/float64(2)) * math.Sin(_c/float64(2))))

	return eqR * (_c - (f/float64(2))*(X+Y))
}

func CalculateDistanceBetweenPointsLambert(a Point, b Point) float64 {
	return CalculateDistanceLambert([2]float64{a.Lat, a.Long}, [2]float64{b.Lat, b.Long})
}

func CalculateDistanceBetweenPoints(a Point, b Point) float64 {
	return CalculateDistance([2]float64{a.Lat, a.Long}, [2]float64{b.Lat, b.Long})
}

func CalculateDistance(a [2]float64, b [2]float64) float64 {
	_diffLong := (math.Pi / 180) * math.Abs(b[1]-a[1])
	_diffLat := (math.Pi / 180) * math.Abs(b[0]-a[0])

	_a1 := math.Pow(math.Sin(_diffLat/2), 2)
	_a2 := math.Cos((math.Pi / 180) * b[0])
	_a3 := math.Cos((math.Pi / 180) * a[0])
	_a4 := math.Pow(math.Sin(_diffLong/2), 2)
	_a := _a1 + _a2*_a3*_a4
	_c := 2 * math.Atan2(math.Sqrt(_a), math.Sqrt(1-_a))

	const R = 6372795
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
