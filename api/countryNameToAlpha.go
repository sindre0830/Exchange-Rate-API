package api

import (
	"encoding/json"
)

// countryAlphacode structure keeps all alphacodes of a given country.
type countryAlphacode []struct {
	Alpha3Code string `json:"alpha3Code"`
}
// handlerCountryNameToAlpha handles converting country name to alphacode.
func handlerCountryNameToAlpha(country string) (string, error) {
	//request country alpha code (3 characters long)
	var inpData countryAlphacode
	err := getCountryAlphaCode(&inpData, country)
	//branch if there is an error
	if err != nil {
		return "", err
	}
	//filter through the inputed data and generate data for output
	outData := inpData[0].Alpha3Code
	return outData, err
}
// getCountryAlphaCode request alphacode of country.
func getCountryAlphaCode(e *countryAlphacode, country string) error {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=alpha3Code"
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
