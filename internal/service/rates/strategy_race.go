package rates

import (
	"context"
	"fmt"
	"fx-service/internal/service/providers"
	c "fx-service/pkg/console"
	"github.com/gofiber/fiber/v2/log"
	"sync"
)

// callProviderRace calls all healthy providers at the same time;
// it waits for the first successful response and cancels other goroutines,
// or returns an error if all providers fail.
func callProviderRace(from string, to interface{}, isMulti bool) (interface{}, *string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.Outf("Race mode: %s -> %v", from, to)

	var (
		wg          sync.WaitGroup
		successChan = make(chan struct {
			result interface{}
			name   string
		})
		errorChan = make(chan error)
		once      sync.Once
	)

	for name, provider := range providers.EnabledProviders {
		wg.Add(1)
		go func(ctx context.Context, name string, provider providers.ProviderInterface) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				c.Out("Context cancelled")
				return
			default:
				c.Outf("Race is calling provider: %s", name)
				result, err := callProvider(provider, from, to, isMulti)
				if err == nil {
					once.Do(func() {
						successChan <- struct {
							result interface{}
							name   string
						}{result, name}
						cancel()
					})
				} else {
					errorChan <- err
					log.Warnf("Provider %s failed: %v\n", name, err)
				}
			}
		}(ctx, name, provider)
	}

	go func() {
		wg.Wait()
		close(successChan)
		close(errorChan)
	}()

	var collectedErrors []error
	for {
		select {
		case success := <-successChan:
			if success.result != nil {
				return success.result, &success.name, nil
			}
		case err := <-errorChan:
			collectedErrors = append(collectedErrors, err)
			if len(collectedErrors) == len(providers.EnabledProviders) {
				return nil, nil, fmt.Errorf("all providers failed: %v", collectedErrors)
			}
		}
	}
}
