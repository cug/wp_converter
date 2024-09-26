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
						panicOnError(err)
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
