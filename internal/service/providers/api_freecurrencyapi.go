package providers

import (
	"encoding/json"
	"fmt"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	"net/http"
	"strings"
)

/**
  Free Currency API

  Website:
  https://app.freecurrencyapi.com

  Base URL:
  https://api.freecurrencyapi.com/v1
*/

type FreeCurrencyApi struct {
	Name                string
	APIKey              string
	Timeout             int
	SupportedCurrencies []string
}

type freeCurrencyApiStatus struct {
	AccountId int64 `json:"account_id"`
}

type freeCurrencyApiResponse struct {
	Data RateList
}

type freeCurrencyApiListResponse struct {
	Data map[string]struct {
		symbol string
	}
}

const freeCurrencyApiBaseURL = "https://api.freecurrencyapi.com/v1"
const freeCurrencyApiStatusEndpoint = "/status?apikey=%s"
const freeCurrencyApiLatestEndpoint = "/latest?apikey=%s&base_currency=%s&currencies=%s"
const freeCurrencyApiListEndpoint = "/currencies?apikey=%s"

func (api *FreeCurrencyApi) updateSupportedCurrencies() error {
	url := fmt.Sprintf(freeCurrencyApiBaseURL+freeCurrencyApiListEndpoint, api.APIKey)

	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return e.FromError(err)
	}
	if status != http.StatusOK {
		return e.Throwf(errNon200, "got non-200 response code: %d", status)
	}

	// Parse the response into our predefined structure
	var response freeCurrencyApiListResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		// Parsing the response should not fail, so we log the error
		return e.FromError(err)
	}

	// Check if the response was successful
	if response.Data == nil || len(response.Data) == 0 {
		return e.Throw(errNoResult, "response does not contain symbols group")
	}

	// Extract the supported currencies from the response
	api.SupportedCurrencies = make([]string, 0, len(response.Data))
	for _, symbol := range response.Data {
		api.SupportedCurrencies = append(api.SupportedCurrencies, symbol.symbol)
	}

	c.Infof("Provider '%s' supports %v currencies", api.Name, len(api.SupportedCurrencies))

	return nil
}

func (api *FreeCurrencyApi) CheckApiKey() bool {
	if api.APIKey == "" {
		return false
	}

	// Format the URL for the get request
	url := fmt.Sprintf(freeCurrencyApiBaseURL+freeCurrencyApiStatusEndpoint, api.APIKey)

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		// Making the request should not fail, so we log the error
		e.FromError(err).Print(0, 0)
		return false
	}
	if status != http.StatusOK {
		c.Warnf("FreeCurrencyApi status check failed. Got non-200 response code: %d", status)
		return false
	}

	// Parse the response into our predefined structure
	var response freeCurrencyApiStatus
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		// Parsing the response should not fail, so we log the error
		e.FromError(err).Print(0, 0)
		return false
	}

	// Check if the account ID is set
	if response.AccountId == 0 {
		c.Warn("FreeCurrencyApi status check failed. Account ID is not set.")
		return false
	}

	// Update supported currencies
	err = api.updateSupportedCurrencies()
	if err != nil {
		e.FromError(err).Print(0, 0)
		return false
	}

	return true
}

func (api *FreeCurrencyApi) GetName() string {
	return api.Name
}

func (api *FreeCurrencyApi) GetRate(from, to string) (float64, error) {
	// set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// Format the URL for the get request
	url := fmt.Sprintf(freeCurrencyApiBaseURL+freeCurrencyApiLatestEndpoint, api.APIKey, from, to)

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return 0, e.FromError(err).SetFields(ef)
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf(api.Name+" got non-200 response code: %d", status)
		return 0, e.Throw(errNon200, msg).SetFields(ef)
	}

	// Parse the response into our predefined structure
	var response freeCurrencyApiResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		c.Warnf("Failed to parse response JSON from %s API", api.Name)
		return 0, e.FromError(err).SetFields(ef)
	}

	result, exists := response.Data[to]
	if !exists || result == 0 {
		// We got a successful response, but didn't find the data
		return 0, e.Throw(errNoResult, "unsupported API response format, does not contain 'to' currency").SetFields(ef)
	}

	return result, nil
}

func (api *FreeCurrencyApi) GetRates(from string, to []string) (RateList, error) {
	// set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// Format the URL for the get request
	url := fmt.Sprintf(freeCurrencyApiBaseURL+freeCurrencyApiLatestEndpoint, api.APIKey, from, strings.Join(to, ","))

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf(api.Name+" got non-200 response code: %d", status)
		return nil, e.Throw(errNon200, msg).SetFields(ef)
	}

	// Parse the response into our predefined structure
	var response freeCurrencyApiResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		c.Warnf("Failed to parse response JSON from %s API", api.Name)
		return nil, e.FromError(err).SetFields(ef)
	}

	if response.Data == nil || len(response.Data) == 0 {
		// We got a successful response, but didn't find the data
		return nil, e.Throw(errNoResult, "unsupported API response format").SetFields(ef)
	}

	result := RateList{}
	for _, currency := range to {
		if rate, ok := response.Data[currency]; ok {
			result[currency] = rate
		} else {
			return nil, e.Throw(errNoResult, "API response is missing currency: "+currency).SetFields(ef)
		}
	}

	return result, nil
}

func (api *FreeCurrencyApi) Supports(currency string) bool {
	for _, next := range api.SupportedCurrencies {
		if next == currency {
			return true
		}
	}
	return false
}
