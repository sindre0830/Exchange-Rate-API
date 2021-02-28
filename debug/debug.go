package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorMessage structure keeps all error data.
type ErrorMessage struct {
	StatusCode 		 int    `json:"status_code"`
	Location   		 string `json:"location"`
	RawError   		 string `json:"raw_error"`
	PossibleReason   string `json:"possible_reason"`
}
// ErrorMsg is a global variable.
var ErrorMsg ErrorMessage
// UpdateStatusCode adds status code to error msg.
func UpdateStatusCode(status int) {
	ErrorMsg.StatusCode = status
}
// UpdateErrorMessage adds new information to error msg.
func UpdateErrorMessage(status int, loc string, err string, reason string) {
	if ErrorMsg.StatusCode == 0 {
		ErrorMsg.StatusCode = status
		ErrorMsg.PossibleReason = reason
	} else {
		ErrorMsg.PossibleReason = "Unknown"
	}
	ErrorMsg.Location = loc
	ErrorMsg.RawError = err
}
// PrintErrorInformation prints error msg to user and terminal.
func PrintErrorInformation(w http.ResponseWriter) {
	//declare error variable
	var err error
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorMsg.StatusCode)
	//send error output to user
	err = json.NewEncoder(w).Encode(ErrorMsg)
	//branch if something went wrong with output
	if err != nil {
		fmt.Println("ERROR encoding JSON in PrintErrorInformation()", err)
		return
	}
	//send error output to console
	fmt.Printf("\nstatus_code: %v,\nlocation: %s,\nraw_error: %s,\npossible_reason: %s\n", ErrorMsg.StatusCode, ErrorMsg.Location, ErrorMsg.RawError, ErrorMsg.PossibleReason)
}
