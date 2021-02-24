package api

import (
	"encoding/json"
)

type countryNameToAlpha []struct {
	Alpha3Code string `json:"alpha3Code"`
}

func handlerCountryNameToAlpha(country string) (string, error) {
	//request country alpha code (3 characters long)
	var inpData countryNameToAlpha
	err := getCountryAlphaCodeData(&inpData, country)
	if err != nil {
		return "", err
	}
	//filter through the inputed data and generate data for output
	outData := inpData[0].Alpha3Code
	return outData, err
}

func getCountryAlphaCodeData(e *countryNameToAlpha, country string) error {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=alpha3Code"
	//gets raw output from api
	output, err := requestData(url)
	if err != nil {
		return err
	}
	//convert raw output to json
	err = json.Unmarshal(output, &e)
	return err
}
