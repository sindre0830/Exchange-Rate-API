package fun_test

import (
	"main/fun"
	"testing"
)

// Test_ValidateDates test ValidateDates with multiple data.
func Test_ValidateDates(t *testing.T) {
	//store expected data to check against
	data := map[string]bool {
		"2020-01-01-2020-02-01": true,
		"2020-01-01-2020-02": false,
		"01-01-2020-2020-02-01": false,
		"2020-01-01-01-02-2020": false,
		"2020-01-01-2O20-02-01": false,
		"2020-01-01-2020-02-00": false,
		"2020-02-01-2020-01-01": false,
	}
	//iterate through map and check each key to expected element
	for key, element := range data {
		err := fun.ValidateDates(key)
		//branch if we get an unexpected answer
		if (err != nil && element) || (err == nil && !element) {
			t.Errorf("Expected '%v' but got '%v'. Tested: %v", element, fun.ValidateDates(key), key)
		}
	}
}