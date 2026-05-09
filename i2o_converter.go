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
		Version:     "OsmAnd 4.6.6",
		Creator:     "OsmAnd Maps 4.6.6 (4.6.6.1)",
		BaseNS:      "https://www.topografix.com/GPX/1/1/",
		OsmNS:       "https://osmand.net",
		Namespace:   "https://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		Xsi:         "https://www.w3.org/2001/XMLSchema-instance",
		XsiLocation: "https://www.topografix.com/GPX/1/1/gpx.xsd",
		Waypoints:   waypoints,
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

	// Prepend XML declaration to the marshaled data
	converted := append([]byte(xml.Header), xmlData...)
	converted = append(converted, '\n')

	if err := writeToWriter(w, converted); err != nil {
		fmt.Printf("Error writing output: %v\n", err)
	}
}

// processLines iterates through the CSV data, filters waypoints based on the provided
// coordinate boundaries, and groups them by category. It returns a slice of
// waypoints and a slice of the associated category groups.
func processLines(data []IOPlace, mapBoundaries map[string]float64) ([]OAWpt, []OAGroup) {
	lonMin, lonMax, latMin, latMax, err := coordinateBoundaries(mapBoundaries)
	if err != nil {
		log.Fatalf("Error validating boundaries: %v", err)
	}

	categoryMap := make(map[string]OAGroup)
	var waypoints []OAWpt
	var discardedWaypoints []OAWpt

	for _, p := range data {
		if validateCsvLine(p) {
			currentLineLon, _ := strconv.ParseFloat(p.Lon, 64)
			currentLineLat, _ := strconv.ParseFloat(p.Lat, 64)
			if currentLineLon > lonMin && currentLineLon < lonMax &&
				currentLineLat > latMin && currentLineLat < latMax {
				wp := convertCsvLineToWaypoint(p)
				if validateWaypoint(wp, false) {
					waypoints = append(waypoints, wp)
					if categoryMap[wp.Type].Name == "" {
						categoryMap[wp.Type] = OAGroup{
							Icon:       wp.Extensions.Icon,
							Background: wp.Extensions.Background,
							Color:      wp.Extensions.Color,
							Name:       wp.Type,
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
func convertCsvLineToWaypoint(p IOPlace) OAWpt {
	waypointType := p.Category
	icon, color, background := iconBackgroundColorForType(waypointType)

	// make places, that aren't open have grey symbols
	if p.Open != "Yes" {
		fmt.Println("Setting line color to grey")
		color = "#aaaaaa"
	}

	wp := OAWpt{
		Lat:         p.Lat,
		Lon:         p.Lon,
		Elevation:   p.Altitude,
		Time:        p.DateVerified,
		Name:        p.Name,
		Description: createDescription(p),
		Type:        waypointType,
		Extensions: OAWptExtensions{
			Icon:       icon,
			Background: background,
			Color:      color,
			// TODO: Figure out whether the below are actually needed for anything
			AmenitySubtype: "user_defined_other_postcode",
			AmenityType:    "user_defined_other",
		},
	}
	return wp
}
