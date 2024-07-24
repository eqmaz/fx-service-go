package rates

import (
	"fx-service/internal/service/providers"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	util "fx-service/pkg/helpers"
	"github.com/gofiber/fiber/v2/log"
)

// aggregateSingleProvider aggregates results from all providers for a single currency conversion
func aggregateSingleResult(from string, toCurrency string) (interface{}, *string, error) {
	var (
		totalRate float64
		numRates  int
	)

	providerName := "Aggregate [all]"
	for _, provider := range providers.EnabledProviders {
		result, err := callProvider(provider, from, toCurrency, false)
		if err == nil {
			// Assuming result is a float64 for single currency rate
			rateF64, ok := result.(float64)
			if !ok {
				// Should never happen
				log.Warnf("Provider %s aggregate failed for %s -> %s: invalid result type\n", providerName, from, toCurrency)
				continue
			}
			totalRate += rateF64
			numRates++
		} else {
			log.Warnf("Provider failed for %s -> %s: %v\n", from, toCurrency, err)
		}
	}

	if numRates == 0 {
		return nil, nil, e.Throw("eSaSr30", "all providers failed")
	}

	c.Outf("GetRate - Averaged values from %d providers", numRates)

	// Calculate mean rate (Round to 8 decimal places)
	meanRate := util.Round(totalRate/float64(numRates), 8)

	return meanRate, &providerName, nil
}

// aggregateMultiProvider aggregates results from all providers for multiple currency conversions
func aggregateMultiResult(from string, toCurrencies []string) (interface{}, *string, error) {
	// Initialize map to accumulate rates for each "to" currency
	ratesMap := make(map[string]float64)
	countMap := make(map[string]int)

	providerName := "Aggregate [all]"
	for _, provider := range providers.EnabledProviders {
		result, err := callProvider(provider, from, toCurrencies, true)
		// TODO sanity check for result type
		if err == nil {
			// Assuming result is a map[string]float64 for multi currency rates
			rates := result.(providers.RateList)
			for currency, rate := range rates {
				ratesMap[currency] += rate
				countMap[currency]++
			}
		} else {
			// TODO - log.warn?
			log.Warnf("Provider failed for %s -> %v: %v\n", from, toCurrencies, err)
		}
	}

	if len(ratesMap) == 0 {
		return nil, nil, e.Throw("eSaMr61", "all providers failed")
	}

	c.Outf("GetRates - averaged values from %d providers", len(ratesMap))

	// Calculate mean rates for each "to" currency
	meanRates := make(providers.RateList)
	for currency, totalRate := range ratesMap {
		count := countMap[currency]
		meanRates[currency] = util.Round(totalRate/float64(count), 8)
	}

	return meanRates, &providerName, nil
}

// callAggregateProvider calls all healthy providers and aggregates the results
// Returns the result and provider name (Aggregate), or error
// For single from-to results, it returns the mean of the "to" rate
// For multi from-to results, it returns the mean of the "to" rates, for each "to" currency
func callProviderAggregate(from string, to interface{}, isMulti bool) (interface{}, *string, error) {
	if isMulti {
		return aggregateMultiResult(from, to.([]string))
	}
	return aggregateSingleResult(from, to.(string))
}
