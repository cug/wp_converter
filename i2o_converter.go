package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strconv"
)

// convertIOverlanderToOsmAnd reads CSV data from r, processes it into waypoints
// filtered by mapBoundaries, and writes the resulting OsmAnd-compatible GPX
// data to w.
func convertIOverlanderToOsmAnd(r io.Reader, w io.Writer, mapBoundaries map[string]float64) {
	data := parseCvsData(r)
	waypoints, groups := processLines(data, mapBoundaries)

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
	converted := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n")
	converted = append(converted, xmlData...)
	converted = append(converted, "\n"...)

	writeToWriter(w, converted)
}

// processLines iterates through the CSV data, filters waypoints based on the provided
// coordinate boundaries, and groups them by category. It returns a slice of
// waypoints and a slice of the associated category groups.
func processLines(data [][]string, mapBoundaries map[string]float64) ([]OAWpt, []OAGroup) {
	lonMin, lonMax, latMin, latMax := coordinateBoundaries(mapBoundaries)

	categoryMap := make(map[string]OAGroup)
	var waypoints []OAWpt
	var discardedWaypoints []OAWpt

	// Use first line in CSV field to create a map between column names and their index
	columnIndexMap := columnHeaderIndexMap(data[0])

	for i, line := range data {
		if i > 0 && validateCsvLine(line, columnIndexMap) {
			currentLineLon, _ := strconv.ParseFloat(line[columnIndexMap[csvLon]], 64)
			currentLineLat, _ := strconv.ParseFloat(line[columnIndexMap[csvLat]], 64)
			if currentLineLon > lonMin && currentLineLon < lonMax &&
				currentLineLat > latMin && currentLineLat < latMax {
				wp := convertCsvLineToWaypoint(line, columnIndexMap)
				if validateWaypoint(wp, false) {
					waypoints = append(waypoints, wp)
					if categoryMap[wp.WptType].GName == "" {
						categoryMap[wp.WptType] = OAGroup{
							GIcon:       wp.WptExtensions.WEIcon,
							GBackground: wp.WptExtensions.WEBackground,
							GColor:      wp.WptExtensions.WEColor,
							GName:       wp.WptType,
						}
					}
				} else {
					discardedWaypoints = append(discardedWaypoints, wp)
					log.Println("Discarded Waypoint: ", wp)
				}
			}
		}
	}

	if len(discardedWaypoints) > 0 {
		log.Printf("Due to validation errors, %d waypoints were discarded\n", len(discardedWaypoints))
	}

	var groups []OAGroup
	for category := range categoryMap {
		groups = append(groups, categoryMap[category])
	}
	return waypoints, groups
}

// convertCsvLineToWaypoint transforms a single CSV line into an OAWpt struct,
// mapping CSV columns to GPX fields and assigning icons/colors based on the
// waypoint category.
func convertCsvLineToWaypoint(line []string, columnIndexMap map[string]int) OAWpt {
	waypointType := line[columnIndexMap[csvCategory]]
	icon, color, background := iconBackgroundColorForType(waypointType)

	// make places, that aren't open have grey symbols
	if line[columnIndexMap[csvOpen]] != "Yes" {
		fmt.Println("Setting line color to grey")
		color = "#aaaaaa"
	}

	wp := OAWpt{
		WptLat:      line[columnIndexMap[csvLat]],
		WptLon:      line[columnIndexMap[csvLon]],
		WptElevaton: line[columnIndexMap[csvAltitude]],
		WptTime:     line[columnIndexMap[csvDateVerified]],
		WptName:     line[columnIndexMap[csvName]],
		WptDesc:     createDescription(line, columnIndexMap),
		WptType:     waypointType,
		WptExtensions: OAWptExtensions{
			WEIcon:       icon,
			WEBackground: background,
			WEColor:      color,
			// TODO: Figure out whether the below are actually needed for anything
			WEAmenitySubtype: "user_defined_other_postcode",
			WEAmenityType:    "user_defined_other",
		},
	}
	return wp
}
