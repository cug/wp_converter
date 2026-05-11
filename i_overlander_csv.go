package main

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
)

// Constants for CSV fields, these need to match the column headers
// in the first line of the CSV file, so they can be mapped via
// value for key lookup
const (
	csvId                = "Id"
	csvLocation          = "Location"
	csvName              = "Name"
	csvCategory          = "Category"
	csvDescription       = "Description"
	csvLat               = "Latitude"
	csvLon               = "Longitude"
	csvAltitude          = "Altitude"
	csvDateVerified      = "Date verified"
	csvOpen              = "Open"
	csvElectricity       = "Electricity"
	csvWifi              = "Wifi"
	csvKitchen           = "Kitchen"
	csvParking           = "Parking"
	csvRestaurant        = "Restaurant"
	csvShowers           = "Showers"
	csvWater             = "Water"
	csvToilets           = "Toilets"
	csvBigRig            = "Big rig friendly"
	csvTent              = "Tent friendly"
	csvPets              = "Pet friendly"
	csvSani              = "Sanitation dump station"
	csvOutdoorGear       = "Outdoor gear"
	csvGroceries         = "Groceries"
	csvArtisan           = "Artisan goods"
	csvBakery            = "Bakery"
	csvRarity            = "Rarity in this area"
	csvRepairsVehicle    = "Repairs vehicles"
	csvRepairsMotorcycle = "Repairs motorcycles"
	csvRepairsBicycle    = "Repairs bicycles"
	csvSellsParts        = "Sells parts"
	csvRecyclesBatteries = "Recycles batteries"
	csvRecyclesOil       = "Recycles oil"
	csvBioFuel           = "Bio fuel"
	csvEvCharging        = "Electric vehicle charging"
	csvCompostSawdust    = "Composting sawdust"
	csvRecycleCenter     = "Recycling center"
)

const baseUrlForDesc string = "https://ioverlander.com/places/"

// IOPlace represents a single waypoint record from iOverlander CSV.
type IOPlace struct {
	ID                string
	Location          string
	Name              string
	Category          string
	Description       string
	Lat               string
	Lon               string
	Altitude          string
	DateVerified      string
	Open              string
	Electricity       string
	Wifi              string
	Kitchen           string
	Parking           string
	Restaurant        string
	Showers           string
	Water             string
	Toilets           string
	BigRig            string
	Tent              string
	Pets              string
	Sani              string
	OutdoorGear       string
	Groceries         string
	Artisan           string
	Bakery            string
	Rarity            string
	RepairsVehicle    string
	RepairsMotorcycle string
	RepairsBicycle    string
	SellsParts        string
	RecyclesBatteries string
	RecyclesOil       string
	BioFuel           string
	EvCharging        string
	CompostSawdust    string
	RecycleCenter     string
}

