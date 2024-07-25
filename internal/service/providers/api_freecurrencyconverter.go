package providers

import (
	"encoding/json"
	"fmt"
	util "fx-service/pkg/helpers"
	"io"
	"net/http"
	"strings"
)

// TODO - finish this provider

/**
 * FreeCurrencyConverterAPI is a provider for free.currconv.com
 * Website: https://www.currencyconverterapi.com/
 */

type FreeCurrencyConverterAPI struct {
	Name                string
	APIKey              string
	Timeout             int
	supportedCurrencies []string // TODO
}

const freeCurrencyConverterAPIBaseURL = "https://free.currconv.com/api/v7/convert?q=%s&compact=ultra&apiKey=%s"

func (api *FreeCurrencyConverterAPI) CheckApiKey() bool {
	if api.APIKey == "" {
		return false
	}

	// TODO
	api.supportedCurrencies = []string{"USD", "EUR", "GBP", "JPY", "AUD", "CAD"}
	return false
}

func (api *FreeCurrencyConverterAPI) GetName() string {
	return api.Name
}

func (api *FreeCurrencyConverterAPI) GetRate(from, to string) (float64, error) {
	query := fmt.Sprintf("%s_%s", from, to)
	url := fmt.Sprintf(freeCurrencyConverterAPIBaseURL, query, api.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	if rate, ok := result[query].(float64); ok {
		return rate, nil
	}
	return 0, fmt.Errorf("unsupported API response format")
}

func (api *FreeCurrencyConverterAPI) GetRates(from string, to []string) (RateList, error) {
	query := strings.Join(to, fmt.Sprintf("_%s,", from)) + "_" + from
	url := fmt.Sprintf(freeCurrencyConverterAPIBaseURL, query, api.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	rates := RateList{}
	for _, currency := range to {
		key := fmt.Sprintf("%s_%s", from, currency)
		if rate, ok := result[key].(float64); ok {
			rates[currency] = rate
		} else {
			return nil, fmt.Errorf("unsupported API response format for currency: %s", currency)
		}
	}
	return rates, nil
}

func (api *FreeCurrencyConverterAPI) Supports(currency string) bool {
	return util.SliceContains(api.supportedCurrencies, currency)
}
