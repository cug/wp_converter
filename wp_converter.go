package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type OAGpx struct {
	XMLName    xml.Name        `xml:"gpx"`
	Version    string          `xml:"version,attr"`
	Creator    string          `xml:"creator,attr"`
	Namepace   string          `xml:"xmlns:gpxtpx,attr"`
	Xsi        string          `xml:"xmlns:xsi,attr"`
	XsiLocaton string          `xml:"xsi:schemaLocation,attr"`
	Metadata   OAGpxMetadata   `xml:"metadata"`
	Waypoints  []OAWpt         `xml:"wpt"`
	Extensions OAGpxExtensions `xml:"extensions"`
}

type OAGpxMetadata struct {
	XMLName xml.Name `xml:"metadata"`
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
	XMLName      xml.Name  `xml:"extensions"`
	PointsGroups []OAGroup `xml:"osmand:points_groups"`
}

type OAGroup struct {
	XMLName     xml.Name `xml:"group"`
	GIcon       string   `xml:"icon,attr"`
	GBackground string   `xml:"background,attr"`
	GColor      string   `xml:"color,attr"`
	GName       string   `xml:"name,attr"`
}

type OAWptExtensions struct {
	WEIcon           string `xml:"osmand:icon"`
	WEBackground     string `xml:"osmand:background"`
	WEColor          string `xml:"osmand:color"`
	WEAmenitySubtype string `xml:"osmand:amenity_subtype"`
	WEAmenityType    string `xml:"osmand:amenity_type"`
}

func main() {

	// Open our xmlFile
	xmlFile, err := os.Open(os.Args[1])
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
	}
}
