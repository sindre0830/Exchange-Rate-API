package api

import (
	"encoding/json"
	"fmt"
	"main/log"
	"net/http"
	"time"
)

type diagnosis struct {
	Exchangeratesapi int    `json:"exchangeratesapi"`
	Restcountries    int    `json:"restcountries"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}

func HandlerDiagnosis(w http.ResponseWriter, r *http.Request) {
	var diag diagnosis
	//get exchange rate api status code
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
	//get rest countries api status code
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
	//set header to json
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

var StartTime time.Time

func getUptime() float64 {
	return time.Since(StartTime).Seconds()
}