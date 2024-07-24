package e

// ErrorMap is a map of error codes to error messages
// These errors can then be called by e.FromCode("e12345")
type ErrorMap map[string]string

// catalogue is a map of error codes to error messages which can be used by the application
// For example - {"e12345": "This is an example error"}
var catalogue = ErrorMap{
	// Use e.SetCatalogue() to set the error catalogue
}

// SetCatalogue sets the error catalogue
func SetCatalogue(c ErrorMap) {
	catalogue = c
}