// parseCvsData reads CSV data from the provided io.Reader and returns it as a slice of IOPlace.
func parseCvsData(r io.Reader) []IOPlace {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(records) < 2 {
		return nil
	}

	headerMap := columnHeaderIndexMap(records[0])
	var places []IOPlace

	for _, line := range records[1:] {
		getVal := func(key string) string {
			if idx, ok := headerMap[key]; ok && idx < len(line) {
				return line[idx]
			}
			return ""
		}

		places = append(places, IOPlace{
			ID:                getVal(csvId),
			Location:          getVal(csvLocation),
			Name:              getVal(csvName),
			Category:          getVal(csvCategory),
			Description:       getVal(csvDescription),
			Lat:               getVal(csvLat),
			Lon:               getVal(csvLon),
			Altitude:          getVal(csvAltitude),
			DateVerified:      getVal(csvDateVerified),
			Open:              getVal(csvOpen),
			Electricity:       getVal(csvElectricity),
			Wifi:              getVal(csvWifi),
			Kitchen:           getVal(csvKitchen),
			Parking:           getVal(csvParking),
			Restaurant:        getVal(csvRestaurant),
			Showers:           getVal(csvShowers),
			Water:             getVal(csvWater),
			Toilets:           getVal(csvToilets),
			BigRig:            getVal(csvBigRig),
			Tent:              getVal(csvTent),
			Pets:              getVal(csvPets),
			Sani:              getVal(csvSani),
			OutdoorGear:       getVal(csvOutdoorGear),
			Groceries:         getVal(csvGroceries),
			Artisan:           getVal(csvArtisan),
			Bakery:            getVal(csvBakery),
			Rarity:            getVal(csvRarity),
			RepairsVehicle:    getVal(csvRepairsVehicle),
			RepairsMotorcycle: getVal(csvRepairsMotorcycle),
			RepairsBicycle:    getVal(csvRepairsBicycle),
			SellsParts:        getVal(csvSellsParts),
			RecyclesBatteries: getVal(csvRecyclesBatteries),
			RecyclesOil:       getVal(csvRecyclesOil),
			BioFuel:           getVal(csvBioFuel),
			EvCharging:        getVal(csvEvCharging),
			CompostSawdust:    getVal(csvCompostSawdust),
			RecycleCenter:     getVal(csvRecycleCenter),
		})
	}

	return places
}

// validateCsvLine checks if an IOPlace contains all the required fields and that coordinates are valid floats.
func validateCsvLine(p IOPlace) bool {
	// Probably neither great nor complete, but it's a start and I can
	// add more validation as problem cases arise
	return validateNotEmptyString(p.Name) &&
		validateNotEmptyString(p.Description) &&
		validateNotEmptyString(p.Lat) &&
		validateStringParsesToFloat(p.Lat) &&
		validateNotEmptyString(p.Lon) &&
		validateStringParsesToFloat(p.Lon) &&
		validateNotEmptyString(p.Category)
}

// createDescription constructs a detailed description string for a waypoint, including
// category-specific fields and a link back to the iOverlander website.
func createDescription(p IOPlace) string {
	var sb strings.Builder
	sb.WriteString(p.Description)
	sb.WriteString("\n\n")

	if isValueInList(p.Category, []string{"Informal Campsite", "Established Campground", "Wild Camping"}) {
		fields := []struct {
			name  string
			value string
		}{
			{csvDateVerified, p.DateVerified}, {csvOpen, p.Open}, {csvElectricity, p.Electricity},
			{csvWifi, p.Wifi}, {csvKitchen, p.Kitchen}, {csvParking, p.Parking},
			{csvRestaurant, p.Restaurant}, {csvShowers, p.Showers}, {csvWater, p.Water},
			{csvToilets, p.Toilets}, {csvBigRig, p.BigRig}, {csvTent, p.Tent},
			{csvPets, p.Pets}, {csvSani, p.Sani},
		}
		for _, f := range fields {
			if f.value != "" {
				sb.WriteString(f.name)
				sb.WriteString(": ")
				sb.WriteString(f.value)
				sb.WriteString("\n")
			}
		}
	} else {
		if p.DateVerified != "" {
			sb.WriteString(csvDateVerified)
			sb.WriteString(": ")
			sb.WriteString(p.DateVerified)
			sb.WriteString("\n")
		}
		if p.Open != "" {
			sb.WriteString(csvOpen)
			sb.WriteString(": ")
			sb.WriteString(p.Open)
			sb.WriteString("\n")
		}
	}
	sb.WriteString("\n")
	sb.WriteString(baseUrlForDesc)
	sb.WriteString(p.ID)
	sb.WriteString(" (network required)")
	sb.WriteString("\n")

	return sb.String()
}

// columnHeaderIndexMap creates a mapping from CSV column headers to their respective column indices.
func columnHeaderIndexMap(line []string) map[string]int {
	columnIndexMap := make(map[string]int)
	for i, column := range line {
		columnIndexMap[column] = i
	}
	return columnIndexMap
}
