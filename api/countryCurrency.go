package api

import (
	"encoding/json"
	"log"
)

type countryCurrency struct {
	Currencies []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

func handlerCountryCurrency(country string, flagAlpha bool) string {
	var countryData countryCurrency
	if flagAlpha {
		getCountryCurrencyData(&countryData, country)
	} else {
		country = handlerCountryNameToAlpha(country)
		getCountryCurrencyData(&countryData, country)
	}
	return countryData.Currencies[0].Code
}

func getCountryCurrencyData(e *countryCurrency, country string) {
	url := "https://restcountries.eu/rest/v2/alpha/" + country + "?fields=currencies"

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}