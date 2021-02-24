package log

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorMessage struct {
	StatusCode 		 int    `json:"status_code"`
	Location   		 string `json:"location"`
	RawError   		 string `json:"raw_error"`
	PossibleReason   string `json:"possible_reason"`
}

var ErrorMsg ErrorMessage

func UpdateErrorMessage(status int, loc string, err string, reason string) {
	ErrorMsg.StatusCode = status
	ErrorMsg.Location = loc
	ErrorMsg.RawError = err
	ErrorMsg.PossibleReason = reason
}

func PrintErrorInformation(w http.ResponseWriter) {
	//update header to json and set http code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorMsg.StatusCode)
	//send error output to user
	err := json.NewEncoder(w).Encode(ErrorMsg)
	//branch if something went wrong with output
	if err != nil {
		fmt.Println("ERROR encoding JSON in PrintErrorInformation()", err)
		return
	}
	//send error output to console
	fmt.Printf("\nstatus_code: %v,\nlocation: %s,\nraw_error: %s,\npossible_reason: %s\n", ErrorMsg.StatusCode, ErrorMsg.Location, ErrorMsg.RawError, ErrorMsg.PossibleReason)
}