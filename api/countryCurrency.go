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
	//request country currency code (NOK, USD, EUR, ...)
	var inpData countryCurrency
	//branch if country parameter is an alpha code (NOR, SWE, FIN, ...)
	if flagAlpha {
		getCountryCurrencyData(&inpData, country)
	//branch if country parameter isn't alpha code and request alpha code
	} else {
		country = handlerCountryNameToAlpha(country)
		getCountryCurrencyData(&inpData, country)
	}
	//filter through the inputed data and generate data for output
	outData := inpData.Currencies[0].Code
	return outData
}

func getCountryCurrencyData(e *countryCurrency, country string) {
	url := "https://restcountries.eu/rest/v2/alpha/" + country + "?fields=currencies"

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
