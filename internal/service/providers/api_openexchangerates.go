package providers

import (
	"encoding/json"
	"fmt"

	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	util "fx-service/pkg/helpers"
	"net/http"
	"strings"
)

/**
  Open Exchange Rates API provider

  Website:
  https://openexchangerates.org

  Documentation URL:
  https://docs.openexchangerates.org/reference

  Important info:
  ---------------
  On the free account, all rates returned are in terms of USD
  Therefore, to get the rate between two currencies, we divide the "to" rate by the "from" rate.

  We always request the "from" rate to be returned along with the others, in terms of USD.
  In the case that we're getting multiple rates, we request the "from" rate to be returned for each of the currencies.

*/

type OpenExchangeRates struct {
	Name                string
	AppID               string
	Timeout             int
	supportedCurrencies []string
}

type OpenExchangeRatesListResult map[string]string

type OpenExchangeRateErrorResponse struct {
	Error       bool   `json:"error"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type OpenExchangeRatesResult struct {
	Timestamp int64    `json:"timestamp"`
	Base      string   `json:"base"`
	Rates     RateList `json:"rates"`
}

const openExchangeRatesBaseURL = "https://openexchangerates.org/api"
const openExchangeRatesList = "/currencies.json?show_alternative=false&show_inactive=false"
const openExchangeRatesLatest = "/latest.json?app_id=%s&symbols=%s&show_alternative=false"

func (api *OpenExchangeRates) updateSupportedCurrencies() *e.Exception {
	ef := e.Fields{"api": api.Name}

	// Make the request and validate the response
	url := openExchangeRatesBaseURL + openExchangeRatesList
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return e.FromError(err).SetFields(ef.With("status", status))
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf(api.Name+" got non-200 response code: %d", status)
		return e.Throw(errNon200, msg).SetFields(ef.With("status", status))
	}

	// Parse the response into our predefined structure
	var response OpenExchangeRatesListResult
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return e.FromError(err).SetFields(ef)
	}

	// Ensure the response was successful
	if len(response) == 0 {
		return e.Throw(errNoResult, "response does not contain any symbols").SetFields(ef)
	}

	// Parse the result into our preferred slice of string
	api.supportedCurrencies = util.GetMapKeys(response)

	c.Infof("Provider '%s' supports %v currencies", api.Name, len(api.supportedCurrencies))

	return nil
}

func (api *OpenExchangeRates) doRequest(url string) (*OpenExchangeRatesResult, *e.Exception) {
	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return nil, e.FromError(err).SetFields(e.Fields{"url": url})
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf("Got non-200 response code: %d", status)
		return nil, e.Throw(errNon200, msg).SetFields(e.Fields{"url": url})
	}

	// Parse the response into our predefined structure
	var response OpenExchangeRatesResult
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return nil, e.Throw(errNotJson, "Could not unmarshal response body").SetFields(e.Fields{"url": url})
	}

	return &response, nil
}

func (api *OpenExchangeRates) CheckApiKey() bool {
	if api.AppID == "" {
		return false
	}

	// Grab the list of supported currencies from the provider
	err := api.updateSupportedCurrencies()
	if err != nil {
		err.Print(0, 0)
		return false
	}

	return true
}

func (api *OpenExchangeRates) GetName() string {
	return api.Name
}

func (api *OpenExchangeRates) GetRate(from, to string) (float64, error) {
	// set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// Make the URL
	url := fmt.Sprintf(
		openExchangeRatesBaseURL+openExchangeRatesLatest,
		api.AppID,
		strings.Join([]string{from, to}, ","),
	)

	// Make the request and validate the response
	response, err := api.doRequest(url)
	if err != nil {
		return 0, e.FromError(err).SetFields(ef)
	}

	// Ensure the response was successful
	// We should have the from and to rate in the response
	// The base currency is assumed as USD for the free account
	if response.Rates[from] == 0 {
		return 0, e.Throw(errNoResult, "response does not contain the from rate").SetFields(ef)
	}
	if response.Rates[to] == 0 {
		return 0, e.Throw(errNoResult, "response does not contain the to rate").SetFields(ef)
	}

	// Calculate the rate
	// Both the "from" and "to" rates, in the result are in terms of USD
	// Therefore to get the rate between the two currencies, we divide the "to" rate by the "from" rate
	actualRate := response.Rates[to] / response.Rates[from]

	return actualRate, nil
}

func (api *OpenExchangeRates) GetRates(from string, to []string) (RateList, error) {
	//set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// need to add "from" to the end of "to" list
	// because all rates in the result come in terms of USD, for this provider (on the free account)
	quotesToFetch := append(to, from)
	url := fmt.Sprintf(
		openExchangeRatesBaseURL+openExchangeRatesLatest,
		api.AppID,
		strings.Join(quotesToFetch, ","),
	)

	// Make the request and validate the response
	response, err := api.doRequest(url)
	if err != nil {
		return nil, e.FromError(err).SetFields(ef)
	}

	// Ensure response.Rates contains all the requested currencies
	for _, next := range quotesToFetch {
		if response.Rates[next] == 0 {
			return nil, e.Throw(errNoResult, "response does not contain rate for "+next).SetFields(ef)
		}
	}

	// Calculate the rates, on the basis that the "from" currency is USD
	// and our target "from" currency exists in the response.Rates

	result := make(RateList)
	for _, next := range to {
		if next == from {
			// The result will not want to have the "from" currency in it
			continue
		}
		result[next] = response.Rates[next] / response.Rates[from]
	}

	return result, nil
}

func (api *OpenExchangeRates) Supports(currency string) bool {
	return util.SliceContains(api.supportedCurrencies, currency)
}
