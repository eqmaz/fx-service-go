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

// TODO add documentation at top of file

type FixerApi struct {
	Name                string
	AccessKey           string
	Timeout             int
	supportedCurrencies []string
}

type FixerApiError struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Info string `json:"info"`
}

type FixerApiResponse struct {
	Success bool              `json:"success"`
	Date    string            `json:"date,omitempty"`
	Error   *FixerApiError    `json:"error,omitempty"`
	Rates   RateList          `json:"rates,omitempty"`
	Symbols map[string]string `json:"symbols,omitempty"` // FYI the documentation said "symbols".
}

const fixerBaseUrl = "https://api.apilayer.com/fixer"
const fixerSymbolsUrl = "/symbols"
const fixerLatestUrl = "/latest?base=%s&symbols=%s"

// getHeaders private helper to Create http headers to be sent with the request (includes API key)
func (api *FixerApi) getHeaders() *map[string]string {
	return &map[string]string{"apikey": api.AccessKey}
}

// checkResponseError private helper to check the response shape for errors
func (api *FixerApi) checkResponseError(response FixerApiResponse, ef e.Fields) error {
	if !response.Success {
		util.Dump(response)
		if response.Error != nil {
			switch response.Error.Type {
			case "invalid_access_key":
				c.Warnf("Invalid access key (%s)", api.AccessKey)
				return e.Throw(errApiKy, "invalid access key").SetFields(ef)
			default:
				return e.Throw(errUnhandled, response.Error.Info).SetFields(ef)
			}
		}
		return e.Throw(errUnhandled, "unknown error occurred").SetFields(ef)
	}
	return nil
}

// updateSupportedCurrencies private helper to update the supported currencies list
func (api *FixerApi) updateSupportedCurrencies() error {
	url := fmt.Sprintf(fixerBaseUrl + fixerSymbolsUrl)

	ef := e.Fields{"api": api.Name, "url": url}

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, api.getHeaders())
	if err != nil {
		return e.FromError(err).SetFields(ef.With("status", status))
	}
	if status != http.StatusOK {
		msg := fmt.Sprintf(api.Name+" got non-200 response code: %d", status)
		return e.Throw(errNon200, msg).SetFields(ef.With("status", status))
	}

	// Parse the response into our predefined structure
	var response FixerApiResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return err
	}

	// Check if the response was successful
	ex := api.checkResponseError(response, ef)
	if ex != nil {
		return ex
	}
	if response.Symbols == nil {
		return e.Throw(errNoResult, "response does not contain symbols group").SetFields(ef)
	}

	// Extract the supported currencies from the response
	api.supportedCurrencies = make([]string, 0, len(response.Symbols))
	for currency := range response.Symbols {
		api.supportedCurrencies = append(api.supportedCurrencies, currency)
	}

	c.Infof("Provider '%s' supports %v currencies", api.Name, len(api.supportedCurrencies))

	return nil
}

func (api *FixerApi) CheckApiKey() bool {
	if api.AccessKey == "" {
		return false
	}

	err := api.updateSupportedCurrencies()
	if err != nil {
		e.FromError(err).Print(0, 0)
		return false
	}

	return true
}

func (api *FixerApi) GetName() string {
	return api.Name
}

func (api *FixerApi) GetRate(from, to string) (float64, error) {
	// Format the URL for the get request
	url := fmt.Sprintf(fixerBaseUrl+fixerLatestUrl, from, to)

	// set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to, "url": url}

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, api.getHeaders())
	if err != nil {
		// Some unknown issue with making the request
		return 0, e.FromError(err).SetFields(ef.With("status", status))
	}
	if status != http.StatusOK {
		// response code is not 200
		return 0, e.FromCode("eAGn2c", status).SetFields(ef.With("status", status))
	}

	// Parse the response into our predefined structure
	var response FixerApiResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return 0, e.FromError(err).SetFields(ef)
	}

	// Check if the response was successful
	err = api.checkResponseError(response, ef)
	if err != nil {
		return 0, err
	}

	rate, ok := response.Rates[to]
	if !ok {
		// Quote symbol not found in the response
		return 0, e.FromCode("ePrRnf").SetFields(ef.With("rates", response.Rates))
	}

	return rate, nil
}

func (api *FixerApi) GetRates(from string, to []string) (RateList, error) {
	// Format the URL for the get request
	url := fmt.Sprintf(currencyLayerBaseURL+currencyLayerLiveURL, from, strings.Join(to, ","))

	// set error fields for traceability
	ef := e.Fields{"api": api.Name, "from": from, "to": to, "url": url}

	// Make the request and validate the response
	status, bodyData, err := makeGetRequest(url, api.Timeout, api.getHeaders())
	if err != nil {
		return nil, e.FromError(err).SetFields(ef.With("status", status))
	}
	if status != http.StatusOK {
		return nil, e.FromCode("eAGn2c", status).SetFields(ef.With("status", status))
	}

	// Parse the response into our predefined structure
	var response FixerApiResponse
	err = json.Unmarshal(bodyData, &response)
	if err != nil {
		return nil, e.FromError(err).SetFields(ef)
	}
	// Ensure the response was successful
	err = api.checkResponseError(response, ef)
	if err != nil {
		return nil, err // It's an e.Exception
	}

	// Extract the rates from the response
	rates := make(RateList)
	util.Dump(response)
	for _, quoteCcy := range to {
		if rate, ok := response.Rates[quoteCcy]; ok {
			rates[quoteCcy] = rate
		} else {
			msg := fmt.Sprintf("Currency '%s' not found in response from API", quoteCcy)
			return nil, e.Throw("eFaG205", msg).SetFields(ef.With("missingQuote", quoteCcy))
		}
	}

	return rates, nil
}

func (api *FixerApi) Supports(currency string) bool {
	return util.SliceContains(api.supportedCurrencies, currency)
}
