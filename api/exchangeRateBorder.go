package api

import "fmt"

func HandlerExchangeRateBorder(country string, limit int) {
	//get bordering countries from country name
	countryNeighbours := handlerCountryBorder(country, limit)
	fmt.Println(countryNeighbours)
}
