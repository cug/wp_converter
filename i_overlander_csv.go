package main

import (
	"encoding/csv"
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

func readCvsData(filename string) [][]string {
	f, err := os.Open(filename)
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

	return data
}

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

func columnHeaderIndexMap(line []string) map[string]int {
	columnIndexMap := make(map[string]int)
	for i, column := range line {
		columnIndexMap[column] = i
	}
	return columnIndexMap
}
