package api

import (
	"encoding/json"
)

// countryCurrency structure keeps all information about currencies.
type countryCurrency struct {
	Currencies []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
}
// countryAll structure keeps all information about one or more countries
type countryAll []struct {
	Name           string    `json:"name"`
	TopLevelDomain []string  `json:"topLevelDomain"`
	Alpha2Code     string    `json:"alpha2Code"`
	Alpha3Code     string    `json:"alpha3Code"`
	CallingCodes   []string  `json:"callingCodes"`
	Capital        string    `json:"capital"`
	AltSpellings   []string  `json:"altSpellings"`
	Region         string    `json:"region"`
	Subregion      string    `json:"subregion"`
	Population     int       `json:"population"`
	Latlng         []float64 `json:"latlng"`
	Demonym        string    `json:"demonym"`
	Area           float64   `json:"area"`
	Gini           float64   `json:"gini"`
	Timezones      []string  `json:"timezones"`
	Borders        []string  `json:"borders"`
	NativeName     string    `json:"nativeName"`
	NumericCode    string    `json:"numericCode"`
	Currencies     []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages []struct {
		Iso6391    string `json:"iso639_1"`
		Iso6392    string `json:"iso639_2"`
		Name       string `json:"name"`
		NativeName string `json:"nativeName"`
	} `json:"languages"`
	Translations struct {
		De string `json:"de"`
		Es string `json:"es"`
		Fr string `json:"fr"`
		Ja string `json:"ja"`
		It string `json:"it"`
		Br string `json:"br"`
		Pt string `json:"pt"`
		Nl string `json:"nl"`
		Hr string `json:"hr"`
		Fa string `json:"fa"`
	} `json:"translations"`
	Flag          string `json:"flag"`
	RegionalBlocs []struct {
		Acronym       string        `json:"acronym"`
		Name          string        `json:"name"`
		OtherAcronyms []interface{} `json:"otherAcronyms"`
		OtherNames    []interface{} `json:"otherNames"`
	} `json:"regionalBlocs"`
	Cioc string `json:"cioc"`
}
// handlerCountryCurrency handles getting currency information of a given country.
func handlerCountryCurrency(arrCountry []string, flagAlpha bool) ([]string, error)  {
	//create error variable
	var err error
	//request country currency code (NOK, USD, EUR, ...)
	var inpData countryAll
	//branch if country parameter isn't alpha code and request alpha code
	if !flagAlpha {
		//this is only used for user input so it will never be more then one element in this case
		arrCountry[0], err = handlerCountryNameToAlpha(arrCountry[0])
		//branch if there is an error
		if err != nil {
			return nil, err
		}
	}
	//request currency code
	err = getCountryCurrency(&inpData, arrCountry)
	//branch if there is an error
	if err != nil {
		return nil, err
	}
	//filter through the inputed data and generate data for output
	var outData []string
	for _, country := range inpData {
		currency := country.Currencies[0].Code
		outData = append(outData, currency)
	}
	return outData, err
}
// getCountryCurrency request currency information of a given country.
func getCountryCurrency(e *countryAll, arrCountry []string) error {
	//
	var codes string
	for _, country := range arrCountry {
		codes += country + ";"
	}
	//url to API
	url := "https://restcountries.eu/rest/v2/alpha?codes=" + codes
	//gets raw output from API
	output, err := requestData(url)
	//branch if there is an error
	if err != nil {
		return err
	}
	//convert raw output to JSON
	err = json.Unmarshal(output, &e)
	return err
}
