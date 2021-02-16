package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	} else {
		return countryData[0].Borders[:limit]
	}
}

func getCountryBorderData(e *countryBorder, country string) {
	url := "https://restcountries.eu/rest/v2/name/" + country + "?fields=borders"

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
