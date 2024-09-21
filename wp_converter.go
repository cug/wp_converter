package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var mapBoundaries map[string]float64
var fileToConvert = "none"

func main() {
	readArguments()
	if fileToConvert == "none" {
		log.Fatal("no input file")
	}
	convert()
}

// Provide boundary arguments for latitude and longitude like this:
// wp_converter --latMin=50.00 -i infile.csv > outfile.gpx
// provide the input file, the one downloaded from iOverlander via:
// ... -i infile.csv ...
// write the output to a file like this:
// ./wp_converter -i infile.csv > outfile.gpx
func readArguments() {
	var boundaryArguments = make(map[string]float64)
	validArgumentNames := []string{"lonMin", "lonMax", "latMin", "latMax"}

	for i, a := range os.Args {
		if i > 0 {
			if len(a) > 2 && a[:2] == "--" {
				v := strings.Split(a, "=")
				if len(v) == 2 {
					key := v[0][2:]
					if isValueInList(key, validArgumentNames) {
						// TODO Handle syntax error from parsing float
						boundaryArguments[key], _ = strconv.ParseFloat(v[1], 8)
					}
				}
			} else {
				if a == "-i" {
					fileToConvert = os.Args[i+1]
				}
			}
			// TODO: Handle invalid arguments
			// TODO: Add argument to specify outfile
		}
	}
	mapBoundaries = boundaryArguments
}

func convert() {

	// TODO: Make the groups dynamic based on the data in the file

	lineData := convertLinesToWaypoints()

	gpx := OAGpx{
		Version:    "OsmAnd 4.6.6",
		Creator:    "OsmAnd Maps 4.6.6 (4.6.6.1)",
		BaseNS:     "https://www.topografix.com/GPX/1/1/",
		OsmNS:      "https://osmand.net",
		Namepace:   "https://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		Xsi:        "https://www.w3.org/2001/XMLSchema-instance",
		XsiLocaton: "https://www.topografix.com/GPX/1/1/gpx.xsd",
		Waypoints:  lineData,
		Metadata: OAGpxMetadata{
			Name:   "favorites",
			GMTime: "1970-01-01T08:00:00Z",
		},
		Extensions: OAGpxExtensions{
			PointsGroups: OAPointsGroups{
				Group: []OAGroup{
					{
						GIcon:       "tourism_camp_site",
						GBackground: "circle",
						GColor:      "#ffff0000",
						GName:       "Informal Campsite",
					},
					{
						GIcon:       "tourism_camp_site",
						GBackground: "circle",
						GColor:      "#ffff0000",
						GName:       "",
					},
					{
						GIcon:       "tourism_camp_site",
						GBackground: "circle",
						GColor:      "#ffff0000",
						GName:       "Established Campground",
					},
				},
			},
		},
	}
	xmlData, err := xml.MarshalIndent(gpx, "", " ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// TODO: Write to outfile if specified as an argument
	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n%s\n", xmlData)
}

func convertLinesToWaypoints() []OAWpt {
	// TODO: Make sure to set places that are marked as not open to a different color

	data := readCvsData(fileToConvert)
	lonMin, lonMax, latMin, latMax := coordinateBoundaries()

	var waypoints []OAWpt
	for i, line := range data {
		// TODO: Run line through validation
		if i > 0 {
			// TODO: Handle errors properly
			currentLineLon, _ := strconv.ParseFloat(line[fieldIndexForString(csvLon)], 8)
			currentLineLat, _ := strconv.ParseFloat(line[fieldIndexForString(csvLat)], 8)
			if currentLineLon > lonMin && currentLineLon < lonMax &&
				// TODO: Make sure the values make sense
				currentLineLat > latMin && currentLineLat < latMax {
				// TODO: Make sure the category is in a supported list
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
				// TODO: Discard waypoint if it failed validation
				waypoints = append(waypoints, wp)
			}
		}

	}
	return waypoints
}

func coordinateBoundaries() (float64, float64, float64, float64) {
	// TODO: Move mapBoundaries to an argument, get rid of global
	lonMin, lonMax, latMin, latMax := -180.0, 180.0, -90.0, 90.0
	if mapBoundaries["lonMin"] != 0.0 {
		lonMin = mapBoundaries["lonMin"]
	}
	if mapBoundaries["lonMax"] != 0.0 {
		lonMax = mapBoundaries["lonMax"]
	}
	if mapBoundaries["latMin"] != 0.0 {
		latMin = mapBoundaries["latMin"]
	}
	if mapBoundaries["latMax"] != 0.0 {
		latMax = mapBoundaries["latMax"]
	}

	return lonMin, lonMax, latMin, latMax
}
