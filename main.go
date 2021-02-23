package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/api"
	"net/http"
	"os"
	"strings"
)

func handleHistory(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 6 {
		status := http.StatusBadRequest
		http.Error(w, "Expecting format .../country_name/begin_date-end_date", status)
		return
	}
	date := strings.Split(parts[5], "-")
	startDate := date[0] + "-" + date[1] + "-" + date[2]
	endDate := date[3] + "-" + date[4] + "-" + date[5]
	
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(api.HandlerExchangeHistory(parts[4], startDate, endDate))
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}

// Main program
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/exchange/v1/exchangehistory/", handleHistory)

	/*task2 := api.HandlerExchangeRateBorder("russia", 4)
	fmt.Println(task2)*/

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
