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
