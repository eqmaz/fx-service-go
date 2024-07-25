package application

var errorMap = map[string]string{
	"eNcF01": "No valid config file found",
	"eGaPf1": "All providers have failed",
	"eCRP68": "All providers failed in round-robin mode",
	"ePrRnf": "To-symbol (quote) not found in response from API provider",
	"eAGn2c": "Got non-200 response code from API provider (status: %d)",
}
