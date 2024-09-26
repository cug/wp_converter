package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var infile, outfile, mapBoundaries = readArguments()
	if infile == "none" {
		log.Fatal("no input file")
	}
	convertIOverlanderToOsmAnd(infile, outfile, mapBoundaries)
}

// Provide boundary arguments for latitude and longitude like this:
// wp_converter --latMin=50.00 --latMax=51.0 -i infile.csv -o outfile.gpx
// provide the input file, the one downloaded from iOverlander via:
// ... -i infile.csv -o outfile.gpx ...
// or write the output to a file like this:
// ./wp_converter -i infile.csv > outfile.gpx
func readArguments() (string, string, map[string]float64) {
	var infile, outfile string = "none", "none"
	var boundaryArguments = make(map[string]float64)
	validArgumentNames := []string{"lonMin", "lonMax", "latMin", "latMax"}

	for i, a := range os.Args {
		if i > 0 {
			if len(a) > 2 && a[:2] == "--" {
				v := strings.Split(a, "=")
				if len(v) == 2 {
					key := v[0][2:]
					if isValueInList(key, validArgumentNames) {
						value, err := strconv.ParseFloat(v[1], 8)
						checkForError(err)
						boundaryArguments[key] = value
					} else {
						log.Fatal("Invalid boundary argument ", key)
					}
				}
			} else {
				if a == "-i" {
					infile = os.Args[i+1]
					// skip next argument since we are reading it as a filename
					i++
				} else if a == "-o" {
					outfile = os.Args[i+1]
					// skip next argument since we are reading it as a filename
					i++
				} else {
					// ignore
				}
			}
		}
	}
	return infile, outfile, boundaryArguments
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

	validateCoordinateBoundaries(lonMin, lonMax, latMin, latMax)
	return lonMin, lonMax, latMin, latMax
}
