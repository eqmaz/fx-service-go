package reply

// Result wraps the response data in the standard success response shape
func Result(response interface{}) map[string]interface{} {
	return map[string]interface{}{
		"result": response,
	}
}
