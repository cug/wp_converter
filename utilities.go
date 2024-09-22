// need to read up about package management and then create an actual utilities package
// until then, keeping helpers here
package main

import (
	"fmt"
	"os"
)

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func writeToFile(b []byte, filename string) {
	if filename == "none" {
		fmt.Printf("%s", b)
	} else {
		err := os.WriteFile(filename, b, 0644)
		checkForError(err)
	}
}

func checkForError(e error) {
	if e != nil {
		panic(e)
	}
}
