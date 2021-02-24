package log

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorInformation struct {
	StatusCode 		 int    `json:"status_code"`
	Location   		 string `json:"location"`
	RawError   		 string `json:"raw_error"`
	PossibleSolution string `json:"possible_solution"`
}

var ErrorInfo ErrorInformation

func UpdateErrorInformation(status int, loc string, err string, solution string) {
	ErrorInfo.StatusCode = status
	ErrorInfo.Location = loc
	ErrorInfo.RawError = err
	ErrorInfo.PossibleSolution = solution
}

func PrintErrorInformation(w http.ResponseWriter) {
	//update header to json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorInfo.StatusCode)
	//send error output to user
	err := json.NewEncoder(w).Encode(ErrorInfo)
	//branch if something went wrong with output
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}