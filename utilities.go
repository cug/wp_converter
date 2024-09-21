// need to read up about package management and then create an actual utilities package
// until then, keeping helpers here
package main

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
