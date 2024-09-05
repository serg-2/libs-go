package main

import (
	"fmt"
	"testing"

	. "github.com/serg-2/libs-go/marinelib"
)

func TestDistance(t *testing.T) {
	// Point 1
	p1 := [2]float64{55, 40}
	p2 := [2]float64{30, 90}
	fmt.Printf("Points: %.6f %.6f AND %.6f %.6f\n", p1[0], p1[1], p2[0], p2[1])
	fmt.Printf("Haversine: %.1f\n", CalculateDistance(p1, p2))
	fmt.Printf("Lambert's: %.1f\n", CalculateDistanceLambert(p1, p2))
	fmt.Printf("==============\n")

	p1 = [2]float64{77.1539, -139.398}
	p2 = [2]float64{-77.1804, -139.55}
	fmt.Printf("Points: %.6f %.6f AND %.6f %.6f\n", p1[0], p1[1], p2[0], p2[1])
	fmt.Printf("Haversine: %.1f\n", CalculateDistance(p1, p2))
	fmt.Printf("Lambert's: %.1f\n", CalculateDistanceLambert(p1, p2))
	fmt.Printf("==============\n")

	p1 = [2]float64{77.1539, 120.398}
	p2 = [2]float64{77.1804, 129.55}
	fmt.Printf("Points: %.6f %.6f AND %.6f %.6f\n", p1[0], p1[1], p2[0], p2[1])
	fmt.Printf("Haversine: %.1f\n", CalculateDistance(p1, p2))
	fmt.Printf("Lambert's: %.1f\n", CalculateDistanceLambert(p1, p2))
	fmt.Printf("==============\n")

	p1 = [2]float64{77.1539, -120.398}
	p2 = [2]float64{77.1804, 129.55}
	fmt.Printf("Points: %.6f %.6f AND %.6f %.6f\n", p1[0], p1[1], p2[0], p2[1])
	fmt.Printf("Haversine: %.1f\n", CalculateDistance(p1, p2))
	fmt.Printf("Lambert's: %.1f\n", CalculateDistanceLambert(p1, p2))
	fmt.Printf("==============\n")

	p1 = [2]float64{44.0001, 55.161}
	p2 = [2]float64{44.0002, 55.162}
	fmt.Printf("Points: %.6f %.6f AND %.6f %.6f\n", p1[0], p1[1], p2[0], p2[1])
	fmt.Printf("Haversine: %.1f\n", CalculateDistance(p1, p2))
	fmt.Printf("Lambert's: %.1f\n", CalculateDistanceLambert(p1, p2))
	fmt.Printf("==============\n")

}
