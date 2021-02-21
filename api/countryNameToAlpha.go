package api

import (
	"encoding/json"
	"log"
)

type countryNameToAlpha []struct {
	Alpha3Code string `json:"alpha3Code"`
}

func handlerCountryNameToAlpha(country string) string {
	//request country alpha code (3 characters long)
	var inpData countryNameToAlpha
	getCountryAlphaCodeData(&inpData, country)
	//filter through the inputed data and generate data for output
	outData := inpData[0].Alpha3Code
	return outData
}

func getCountryAlphaCodeData(e *countryNameToAlpha, country string) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=alpha3Code"

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
