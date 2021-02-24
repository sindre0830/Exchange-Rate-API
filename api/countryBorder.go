package api

import (
	"encoding/json"
)

type countryBorder []struct {
	Borders []string `json:"borders"`
}

func handlerCountryBorder(country string) ([]string, error) {
	//request all bordering countries of inputed country
	var inpData countryBorder
	err := getCountryBorderData(&inpData, country)
	if err != nil {
		return nil, err
	}
	//filter through the inputed data and generate data for output
	outData := inpData[0].Borders
	return outData, err
}

func getCountryBorderData(e *countryBorder, country string) error {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=borders"
	//gets raw output from api
	output, err := requestData(url)
	if err != nil {
		return err
	}
	//convert raw output to json
	err = json.Unmarshal(output, &e)
	return err
}
