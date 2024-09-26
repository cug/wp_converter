package main

import (
	"testing"
)

func TestValidateEmptyString(t *testing.T) {
	s1 := "Hello"
	s2 := ""
	var s3 string

	if !validateNotEmptyString(s1) {
		t.Errorf("Failed for non-emtpy string")
	}
	if validateNotEmptyString(s2) {
		t.Errorf("failed for empty string")
	}
	if validateNotEmptyString(s3) {
		t.Errorf("Failed for non initialized string")
	}
}

func TestValidateStringParsesToFloat(t *testing.T) {
	if !validateStringParsesToFloat("3.14159") {
		t.Errorf("3.14159 wrongly did not pass validation")
	}
	if !validateStringParsesToFloat("0") {
		t.Errorf("0 did not pass validation")
	}
	if !validateStringParsesToFloat("-3.14159") {
		t.Errorf("Negative Pi did not parse to float")
	}
	if validateStringParsesToFloat("foobar") {
		t.Errorf("foobar passed to float")
	}
}

func TestValidateCoordinateBoundaries(t *testing.T) {
	tooSmallLon := -181.0
	tooLargeLon := 181.0
	tooSmallLat := -90.5
	tooLargeLat := 90.5

	minLon := -100.0
	maxLon := 100.0
	minLat := -45.0
	maxLat := 45.0

	var r bool
	var m string

	r, m = validateCoordinateBoundaries(tooSmallLon, maxLon, minLat, maxLat)
	outputTestResult(false, r, m, t)

	r, m = validateCoordinateBoundaries(minLon, tooLargeLon, minLat, maxLat)
	outputTestResult(false, r, m, t)

	r, m = validateCoordinateBoundaries(minLon, maxLon, tooSmallLat, maxLat)
	outputTestResult(false, r, m, t)

	r, m = validateCoordinateBoundaries(minLon, maxLon, minLat, tooLargeLat)
	outputTestResult(false, r, m, t)

	r, m = validateCoordinateBoundaries(maxLon, minLon, maxLat, minLat)
	outputTestResult(false, r, m, t)

	// and finally a good one
	r, m = validateCoordinateBoundaries(minLon, maxLon, minLat, maxLat)
	outputTestResult(true, r, m, t)
}

func outputTestResult(expected bool, r bool, m string, t *testing.T) {
	if r != expected {
		t.Error(m)
		r, m = true, ""
	}
}
