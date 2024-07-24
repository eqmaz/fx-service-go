package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TODO add documentation at top of file

type FixerApi struct {
	Name      string
	AccessKey string
	Timeout   int
}

const fixerBaseURL = "https://data.fixer.io/api/latest?access_key=%s&symbols=%s"

var fixerSupportedCurrencies = []string{"USD", "EUR", "GBP", "JPY", "AUD", "CAD"}

func (api *FixerApi) CheckApiKey() bool {
	return api.AccessKey != ""
}

func (api *FixerApi) GetName() string {
	return api.Name
}

func (api *FixerApi) GetRate(from, to string) (float64, error) {
	url := fmt.Sprintf(fixerBaseURL, api.AccessKey, to)
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

	if rates, ok := result["rates"].(map[string]interface{}); ok {
		if rate, ok := rates[to].(float64); ok {
			return rate, nil
		}
	}
	return 0, fmt.Errorf("unsupported API response format")
}

func (api *FixerApi) GetRates(from string, to []string) (RateList, error) {
	url := fmt.Sprintf(fixerBaseURL, api.AccessKey, strings.Join(to, ","))
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
	if ratesMap, ok := result["rates"].(map[string]interface{}); ok {
		for _, currency := range to {
			if rate, ok := ratesMap[currency].(float64); ok {
				rates[currency] = rate
			} else {
				return nil, fmt.Errorf("unsupported API response format for currency: %s", currency)
			}
		}
	}
	return rates, nil
}

func (api *FixerApi) Supports(currency string) bool {
	for _, c := range fixerSupportedCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}
