package api

import (
	"encoding/json"
	"log"
)

type exchangeRate struct {
	Rates map[string]float32 `json:"rates"`
	Base  string `json:"base"`
	Date  string `json:"date"`
}

type ExchangeRateBorder struct {
	Rates map[string]countryCurrencyRate `json:"rates"`
	Base string `json:"base"`
}

type countryCurrencyRate struct {
	Currency string `json:"currency`
	Rate float32 `json:"rate"`
}

func HandlerExchangeRateBorder(country string, limit int) ExchangeRateBorder {
	//get bordering countries from requested country
	arrNeighbourCode := handlerCountryBorder(country)
	//get the currencies of the bordering countries
	var arrNeighbourCurrency []string
	for _, neighbour := range arrNeighbourCode {
		arrNeighbourCurrency = append(arrNeighbourCurrency, handlerCountryCurrency(neighbour, true))
	}
	//get the currency of the requested country and set it as base
	baseCurrency := handlerCountryCurrency(country, false)
	//request all available currency data
	var inpData exchangeRate
	getExchangeRateBorderData(&inpData, baseCurrency)
	//filter through the inputed data and generate data for output
	var outData ExchangeRateBorder
	filterExchangeRateBorder(&inpData, &outData, arrNeighbourCode, arrNeighbourCurrency, limit)
	
	return outData
}

func filterExchangeRateBorder(inpE *exchangeRate, outE *ExchangeRateBorder, arrNeighbourCode []string, arrNeighbourCurrency []string, limit int) {
	//update output
	outE.Base = inpE.Base
	//initialize map in struct (could be done in a constructor)
	outE.Rates = make(map[string]countryCurrencyRate)
	//initialize a buffer struct to add in map
	var bufferStruct countryCurrencyRate
	//iterate through array of neighbours currency
	for i, targetCurrency := range arrNeighbourCurrency {
		//branch if target is the same as the base currency and set rate to 1
		if targetCurrency == inpE.Base {
			//update buffer
			bufferStruct.Currency = inpE.Base
			bufferStruct.Rate = 1
			//update output
			outE.Rates[arrNeighbourCode[i]] = bufferStruct
			//check if limit is hit
			if (limit > 0) && (len(outE.Rates) >= limit) {
				return
			}
		} else {
			//iterate through map of all available currencies
			for currencyName, currencyRate := range inpE.Rates {
				//branch if target currency is available
				if targetCurrency == currencyName {
					//update buffer
					bufferStruct.Currency = currencyName
					bufferStruct.Rate = currencyRate
					//update output
					outE.Rates[arrNeighbourCode[i]] = bufferStruct
					if (limit > 0) && (len(outE.Rates) >= limit) {
						return
					}
				}
			}
		}
	}
}

func getExchangeRateBorderData(e *exchangeRate, baseCurrency string) {
	url := "https://api.exchangeratesapi.io/latest?base=" + baseCurrency

	output := requestData(url)
	
	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
