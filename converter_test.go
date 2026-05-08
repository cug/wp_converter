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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lonMin, lonMax, latMin, latMax := coordinateBoundaries(tc.boundaries)
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
		WptName: "Test Place",
		WptDesc: "Test Description",
		WptLat:  "50.0",
		WptLon:  "10.0",
		WptExtensions: OAWptExtensions{
			WEIcon:  "icon",
			WEColor: "color",
		},
	}

	t.Run("Valid waypoint", func(t *testing.T) {
		if !validateWaypoint(validWp, false) {
			t.Error("validateWaypoint should return true for valid waypoint")
		}
	})

	t.Run("Invalid name", func(t *testing.T) {
		wp := validWp
		wp.WptName = ""
		if validateWaypoint(wp, false) {
			t.Error("validateWaypoint should return false for empty name")
		}
	})

	t.Run("Invalid lat", func(t *testing.T) {
		wp := validWp
		wp.WptLat = "not-a-float"
		if validateWaypoint(wp, false) {
			t.Error("validateWaypoint should return false for invalid latitude")
		}
	})

	t.Run("Enforce supported types - valid", func(t *testing.T) {
		wp := validWp
		wp.WptType = "Water"
		if !validateWaypoint(wp, true) {
			t.Error("validateWaypoint should return true for supported type")
		}
	})

	t.Run("Enforce supported types - invalid", func(t *testing.T) {
		wp := validWp
		wp.WptType = "Unknown Type"
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

func TestDescriptionFieldsForCategory(t *testing.T) {
	t.Run("Campground category", func(t *testing.T) {
		fields := descriptionFieldsForCategory("Established Campground")
		if len(fields) < 5 {
			t.Errorf("Expected more fields for campground, got %d", len(fields))
		}
	})

	t.Run("Default category", func(t *testing.T) {
		fields := descriptionFieldsForCategory("Other")
		if len(fields) != 2 {
			t.Errorf("Expected 2 fields for default category, got %d", len(fields))
		}
	})
}

func TestCreateDescription(t *testing.T) {

	m := map[string]int{"Id": 0, "Description": 1, "Category": 2, "Date verified": 3, "Open": 4}
	line := []string{"123", "Nice place", "Other", "2024-01-01", "Yes"}

	desc := createDescription(line, m)
	if !strings.Contains(desc, "Nice place") || !strings.Contains(desc, "123") {
		t.Errorf("createDescription output missing info: %q", desc)
	}
}

func TestValidateCsvLine(t *testing.T) {
	m := map[string]int{
		"Name": 0, "Description": 1, "Latitude": 2, "Longitude": 3, "Category": 4,
	}

	t.Run("Valid line", func(t *testing.T) {
		line := []string{"Place", "Desc", "50.0", "10.0", "Campground"}
		if !validateCsvLine(line, m) {
			t.Error("validateCsvLine should return true for valid line")
		}
	})

	t.Run("Invalid lat", func(t *testing.T) {
		line := []string{"Place", "Desc", "abc", "10.0", "Campground"}
		if validateCsvLine(line, m) {
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
	boundaries := map[string]float64{"latMin": 40, "latMax": 60, "lonMin": -130, "lonMax": -110}

	wpts, groups := processLines(data, boundaries)

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

	if !strings.Contains(result, "osmand:points_groups") {
		t.Errorf("Result missing OsmAnd extensions: %q", result)
	}
}

func TestParseArguments(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expInfile     string
		expOutfile    string
		expBoundaries map[string]float64
	}{
		{
			"No arguments",
			[]string{"prog"},
			"none",
			"none",
			map[string]float64{},
		},
		{
			"Input and output",
			[]string{"prog", "-i", "in.csv", "-o", "out.gpx"},
			"in.csv",
			"out.gpx",
			map[string]float64{},
		},
		{
			"Boundaries only",
			[]string{"prog", "--latMin=50.0", "--latMax=60.0", "--lonMin=-10.0", "--lonMax=10.0"},
			"none",
			"none",
			map[string]float64{"latMin": 50.0, "latMax": 60.0, "lonMin": -10.0, "lonMax": 10.0},
		},
		{
			"Mixed arguments",
			[]string{"prog", "-i", "in.csv", "--latMin=50.0", "-o", "out.gpx", "--lonMax=10.0"},
			"in.csv",
			"out.gpx",
			map[string]float64{"latMin": 50.0, "lonMax": 10.0},
		},
		{
			"Missing values for flags",
			[]string{"prog", "-i", "-o"},
			"none",
			"none",
			map[string]float64{},
		},
		{
			"Incomplete flags",
			[]string{"prog", "-i"},
			"none",
			"none",
			map[string]float64{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			infile, outfile, boundaries := parseArguments(tc.args)
			if infile != tc.expInfile {
				t.Errorf("infile = %q; want %q", infile, tc.expInfile)
			}
			if outfile != tc.expOutfile {
				t.Errorf("outfile = %q; want %q", outfile, tc.expOutfile)
			}
			if len(boundaries) != len(tc.expBoundaries) {
				t.Errorf("boundaries length = %d; want %d", len(boundaries), len(tc.expBoundaries))
			}
			for k, v := range tc.expBoundaries {
				if boundaries[k] != v {
					t.Errorf("boundary %s = %f; want %f", k, boundaries[k], v)
				}
			}
		})
	}
}
