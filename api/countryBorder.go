package api

import (
	"encoding/json"
)

// countryBorder structure keeps all information about bordering countries.
type countryBorder []struct {
	Borders []string `json:"borders"`
}
// handlerCountryBorder handles getting bordering countries of a given country.
func handlerCountryBorder(country string) ([]string, error) {
	//request all bordering countries of inputed country
	var inpData countryBorder
	err := getCountryBorderData(&inpData, country)
	//branch if there is an error
	if err != nil {
		return nil, err
	}
	//filter through the inputed data and generate data for output
	outData := inpData[0].Borders
	return outData, err
}
// getCountryBorderData request bordering countries of a given country.
func getCountryBorderData(e *countryBorder, country string) error {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=borders"
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
