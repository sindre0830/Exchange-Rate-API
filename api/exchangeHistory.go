package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type exchangeHistory struct {
	Rates    map[string](map[string]float32) `json:"rates"`
	Start_at string                        `json:"start_at"`
	Base     string                        `json:"base"`
	End_at   string                        `json:"end_at"`
}

// runs getInfo
func HandlerExchangeHistory(country string, startDate string, endDate string) {
	currency := HandlerCountryCurrency(country)
	var exchanges exchangeHistory
	getExchangeHistoryData(&exchanges, startDate, endDate)
	fmt.Println(exchanges.Rates["2020-01-02"][currency])
}

func getExchangeHistoryData(e *exchangeHistory, startDate string, endDate string) {
	url := "https://api.exchangeratesapi.io/history?start_at=" + startDate + "&end_at=" + endDate

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
