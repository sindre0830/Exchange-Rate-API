package main

import (
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

	api.HandlerExchangeHistory("norway", "2020-01-01", "2020-01-10")

	log.Fatal(http.ListenAndServe(":" + port, nil))
}
