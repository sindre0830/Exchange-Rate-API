package api

import (
	"encoding/json"
	"main/fun"
	"main/log"
	"net/http"
	"strings"
)

// exchangeHistory structure keeps all rates between two dates based on a base currency.
type exchangeHistory struct {
	Rates    map[string](map[string]float32) `json:"rates"`
	StartAt  string                          `json:"start_at"`
	Base     string                          `json:"base"`
	EndAt    string                          `json:"end_at"`
}
// HandlerExchangeHistory handles getting input from URL and requesting exchange history to output to user.
func HandlerExchangeHistory(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/'
	arrURL := strings.Split(r.URL.Path, "/")
	//branch if there is an error
	if len(arrURL) != 6 {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> Checking length of URL",
			"url validation: either too many or too few arguments in url path",
			"Path format. Expected format: '.../country/start_at-end_at'. Example: '.../norway/2020-01-20-2021-02-01'.",
		)
		log.PrintErrorInformation(w)
		return
	}
	//set country variable
	country := arrURL[4]
	//request base currency code from country name
	currency, err := handlerCountryCurrency(country, false)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> handlerCountryCurrency() -> Getting base currency from requested country in URL",
			err.Error(),
			"Not valid country. Expected format: '.../country/...'. Example: '.../norway/...'.",
		)
		log.PrintErrorInformation(w)
		return
	}
	//get dates from url
	dates := arrURL[5]
	//validate date input
	err = fun.ValidateDates(dates)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> Checking if inputed dates are valid",
			err.Error(),
			"Date format. Expected format: '.../start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../2020-01-20-2021-02-01'",
		)
		log.PrintErrorInformation(w)
		return
	}
	//set start- and end date variables
	startDate := dates[:10]
	endDate := dates[11:]
	//request all exchange history between two dates
	var inpData exchangeHistory
	err = getExchangeHistoryData(&inpData, startDate, endDate)
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeHistory() -> getExchangeHistoryData() -> Getting all rates between two dates",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	//branch if start date is empty in input. This is an indicator to invalid date. 2020-60-42 == false
	if inpData.StartAt == "" {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeHistory() -> getExchangeHistoryData() -> Getting all rates between two dates",
			"date validation: empty input from API",
			"Date is not valid, check if start or end date is a valid date. Expected format: '.../start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../2020-01-20-2021-02-01'",
		)
		log.PrintErrorInformation(w)
		return
	}
	//filter through the inputed data and generate data for output
	var outData exchangeHistory
	filterExchangeHistory(&inpData, &outData, currency, startDate, endDate)
	//set header to json
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(outData)
	//branch if something went wrong with output
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeHistory() -> Sending output to user",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
	}
}
// filterExchangeHistory filters through input and send data to output
func filterExchangeHistory(inpData *exchangeHistory, outData *exchangeHistory, currency string, startDate string, endDate string) {
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
// getExchangeHistoryData request all rates between two dates.
func getExchangeHistoryData(e *exchangeHistory, startDate string, endDate string) error {
	url := "https://api.exchangeratesapi.io/history?start_at=" + startDate + "&end_at=" + endDate
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
