package main

import "strconv"

func validateNotEmptyString(s string) bool {
	return s != ""
}

func validateStringParsesToFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
