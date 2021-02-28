package api

import (
	"io/ioutil"
	"main/debug"
	"net/http"
	"time"
)

// requestData gets raw data from API's
func requestData(url string) ([]byte, error) {
	//declare error variable
	var err error
	//create new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	//branch if there is an error
	if err != nil || req.Response.StatusCode != 200 {
		debug.UpdateStatusCode(req.Response.StatusCode)
		return nil, err
	}
	//timeout after 2 seconds
	apiClient := http.Client{
		Timeout: time.Second * 2,
	}
	//get response
	res, err := apiClient.Do(req)
	//branch if there is an error
	if err != nil {
		return nil, err
	}
	//read output
	output, err := ioutil.ReadAll(res.Body)
	//branch if there is an error
	if err != nil {
		return nil, err
	}
	return output, err
}
