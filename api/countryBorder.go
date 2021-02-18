package api

import (
	"encoding/json"
	"log"
)

type countryBorder []struct {
	Borders []string `json:"borders"`
}

func handlerCountryBorder(country string, limit int) []string {
	var countryData countryBorder
	getCountryBorderData(&countryData, country)
	//branch if limit is less than 1 or limit is higher than array length and return all the bordering countries (when limit parameter is not used)
	if (limit < 1) || (limit >= len(countryData[0].Borders)) {
		return countryData[0].Borders
	}
	return countryData[0].Borders[:limit]
}

func getCountryBorderData(e *countryBorder, country string) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=borders"

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
