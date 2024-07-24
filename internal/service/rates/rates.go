package rates

import (
	"fx-service/internal/service/providers"
	"fx-service/internal/service/ratecache"
	"fx-service/pkg/config"
	"fx-service/pkg/e"
	util "fx-service/pkg/helpers"
)

type RateGetter interface {
	GetRate(from, to string, mode config.Mode) (float64, bool, error)
	GetRates(from string, to []string, mode config.Mode) (map[string]float64, error)
}

type GetRateResult struct {
	Base      string
	Quote     string
	Rate      float64
	WasCached bool
	Provider  *string
}

type GetRatesResult struct {
	Base      string
	Quotes    []string
	Rates     providers.RateList
	WasCached bool
	Provider  *string
}

// GetRate obtains the rate for the given currency pair.
// Returns the rate, a boolean indicating if the rate was found in the cache, or an error.
func GetRate(from, to string, mode config.Mode) (*GetRateResult, error) {
	result := GetRateResult{
		Base:  from,
		Quote: to,
	}
	// Check if we have the rate in the cache
	if rate := ratecache.GetInstance().Get(from, to); rate != nil {
		result.Rate = *rate
		result.WasCached = true
		return &result, nil
	}

	// Get the rate from the provider
	rate, providerName, err := runAPIStrategy(from, to, mode, false)
	if err != nil {
		return nil, err
	}

	// For single rates, the result from the API calling strategy is a float64
	rateF64 := rate.(float64)
	result.Rate = rateF64
	result.Provider = providerName

	// Update the cache, asynchronously
	defer func() {
		go func() {
			ratecache.GetInstance().Set(from, to, rateF64)
		}()
	}()

	return &result, nil
}

// GetRates obtains multiple quotes for the given currency rate
func GetRates(from string, toList []string, mode config.Mode) (*GetRatesResult, error) {
	var ratesToGet []string
	result := GetRatesResult{
		Base:   from,
		Quotes: toList,
		Rates:  make(providers.RateList),
	}

	// Check which combinations we have in the cache
	cache := ratecache.GetInstance()
	for _, toCurrency := range toList {
		if rate := cache.Get(from, toCurrency); rate != nil {
			// Found it in the cache
			result.Rates[toCurrency] = *rate
		} else {
			// Not in the cache - we'll need to get this from the API provider(s)
			ratesToGet = append(ratesToGet, toCurrency)
		}
	}

	// If we have all the rates in the cache, return the result
	if len(ratesToGet) == 0 {
		result.WasCached = true
		return &result, nil
	}

	// Get the rates from the provider, for the ones we don't have in the cache
	strategyResult, providerName, err := runAPIStrategy(from, ratesToGet, mode, true)
	if err != nil {
		return nil, err
	}

	// Ensure strategyResult is a RateList type (should always be the case)
	apiRatesResult, ok := strategyResult.(providers.RateList)
	if !ok {
		// Sanity check - this should never happen
		actualType := util.GetType(strategyResult)
		return nil, e.Throw("eRgr71", "runApiStrategy returned an invalid type. Expected 'RateList'; got: "+actualType)
	}

	// Update the cache asynchronously
	defer func() {
		go func() {
			rc := ratecache.GetInstance()
			for currency, newRate := range apiRatesResult {
				rc.Set(from, currency, newRate)
			}
		}()
	}()

	// Combine the rates we just got from the API provider with the ones we already had in the cache
	for currency, rate := range apiRatesResult {
		result.Rates[currency] = rate
	}

	result.Provider = providerName
	return &result, nil
}
