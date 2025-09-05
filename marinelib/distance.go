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
	// Reduced Latitudes
	rLatARad := getReducedLat(degToRad(a[0]))
	rLatBRad := getReducedLat(degToRad(b[0]))

	_c := getCentralAngleRad(
		rLatARad,
		degToRad(a[1]),
		rLatBRad,
		degToRad(b[1]),
	)
	
	P := (rLatARad + rLatBRad) / 2
	Q := (rLatBRad - rLatARad) / 2
	X := (_c - math.Sin(_c)) * ((math.Sin(P) * math.Sin(P) * math.Cos(Q) * math.Cos(Q)) / (math.Cos(_c/float64(2)) * math.Cos(_c/float64(2))))
	Y := (_c + math.Sin(_c)) * ((math.Sin(Q) * math.Sin(Q) * math.Cos(P) * math.Cos(P)) / (math.Sin(_c/float64(2)) * math.Sin(_c/float64(2))))

	return eqR * (_c - (f/2)*(X+Y))
}

// point [lat lon]
func CalculateDistance(point1 [2]float64, point2 [2]float64) float64 {
	_c := getCentralAngleRad(
		degToRad(point1[0]),
		degToRad(point1[1]),
		degToRad(point2[0]),
		degToRad(point2[1]),
	)
	return R * _c
}

func CalculateDistanceBetweenPointsLambert(a Point, b Point) float64 {
	return CalculateDistanceLambert([2]float64{a.Lat, a.Long}, [2]float64{b.Lat, b.Long})
}

func CalculateDistanceBetweenPoints(a Point, b Point) float64 {
	return CalculateDistance([2]float64{a.Lat, a.Long}, [2]float64{b.Lat, b.Long})
}
