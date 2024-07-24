package config

import (
	"fx-service/pkg/helpers"
	"strings"
)

// ProviderConfig structure for rate provider API configurations
type ProviderConfig struct {
	Enabled    bool     `json:"enabled"`
	Key        string   `json:"key"`
	Currencies []string `json:"currencies"`
	Priority   uint     `json:"priority"`
}

// RateLimiterConfig structure for rate limiter configurations
type RateLimiterConfig struct {
	Enabled     bool `json:"enabled"`
	MaxRequests int  `json:"maxRequests"`
	Timeframe   int  `json:"timeframe"`
}

// Config - main (parent) struct for app configs
type Config struct {
	CurrenciesEnabled       []string                  `json:"currenciesEnabled"`
	CurrenciesCaseSensitive bool                      `json:"currenciesCaseSensitive"`
	APITimeout              int                       `json:"apiTimeout"`
	RateLimiter             RateLimiterConfig         `json:"rateLimiter"`
	CacheExpirySec          int                       `json:"cacheExpirySec"`
	ShowProvider            bool                      `json:"showProvider"`
	Mode                    Mode                      `json:"mode"`
	Router                  string                    `json:"router"`
	Port                    uint64                    `json:"port"`
	Providers               map[string]ProviderConfig `json:"providers"`
}

// CurrenciesToUppercase converts all currencies, from the config, to uppercase
// All comparisons are done in uppercase
func (c *Config) CurrenciesToUppercase() {
	for i, currency := range c.CurrenciesEnabled {
		c.CurrenciesEnabled[i] = strings.ToUpper(currency)
	}
}

// IsCurrencySupported checks if the given currency is supported
func (c *Config) IsCurrencySupported(currency string) bool {
	return helpers.SliceContains(c.CurrenciesEnabled, currency)
}
