package api

import (
	"encoding/json"
	"fmt"
	"main/log"
	"net/http"
	"time"
)

// diagnosis structure keeps version, uptime and status codes on used API's.
type diagnosis struct {
	Exchangeratesapi int    `json:"exchangeratesapi"`
	Restcountries    int    `json:"restcountries"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}
// StartTime is a variable declared at the start of the program to calculate uptime.
var StartTime time.Time
// getUptime calculates uptime based on start time and current time.
func getUptime() float64 {
	return time.Since(StartTime).Seconds()
}
// HandlerDiagnosis request input to diagnosis structure and writes to user.
func HandlerDiagnosis(w http.ResponseWriter, r *http.Request) {
	var diag diagnosis
	//declare error variable
	var err error
	//get exchange rate API status code
	resp, err := http.Get("https://api.exchangeratesapi.io/latest")
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerDiagnosis() -> Getting status code from 'https://api.exchangeratesapi.io/latest'",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	diag.Exchangeratesapi = resp.StatusCode
	//get rest countries API status code
	resp, err = http.Get("https://restcountries.eu/rest/v2/all")
	//branch if there is an error
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerDiagnosis() -> Getting status code from 'https://restcountries.eu/rest/v2/all'",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
		return
	}
	diag.Restcountries = resp.StatusCode
	//set version
	diag.Version = "v1"
	//get uptime
	diag.Uptime = fmt.Sprintf("%f", getUptime())
	//set header to JSON
	w.Header().Set("Content-Type", "application/json")
	//send output to user
	err = json.NewEncoder(w).Encode(diag)
	//branch if something went wrong with output
	if err != nil {
		log.UpdateErrorMessage(
			http.StatusInternalServerError, 
			"HandlerDiagnosis() -> Sending output to user",
			err.Error(),
			"Unknown",
		)
		log.PrintErrorInformation(w)
	}
}
