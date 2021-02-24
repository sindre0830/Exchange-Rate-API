package api

import (
	"encoding/json"
	"main/log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	Currency string `json:"currency"`
	Rate float32 `json:"rate"`
}

func HandlerExchangeRateBorder(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/'
	arrURL := strings.Split(r.URL.Path, "/")
	//branch if there is an error
	if len(arrURL) != 5 {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeRateBorder() -> Checking length of URL",
			"Either too many or too few arguments in path.",
			"Path format. Expected format: '.../country?limit=num' ('?limit=num' is optional). Example: '.../norway?limit=2'.",
		)
		log.PrintErrorInformation(w)
		return
	}
	//set country variable
	country := arrURL[4]
	//get the currency of the requested country and set it as base
	baseCurrency, err := handlerCountryCurrency(country, false)
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
	err = getExchangeRateBorderData(&inpData, baseCurrency)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeRateBorder() -> getExchangeRateBorderData() -> Getting latest rates based on base currency",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	//set default limit to 0 (no limit)
	limit := 0
	//get all parameters from URL
	arrURLParameters, err := url.ParseQuery(r.URL.RawQuery)
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
	if len(arrURLParameters) > 0 {
		//branch if field 'limit' exist
		if elemURLParameters, ok := arrURLParameters["limit"]; ok {
			//set new limit according to URL parameter
			limit, err = strconv.Atoi(elemURLParameters[0])
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
				"HandlerExchangeRateBorder() -> Checking if limit field exist",
				"Fields in URL used, but doesn't contain 'limit'.",
				"Wrong field, or typo. Expected format: '...?limit=num'. Example: '...limit=2'.",
			)
			log.PrintErrorInformation(w)
			return
		}
	}
	//get bordering countries from requested country
	arrNeighbourCode, err := handlerCountryBorder(country)
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
	var arrNeighbourCurrency []string
	for _, neighbour := range arrNeighbourCode {
		neighbourCurrency, err := handlerCountryCurrency(neighbour, true)
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
		arrNeighbourCurrency = append(arrNeighbourCurrency, neighbourCurrency)
	}
	//filter through the inputed data and generate data for output
	var outData exchangeRateBorder
	filterExchangeRateBorder(&inpData, &outData, arrNeighbourCode, arrNeighbourCurrency, limit)
	//set header to json
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

func filterExchangeRateBorder(inpData *exchangeRate, outData *exchangeRateBorder, arrNeighbourCode []string, arrNeighbourCurrency []string, limit int) {
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
	//branch if there is an error
	if err != nil {
		return err
	}
	//convert raw output to json
	err = json.Unmarshal(output, &e)
	return err
}
