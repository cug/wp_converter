package main

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func TestIsValueInList(t *testing.T) {
	list := []string{"apple", "banana", "cherry"}
	tests := []struct {
		value    string
		expected bool
	}{
		{"apple", true},
		{"banana", true},
		{"cherry", true},
		{"date", false},
		{"", false},
	}

	for _, tc := range tests {
		if res := isValueInList(tc.value, list); res != tc.expected {
			t.Errorf("isValueInList(%q, %v) = %v; want %v", tc.value, list, res, tc.expected)
		}
	}
}

func TestCoordinateBoundaries(t *testing.T) {
	tests := []struct {
		name       string
		boundaries map[string]float64
		expected   [4]float64 // lonMin, lonMax, latMin, latMax
	}{
		{
			"Default boundaries",
			map[string]float64{},
			[4]float64{-180.0, 180.0, -90.0, 90.0},
		},
		{
			"Partial boundaries",
			map[string]float64{"latMin": 50.0, "latMax": 60.0},
			[4]float64{-180.0, 180.0, 50.0, 60.0},
		},
		{
			"All boundaries",
			map[string]float64{"lonMin": -10, "lonMax": 10, "latMin": 40, "latMax": 50},
			[4]float64{-10.0, 10.0, 40.0, 50.0},
		},
		{
			"Zero longitude boundaries (prime meridian)",
			map[string]float64{"lonMin": 0.0, "lonMax": 10.0},
			[4]float64{0.0, 10.0, -90.0, 90.0},
		},
		{
			"Zero latitude boundary (equator)",
			map[string]float64{"latMin": 0.0, "latMax": 45.0},
			[4]float64{-180.0, 180.0, 0.0, 45.0},
		},
		{
			"All zero boundaries",
			map[string]float64{"lonMin": 0.0, "lonMax": 0.0, "latMin": 0.0, "latMax": 0.0},
			[4]float64{0.0, 0.0, 0.0, 0.0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lonMin, lonMax, latMin, latMax, err := coordinateBoundaries(tc.boundaries)
			if err != nil {
				t.Fatalf("coordinateBoundaries failed: %v", err)
			}
			if lonMin != tc.expected[0] || lonMax != tc.expected[1] || latMin != tc.expected[2] || latMax != tc.expected[3] {
				t.Errorf("coordinateBoundaries(%v) = %v, %v, %v, %v; want %v", tc.boundaries, lonMin, lonMax, latMin, latMax, tc.expected)
			}
		})
	}
}

func TestSupportedPOITypes(t *testing.T) {
	types := supportedPOITypes()
	if len(types) == 0 {
		t.Error("supportedPOITypes() returned an empty list")
	}
}

func TestIconBackgroundColorForType(t *testing.T) {
	tests := []struct {
		typ      string
		expIcon  string
		expColor string
		expBg    string
	}{
		{"Established Campground", "tourism_camp_site", "#339933", "circle"},
		{"Water", "amenity_drinking_water", "#0099ff", "circle"},
		{"Unknown", "tourism_viewpoint", "#ffff80ff", "star"},
	}

	for _, tc := range tests {
		icon, color, bg := iconBackgroundColorForType(tc.typ)
		if icon != tc.expIcon || color != tc.expColor || bg != tc.expBg {
			t.Errorf("iconBackgroundColorForType(%q) = %q, %q, %q; want %q, %q, %q", tc.typ, icon, color, bg, tc.expIcon, tc.expColor, tc.expBg)
		}
	}
}

