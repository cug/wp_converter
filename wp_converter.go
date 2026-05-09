package main

import (
	"flag"
	"io"
	"log"
	"os"
)

// main is the entry point of the application. It parses command-line arguments,
// opens the input file and prepares the output destination, then triggers the
// conversion process from iOverlander CSV to OsmAnd GPX.
func main() {
	var infile, outfile string
	var lonMin, lonMax, latMin, latMax float64

	flag.StringVar(&infile, "i", "", "Input CSV file")
	flag.StringVar(&outfile, "o", "", "Output GPX file (default: stdout)")
	flag.Float64Var(&lonMin, "lonMin", -180.0, "Minimum longitude boundary")
	flag.Float64Var(&lonMax, "lonMax", 180.0, "Maximum longitude boundary")
	flag.Float64Var(&latMin, "latMin", -90.0, "Minimum latitude boundary")
	flag.Float64Var(&latMax, "latMax", 90.0, "Maximum latitude boundary")
	flag.Parse()

	if infile == "" {
		log.Fatal("no input file provided. Use -i <filename>")
	}

	fIn, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer fIn.Close()

	var w io.Writer
	if outfile == "" {
		w = os.Stdout
	} else {
		fOut, err := os.Create(outfile)
		if err != nil {
			log.Fatal(err)
		}
		defer fOut.Close()
		w = fOut
	}

	mapBoundaries := map[string]float64{
		"lonMin": lonMin,
		"lonMax": lonMax,
		"latMin": latMin,
		"latMax": latMax,
	}

	// Validate boundaries before starting conversion
	if _, _, _, _, err := coordinateBoundaries(mapBoundaries); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	convertIOverlanderToOsmAnd(fIn, w, mapBoundaries)
}
