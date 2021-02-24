package api

import (
	"encoding/json"
)

type countryCurrency struct {
	Currencies []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

func handlerCountryCurrency(country string, flagAlpha bool) (string, error)  {
	//create error variable
	var err error
	//request country currency code (NOK, USD, EUR, ...)
	var inpData countryCurrency
	//branch if country parameter is an alpha code (NOR, SWE, FIN, ...)
	if flagAlpha {
		err = getCountryCurrencyData(&inpData, country)
	//branch if country parameter isn't alpha code and request alpha code
	} else {
		country, err := handlerCountryNameToAlpha(country)
		if err != nil {
			return "", err
		}
		err = getCountryCurrencyData(&inpData, country)
	}
	//filter through the inputed data and generate data for output
	outData := inpData.Currencies[0].Code
	return outData, err
}

func getCountryCurrencyData(e *countryCurrency, country string) error {
	//url to api
	url := "https://restcountries.eu/rest/v2/alpha/" + country + "?fields=currencies"
	//gets raw output from api
	output, err := requestData(url)
	if err != nil {
		return err
	}
	//convert raw output to json
	err = json.Unmarshal(output, &e)
	return err
}
