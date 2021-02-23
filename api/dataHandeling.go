package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func requestData(url string) ([]byte, error) {
	//create new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("ERROR in creating request", err)
		return nil, err
	}
	//timeout after 2 seconds
	apiClient := http.Client{
		Timeout: time.Second * 2,
	}
	//get response
	res, err := apiClient.Do(req)
	if err != nil {
		fmt.Println("ERROR in response", err)
		return nil, err
	}
	//read output
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("ERROR when reading response", err)
		return nil, err
	}
	return output, err
}
