package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type OAGpx struct {
	XMLName    xml.Name        `xml:"gpx"`
	Version    string          `xml:"version,attr"`
	Creator    string          `xml:"creator,attr"`
	BaseNS     string          `xml:"xmlns,attr"`
	OsmNS      string          `xml:"xmlns:osmand,attr"`
	Namepace   string          `xml:"xmlns:gpxtpx,attr"`
	Xsi        string          `xml:"xmlns:xsi,attr"`
	XsiLocaton string          `xml:"xsi:schemaLocation,attr"`
	Metadata   OAGpxMetadata   `xml:"metadata"`
	Waypoints  []OAWpt         `xml:"wpt"`
	Extensions OAGpxExtensions `xml:"extensions"`
}

type OAGpxMetadata struct {
	XMLName xml.Name `xml:"metadata"`
	Name    string   `xml:"name"`
	GMTime  string   `xml:"time"`
}

type OAWpt struct {
	XMLName       xml.Name        `xml:"wpt"`
	WptLat        string          `xml:"lat,attr"`
	WptLon        string          `xml:"lon,attr"`
	WptElevaton   string          `xml:"ele"`
	WptTime       string          `xml:"time"`
	WptName       string          `xml:"name"`
	WptDesc       string          `xml:"desc"`
	WptType       string          `xml:"type"`
	WptExtensions OAWptExtensions `xml:"extensions"`
}

type OAGpxExtensions struct {
	XMLName      xml.Name       `xml:"extensions"`
	PointsGroups OAPointsGroups `xml:"osmand:points_groups"`
}

type OAPointsGroups struct {
	XMLName xml.Name  `xml:"osmand:points_groups"`
	Group   []OAGroup `xml:"osmand:group"`
}

type OAGroup struct {
	XMLName     xml.Name `xml:"osmand:group"`
	GIcon       string   `xml:"icon,attr"`
	GBackground string   `xml:"background,attr"`
	GColor      string   `xml:"color,attr"`
	GName       string   `xml:"name,attr"`
}

type OAWptExtensions struct {
	XMLName          xml.Name `xml:"extensions"`
	WEIcon           string   `xml:"osmand:icon"`
	WEBackground     string   `xml:"osmand:background"`
	WEColor          string   `xml:"osmand:color"`
	WEAmenitySubtype string   `xml:"osmand:amenity_subtype"`
	WEAmenityType    string   `xml:"osmand:amenity_type"`
}

func main() {
	if os.Args[1] == "read" {
		read()
	} else {
		write()
	}
}

func read() {
	// Open our xmlFile
	xmlFile, err := os.Open(os.Args[2])
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened xml file")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(io.Reader(xmlFile))
	var gpx OAGpx

	uErr := xml.Unmarshal(byteValue, &gpx)
	if uErr != nil {
		fmt.Println(uErr)
		os.Exit(-1)
	}

	for i := 0; i < len(gpx.Waypoints); i++ {
		fmt.Println("Waypoint Name: " + gpx.Waypoints[i].WptName)
		fmt.Println("Waypoint Desc: " + gpx.Waypoints[i].WptDesc)
		fmt.Println("Waypoint Position: " + gpx.Waypoints[i].WptLat + ", " + gpx.Waypoints[i].WptLon)
		fmt.Println("Icon: " + gpx.Waypoints[i].WptExtensions.WEIcon)
		fmt.Println("Background: " + gpx.Waypoints[i].WptExtensions.WEBackground)
	}
}

func write() {

	gpx := OAGpx{
		Version:    "OsmAnd 4.6.6",
		Creator:    "OsmAnd Maps 4.6.6 (4.6.6.1)",
		BaseNS:     "http://www.topografix.com/GPX/1/1",
		OsmNS:      "https://osmand.net",
		Namepace:   "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		Xsi:        "http://www.w3.org/2001/XMLSchema-instance",
		XsiLocaton: "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd",
		Waypoints:  readCsvFile(),
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

func readCsvFile() []OAWpt {
	f, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return convertLinesToWaypoints(data)
}

func convertLinesToWaypoints(data [][]string) []OAWpt {

	// 0: 	Id,
	// 1: 	Location,
	// 2: 	Name,
	// 3: 	Category,
	// 4:	Description,
	// 5:	Latitude,
	// 6:	Longitude,
	// 7:	Altitude,
	// 8:	Date verified,
	// 9:	Open,
	// 10:	Electricity,
	// 11:	Wifi,
	// 12:	Kitchen,
	// 13:	Parking,
	// 14:	Restaurant,
	// 15:	Showers,
	// 16:	Water,
	// 17:	Toilets,
	// 18:	Big rig friendly,
	// 19:	Tent friendly,
	// 20:	Pet friendly,
	// 21:	Sanitation dump station,
	// 22:	Outdoor gear,
	// 23:	Groceries,
	// 24:	Artisan goods,
	// 25:	Bakery,
	// 26:	Rarity in this area,
	// 27:	Repairs vehicles,
	// 28:	Repairs motorcycles,
	// 29:	Repairs bicycles,
	// 30:	Sells parts,
	// 31:	Recycles batteries,
	// 32:	Recycles oil,
	// 33:	Bio fuel,
	// 34:	Electric vehicle charging,
	// 35:	Composting sawdust,
	// 36:	Recycling center

	// 16939,"R. Mte. Cardoso 12-14, 3090, Portugal",Costa De Lavos Service area ,Established Campground,"Service area with toilets, showers, (only during the season) water and place to dispose the gray water. Right on the beach. Nice village to see, Casa dos Pescadores to be visited.",40.087700,-8.872940,17.5283203125,2022-07-17 00:00:00 UTC,Yes,No,No,No,,No,Cold,Non-Potable,Running Water,Yes,No,Yes,Yes,,,,,,,,,,,,,,,

	var waypoints []OAWpt
	for i, line := range data {
		if i > 0 {
			wp := OAWpt{
				WptLat:      line[5],
				WptLon:      line[6],
				WptElevaton: line[7],
				WptTime:     line[8],
				WptName:     line[2],
				WptDesc:     line[4],
				WptType:     line[3],
				WptExtensions: OAWptExtensions{
					WEIcon:           "tourism_camp_site",
					WEBackground:     "circle",
					WEColor:          "#ffff80ff",
					WEAmenitySubtype: "user_defined_other_postcode",
					WEAmenityType:    "user_defined_other",
				},
			}
			waypoints = append(waypoints, wp)
		}

	}
	return waypoints
}
