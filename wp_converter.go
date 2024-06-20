package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// TODO: read args here and then call the appropriate function or
	// fail the call to the app

	write()
}

func write() {

	// TODO: Make the groups dynamic based on the data in the file

	gpx := OAGpx{
		Version:    "OsmAnd 4.6.6",
		Creator:    "OsmAnd Maps 4.6.6 (4.6.6.1)",
		BaseNS:     "http://www.topografix.com/GPX/1/1",
		OsmNS:      "https://osmand.net",
		Namepace:   "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		Xsi:        "http://www.w3.org/2001/XMLSchema-instance",
		XsiLocaton: "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd",
		Waypoints:  readWaypoints(),
		Metadata: OAGpxMetadata{
			Name:   "favories",
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

	return convertLinesToWaypoints(readCvsData(os.Args[1]))
}

func convertLinesToWaypoints(data [][]string) []OAWpt {

	// TODO: Make sure to set places that are marked as not open to a different color

	// TODO: Make this xargs
	lonMax := -114.0
	lonMin := -168.0

	var waypoints []OAWpt
	for i, line := range data {
		if i > 0 {
			currentLineLon, _ := strconv.ParseFloat(line[fieldIndexForString(csvLon)], 8)
			if currentLineLon > lonMin && currentLineLon < lonMax {
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
						// TODO: Make these dynamic
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
