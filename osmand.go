package main

import (
	"encoding/xml"
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

func iconBackgroundColorForType(t string) (string, string, string) {
	// TODO: Make this configurable and more flexible
	var icon, background, color string
	switch t {
	case "Established Campground":
		icon = "tourism_camp_site"
		background = "circle"
		color = "#339933"
	case "Informal Campsite":
		icon = "tourism_camp_site"
		background = "circle"
		color = "#79d279"
	case "Wild Camping":
		icon = "tourism_camp_site"
		background = "circle"
		color = "#00ff00"
	case "Water":
		icon = "amenity_drinking_water"
		background = "circle"
		color = "#0099ff"
	case "Mechanic and Parts":
		icon = "shop_car_repair"
		background = "circle"
		color = "#9999ff"
	case "Shopping":
		icon = "shop_supermarket"
		background = "circle"
		color = "#339933"
	case "Laundromat":
		icon = "tourism_viewpoint"
		background = "circle"
		color = "#339933"
	case "Fuel Station":
		icon = "fuel"
		background = "circle"
		color = "#339933"
	default:
		// TODO: make this a proper default, e.g. a star or so
		icon = "tourism_viewpoint"
		background = "star"
		color = "#ffff80ff"
	}
	return icon, color, background
}

func validateWaypoint(wp OAWpt) bool {
	// This is likely not complete, but it's a start, better than nothing
	return validateNotEmptyString(wp.WptName) &&
		validateNotEmptyString(wp.WptDesc) &&
		validateNotEmptyString(wp.WptLat) &&
		validateNotEmptyString(wp.WptLon) &&
		validateStringParsesToFloat(wp.WptLon) &&
		validateStringParsesToFloat(wp.WptLat) &&
		validateNotEmptyString(wp.WptExtensions.WEIcon) &&
		validateNotEmptyString(wp.WptExtensions.WEColor)
}
