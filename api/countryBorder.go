package api

import (
	"encoding/json"
)

// countryBorder structure keeps all information about bordering countries.
type countryBorder struct {
	Borders []string `json:"borders"`
}
// handlerCountryBorder handles getting bordering countries of a given country.
func handlerCountryBorder(country string) ([]string, error) {
	//create error variable
	var err error
	//request all bordering countries of inputed country
	var inpData countryBorder
	err = getCountryBorder(&inpData, country)
	//branch if there is an error
	if err != nil {
		return nil, err
	}
	//filter through the inputed data and generate data for output
	outData := inpData.Borders
	return outData, err
}
// getCountryBorder request bordering countries of a given country.
func getCountryBorder(e *countryBorder, country string) error {
	//declare error variable
	var err error
	//url to API
	url := "https://restcountries.eu/rest/v2/alpha/" + country + "?fields=borders"
	//gets raw output from API
	output, err := requestData(url)
	//branch if there is an error
	if err != nil {
		return err
	}
	//convert raw output to JSON
	err = json.Unmarshal(output, &e)
	return err
}
