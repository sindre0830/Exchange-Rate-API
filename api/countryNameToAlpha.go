package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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