package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
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

// readCvsData opens a CSV file by name and returns its contents as a slice of string slices.
func readCvsData(filename string) [][]string {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	return parseCvsData(f)
}

// parseCvsData reads CSV data from the provided io.Reader and returns it as a slice of string slices.
func parseCvsData(r io.Reader) [][]string {
	csvReader := csv.NewReader(r)
	csvReader.FieldsPerRecord = -1
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// validateCsvLine checks if a CSV line contains all the required fields and that coordinates are valid floats.
func validateCsvLine(line []string, columnIndexMap map[string]int) bool {
	// Probably neither great nor complete, but it's a start and I can
	// add more validation as problem cases arise
	return validateNotEmptyString(line[columnIndexMap[csvName]]) &&
		validateNotEmptyString(line[columnIndexMap[csvDescription]]) &&
		validateNotEmptyString(line[columnIndexMap[csvLat]]) &&
		validateStringParsesToFloat(line[columnIndexMap[csvLat]]) &&
		validateNotEmptyString(line[columnIndexMap[csvLon]]) &&
		validateStringParsesToFloat(line[columnIndexMap[csvLon]]) &&
		validateNotEmptyString(line[columnIndexMap[csvCategory]])
}

// descriptionFieldsForCategory returns a list of CSV column names that should be included in the
// waypoint description based on the category of the point of interest.
func descriptionFieldsForCategory(category string) []string {
	if isValueInList(category, []string{"Informal Campsite", "Established Campground", "Wild Camping"}) {
		return []string{
			csvDateVerified, csvOpen, csvElectricity, csvWifi, csvKitchen, csvParking,
			csvRestaurant, csvShowers, csvWater, csvToilets, csvBigRig, csvTent, csvPets, csvSani,
		}
	}
	// Default values
	return []string{csvDateVerified, csvOpen}
}

// createDescription constructs a detailed description string for a waypoint, including
// category-specific fields and a link back to the iOverlander website.
func createDescription(line []string, columnIndexMap map[string]int) string {
	var desc string
	desc = line[columnIndexMap[csvDescription]] + "\n\n"

	fieldListForCategory := descriptionFieldsForCategory(line[columnIndexMap[csvCategory]])
	for _, f := range fieldListForCategory {
		if line[columnIndexMap[f]] != "" {
			desc += f + ": " + line[columnIndexMap[f]] + "\n"
		}
	}
	desc += "\n" + baseUrlForDesc + line[columnIndexMap[csvId]] + " (network required)" + "\n"

	return desc
}

// columnHeaderIndexMap creates a mapping from CSV column headers to their respective column indices.
func columnHeaderIndexMap(line []string) map[string]int {
	columnIndexMap := make(map[string]int)
	for i, column := range line {
		columnIndexMap[column] = i
	}
	return columnIndexMap
}
