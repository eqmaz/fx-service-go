package rates

import (
	"fx-service/internal/service/providers"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
)

// callProviderFirst calls the first healthy provider that is available
func callProviderFirst(from string, to interface{}, isMulti bool) (interface{}, *string, error) {
	count := 0
	for name, provider := range providers.EnabledProviders {
		count++
		result, err := callProvider(provider, from, to, isMulti)
		if err != nil {
			c.Warnf("Provider '%s' failed: %v", name, err.Error())
			e.FromError(err).SetField("strategy", "first").Print(0, 0)
			continue
		}

		// From here, we have a successful result
		// if more than 1 provider was tried, log it
		if count > 1 {
			// TODO use preferred logger
			c.Warnf("First mode succeeded, but tried more than 1 provider (%v)", count)
		}
		return result, &name, nil

	}
	return nil, nil, e.Throwf("eCpf34", "all providers failed, tried %v enabled providers", count)
}
