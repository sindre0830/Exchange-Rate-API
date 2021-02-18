package api

import (
	"encoding/json"
	"log"
)

type countryNameToAlpha []struct {
	Alpha3Code string `json:"alpha3Code"`
}

func handlerCountryNameToAlpha(country string) string {
	var countryData countryNameToAlpha
	getCountryAlphaCodeData(&countryData, country)
	return countryData[0].Alpha3Code
}

func getCountryAlphaCodeData(e *countryNameToAlpha, country string) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=alpha3Code"

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}