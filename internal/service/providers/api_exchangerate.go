package providers

import (
	"encoding/json"
	"fmt"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	"net/http"
)

/**
  ExchangeRate-API service

  Documentation here:
  https://www.exchangerate-api.com/docs

  Terms of use here:
  "https://www.exchangerate-api.com/terms",
*/

const exchangeRateAPIBaseURL = "https://v6.exchangerate-api.com/v6/%s"
const exchangeRateAPISingle = "/pair/%s/%s/1"
const exchangeRateAPIMulti = "/latest/%s"
const exchangeRateAPIList = "/codes"

type ExchangeRateApi struct {
	Name                string
	APIKey              string
	Timeout             int
	SupportedCurrencies []string
}

type exchangeRateAPIResponse struct {
	Result          string     `json:"result"`                     // "success" field in the response
	ErrorType       string     `json:"error-type"`                 // "error-type" field in the response
	BaseCode        string     `json:"base_code"`                  // "source" field in the response
	TargetCode      string     `json:"target_code"`                // "target" field in the response
	ConversionRate  float64    `json:"conversion_rate"`            // "rate" field in the response
	ConversionRates RateList   `json:"conversion_rates,omitempty"` // "rates" field in the response is a map of currency codes to rates
	SupportedCodes  [][]string `json:"supported_codes,omitempty"`  // "supported_codes" field in the response (array of [code, name])
}

func (api *ExchangeRateApi) getSupportedCurrencies() error {
	ef := e.Fields{"api": api.Name}

	// Make request to get the list of supported currencies
	url := fmt.Sprintf(exchangeRateAPIBaseURL+exchangeRateAPIList, api.APIKey)
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return e.FromError(err).SetFields(ef)
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf("ExchangeRate-API got non-200 response code: %d", status)
		return e.Throw(errNon200, msg).SetFields(ef.With("status", status))
	}

	// Parse the response
	var response exchangeRateAPIResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return e.FromError(err).SetFields(ef)
	}

	// Check if the response was successful - "Result" usually comes with a "success" value
	if response.Result != "success" {
		return e.Throw(response.ErrorType, "ExchangeRate-API response was not successful").SetFields(ef)
	}

	// Check if the response contains supported codes
	if len(response.SupportedCodes) == 0 {
		return e.Throw(errNoResult, "ExchangeRate-API response did not contain supported codes").SetFields(ef)
	}

	// Extract the supported currencies from the response
	api.SupportedCurrencies = make([]string, 0, len(response.SupportedCodes))
	for _, code := range response.SupportedCodes {
		if len(code) == 0 {
			// Should never happen
			continue
		}
		supportedCode := code[0]
		api.SupportedCurrencies = append(api.SupportedCurrencies, supportedCode)
	}

	// Check length of supported currencies
	if len(api.SupportedCurrencies) == 0 {
		// Should never happen
		return e.Throw(errNoResult, "ExchangeRate-API response did not contain supported codes").SetFields(ef)
	}

	c.Infof("Provider '%s' supports %v currencies", api.Name, len(api.SupportedCurrencies))

	return nil
}

func (api *ExchangeRateApi) CheckApiKey() bool {
	if api.APIKey == "" {
		c.Warn(api.Name + " API key is not set")
		return false
	}

	err := api.getSupportedCurrencies()
	if err != nil {
		c.Warnf("Failed to update supported currencies for provider '%s': %s", api.Name, err)
		e.FromError(err).Print(0, 0)
		return false
	}

	return true
}

func (api *ExchangeRateApi) GetName() string {
	return api.Name
}

func (api *ExchangeRateApi) GetRate(from, to string) (float64, error) {
	// Error fields, for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// Build the URL
	url := fmt.Sprintf(exchangeRateAPIBaseURL+exchangeRateAPISingle, api.APIKey, from, to)

	// Make the request
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return 0, err
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf("ExchangeRate-API got non-200 response code: %d", status)
		return 0, e.Throw(errNon200, msg).SetFields(ef.With("status", status))
	}

	// Parse the API response into a formal struct
	var response exchangeRateAPIResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return 0, err
	}

	// Check if the response was successful
	// "Result" usually comes with a "success" value
	if response.Result != "success" {
		return 0, e.Throw(response.ErrorType, "ExchangeRate-API response was not successful")
	}

	// Check conversion rate is set
	if response.ConversionRate == 0 {
		return 0, e.Throw("", "ExchangeRate-API Conversion rate is not set in response from API")
	}

	return response.ConversionRate, nil
}

func (api *ExchangeRateApi) GetRates(from string, to []string) (RateList, error) {
	// Error context fields, for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// Build the URL
	url := fmt.Sprintf(exchangeRateAPIBaseURL+exchangeRateAPIMulti, api.APIKey, from)

	// Make the request
	status, bodyData, err := makeGetRequest(url, api.Timeout, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf(api.Name+" got non-200 response code: %d", status)
		return nil, e.Throw(errNon200, msg).SetFields(ef.With("status", status))
	}

	// Parse the API response into a formal struct
	var response exchangeRateAPIResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return nil, err
	}

	// Check if the response was successful
	// "Result" usually comes with a "success" value
	if response.Result != "success" {
		msg := "ExchangeRate-API response was not successful"
		return nil, e.Throw(response.ErrorType, msg).SetFields(ef)
	}

	// Check ConversionRates is set
	if len(response.ConversionRates) == 0 {
		msg := "ExchangeRate-API Conversion rates are not set in response from API"
		return nil, e.Throw(errNoResult, msg).SetFields(ef)
	}

	// Try to build the RateList, from the unmarshalled response
	result := make(RateList)
	for _, currency := range to {
		if rate, ok := response.ConversionRates[currency]; ok {
			result[currency] = rate
		} else {
			return nil, e.Throwf("unsupported API response format for currency: %s", currency).SetFields(ef)
		}
	}

	return result, nil
}

func (api *ExchangeRateApi) Supports(currency string) bool {
	for _, next := range api.SupportedCurrencies {
		if next == currency {
			return true
		}
	}
	return false
}
