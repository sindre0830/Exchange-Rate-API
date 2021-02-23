package api

import (
	"encoding/json"
	"log"
)

type ExchangeHistory struct {
	Rates    map[string](map[string]float32) `json:"rates"`
	StartAt string                          `json:"start_at"`
	Base     string                          `json:"base"`
	EndAt   string                          `json:"end_at"`
}

func HandlerExchangeHistory(country string, startDate string, endDate string) ExchangeHistory {
	//request base currency code from country name
	currency := handlerCountryCurrency(country, false)
	//request all exchange history between two dates
	var inpData ExchangeHistory
	getExchangeHistoryData(&inpData, startDate, endDate)
	//filter through the inputed data and generate data for output
	var outData ExchangeHistory
	filterExchangeHistory(&inpData, &outData, currency, startDate, endDate)
	return outData
}

func filterExchangeHistory(inpData *ExchangeHistory, outData *ExchangeHistory, currency string, startDate string, endDate string) {
	//initialize map in struct (could be done in a constructor)
	outData.Rates = make(map[string]map[string]float32)
	//iterate through input structure and add only the values where the currency code is equal to the request
	for date, mapCurrencies := range inpData.Rates {
		for currencyCode, currencyValue := range mapCurrencies {
			//branch if key in map equal requested currency
			if currencyCode == currency {
				//initialize a buffer map to add in Rates map
				buffer := make(map[string]float32)
				buffer[currencyCode] = currencyValue
				outData.Rates[date] = buffer
			}
		}
    }
	//copy data from input structure to output structure
	outData.StartAt = inpData.StartAt
	outData.Base = inpData.Base
	outData.EndAt = inpData.EndAt
}

func getExchangeHistoryData(e *ExchangeHistory, startDate string, endDate string) {
	url := "https://api.exchangeratesapi.io/history?start_at=" + startDate + "&end_at=" + endDate

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
