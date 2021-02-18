package api

import (
	"encoding/json"
	"log"
)

type countryBorder []struct {
	Borders []string `json:"borders"`
}

func handlerCountryBorder(country string) []string {
	var countryData countryBorder
	getCountryBorderData(&countryData, country)
	return countryData[0].Borders
}

func getCountryBorderData(e *countryBorder, country string) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=borders"

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
