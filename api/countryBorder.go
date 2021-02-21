package api

import (
	"encoding/json"
	"log"
)

type countryBorder []struct {
	Borders []string `json:"borders"`
}

func handlerCountryBorder(country string) []string {
	//request all bordering countries of inputed country
	var inpData countryBorder
	getCountryBorderData(&inpData, country)
	//filter through the inputed data and generate data for output
	outData := inpData[0].Borders
	return outData
}

func getCountryBorderData(e *countryBorder, country string) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=borders"

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
