package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/api"
	"net/http"
	"os"
	"time"
)

type diagnosis struct {
	Exchangeratesapi int `json:"exchangeratesapi"`
	Restcountries int `json:"restcountries"`
	Version string `json:"version"`
	Uptime string `json:"uptime"`
}

func handlerDiagnosis(w http.ResponseWriter, r *http.Request) {
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

var startTime time.Time

func getUptime() float64 {
    return time.Since(startTime).Seconds()
}

// Main program
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	startTime = time.Now()

	http.HandleFunc("/exchange/v1/exchangehistory/", api.HandlerExchangeHistory)

	http.HandleFunc("/exchange/v1/exchangeborder/", api.HandlerExchangeRateBorder)

	http.HandleFunc("/exchange/v1/diag/", handlerDiagnosis)

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
