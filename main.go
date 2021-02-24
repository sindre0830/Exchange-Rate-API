package main

import (
	"log"
	"main/api"
	"net/http"
	"os"
	"time"
)

// Main program
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	api.StartTime = time.Now()

	http.HandleFunc("/exchange/v1/exchangehistory/", api.HandlerExchangeHistory)

	http.HandleFunc("/exchange/v1/exchangeborder/", api.HandlerExchangeRateBorder)

	http.HandleFunc("/exchange/v1/diag/", api.HandlerDiagnosis)

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
