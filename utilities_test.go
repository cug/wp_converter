package main

import (
	"testing"
)

func TestIsValueInList(t *testing.T) {
	list := &[]string{"a", "b", "c"}
	vInlist := "a"
	vNotInList := "d"
	if !isValueInList(vInlist, list) {
		t.Errorf("failed for value that is in list")
	}
	if isValueInList(vNotInList, list) {
		t.Errorf("returned true for value that isn't in list")
	}
}
