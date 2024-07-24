package reply

// Error wraps the error response data in the standard error response shape
func Error(response interface{}) map[string]interface{} {
	res := map[string]interface{}{
		"result": nil,
		"error":  response,
	}
	return res
}
