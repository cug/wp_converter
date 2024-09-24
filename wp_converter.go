package main

import (
	"encoding/xml"
	"fmt"
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
	convert(infile, outfile, mapBoundaries)
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

func convert(infile string, outfile string, mapBoundaries map[string]float64) {
	waypoints, groups := convertLines(infile, mapBoundaries)

	gpx := OAGpx{
		Version:    "OsmAnd 4.6.6",
		Creator:    "OsmAnd Maps 4.6.6 (4.6.6.1)",
		BaseNS:     "https://www.topografix.com/GPX/1/1/",
		OsmNS:      "https://osmand.net",
		Namepace:   "https://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		Xsi:        "https://www.w3.org/2001/XMLSchema-instance",
		XsiLocaton: "https://www.topografix.com/GPX/1/1/gpx.xsd",
		Waypoints:  waypoints,
		Metadata: OAGpxMetadata{
			Name:   "favorites",
			GMTime: "1970-01-01T08:00:00Z",
		},
		Extensions: OAGpxExtensions{
			PointsGroups: OAPointsGroups{
				Group: groups,
			},
		},
	}
	xmlData, err := xml.MarshalIndent(gpx, "", " ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Not very elegant, but it works, maybe I'll learn a better way later
	var converted = []byte("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n")
	converted = append(converted, xmlData...)
	converted = append(converted, "\n"...)

	writeToFile(converted, outfile)
}

func convertLines(infile string, mapBoundaries map[string]float64) ([]OAWpt, []OAGroup) {
	// TODO: Make sure to set places that are marked as not open to a different color

	data := readCvsData(infile)
	lonMin, lonMax, latMin, latMax := coordinateBoundaries(mapBoundaries)

	var categoryMap = make(map[string]OAGroup)
	var waypoints []OAWpt
	var discardedWaypoints []OAWpt

	for i, line := range data {
		if i > 0 && validateCsvLine(line) {
			// Currently I'm relying on above line validation, therefore not handling
			// the parse errors
			currentLineLon, _ := strconv.ParseFloat(line[fieldIndexForString(csvLon)], 8)
			currentLineLat, _ := strconv.ParseFloat(line[fieldIndexForString(csvLat)], 8)
			if currentLineLon > lonMin && currentLineLon < lonMax &&
				currentLineLat > latMin && currentLineLat < latMax {
				wp := convertCsvLineToWaypoint(line)
				if validateWaypoint(wp) {
					waypoints = append(waypoints, wp)
					// TODO: Check the map whether the entry already exists

					// Create a map of categories to dynamically fill the OsmAnd
					// Point Group
					categoryMap[wp.WptType] = OAGroup{
						GIcon:       wp.WptExtensions.WEIcon,
						GBackground: wp.WptExtensions.WEBackground,
						GColor:      wp.WptExtensions.WEColor,
						GName:       wp.WptType,
					}
				} else {
					discardedWaypoints = append(discardedWaypoints, wp)
					log.Println("Discarded Waypoint: ", wp)
				}
			}
		}
	}

	log.Printf("Due to validation errors, %d waypoints were discarded\n", len(discardedWaypoints))

	var groups []OAGroup
	for category := range categoryMap {
		groups = append(groups, categoryMap[category])
	}
	return waypoints, groups
}

func convertCsvLineToWaypoint(line []string) OAWpt {
	waypointType := line[fieldIndexForString(csvCategory)]
	icon, color, background := iconBackgroundColorForType(waypointType)
	wp := OAWpt{
		WptLat:      line[fieldIndexForString(csvLat)],
		WptLon:      line[fieldIndexForString(csvLon)],
		WptElevaton: line[fieldIndexForString(csvAltitude)],
		WptTime:     line[fieldIndexForString(csvDateVerified)],
		WptName:     line[fieldIndexForString(csvName)],
		WptDesc:     createDescription(line),
		WptType:     waypointType,
		WptExtensions: OAWptExtensions{
			WEIcon:           icon,
			WEBackground:     background,
			WEColor:          color,
			WEAmenitySubtype: "user_defined_other_postcode",
			WEAmenityType:    "user_defined_other",
		},
	}
	return wp
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

func createDescription(line []string) string {
	// TODO: use category field to determine which fields to combine

	var campsiteFields = []string{
		csvDateVerified, csvOpwn, csvElectricity, csvWifi, csvKitchen, csvParking,
		csvRestaurant, csvShowers, csvWater, csvToilets, csvBigRig, csvTent, csvPets, csvSani,
	}
	var desc string
	// TODO: Handle empty fields properly
	desc = line[fieldIndexForString(csvDescription)] + "\n\n"

	for _, f := range campsiteFields {
		// TODO: Handle empty fields properly
		desc += f + ": " + line[fieldIndexForString(f)] + "\n"
	}

	return desc
}
