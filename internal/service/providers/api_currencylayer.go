package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	util "fx-service/pkg/helpers"
)

/**
  "Currency Data API" provided by API Layer (apilayer.com)
  "Reliable Exchange Rates & Currency Conversion for your Business."

  Documentation website:
  https://apilayer.com/marketplace/currency_data-api#endpoints
*/

type CurrencyLayer struct {
	Name                string
	AccessKey           string
	Timeout             int
	supportedCurrencies []string
}

type CurrencyLayerError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

type CurrencyLayerResponse struct {
	Success bool                `json:"success"`
	Error   *CurrencyLayerError `json:"error,omitempty"`
	Quotes  RateList            `json:"quotes,omitempty"`
	Symbols map[string]string   `json:"currencies,omitempty"` // FYI the documentation said "symbols".
}

const currencyLayerBaseURL = "https://api.apilayer.com/currency_data"
const currencyLayerLiveURL = "/live?source=%s&currencies=%s"
const currencyLayerListURL = "/list"

// getHeaders private helper to Create http headers to be sent with the request (includes API key)
func (api *CurrencyLayer) getHeaders() *map[string]string {
	return &map[string]string{"apikey": api.AccessKey}
}

// checkResponseError private helper to check the response shape for errors
func (api *CurrencyLayer) checkResponseError(response CurrencyLayerResponse, ef e.Fields) error {
	if !response.Success {
		if response.Error != nil {
			switch response.Error.Type {
			case "invalid_access_key":
				c.Outf("Invalid access key for CurrencyLayer API (%s)", api.AccessKey)
				return e.Throw(errApiKy, "invalid access key").SetFields(ef)
			default:
				return e.Throw(errUnhandled, response.Error.Info).SetFields(ef)
			}
		}
		return e.Throw(errUnhandled, "unknown error occurred in "+api.Name).SetFields(ef)
	}
	return nil
}

func (api *CurrencyLayer) updateSupportedCurrencies() error {
	ef := e.Fields{"api": api.Name}

	url := fmt.Sprintf(currencyLayerBaseURL + currencyLayerListURL)

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, api.getHeaders())
	if err != nil {
		return e.FromError(err).SetFields(ef.With("status", status))
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf(api.Name+" got non-200 response code: %d", status)
		return e.Throw(errNon200, msg).SetFields(e.Fields{"status": status, "url": url, "api": api.Name})
	}

	// Parse the response into our predefined structure
	var response CurrencyLayerResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return err
	}

	// Check if the response was successful
	err = api.checkResponseError(response, ef)
	if err != nil {
		return err
	}
	if response.Symbols == nil {
		return e.Throw(errNoResult, "response does not contain symbols group").SetFields(ef)
	}

	// Extract the supported currencies from the response
	api.supportedCurrencies = make([]string, 0, len(response.Symbols))
	for currency := range response.Symbols {
		//currency = strings.ToUpper(currency) // To ensure the currency is uppercase
		api.supportedCurrencies = append(api.supportedCurrencies, currency)
	}

	c.Infof("Provider '%s' supports %v currencies", api.Name, len(response.Symbols))

	return nil
}

func (api *CurrencyLayer) GetName() string {
	return api.Name
}

func (api *CurrencyLayer) CheckApiKey() bool {
	if api.AccessKey == "" {
		c.Warn(api.Name + " API key is not set")
		return false
	}

	err := api.updateSupportedCurrencies()
	if err != nil {
		c.Warnf("Failed to update supported currencies for provider '%s': %s", api.Name, err)
		e.FromError(err).Print(0, 0)
		return false
	}

	return true
}

func (api *CurrencyLayer) GetRate(from, to string) (float64, error) {
	// set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// Format the URL for the get request
	url := fmt.Sprintf(currencyLayerBaseURL+currencyLayerLiveURL, from, to)

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, api.getHeaders())
	if err != nil {
		return 0, e.FromError(err).SetFields(ef.With("status", status))
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf("CurrencyLayer API got non-200 response code: %d", status)
		return 0, e.Throw(errNon200, msg).SetFields(ef.With("status", status))
	}

	// Parse the response into our predefined structure
	var response CurrencyLayerResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return 0, err
	}

	// Check if the response was successful
	err = api.checkResponseError(response, ef)
	if err != nil {
		return 0, err
	}

	// With CurrencyLayer, the response looks like
	// {"success":true,"quotes":{"USDEUR":0.85,"USDGBP":0.75,"USDJPY":110.0,"USDAUD":1.35,"USDCAD":1.25}}
	rateKey := fmt.Sprintf("%s%s", from, to)
	rate, ok := response.Quotes[rateKey]
	if !ok {
		return 0, e.Throwf("eClGr144", "rate not found").SetFields(ef)
	}

	return rate, nil
}

func (api *CurrencyLayer) GetRates(from string, to []string) (RateList, error) {
	// set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to}

	// Format the URL for the get request
	url := fmt.Sprintf(currencyLayerBaseURL+currencyLayerLiveURL, from, strings.Join(to, ","))

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, api.getHeaders())
	if err != nil {
		return nil, e.FromError(err).SetFields(ef.With("status", status))
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf("CurrencyLayer API got non-200 response code: %d", status)
		return nil, e.Throw(errNon200, msg).SetFields(ef.With("status", status))
	}

	// Parse the response into our predefined structure
	var response CurrencyLayerResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return nil, e.FromError(err).SetFields(ef)
	}

	// Ensure the response was successful
	err = api.checkResponseError(response, ef)
	if err != nil {
		return nil, err
	}

	// Extract the rates from the response
	rates := make(RateList)
	for _, currency := range to {
		key := fmt.Sprintf("%s%s", from, currency)
		if rate, ok := response.Quotes[key]; ok {
			rates[currency] = rate
		} else {
			msg := fmt.Sprintf("Currency '%s' was not found in response from Currency Layer API", currency)
			return nil, e.Throw("eClGr133", msg).SetFields(ef)
		}
	}

	return rates, nil
}

func (api *CurrencyLayer) Supports(currency string) bool {
	return util.SliceContains(api.supportedCurrencies, currency)
}
