package main

import (
	"fmt"
	"log"
	"main/api"
	"net/http"
	"os"
)

// Main program
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	data1 := api.HandlerExchangeHistory("norway", "2020-01-01", "2020-01-10")
	fmt.Println(data1)

	data2 := api.HandlerExchangeRateBorder("russia", 4)
	fmt.Println(data2)

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
