package api

import (
	"encoding/json"
	"fmt"
	"log"
)

type exchangeRate struct {
	Rates map[string]float32 `json:"rates"`
	Base  string `json:"base"`
	Date  string `json:"date"`
}

type exchangeRateBorder struct {
	Rates map[string]countryCurrencyRate `json:"rates"`
	Base string `json:"base"`
}

type countryCurrencyRate struct {
	Currency string `json:"currency`
	Rate float32 `json:"rate"`
}

func HandlerExchangeRateBorder(country string, limit int) {
	//get bordering countries from requested country
	arrNeighbourCode := handlerCountryBorder(country, limit)
	fmt.Println(arrNeighbourCode)
	//get the currencies of the bordering countries
	var arrNeighbourCurrency []string
	for _, neighbour := range arrNeighbourCode {
		arrNeighbourCurrency = append(arrNeighbourCurrency, handlerCountryCurrency(neighbour, true))
	}
	fmt.Println(arrNeighbourCurrency)
	//get the currency of the requested country and set it as base
	baseCurrency := handlerCountryCurrency(country, false)
	//request all available currency data
	var inpData exchangeRate
	getExchangeRateBorderData(&inpData, baseCurrency)
	//filter through the inputed data and generate data for output
	var outData exchangeRateBorder
	filterExchangeRateBorder(&inpData, &outData, arrNeighbourCode, arrNeighbourCurrency)
		fmt.Println(outData)
}

func filterExchangeRateBorder(inpE *exchangeRate, outE *exchangeRateBorder, arrNeighbourCode []string, arrNeighbourCurrency []string) {
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
				}
			}
		}
	}
	//update output
	outE.Base = inpE.Base
}

func getExchangeRateBorderData(e *exchangeRate, baseCurrency string) {
	url := "https://api.exchangeratesapi.io/latest?base=" + baseCurrency

	output := requestData(url)
	
	jsonErr := json.Unmarshal(output, &e)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
