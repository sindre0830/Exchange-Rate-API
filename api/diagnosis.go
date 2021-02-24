package api

import (
	"encoding/json"
	"fmt"
	"log"
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
	if err != nil {
		log.Fatal(err)
	}
	diag.Exchangeratesapi = resp.StatusCode
	//get rest countries api status code
	resp, err = http.Get("https://restcountries.eu/rest/v2/all")
	if err != nil {
		log.Fatal(err)
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
		fmt.Println("ERROR encoding JSON", err)
	}
}

var StartTime time.Time

func getUptime() float64 {
	return time.Since(StartTime).Seconds()
}