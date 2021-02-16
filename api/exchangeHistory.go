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
	Start_at string                          `json:"start_at"`
	Base     string                          `json:"base"`
	End_at   string                          `json:"end_at"`
}

// runs getInfo
func HandlerExchangeHistory(country string, startDate string, endDate string) {
	//get currency code from country name
	currency := handlerCountryCurrency(country)
	//get all exchange history between two dates
	var inpExchanges exchangeHistory
	getExchangeHistoryData(&inpExchanges, startDate, endDate)
	//filter out exchange history to one specific currency
	var outExchanges exchangeHistory
	filterExchangeHistory(&inpExchanges, &outExchanges, currency, startDate, endDate)

	fmt.Println(outExchanges)
}

func filterExchangeHistory(inpE *exchangeHistory, outE *exchangeHistory, currency string, startDate string, endDate string) {
	//initializer map in struct (could be done in a constructor)
	outE.Rates = make(map[string]map[string]float32)
	//iterate through input structen and adds only the values where the currency code is equal to the request
	for date, mapCurrencies := range inpE.Rates {
		for currencyCode, currencyValue := range mapCurrencies {
			//branch if key in map equal requested currency
			if(currencyCode == currency) {
				//initialize a buffer map to add in Rates map
				buffer := make(map[string]float32)
				buffer[currencyCode] = currencyValue
				outE.Rates[date] = buffer
			}
		}
    }
	//copy data from input structure to output structure
	outE.Start_at = inpE.Start_at
	outE.Base = inpE.Base
	outE.End_at = inpE.End_at
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
