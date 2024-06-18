package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
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

type OAGpxExtensions struct {
	XMLName      xml.Name       `xml:"extensions"`
	PointsGroups OAPointsGroups `xml:"osmand:points_groups"`
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

type OAWptExtensions struct {
	XMLName          xml.Name `xml:"extensions"`
	WEIcon           string   `xml:"osmand:icon"`
	WEBackground     string   `xml:"osmand:background"`
	WEColor          string   `xml:"osmand:color"`
	WEAmenitySubtype string   `xml:"osmand:amenity_subtype"`
	WEAmenityType    string   `xml:"osmand:amenity_type"`
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

// Constants for CSV fields
const id = "Id"
const location = "Location"
const name = "Name"
const category = "Category"
const description = "Description"
const lat = "Latitude"
const lon = "Longitude"
const altitude = "Altitude"
const dateVerified = "Date verified"
const open = "Open"
const electricity = "Electricity"
const wifi = "WiFi"
const kitchen = "Kitchen"
const parking = "Parking"
const restaurant = "Restaurant"
const showers = "Showers"
const water = "Water"
const wc = "Toilets"
const bigRig = "Big rig friendly"
const tent = "Tent friendly"
const pets = "Pet friendly"
const sani = "Sanitation dump station"
const outdoorGear = "Outdoor gear"
const groceries = "Groceries"
const artisan = "Artisan goods"
const bakery = "Bakery"
const rarity = "Rarity in this area"
const repVehicle = "Repairs vehicle"
const repMoto = "Repairs motorcycles"
const repBicycle = "Repairs bicycles"
const sellParts = "Sells parts"
const recBatt = "Recycles batteries"
const recOil = "Recycles oil"
const bioFuel = "Bio fuel"
const evCharging = "Electric vehicle charging"
const compostSawdust = "Composting sawdust"
const recCenter = "Recycling center"

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
	// TODO: Clean this up into separate functions to open file, read CSV and then
	// decode into useful data

	f, err := os.Open(os.Args[1])
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

func field(s string) int {
	// TODO: Make sure, all exports have same fields
	var fields = map[string]int{
		id:             0,
		location:       1,
		name:           2,
		category:       3,
		description:    4,
		lat:            5,
		lon:            6,
		altitude:       7,
		dateVerified:   8,
		open:           9,
		electricity:    10,
		wifi:           11,
		kitchen:        12,
		parking:        13,
		restaurant:     14,
		showers:        15,
		water:          16,
		wc:             17,
		bigRig:         18,
		tent:           19,
		pets:           20,
		sani:           21,
		outdoorGear:    22,
		groceries:      23,
		artisan:        24,
		bakery:         25,
		rarity:         26,
		repVehicle:     27,
		repMoto:        28,
		repBicycle:     29,
		sellParts:      30,
		recBatt:        31,
		recOil:         32,
		bioFuel:        33,
		evCharging:     34,
		compostSawdust: 35,
		recCenter:      36,
	}

	index, exists := fields[s]
	if exists {
		return index
	} else {
		// TODO: Need error handling for this, otherwise app will just crash
		return -1
	}
}

func convertLinesToWaypoints(data [][]string) []OAWpt {

	// TODO: Make sure to set places that are marked as not open to a different color

	// TODO: Make this xargs
	lonMax := -114.0
	lonMin := -168.0

	var waypoints []OAWpt
	for i, line := range data {
		if i > 0 {
			currentLineLon, _ := strconv.ParseFloat(line[field(lon)], 8)
			if currentLineLon > lonMin && currentLineLon < lonMax {
				wp := OAWpt{
					WptLat:      line[field(lat)],
					WptLon:      line[field(lon)],
					WptElevaton: line[field(altitude)],
					WptTime:     line[field(dateVerified)],
					WptName:     line[field(name)],
					WptDesc:     createDescription(line),
					WptType:     line[field(category)],
					WptExtensions: OAWptExtensions{
						// TODO: Make these dynamic
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

	}
	return waypoints
}

func createDescription(line []string) string {
	// TODO: use category field to determine which fields to combine

	var campsiteFields = []string{
		dateVerified, open, electricity, wifi, kitchen, parking,
		restaurant, showers, water, wc, bigRig, tent, pets, sani,
	}
	var desc string
	desc = line[field(description)] + "\n\n"

	for _, f := range campsiteFields {
		desc += f + ": " + line[field(f)] + "\n"
	}

	return desc
}
