package config

// defaultConfigMap contains the default configuration values.
// If no config.json file exists, these values will be used.
// The config.json file values can override these values if they are present.
var defaultConfigMap = map[string]interface{}{
	"currenciesEnabled": []string{ // Enabled currency codes (regardless if providers support more)
		"USD", "EUR", "GBP", "JPY", "AUD", "CAD",
	},
	"currenciesCaseSensitive": false,
	"apiTimeout":              10,      // 10 seconds
	"cacheExpirySec":          60 * 60, // 1 hour, in seconds
	"showProvider":            false,   // Whether to display the provider name in each response
	"RateLimiter": map[string]interface{}{ // Rate limit configuration (requests to us)
		"Enabled":     true, // Whether rate limiting is enabled
		"MaxRequests": 10,   // Maximum number of requests within the timeframe period
		"Timeframe":   30,   // Timeframe period in seconds for the rate limit
	},
	"Mode":   "random", // The strategy to fetch exchange rates from different providers
	"Router": "Fiber",  // The http router framework to use for the API
	"Port":   8080,     // The port to listen on for incoming HTTP requests
}