func TestValidateWaypoint(t *testing.T) {
	validWp := OAWpt{
		Name:        "Test Place",
		Description: "Test Description",
		Lat:         "50.0",
		Lon:         "10.0",
		Extensions: OAWptExtensions{
			Icon:  "icon",
			Color: "color",
		},
	}

	t.Run("Valid waypoint", func(t *testing.T) {
		if !validateWaypoint(validWp, false) {
			t.Error("validateWaypoint should return true for valid waypoint")
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		wp := validWp
		wp.Name = ""
		if validateWaypoint(wp, false) {
			t.Error("validateWaypoint should return false for empty name")
		}
	})

	t.Run("Invalid lat", func(t *testing.T) {
		wp := validWp
		wp.Lat = "not-a-float"
		if validateWaypoint(wp, false) {
			t.Error("validateWaypoint should return false for invalid latitude")
		}
	})

	t.Run("Enforce supported types - valid", func(t *testing.T) {
		wp := validWp
		wp.Type = "Water"
		if !validateWaypoint(wp, true) {
			t.Error("validateWaypoint should return true for supported type")
		}
	})

	t.Run("Enforce supported types - invalid", func(t *testing.T) {
		wp := validWp
		wp.Type = "Unknown Type"
		if validateWaypoint(wp, true) {
			t.Error("validateWaypoint should return false for unsupported type")
		}
	})
}

func TestColumnHeaderIndexMap(t *testing.T) {
	header := []string{"Id", "Location", "Name"}
	m := columnHeaderIndexMap(header)
	if m["Id"] != 0 || m["Location"] != 1 || m["Name"] != 2 {
		t.Errorf("columnHeaderIndexMap failed: %v", m)
	}
}

func TestCreateDescription(t *testing.T) {
	p := IOPlace{
		ID:           "123",
		Description:  "Nice place",
		Category:     "Other",
		DateVerified: "2024-01-01",
		Open:         "Yes",
	}

	desc := createDescription(p)
	if !strings.Contains(desc, "Nice place") || !strings.Contains(desc, "123") {
		t.Errorf("createDescription output missing info: %q", desc)
	}
}

func TestValidateCsvLine(t *testing.T) {
	t.Run("Valid line", func(t *testing.T) {
		p := IOPlace{
			Name:        "Place",
			Description: "Desc",
			Lat:         "50.0",
			Lon:         "10.0",
			Category:    "Campground",
		}
		if !validateCsvLine(p) {
			t.Error("validateCsvLine should return true for valid line")
		}
	})

	t.Run("Invalid lat", func(t *testing.T) {
		p := IOPlace{
			Name:        "Place",
			Description: "Desc",
			Lat:         "abc",
			Lon:         "10.0",
			Category:    "Campground",
		}
		if validateCsvLine(p) {
			t.Error("validateCsvLine should return false for invalid latitude")
		}
	})
}

func getSampleCSVData() [][]string {
	return [][]string{
		{"Id", "Location", "Name", "Category", "Description", "Latitude", "Longitude", "Altitude", "Date verified", "Open", "Electricity", "Wifi", "Kitchen", "Parking", "Restaurant", "Showers", "Water", "Toilets", "Big rig friendly", "Tent friendly", "Pet friendly", "Sanitation dump station", "Outdoor gear", "Groceries", "Artisan goods", "Bakery", "Rarity in this area", "Repairs vehicles", "Repairs motorcycles", "Repairs bicycles", "Sells parts", "Recycles batteries", "Recycles oil", "Bio fuel", "Electric vehicle charging", "Composting sawdust", "Recycling center"},
		{"1", "Loc1", "Name1", "Established Campground", "Desc1", "50.0", "-120.0", "0", "2024-01-01", "Yes", "Yes", "No", "No", "", "No", "No", "No", "No", "Yes", "Yes", "Yes", "No", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		{"2", "Loc2", "Name2", "Water", "Desc2", "51.0", "-121.0", "0", "2024-01-01", "Yes", "", "", "", "", "", "", "Yes", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	}
}

func TestProcessLines(t *testing.T) {
	data := getSampleCSVData()
	// convert raw CSV sample to IOPlace slice
	var places []IOPlace
	headerMap := columnHeaderIndexMap(data[0])
	for _, line := range data[1:] {
		places = append(places, IOPlace{
			ID:           line[headerMap[csvId]],
			Name:         line[headerMap[csvName]],
			Category:     line[headerMap[csvCategory]],
			Description:  line[headerMap[csvDescription]],
			Lat:          line[headerMap[csvLat]],
			Lon:          line[headerMap[csvLon]],
			Open:         line[headerMap[csvOpen]],
			DateVerified: line[headerMap[csvDateVerified]],
		})
	}
	boundaries := map[string]float64{"latMin": 40, "latMax": 60, "lonMin": -130, "lonMax": -110}

	wpts, groups := processLines(places, boundaries)

	if len(wpts) != 2 {
		t.Errorf("Expected 2 waypoints, got %d", len(wpts))
	}
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
}

func TestConvertIOverlanderToOsmAnd(t *testing.T) {
	data := getSampleCSVData()
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	for _, line := range data {
		writer.Write(line)
	}
	writer.Flush()

	r := strings.NewReader(buf.String())

	w := &bytes.Buffer{}
	boundaries := map[string]float64{"latMin": 40, "latMax": 60, "lonMin": -130, "lonMax": -110}

	convertIOverlanderToOsmAnd(r, w, boundaries)

	result := w.String()
	if !strings.Contains(result, "Name1") || !strings.Contains(result, "Name2") {
		t.Errorf("Result missing waypoint names: %q", result)
	}
	if !strings.Contains(result, "osmand:points_groups") {
		t.Errorf("Result missing OsmAnd extensions: %q", result)
	}
}

// TestParseArguments is removed as parseArguments was replaced by the flag package.
