// need to read up about package management and then create an actual utilities package
// until then, keeping helpers here
package main

import (
	"io"
	"log"
	"os"
	"slices"
)

// isValueInList checks if a given string exists within a slice of strings.
func isValueInList(value string, list []string) bool {
	return slices.Contains(list, value)
}

// writeToWriter writes a byte slice to the provided io.Writer and panics on error.
func writeToWriter(w io.Writer, b []byte) {
	_, err := w.Write(b)
	panicOnError(err)
}

// writeToFile writes a byte slice to a file. If filename is "none", it writes to standard output.
func writeToFile(b []byte, filename string) {
	if filename == "none" {
		writeToWriter(os.Stdout, b)
	} else {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		panicOnError(err)
		defer f.Close()
		writeToWriter(f, b)
	}
}

// panicOnError panics if the provided error is not nil.
func panicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

// coordinateBoundaries resolves the geographic bounding box from the provided map.
// It uses default world boundaries if specific boundaries are not provided or are 0.0.
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
