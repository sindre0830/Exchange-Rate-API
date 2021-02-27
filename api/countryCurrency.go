package api

import (
	"encoding/json"
)

// countryCurrency structure keeps all information about currencies.
type countryCurrency struct {
	Currencies []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}
// handlerCountryCurrency handles getting currency information of a given country.
func handlerCountryCurrency(country string, flagAlpha bool) (string, error)  {
	//create error variable
	var err error
	//request country currency code (NOK, USD, EUR, ...)
	var inpData countryCurrency
	//branch if country parameter isn't alpha code and request alpha code
	if !flagAlpha {
		country, err = handlerCountryNameToAlpha(country)
		//branch if there is an error
		if err != nil {
			return "", err
		}
	}
	//request currency code
	err = getCountryCurrency(&inpData, country)
	//branch if there is an error
	if err != nil {
		return "", err
	}
	//filter through the inputed data and generate data for output
	outData := inpData.Currencies[0].Code
	return outData, err
}
// getCountryCurrency request currency information of a given country.
func getCountryCurrency(e *countryCurrency, country string) error {
	//url to API
	url := "https://restcountries.eu/rest/v2/alpha/" + country + "?fields=currencies"
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
