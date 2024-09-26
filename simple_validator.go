package main

import (
	"strconv"
)

func validateNotEmptyString(s string) bool {
	return s != ""
}

func validateStringParsesToFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func validateCoordinateBoundaries(lonMin float64, lonMax float64, latMin float64, latMax float64) (bool, string) {
	var message = ""
	if lonMin > lonMax {
		message += "\nlonMin > lonMax"
	}
	if lonMin < -180.0 || lonMin > 180.0 {
		message += "\nlonMin out of bounds"
	}
	if lonMax < -180.0 || lonMax > 180.0 {
		message += "\nlonMax out of bounds"
	}
	if latMin > latMax {
		message += "\nlatMin > latMax"
	}
	if latMin < -90.0 || latMin > 90.0 {
		message += "\nlatMin out of bounds"
	}
	if latMax < -90.0 || latMax > 90.0 {
		message += "\nlatMax out of bounds"
	}
	return message == "", message
}
