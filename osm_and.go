package main

import (
	"encoding/xml"
)

type OAGpx struct {
	XMLName     xml.Name        `xml:"gpx"`
	Version     string          `xml:"version,attr"`
	Creator     string          `xml:"creator,attr"`
	BaseNS      string          `xml:"xmlns,attr"`
	OsmNS       string          `xml:"xmlns:osmand,attr"`
	Namespace   string          `xml:"xmlns:gpxtpx,attr"`
	Xsi         string          `xml:"xmlns:xsi,attr"`
	XsiLocation string          `xml:"xsi:schemaLocation,attr"`
	Metadata    OAGpxMetadata   `xml:"metadata"`
	Waypoints   []OAWpt         `xml:"wpt"`
	Extensions  OAGpxExtensions `xml:"extensions"`
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
	XMLName     xml.Name        `xml:"wpt"`
	Lat         string          `xml:"lat,attr"`
	Lon         string          `xml:"lon,attr"`
	Elevation   string          `xml:"ele"`
	Time        string          `xml:"time"`
	Name        string          `xml:"name"`
	Description string          `xml:"desc"`
	Type        string          `xml:"type"`
	Extensions  OAWptExtensions `xml:"extensions"`
}

type OAWptExtensions struct {
	XMLName        xml.Name `xml:"extensions"`
	Icon           string   `xml:"osmand:icon"`
	Background     string   `xml:"osmand:background"`
	Color          string   `xml:"osmand:color"`
	AmenitySubtype string   `xml:"osmand:amenity_subtype"`
	AmenityType    string   `xml:"osmand:amenity_type"`
}

type OAPointsGroups struct {
	XMLName xml.Name  `xml:"osmand:points_groups"`
	Group   []OAGroup `xml:"osmand:group"`
}

type OAGroup struct {
	XMLName    xml.Name `xml:"osmand:group"`
	Icon       string   `xml:"icon,attr"`
	Background string   `xml:"background,attr"`
	Color      string   `xml:"color,attr"`
	Name       string   `xml:"name,attr"`
}

// supportedPOITypes returns a list of POI categories that are officially supported
// by the converter with specific icon and color mappings.
func supportedPOITypes() []string {
	return []string{
		"Established Campground",
		"Informal Campsite",
		"Wild Camping",
		"Water",
		"Mechanic and Parts",
		"Shopping",
		"Laundromat",
		"Fuel Station",
	}
}

type WaypointStyle struct {
	Icon       string
	Background string
	Color      string
}

var styleMap = map[string]WaypointStyle{
	"Established Campground": {"tourism_camp_site", "circle", "#339933"},
	"Informal Campsite":      {"tourism_camp_site", "circle", "#79d279"},
	"Wild Camping":           {"tourism_camp_site", "circle", "#00ff00"},
	"Water":                  {"amenity_drinking_water", "circle", "#0099ff"},
	"Mechanic and Parts":     {"shop_car_repair", "circle", "#9999ff"},
	"Shopping":               {"shop_supermarket", "circle", "#339933"},
	"Laundromat":             {"tourism_viewpoint", "circle", "#339933"},
	"Fuel Station":           {"fuel", "circle", "#339933"},
}

var defaultStyle = WaypointStyle{
	Icon:       "tourism_viewpoint",
	Background: "star",
	Color:      "#ffff80ff",
}

// iconBackgroundColorForType maps a POI category to its corresponding OsmAnd
// icon name, color hex code, and background shape.
func iconBackgroundColorForType(t string) (string, string, string) {
	style, ok := styleMap[t]
	if !ok {
		style = defaultStyle
	}
	return style.Icon, style.Color, style.Background
}

// validateWaypoint checks if a waypoint has all the required fields and valid coordinates.
// If enforceSupportedTypes is true, it also verifies that the waypoint category is in the supported list.
func validateWaypoint(wp OAWpt, enforceSupportedTypes bool) bool {
	// This is likely not complete, but it's a start, better than nothing
	var wpInSupportedTypes bool = true
	if enforceSupportedTypes {
		wpInSupportedTypes = isValueInList(wp.Type, supportedPOITypes())
	}
	return validateNotEmptyString(wp.Name) &&
		validateNotEmptyString(wp.Description) &&
		validateNotEmptyString(wp.Lat) &&
		validateStringParsesToFloat(wp.Lat) &&
		validateNotEmptyString(wp.Lon) &&
		validateStringParsesToFloat(wp.Lon) &&
		validateNotEmptyString(wp.Extensions.Icon) &&
		validateNotEmptyString(wp.Extensions.Color) &&
		wpInSupportedTypes
}
