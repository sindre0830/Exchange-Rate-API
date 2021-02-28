package api

import (
	"encoding/json"
	"main/debug"
	"main/fun"
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
	//create error variable
	var err error
	//split URL path by '/'
	arrURL := strings.Split(r.URL.Path, "/")
	//branch if there is an error
	if len(arrURL) != 6 {
		debug.UpdateErrorMessage(
			http.StatusBadRequest,
			"HandlerExchangeHistory() -> Checking length of URL",
			"url validation: either too many or too few arguments in url path",
			"Path format. Expected format: '.../country/start_at-end_at'. Example: '.../norway/2020-01-20-2021-02-01'.",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//set country variable
	var country []string
	country = append(country, arrURL[4])
	//convert country name to its alphacode
	country[0], err = handlerCountryNameToAlpha(country[0])
	//branch if there is an error
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> handlerCountryNameToAlpha() -> Converting country name to its alphacode",
			err.Error(),
			"Not valid country. Expected format: '.../country'. Example: '.../norway'.",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//request currency code from country name
	currency, _, err := handlerCountryCurrency(country)
	//branch if there is an error
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> handlerCountryCurrency() -> Getting base currency from requested country in URL",
			err.Error(),
			"Not valid country. Expected format: '.../country/...'. Example: '.../norway/...'.",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//get dates from URL
	dates := arrURL[5]
	//validate date input
	err = fun.ValidateDates(dates)
	//branch if there is an error
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> Checking if inputed dates are valid",
			err.Error(),
			"Date format. Expected format: '.../start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../2020-01-20-2021-02-01'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//set start- and end date variables
	startDate := dates[:10]
	endDate := dates[11:]
	//request all exchange history between two dates of a given currency
	var inpData exchangeHistory
	err = getExchangeHistory(&inpData, startDate, endDate, currency[0])
	//branch if there is an error
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeHistory() -> getExchangeHistory() -> Getting all rates between two dates",
			err.Error(),
			"Unknown",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//branch if start date is empty in input, this is likely caused by an invalid date. 2020-60-42 == false
	if inpData.StartAt == "" {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeHistory() -> getExchangeHistory() -> Getting all rates between two dates",
			"date validation: empty input from API",
			"Date is not valid, check if start or end date is a valid date. Expected format: '.../start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../2020-01-20-2021-02-01'",
		)
		debug.PrintErrorInformation(w)
		return
	}
	//since the input is already filtered, we can send it unedited to output (made new variable for consistency)
	outData := inpData
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(outData)
	//branch if something went wrong with output
	if err != nil {
		debug.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerExchangeHistory() -> Sending output to user",
			err.Error(),
			"Unknown",
		)
		debug.PrintErrorInformation(w)
	}
}
// getExchangeHistory request all rates between two dates of a given currency.
func getExchangeHistory(e *exchangeHistory, startDate string, endDate string, currency string) error {
	//declare error variable
	var err error
	//url to API
	url := "https://api.exchangeratesapi.io/history?start_at=" + startDate + "&end_at=" + endDate + "&symbols=" + currency
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
