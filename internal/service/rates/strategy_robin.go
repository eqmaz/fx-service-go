package rates

import (
	"fx-service/internal/service/providers"
	c "fx-service/pkg/console"
	"fx-service/pkg/e"
	"sync"
)

const (
	errAllFailed = "eCRP68"
)

var (
	rrs     *roundRobinState
	rrsOnce sync.Once
)

// roundRobinState manages the state for round-robin provider calling
type roundRobinState struct {
	mu        sync.Mutex
	providers []providers.ProviderInterface
	nextIndex int
}

func getRobinState() *roundRobinState {
	rrsOnce.Do(func() {
		rrs = &roundRobinState{
			providers: make([]providers.ProviderInterface, 0, len(providers.EnabledProviders)),
		}
		// We need to convert the EnabledProviders map to a slice to ensure the order is consistent
		for _, provider := range providers.EnabledProviders {
			rrs.providers = append(rrs.providers, provider)
		}
	})
	return rrs
}

// callProviderRoundRobin calls the next healthy provider in a round-robin fashion.
// It locks the mutex to ensure thread safety when accessing shared state.
func callProviderRoundRobin(from string, to interface{}, isMulti bool) (interface{}, *string, error) {
	// Lazy initialization of providers slice if not already initialized
	rr := getRobinState()

	// Acquire lock to ensure exclusive access to shared state
	rr.mu.Lock()
	defer rr.mu.Unlock()

	// Iterate through the providers in a round-robin fashion
	// We stop when we reach a healthy provider or when we've tried all providers
	for i := 0; i < len(rr.providers); i++ {
		index := (rr.nextIndex + i) % len(rr.providers) // Calculate current provider index
		provider := rr.providers[index]                 // Select provider at the calculated index
		rr.nextIndex = (index + 1) % len(rr.providers)  // Update nextIndex for next iteration

		// Call the provider to fetch the result
		result, err := callProvider(provider, from, to, isMulti)
		if err == nil {
			// Return result if provider call is successful
			providerName := provider.GetName()
			return result, &providerName, nil
		} else {
			// TODO - get preferred logger in here somehow
			c.Warnf("Provider %T failed during round robin call: %v\n", provider, err)
			//log.Printf("Provider %T failed: %v\n", provider, err) // Log failure for the current provider
		}
	}

	return nil, nil, e.FromCode(errAllFailed)
}
