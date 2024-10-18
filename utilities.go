// need to read up about package management and then create an actual utilities package
// until then, keeping helpers here
package main

import (
	"fmt"
	"log"
	"os"
)

func isValueInList(value string, list *[]string) bool {
	for _, v := range *list {
		if v == value {
			return true
		}
	}
	return false
}

func writeToFile(b []byte, filename string) {
	if filename == "none" {
		fmt.Printf("%s", b)
	} else {
		err := os.WriteFile(filename, b, 0644)
		panicOnError(err)
	}
}

func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func coordinateBoundaries(boundaries map[string]float64) (float64, float64, float64, float64) {
	lonMin, lonMax, latMin, latMax := -180.0, 180.0, -90.0, 90.0
	if boundaries["lonMin"] != 0.0 {
		lonMin = boundaries["lonMin"]
	}
	if boundaries["lonMax"] != 0.0 {
		lonMax = boundaries["lonMax"]
	}
	if boundaries["latMin"] != 0.0 {
		latMin = boundaries["latMin"]
	}
	if boundaries["latMax"] != 0.0 {
		latMax = boundaries["latMax"]
	}

	r, message := validateCoordinateBoundaries(lonMin, lonMax, latMin, latMax)
	if !r {
		log.Fatal("Coordinates invalid\n")
		log.Fatal(message)
	}
	return lonMin, lonMax, latMin, latMax
}
