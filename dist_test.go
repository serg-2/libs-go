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
	testOne(p1, p2)

	p1 = [2]float64{77.1539, -139.398}
	p2 = [2]float64{-77.1804, -139.55}
	testOne(p1, p2)

	p1 = [2]float64{77.1539, 120.398}
	p2 = [2]float64{77.1804, 129.55}
	testOne(p1, p2)

	p1 = [2]float64{77.1539, -120.398}
	p2 = [2]float64{77.1804, 129.55}
	testOne(p1, p2)

	p1 = [2]float64{44.0001, 55.161}
	p2 = [2]float64{44.0002, 55.162}
	testOne(p1, p2)

	p1 = [2]float64{14.467433, -60.866204}
	p2 = [2]float64{14.387486, -60.665988}
	testOne(p1, p2)
}

func testOne(p1 [2]float64, p2 [2]float64) {
	fmt.Printf("Points: %.6f lat %.6f lon AND %.6f lat %.6f lon\n", p1[0], p1[1], p2[0], p2[1])
	dist1 := CalculateDistance(p1, p2)
	dist2 := CalculateDistanceLambert(p1, p2)
	fmt.Printf("Haversine: %.1f\n", dist1)
	fmt.Printf("Lambert's: %.1f\n", dist2)
	fmt.Printf("Difference: %.1f%%\n", (dist1/dist2)*100)

	fmt.Printf("==============\n")
}
