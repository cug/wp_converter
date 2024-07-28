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
	write()
}

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
						boundaryArguments[key], _ = strconv.ParseFloat(v[1], 8)
					}
				}
			} else {
				if a == "-i" {
					fileToConvert = os.Args[i+1]
				}
			}
		}
	}
	mapBoundaries = boundaryArguments
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func write() {

	// TODO: Make the groups dynamic based on the data in the file

	gpx := OAGpx{
		Version:    "OsmAnd 4.6.6",
		Creator:    "OsmAnd Maps 4.6.6 (4.6.6.1)",
		BaseNS:     "https://www.topografix.com/GPX/1/1/",
		OsmNS:      "https://osmand.net",
		Namepace:   "https://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		Xsi:        "https://www.w3.org/2001/XMLSchema-instance",
		XsiLocaton: "https://www.topografix.com/GPX/1/1/gpx.xsd",
		Waypoints:  readWaypoints(),
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
	fmt.Printf("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n%s\n", xmlData)
}

func readWaypoints() []OAWpt {
	return convertLinesToWaypoints(readCvsData(fileToConvert))
}

func convertLinesToWaypoints(data [][]string) []OAWpt {
	// TODO: Make sure to set places that are marked as not open to a different color

	lonMin, lonMax, latMin, latMax := coordinateBoundaries()

	var waypoints []OAWpt
	for i, line := range data {
		if i > 0 {
			currentLineLon, _ := strconv.ParseFloat(line[fieldIndexForString(csvLon)], 8)
			currentLineLat, _ := strconv.ParseFloat(line[fieldIndexForString(csvLat)], 8)
			if currentLineLon > lonMin && currentLineLon < lonMax &&
				currentLineLat > latMin && currentLineLat < latMax {
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
				waypoints = append(waypoints, wp)
			}
		}

	}
	return waypoints
}

func coordinateBoundaries() (float64, float64, float64, float64) {
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
