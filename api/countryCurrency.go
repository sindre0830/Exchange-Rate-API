package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type countryCurrency []struct {
	Currencies []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}

func HandlerCountryCurrency(country string) string {
	var currency countryCurrency
	getCountryCurrencyData(&currency, country)
	return currency[0].Currencies[0].Code
}

func getCountryCurrencyData(e *countryCurrency, country string) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=currencies"

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