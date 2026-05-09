// need to read up about package management and then create an actual utilities package
// until then, keeping helpers here
package main

import (
	"fmt"
	"io"
	"os"
	"slices"
)

// isValueInList checks if a given string exists within a slice of strings.
func isValueInList(value string, list []string) bool {
	return slices.Contains(list, value)
}

// writeToWriter writes a byte slice to the provided io.Writer.
func writeToWriter(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}

// writeToFile writes a byte slice to a file. If filename is empty, it writes to standard output.
func writeToFile(b []byte, filename string) error {
	if filename == "" {
		return writeToWriter(os.Stdout, b)
	} else {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		return writeToWriter(f, b)
	}
}

// Removed panicOnError in favor of returning errors to the caller.

// coordinateBoundaries resolves the geographic bounding box from the provided map.
// It uses default world boundaries if specific boundaries are not provided or are 0.0.
func coordinateBoundaries(boundaries map[string]float64) (float64, float64, float64, float64, error) {
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
		return 0, 0, 0, 0, fmt.Errorf("coordinates invalid: %s", message)
	}
	return lonMin, lonMax, latMin, latMax, nil
}
