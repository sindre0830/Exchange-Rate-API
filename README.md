# [Assignment 1](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Assignment-1)

### Info
- Root path: https://sindre-assignment-1.herokuapp.com/exchange/v1/
- While the submission is individual, I have discussed the tasks with Rickard Loland and Susanne Edvardsen. We have also helped each other with problems that occurred during development ([rubber ducking](https://en.wikipedia.org/wiki/Rubber_duck_debugging) mostly).
- I have used these REST web services to build my service:
    - https://exchangeratesapi.io/
    - https://restcountries.eu/

### How to use

There are 3 endpoints that you can append to the root path.

1. Exchange Rate History of a Currency for a Given Country
    ```
    Method: GET
    Path: .../exchangehistory/{:country_name}/{:begin_date-end_date}
    ```
    - {:country_name} refers to the English name for the country as supported by https://restcountries.eu/.
    - {:begin_date-end_date} indicates the begin date (i.e., the earliest date to be reported) of the exchange rate and the end date (i.e., the latest date of the range) of the period over which exchange rates are reported. Date format: YYYY-MM-DD
        ```
        Example: https://sindre-assignment-1.herokuapp.com/exchange/v1/exchangehistory/norway/2020-01-01-2020-02-01
        ```

2. Current Exchange Rate for Bordering Countries
    ```
    Method: GET
    Path: .../exchangeborder/{:country_name}{?limit={:number}}
    ```
    - {:country_name} refers to the English name for the country as supported by https://restcountries.eu/.
    - {?limit={:number}} is an optional parameter that limits the number of bordering countries for which currencies are reported.
        ```
        Example: https://sindre-assignment-1.herokuapp.com/exchange/v1/exchangeborder/norway?limit=2
        ```

3. Diagnostics Interface
    ```
    Method: GET
    Path: .../diag/
    ```
    - Outputs status of each API used by the program as well as version and uptime of the program.
        ```
        Example: https://sindre-assignment-1.herokuapp.com/exchange/v1/diag/
        ```

## Notes

#### Design decisions

The REST service, https://exchangeratesapi.io, allows for symbol filtering but gives an error if the given symbol doesn't exist in their database. While I could get better performance by filtering currencies when getting all the bordering currencies, I've decided to get all rates, then filter myself. That way the user can request bordering currencies of Russia and get all available currencies that exist, instead of an error.
The REST service, https://restcountries.eu/, allows for alpha code filtering which would be beneficial when getting all the bordering countries' currencies but doesn't seem to allow for field filtering as well. This means that I have to decide on either storing a massive amount of data that I won't use for each bordering country or request more than once (one request per bordering country). I've decided on the latter, but I'm very conflicted on what would result in better performance for my service and the one I request from.

#### Structure

I decided to go with a simple folder structure that mimics the name of the package. I.e. The **api** folder contains the **api** package. 

All the pure functions will be in the **fun** package. Error handling will be in the **log** package. Testing will be in the **testing** folder where each package gets a testing file (Only testing the **fun** package for now).

#### Error handling
###### Location:    '.../projectfolder/log/...'
###### package:     'log'

I decided to handle errors by sending 4 variables to help with debugging. This error message will be in the form of a JSON structure and will be sent to the client and the console.
```go
type ErrorMessage struct {
    StatusCode       int    `json:"status_code"`
    Location         string `json:"location"`
    RawError         string `json:"raw_error"`
    PossibleReason   string `json:"possible_reason"`
}
```
The status code will be sent to both the header and the structure for ease of debugging. The 'location' tells us which function this happened. 'raw_error' is just an error message that explains what failed. 'possible_reason' is used when the status code is 400 and is used to showcase proper formating of URL.
Example:
```go
//branch if there is an error
if err != nil {
    log.UpdateErrorMessage(
        http.StatusBadRequest, 
        "HandlerExchangeHistory() -> handlerCountryCurrency() -> Getting base currency from requested country in URL",
        err.Error(),
        "Not valid country. Expected format: '.../country/...'. Example: '.../norway/...'.",
    )
    log.PrintErrorInformation(w)
    return
}
```
All of my error handlings will happen in my custom package 'log'.

#### Testing
###### Location:    '.../projectfolder/testing/...'
###### Coverage:    'fun_test'

Not sure how I could do any testing on API's so I decided to only test on pure functions, which turned out to only be 1 (which was relevant). All the testing will be done on the **fun** package since that's where the pure functions go.

##### How to

For Visual Studio Code with Golang extension:
1. Open testing file in IDE (In the testing folder)
2. Click the ```run test``` label for any function that you want to test
