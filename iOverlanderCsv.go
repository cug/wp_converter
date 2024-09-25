package main

import (
	"encoding/csv"
	"log"
	"os"
)

// Constants for CSV fields, these need to match the column headers
// in the first line of the CSV file, so they can be mapped via
// value for key lookup
const csvId = "Id"
const csvLocation = "Location"
const csvName = "Name"
const csvCategory = "Category"
const csvDescription = "Description"
const csvLat = "Latitude"
const csvLon = "Longitude"
const csvAltitude = "Altitude"
const csvDateVerified = "Date verified"
const csvOpen = "Open"
const csvElectricity = "Electricity"
const csvWifi = "Wifi"
const csvKitchen = "Kitchen"
const csvParking = "Parking"
const csvRestaurant = "Restaurant"
const csvShowers = "Showers"
const csvWater = "Water"
const csvToilets = "Toilets"
const csvBigRig = "Big rig friendly"
const csvTent = "Tent friendly"
const csvPets = "Pet friendly"
const csvSani = "Sanitation dump station"
const csvOutdoorGear = "Outdoor gear"
const csvGroceries = "Groceries"
const csvArtisan = "Artisan goods"
const csvBakery = "Bakery"
const csvRarity = "Rarity in this area"
const csvRepairsVehicle = "Repairs vehicles"
const csvRepairsMotorcycle = "Repairs motorcycles"
const csvRepairsBicycle = "Repairs bicycles"
const csvSellsParts = "Sells parts"
const csvRecyclesBatteries = "Recycles batteries"
const csvRecyclesOil = "Recycles oil"
const csvBioFuel = "Bio fuel"
const csvEvCharging = "Electric vehicle charging"
const csvCompostSawdust = "Composting sawdust"
const csvRecycleCenter = "Recycling center"

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
