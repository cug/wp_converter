package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// main is the entry point of the application. It parses command-line arguments,
// opens the input file and prepares the output destination, then triggers the
// conversion process from iOverlander CSV to OsmAnd GPX.
func main() {
	infile, outfile, mapBoundaries := readArguments()
	if infile == "none" {
		log.Fatal("no input file")
	}

	fIn, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer fIn.Close()

	var w io.Writer
	if outfile == "none" {
		w = os.Stdout
	} else {
		fOut, err := os.Create(outfile)
		if err != nil {
			log.Fatal(err)
		}
		defer fOut.Close()
		w = fOut
	}

	convertIOverlanderToOsmAnd(fIn, w, mapBoundaries)
}

// Provide boundary arguments for latitude and longitude like this:
// wp_converter --latMin=50.00 --latMax=51.0 -i infile.csv -o outfile.gpx
// provide the input file, the one downloaded from iOverlander via:
// ... -i infile.csv -o outfile.gpx ...
// or write the output to a file like this:
// ./wp_converter -i infile.csv > outfile.gpx
// readArguments reads the command-line arguments from os.Args and parses them.
// It returns the input filename, output filename, and a map of geographic boundaries.
func readArguments() (string, string, map[string]float64) {
	return parseArguments(os.Args)
}

// parseArguments takes a slice of strings (typically os.Args) and extracts the
// input file (-i), output file (-o), and boundary constraints (--latMin, --latMax, --lonMin, --lonMax).
// It returns "none" for filenames if not provided.
func parseArguments(args []string) (string, string, map[string]float64) {
	infile, outfile := "none", "none"
	boundaryArguments := make(map[string]float64)
	validArgumentNames := []string{"lonMin", "lonMax", "latMin", "latMax"}

	for i := 1; i < len(args); i++ {
		a := args[i]
		if len(a) > 2 && a[:2] == "--" {
			v := strings.Split(a, "=")
			if len(v) == 2 {
				key := v[0][2:]
				if isValueInList(key, validArgumentNames) {
					value, err := strconv.ParseFloat(v[1], 64)
					panicOnError(err)
					boundaryArguments[key] = value
				} else {
					log.Fatal("Invalid boundary argument ", key)
				}
			}
		} else {
			if a == "-i" && i+1 < len(args) && args[i+1][0] != '-' {
				infile = args[i+1]
				i++
			} else if a == "-o" && i+1 < len(args) && args[i+1][0] != '-' {
				outfile = args[i+1]
				i++
			}
		}
	}
	return infile, outfile, boundaryArguments
}
