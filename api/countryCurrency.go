package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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

	// Create new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Errorf("Error in creating request:", err.Error())
	}

	apiClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	// Get response
	res, err := apiClient.Do(req)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	// Read output
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}