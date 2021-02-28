# [Assignment 1](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Assignment-1)

### Info
- Author: Sindre Eiklid (sindreik@stud.ntnu.no)
    - While the submission is individual, I have discussed the tasks with Rickard Loland and Susanne Edvardsen. We have also helped each other with problems that occurred during development ([rubber ducking](https://en.wikipedia.org/wiki/Rubber_duck_debugging) mostly).
- Root path: https://sindre-assignment-1.herokuapp.com/exchange/v1/
- I have used these REST web services to build my service:
    - https://exchangeratesapi.io/
    - https://restcountries.eu/
### Usage

There are 3 endpoints that you can append to the root path.

1. Exchange rate history of a currency for a given country
    ```
    Method: GET
    Path: .../exchangehistory/{:country_name}/{:begin_date-end_date}
    ```
    - {:country_name} refers to the English name for the country as supported by https://restcountries.eu/.
    - {:begin_date-end_date} indicates the start date (i.e., the earliest date to be reported) of the exchange rate and the end date (i.e., the latest date of the range) of the period over which exchange rates are reported. Date format: YYYY-MM-DD
        ```
        Example: https://sindre-assignment-1.herokuapp.com/exchange/v1/exchangehistory/norway/2020-01-01-2020-02-01
        ```

2. Current exchange rates for bordering countries
    ```
    Method: GET
    Path: .../exchangeborder/{:country_name}{?limit={:number}}
    ```
    - {:country_name} refers to the English name for the country as supported by https://restcountries.eu/.
    - {?limit={:number}} is an optional parameter that limits the number of bordering countries for which currencies are reported.
        ```
        Example: https://sindre-assignment-1.herokuapp.com/exchange/v1/exchangeborder/norway?limit=2
        ```

3. Diagnostics interface
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

The REST service, https://exchangeratesapi.io, allows for symbol filtering but gives an error if the given symbol doesn't exist in their database. While I could get better performance by filtering currencies when getting all the bordering currencies, I've decided to get all rates, then filter myself. That way, the user can request all the bordering currencies of Russia and get all of the available currencies instead of just an error.

The REST service, https://restcountries.eu/, allows for alpha code filtering, which is beneficial when getting all the bordering currencies but doesn't seem to allow for field filtering and code filtering at the same time. That means that I have to decide on either multiple requests per bordering country or one for each bordering country which would store a massive amount of data. I decided on the latter since requests are expensive for the REST service, but if you would like to check out my previous solution, you can go [here](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021-workspace/sindre0830/assignment-1/-/blob/9d4bf54371ce811aae325cad24a74dc5c549d641/api/countryCurrency.go).

In exchangeHistory I decided to remove my error handling if the output is empty. This is because it is either an invalid date (I.e. 2020-40-01) or the country doesn't exist in the exchange rate database (I.e. Mongolia). While I did error handling before, I decided that it's better to receive an empty struct. This makes my service more dynamic and easier to use in my opinion.

#### Structure

I decided to go with a simple folder structure that mimics the name of the package. I.e. The **api** folder contains the **api** package. 

All the pure functions will be in the **fun** package. Error handling will be in the **debug** package. Testing will be in the **testing** folder where each package gets a testing file (Only testing the **fun** package for now).

#### Error handling
```
Location:    '.../projectfolder/debug/...'
package:     'debug'
````

I decided to handle errors by sending 4 variables to help with debugging. This error message will be in the form of a JSON structure and will be sent to the client and the console.
```go
type ErrorMessage struct {
    StatusCode       int    `json:"status_code"`
    Location         string `json:"location"`
    RawError         string `json:"raw_error"`
    PossibleReason   string `json:"possible_reason"`
}
```
The status code will be sent to both the header and the structure for ease of debugging. The 'location' tells us which function this error occured. 'raw_error' is just an error message that explains what failed. 'possible_reason' is used to give a possible reason for this error.

Example:
```go
//branch if there is an error
if err != nil {
    debug.UpdateErrorMessage(
        http.StatusBadRequest, 
        "HandlerExchangeHistory() -> handlerCountryCurrency() -> Getting base currency from requested country in URL",
        err.Error(),
        "Not valid country. Expected format: '.../country/...'. Example: '.../norway/...'.",
    )
    debug.PrintErrorInformation(w)
    return
}
```

The 'possible_reason' variable might not always be correct. I.e. if my service is timed out during a request, it will still show the same 'possible_reason' as when the country name is incorrect. While I could remove the entire variable and not have this mistake, I've decided to keep it since it's helpful for the client.
Another solution to this would be to keep the status code of each request and check if 'possible_reason' is needed (I.e. status code 500 from REST service would not print 'possible_reason'), but I didn't have time to implement this. 

#### Testing
```
Location:    '.../projectfolder/testing/...'
Coverage:    'fun_test'
```

Not sure how I could do any testing on API's so I decided to only test on pure functions, which turned out to only be 1 (which was relevant). All the testing will be done on the **fun** package since that's where the pure functions go.

##### Usage

For Visual Studio Code with Golang extension:
1. Open testing file in IDE (In the testing folder)
2. Click the ```run test``` label for any function that you want to test
