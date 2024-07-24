package rates

import (
	"fx-service/internal/service/providers"
	"fx-service/pkg/config"
	"fx-service/pkg/e"
)

// callProvider calls GetRate, GetRates, or any other future method required on the specific provider
func callProvider(provider providers.ProviderInterface, from string, to interface{}, isMulti bool) (interface{}, error) {
	var result interface{}
	var err error
	if isMulti {
		// calling GetRates, for multi-currency result

		toList, ok := to.([]string) // Just make sure it's a slice of strings
		if !ok {
			// This should never happen, as the API should have validated the input
			return nil, e.Throw("eScp20", "invalid type for 'to' parameter, expected []string")
		}
		result, err = provider.GetRates(from, toList)
	} else {
		// calling GetRate, for single-currency result

		toCcy, ok := to.(string) // Just make sure it's a string
		if !ok {
			// This should never happen, as the API should have validated the input
			return nil, e.Throw("eScp30", "invalid type for 'to' parameter, expected string")
		}
		result, err = provider.GetRate(from, toCcy)
	}
	if err != nil {
		return nil, err // Just pass the error through, the handler will deal with it
	}

	return result, nil
}

// runAPIStrategy runs the API strategy based on the mode of calling API providers
func runAPIStrategy(from string, to interface{}, mode config.Mode, isMulti bool) (interface{}, *string, error) {
	switch mode {
	case config.Random:
		return callProviderRandom(from, to, isMulti)
	case config.Robin:
		return callProviderRoundRobin(from, to, isMulti)
	case config.Priority:
		return callPriorityOrder(from, to, isMulti)
	case config.Aggregate:
		return callProviderAggregate(from, to, isMulti)
	case config.Race:
		return callProviderRace(from, to, isMulti)
	default:
		return callProviderFirst(from, to, isMulti)
	}
}
