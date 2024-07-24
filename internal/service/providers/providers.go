package providers

import (
	"fx-service/pkg/config"
	c "fx-service/pkg/console"
	"os"
	"sync"
)

const (
	errUnhandled = "unhandled"
	errApiKy     = "apiKey"
	errNon200    = "non200"
	errNoResult  = "noResult"
	errNotJson   = "notJson"
)

type RateList map[string]float64

type ProviderInterface interface {
	CheckApiKey() bool
	GetName() string
	GetRate(from, to string) (float64, error)
	GetRates(from string, to []string) (RateList, error)
	Supports(currency string) bool
}

// InstalledProviders is a map of provider names to their structs. These are the available providers
// Do not confuse with EnabledProviders, which are the providers that are enabled and ready to use
var InstalledProviders = map[string]ProviderInterface{
	"CurrencyLayer":            &CurrencyLayer{},
	"ExchangeRateAPI":          &ExchangeRateApi{},
	"FixerApi":                 &FixerApi{},
	"FreeCurrencyApi":          &FreeCurrencyApi{},
	"OpenExchangeRates":        &OpenExchangeRates{},
	"FreeCurrencyConverterAPI": &FreeCurrencyConverterAPI{},
}

// EnabledProviders is a map of enabled providers. These have been initialized and are ready to use
var EnabledProviders = make(map[string]ProviderInterface)

// providerConstructors is a map of provider names to their constructor functions
// For use in initializing the providers only
var providerConstructors = map[string]func(apiKey string, timeout int) ProviderInterface{
	"CurrencyLayer":            NewCurrencyLayer,
	"ExchangeRateApi":          NewExchangeRateApi,
	"FixerApi":                 NewFixerApi,
	"FreeCurrencyApi":          NewFreeCurrencyApi,
	"OpenExchangeRates":        NewOpenExchangeRates,
	"FreeCurrencyConverterAPI": NewFreeCurrencyConverterApi,
}

// ProviderPriority is a map of provider names to their priority order, as per the Json Config
// If no priority is specified, the default is 0
// The lower the number, the higher the priority (but 0 means no (any) priority)
var ProviderPriority = make(map[string]uint)

// NewCurrencyLayer constructs a new CurrencyLayer provider
func NewCurrencyLayer(apiKey string, timeout int) ProviderInterface {
	return &CurrencyLayer{
		Name:      "Currency Data API on API Layer",
		AccessKey: apiKey,
		Timeout:   timeout,
	}
}

// NewExchangeRateApi constructs a new ExchangeRateAPI provider
func NewExchangeRateApi(apiKey string, timeout int) ProviderInterface {
	return &ExchangeRateApi{
		Name:    "ExchangeRate-API (exchangerate-api.com)",
		APIKey:  apiKey,
		Timeout: timeout,
	}
}

// NewFixerApi constructs a new FixerApi provider
func NewFixerApi(apiKey string, timeout int) ProviderInterface {
	return &FixerApi{
		Name:      "Fixer API on API Layer",
		AccessKey: apiKey,
		Timeout:   timeout,
	}
}

// NewFreeCurrencyApi constructs a new FreeCurrencyAPI provider
func NewFreeCurrencyApi(apiKey string, timeout int) ProviderInterface {
	return &FreeCurrencyApi{
		Name:    "FreecurrencyAPI",
		APIKey:  apiKey,
		Timeout: timeout,
	}
}

// NewFreeCurrencyConverterApi constructs a new FreeCurrencyConverterAPI provider
func NewFreeCurrencyConverterApi(apiKey string, timeout int) ProviderInterface {
	return &FreeCurrencyConverterAPI{
		Name:    "The Free Currency Converter API",
		APIKey:  apiKey,
		Timeout: timeout,
	}
}

// NewOpenExchangeRates constructs a new OpenExchangeRates provider
func NewOpenExchangeRates(apiKey string, timeout int) ProviderInterface {
	return &OpenExchangeRates{
		Name:    "Open Exchange Rates (openexchangerates.org)",
		AppID:   apiKey,
		Timeout: timeout,
	}
}

// initProvider initializes a single provider
func initProvider(name string, providerConfig config.ProviderConfig, timeout int, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	if !providerConfig.Enabled {
		c.Warnf(" -> Provider %s is disabled", name)
		return
	}
	if providerConfig.Key == "" {
		c.Warnf(" -> API key for provider %s is missing", name)
		return
	}

	makeProvider, exists := providerConstructors[name]
	if !exists {
		// This should never happen
		c.Warnf(" -> No constructor found for provider %s", name)
		return
	}

	nextProvider := makeProvider(providerConfig.Key, timeout)
	if !nextProvider.CheckApiKey() {
		c.Warnf(" -> API key for provider %s is invalid", name)
		return
	}

	// If we get here, we're good, so we append the provider to the list of enabled providers
	mu.Lock()
	EnabledProviders[name] = nextProvider
	ProviderPriority[name] = providerConfig.Priority
	mu.Unlock()
	c.Successf("Provider '%s' is enabled", name)
}

// InitProviders initializes the API exchange rate providers in parallel. Performs various checks for each.
func InitProviders(providers *map[string]config.ProviderConfig, timeout int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, providerConfig := range *providers {
		wg.Add(1)
		go initProvider(name, providerConfig, timeout, &wg, &mu)
	}

	wg.Wait()

	// Check that we have at least one provider enabled
	if len(EnabledProviders) == 0 {
		c.Warnf("No providers enabled. Cannot continue\n")
		os.Exit(1)
	} else {
		c.Outf("Enabled %v out of %v providers\n", len(EnabledProviders), len(*providers))
	}
}
