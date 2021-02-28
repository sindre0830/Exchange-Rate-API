package api

import (
	"encoding/json"
	"fmt"
	"main/log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// exchangeRate structure keeps all current rates based on a base currency.
type exchangeRate struct {
	Rates map[string]float32 `json:"rates"`
	Base  string 			 `json:"base"`
	Date  string 			 `json:"date"`
}
// exchangeRateBorder structure keeps a map of neighbouring countries based on a base currency.
// 	@see countryCurrencyRate struct
type exchangeRateBorder struct {
	Rates map[string]countryCurrencyRate `json:"rates"`
	Base  string 						 `json:"base"`
}
// countryCurrencyRate structure keeps current rates for neighbouring countries.
type countryCurrencyRate struct {
	Currency string  `json:"currency"`
	Rate 	 float32 `json:"rate"`
}
// HandlerExchangeRateBorder handles getting input from URL and requesting exchange rate to bordering countries to output to user.
func HandlerExchangeRateBorder(w http.ResponseWriter, r *http.Request) {
	//create error variable
	var err error
	//split URL path by '/'
	arrURL := strings.Split(r.URL.Path, "/")
	//branch if there is an error
	if len(arrURL) != 5 {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeRateBorder() -> Checking length of URL",
			"url validation: either too many or too few arguments in url path",
			"Path format. Expected format: '.../country?limit=num' ('?limit=num' is optional). Example: '.../norway?limit=2'.",
		)
		log.PrintErrorInformation(w)
		return
	}
	//set country variable
	var country []string
	country = append(country, arrURL[4])
	fmt.Printf("country before alpha: %s\n", country[0])
	//convert country name to its alphacode
	country[0], err = handlerCountryNameToAlpha(country[0])
	fmt.Printf("country after alpha: %s\n", country[0])
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeRateBorder() -> handlerCountryNameToAlpha() -> Converting country name to its alphacode",
			err.Error(),
			"Not valid country. Expected format: '.../country'. Example: '.../norway'.",
		)
		log.PrintErrorInformation(w)
		return
	}
	//get the currency of the requested country and set it as base
	baseCurrency, _, err := handlerCountryCurrency(country)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeRateBorder() -> handlerCountryCurrency() -> Getting base currency from requested country in URL",
			err.Error(),
			"Not valid country. Expected format: '.../country'. Example: '.../norway'.",
		)
		log.PrintErrorInformation(w)
		return
	}
	//request all available currency data
	var inpData exchangeRate
	err = getExchangeRateBorder(&inpData, baseCurrency[0])
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeRateBorder() -> getExchangeRateBorder() -> Getting latest rates based on base currency",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	//set default limit to 0 (no limit)
	limit := 0
	//get all parameters from URL
	arrPathParameters, err := url.ParseQuery(r.URL.RawQuery)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeRateBorder() -> Getting URL field (...?limit=num)",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	//branch if any parameters exist
	if len(arrPathParameters) > 0 {
		//branch if field 'limit' exist
		if targetParameters, ok := arrPathParameters["limit"]; ok {
			//set new limit according to URL parameter
			limit, err = strconv.Atoi(targetParameters[0])
			//branch if there is an error
			if err != nil {
				log.UpdateErrorMessage(
					http.StatusBadRequest, 
					"HandlerExchangeRateBorder() -> Converting limit field to integer",
					err.Error(),
					"Limit value is not a number. Expected format: '...?limit=num'. Example: '...?limit=2'.",
				)
				log.PrintErrorInformation(w)
				return
			}
		//branch if there is an error
		} else {
			log.UpdateErrorMessage(
				http.StatusBadRequest, 
				"HandlerExchangeRateBorder() -> Validating path parameters",
				"path validation: fields in URL used, but doesn't contain 'limit'",
				"Wrong field, or typo. Expected format: '...?limit=num'. Example: '...limit=2'.",
			)
			log.PrintErrorInformation(w)
			return
		}
	}
	//get bordering countries from requested country
	arrNeighbourCode, err := handlerCountryBorder(country[0])
	fmt.Printf("arrNeighbourCode: %v\n", arrNeighbourCode)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeRateBorder() -> handlerCountryBorder() -> Getting neighbouring countries",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	//get the currencies of the bordering countries
	arrNeighbourCurrency, arrNeighbourCode, err := handlerCountryCurrency(arrNeighbourCode)
	fmt.Printf("arrNeighbourCurrency: %v\n", arrNeighbourCurrency)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeRateBorder() -> handlerCountryCurrency() -> Getting neighbouring countries currencies",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	//filter through the inputed data and generate data for output
	var outData exchangeRateBorder
	filterExchangeRateBorder(&inpData, &outData, arrNeighbourCode, arrNeighbourCurrency, limit)
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(outData)
	//branch if something went wrong with output
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeRateBorder() -> Sending output to user",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
	}
}
// filterExchangeRateBorder filters through input and send data to output.
func filterExchangeRateBorder(inpData *exchangeRate, outData *exchangeRateBorder, arrNeighbourCode []string, arrNeighbourCurrency []string, limit int) {
	//update output
	outData.Base = inpData.Base
	//initialize map in struct
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
// getExchangeRateBorder request all current rates based on base currency.
func getExchangeRateBorder(e *exchangeRate, baseCurrency string) error {
	//declare error variable
	var err error
	//url to API
	url := "https://api.exchangeratesapi.io/latest?base=" + baseCurrency
	//gets raw output from API
	output, err := requestData(url)
	//branch if there is an error
	if err != nil {
		return err
	}
	//convert raw output to JSON
	err = json.Unmarshal(output, &e)
	return err
}
