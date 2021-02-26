package main

import (
	"log"
	"main/api"
	"net/http"
	"os"
	"time"
)

// Main program.
func main() {
	//get port
	port := os.Getenv("PORT")
	//branch if there isn't a port and set it to 8080
	if port == "" {
		port = "8080"
	}
	//set varible to current time (For checking uptime)
	api.StartTime = time.Now()
	//handle exchange history
	http.HandleFunc("/exchange/v1/exchangehistory/", api.HandlerExchangeHistory)
	//handle exchange rate to bordering countries
	http.HandleFunc("/exchange/v1/exchangeborder/", api.HandlerExchangeRateBorder)
	//handle program diagnosis
	http.HandleFunc("/exchange/v1/diag/", api.HandlerDiagnosis)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
