package api

import (
	"encoding/json"
	"log"
)

type ExchangeHistory struct {
	Rates    map[string](map[string]float32) `json:"rates"`
	Start_at string                          `json:"start_at"`
	Base     string                          `json:"base"`
	End_at   string                          `json:"end_at"`
}

// runs getInfo
func HandlerExchangeHistory(country string, startDate string, endDate string) ExchangeHistory {
	//get currency code from country name
	currency := handlerCountryCurrency(country, false)
	//get all exchange history between two dates
	var inpExchanges ExchangeHistory
	getExchangeHistoryData(&inpExchanges, startDate, endDate)
	//filter out exchange history to one specific currency
	var outExchanges ExchangeHistory
	filterExchangeHistory(&inpExchanges, &outExchanges, currency, startDate, endDate)

	return outExchanges
}

func filterExchangeHistory(inpE *ExchangeHistory, outE *ExchangeHistory, currency string, startDate string, endDate string) {
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

func getExchangeHistoryData(e *ExchangeHistory, startDate string, endDate string) {
	url := "https://api.exchangeratesapi.io/history?start_at=" + startDate + "&end_at=" + endDate

	output := requestData(url)

	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
