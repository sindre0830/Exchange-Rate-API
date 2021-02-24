package api

import (
	"encoding/json"
	"main/log"
	"net/http"
	"strconv"
	"strings"
)

type exchangeHistory struct {
	Rates    map[string](map[string]float32) `json:"rates"`
	StartAt string                          `json:"start_at"`
	Base     string                          `json:"base"`
	EndAt   string                          `json:"end_at"`
}

func HandlerExchangeHistory(w http.ResponseWriter, r *http.Request) {
	//split URL path by '/'
	arrURL := strings.Split(r.URL.Path, "/")
	//branch if there is an error
	if len(arrURL) != 6 {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> Checking length of URL",
			"Either too many or too few arguments in path.",
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
	//split date by '-' for format checking
	arrDate := strings.Split(dates, "-")
	//check if date format is invalid
	var invalidDateFlag bool
	//check if date has correct amount of elements
	invalidDateFlag = (len(arrDate) != 6) || (len(dates) != 21)
	//check if start date is using correct format YYYY-MM-DD
	invalidDateFlag = invalidDateFlag || ((len(arrDate[0]) != 4) || (len(arrDate[1]) != 2) || (len(arrDate[2]) != 2))
	//check if end date is using correct format YYYY-MM-DD
	invalidDateFlag = invalidDateFlag || ((len(arrDate[3]) != 4) || (len(arrDate[4]) != 2) || (len(arrDate[5]) != 2))
	//branch if date is valid so far
	if !invalidDateFlag {
		//check if all date elements are integers and larger than 0. 'hehe-01-00' == false
		for _, elemDate := range arrDate {
			elemDateNum, err := strconv.Atoi(elemDate)
			if err != nil || elemDateNum < 1 {
				invalidDateFlag = true
				break
			}
		}
	}
	//check if end date is larger or equal than start date
	invalidDateFlag = invalidDateFlag || (dates[:10] > dates[11:])
	//branch if there is an error
	if invalidDateFlag {
		log.UpdateErrorMessage(
			http.StatusBadRequest, 
			"HandlerExchangeHistory() -> Checking if inputed dates are valid",
			"Wrong format or the start date is larger than the end date.",
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
			"Start date is empty.",
			"Date is not valid, check if day or month is larger than valid number. Expected format: '.../start_at-end_at' (YYYY-MM-DD-YYYY-MM-DD). Example: '.../2020-01-20-2021-02-01'",
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
