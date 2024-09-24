package main

import (
	"encoding/csv"
	"log"
	"os"
)

// TODO: Is there a better way to define these?

// Constants for CSV fields
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
const csvWifi = "WiFi"
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
const csvRepairsVehicle = "Repairs vehicle"
const csvRepairsMotorcycle = "Repairs motorcycles"
const csvRepairsBicycle = "Repairs bicycles"
const csvSellsParts = "Sells parts"
const csvRecyclesBatteries = "Recycles batteries"
const csvRecyclesOil = "Recycles oil"
const csvBioFuel = "Bio fuel"
const csvEvCharging = "Electric vehicle charging"
const csvCompostSawdust = "Composting sawdust"
const csvRecycleCenter = "Recycling center"

func fieldIndexForString(s string) int {
	// TODO: Make sure, all exports have same fields
	var fields = map[string]int{
		csvId:                0,
		csvLocation:          1,
		csvName:              2,
		csvCategory:          3,
		csvDescription:       4,
		csvLat:               5,
		csvLon:               6,
		csvAltitude:          7,
		csvDateVerified:      8,
		csvOpen:              9,
		csvElectricity:       10,
		csvWifi:              11,
		csvKitchen:           12,
		csvParking:           13,
		csvRestaurant:        14,
		csvShowers:           15,
		csvWater:             16,
		csvToilets:           17,
		csvBigRig:            18,
		csvTent:              19,
		csvPets:              20,
		csvSani:              21,
		csvOutdoorGear:       22,
		csvGroceries:         23,
		csvArtisan:           24,
		csvBakery:            25,
		csvRarity:            26,
		csvRepairsVehicle:    27,
		csvRepairsMotorcycle: 28,
		csvRepairsBicycle:    29,
		csvSellsParts:        30,
		csvRecyclesBatteries: 31,
		csvRecyclesOil:       32,
		csvBioFuel:           33,
		csvEvCharging:        34,
		csvCompostSawdust:    35,
		csvRecycleCenter:     36,
	}

	index, exists := fields[s]
	if exists {
		return index
	} else {
		// TODO: Need error handling for this, otherwise app will just crash
		return -1
	}
}

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

func validateCsvLine(line []string) bool {
	// Probably neither great nor complete, but it's a start and I can
	// add more validation as problem cases arise
	return validateNotEmptyString(line[fieldIndexForString(csvName)]) &&
		validateNotEmptyString(line[fieldIndexForString(csvDescription)]) &&
		validateNotEmptyString(line[fieldIndexForString(csvLat)]) &&
		validateStringParsesToFloat(line[fieldIndexForString(csvLat)]) &&
		validateNotEmptyString(line[fieldIndexForString(csvLon)]) &&
		validateStringParsesToFloat(line[fieldIndexForString(csvLon)]) &&
		validateNotEmptyString(line[fieldIndexForString(csvCategory)])
}
