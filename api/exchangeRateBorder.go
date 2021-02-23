package api

import (
	"encoding/json"
	"fmt"
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
	Currency string `json:"currency"`
	Rate float32 `json:"rate"`
}

func HandlerExchangeRateBorder(country string, limit int) ExchangeRateBorder {
	//get bordering countries from requested country
	arrNeighbourCode, err := handlerCountryBorder(country)
	if err != nil {
		//handle error
	}
	//get the currencies of the bordering countries
	var arrNeighbourCurrency []string
	for _, neighbour := range arrNeighbourCode {
		neighbourCurrency, err := handlerCountryCurrency(neighbour, true)
		if err != nil {
			//handle error
		}
		arrNeighbourCurrency = append(arrNeighbourCurrency, neighbourCurrency)
	}
	//get the currency of the requested country and set it as base
	baseCurrency, err := handlerCountryCurrency(country, false)
	if err != nil {
		//handle error
	}
	//request all available currency data
	var inpData exchangeRate
	err = getExchangeRateBorderData(&inpData, baseCurrency)
	if err != nil {
		//handle error
	}
	//filter through the inputed data and generate data for output
	var outData ExchangeRateBorder
	filterExchangeRateBorder(&inpData, &outData, arrNeighbourCode, arrNeighbourCurrency, limit)
	return outData
}

func filterExchangeRateBorder(inpData *exchangeRate, outData *ExchangeRateBorder, arrNeighbourCode []string, arrNeighbourCurrency []string, limit int) {
	//update output
	outData.Base = inpData.Base
	//initialize map in struct (could be done in a constructor)
	outData.Rates = make(map[string]countryCurrencyRate)
	//initialize a buffer struct to add in map
	var bufferStruct countryCurrencyRate
	//iterate through array of neighbours currency
	for i, targetCurrency := range arrNeighbourCurrency {
		//branch if target is the same as the base currency and set rate to 1
		if targetCurrency == inpData.Base {
			//update buffer
			bufferStruct.Currency = inpData.Base
			bufferStruct.Rate = 1
			//update output
			outData.Rates[arrNeighbourCode[i]] = bufferStruct
			//check if limit is hit
			if (limit > 0) && (len(outData.Rates) >= limit) {
				return
			}
		} else {
			//iterate through map of all available currencies
			for currencyName, currencyRate := range inpData.Rates {
				//branch if target currency is available
				if targetCurrency == currencyName {
					//update buffer
					bufferStruct.Currency = currencyName
					bufferStruct.Rate = currencyRate
					//update output
					outData.Rates[arrNeighbourCode[i]] = bufferStruct
					//check if limit is hit
					if (limit > 0) && (len(outData.Rates) >= limit) {
						return
					}
				}
			}
		}
	}
}

func getExchangeRateBorderData(e *exchangeRate, baseCurrency string) error {
	url := "https://api.exchangeratesapi.io/latest?base=" + baseCurrency
	//gets raw output from api
	output, err := requestData(url)
	if err != nil {
		return err
	}
	//convert raw output to json
	err = json.Unmarshal(output, &e)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
	return err
}
