package log

import (
	"encoding/json"
	"fmt"
	"net/http"
)


type ErrorInformation struct {
	StatusCode int    `json:"status_code"`
	RawError   string `json:"raw_error"`
	Location   string `json:"location"`
}

var ErrorInfo ErrorInformation

func UpdateErrorInformation(status int, err string, loc string) {
	ErrorInfo.StatusCode = status
	ErrorInfo.Location = loc
	ErrorInfo.RawError = err
}

func PrintErrorInformation(w http.ResponseWriter) {
	//set header to json
	w.Header().Set("Content-Type", "application/json")
	//send error output to user
	err := json.NewEncoder(w).Encode(ErrorInfo)
	//branch if something went wrong with output
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}