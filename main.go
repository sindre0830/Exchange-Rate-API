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

func handleDiag(w http.ResponseWriter, r *http.Request) {
	var diag diagnosis

	resp, err := http.Get("https://api.exchangeratesapi.io/latest")
    if err != nil {
        log.Fatal(err)
    }
	diag.Exchangeratesapi = resp.StatusCode
	

	resp, err = http.Get("https://restcountries.eu/rest/v2/all")
    if err != nil {
        log.Fatal(err)
    }
	diag.Restcountries = resp.StatusCode

	diag.Version = "v1"

	diag.Uptime = fmt.Sprintf("%f", upTime())


	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(diag)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}

var startTime time.Time

func upTime() float64 {
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

	http.HandleFunc("/exchange/v1/diag/", handleDiag)

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
