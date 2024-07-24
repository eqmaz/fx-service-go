package rates

import (
	"errors"
	"fx-service/internal/service/providers"
	util "fx-service/pkg/helpers"
	"github.com/gofiber/fiber/v2/log"
)

// callProviderRandom calls a random provider that is available and healthy
// returns the rate result, provider name, or an error
func callProviderRandom(from string, to interface{}, isMulti bool) (interface{}, *string, error) {
	providersTried := make(map[string]bool)
	providersNotTried := util.GetMapKeys(providers.EnabledProviders)

	for len(providersNotTried) > 0 {
		// get a random provider
		nextIndex := util.GetRandomSliceIndex(providersNotTried)
		providerName := providersNotTried[nextIndex]
		util.RemoveSliceElement(providersNotTried, nextIndex)
		provider := providers.EnabledProviders[providerName]
		providersTried[providerName] = true

		result, err := callProvider(provider, from, to, isMulti)
		if err == nil {
			// if more than 1 provider was tried, log it
			if len(providersTried) > 1 {
				log.Warn("Random mode tried providers: ", providersTried) // todo
			}

			return result, &providerName, nil
		}
	}

	// If we get here, we've exhausted all providers available
	return nil, nil, errors.New("all providers failed")
}
