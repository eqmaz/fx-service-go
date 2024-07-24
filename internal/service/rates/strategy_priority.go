package rates

import (
	"fx-service/internal/service/providers"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	"sort"
	"sync"
)

type posState struct {
	providers []providers.ProviderInterface
}

var (
	pos     *posState
	posOnce sync.Once
)

// sortProvidersByPriority sorts providers by their priority order and returns the sorted slice
// Priority of 0 means no priority (they will go last, in a non-guaranteed order)
// When priorities are equal, the order between them is not guaranteed either.
// The highest priority is 1, the next highest is 2, and so on.
func sortProvidersByPriority() []providers.ProviderInterface {
	providerCount := len(providers.EnabledProviders)
	providersWithPriority := make([]struct {
		provider providers.ProviderInterface
		priority uint
	}, 0, providerCount)

	for name, provider := range providers.EnabledProviders {
		priority, exists := providers.ProviderPriority[name]
		if !exists {
			priority = 0
		}
		providersWithPriority = append(providersWithPriority, struct {
			provider providers.ProviderInterface
			priority uint
		}{
			provider: provider,
			priority: priority,
		})
	}

	// Sort the providers by priority (ascending)
	sort.SliceStable(providersWithPriority, func(i, j int) bool {
		return providersWithPriority[i].priority < providersWithPriority[j].priority
	})

	// Extract sorted providers
	sortedProviders := make([]providers.ProviderInterface, 0, providerCount)
	for _, pw := range providersWithPriority {
		sortedProviders = append(sortedProviders, pw.provider)
	}

	return sortedProviders
}

// getPosState returns the singleton priority order state
func getPosState() *posState {
	posOnce.Do(func() {
		pos = &posState{
			providers: sortProvidersByPriority(),
		}
	})
	return pos
}

// callPriorityOrder calls the providers in priority order until one returns a result
func callPriorityOrder(from string, to interface{}, isMulti bool) (interface{}, *string, error) {
	state := getPosState()

	// Iterate through the providers in priority order
	for i, provider := range state.providers {
		result, err := callProvider(provider, from, to, isMulti)
		if err == nil {
			c.Outf("Priority Order %d - Provider %s succeeded", i, provider.GetName())
			providerName := provider.GetName()
			return result, &providerName, nil
		}

		// Log the error and continue to the next provider
		// TODO - use logger, not the console output
		c.Warnf("Provider failed for %s -> %v: %v\n", from, to, err)
	}

	return nil, nil, e.FromCode("eGaPf1")
}
