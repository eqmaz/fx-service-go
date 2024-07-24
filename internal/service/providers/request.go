package providers

import (
	"fmt"
	"fx-service/pkg/e"
	"io"
	"net/http"
	"strings"
	"time"
)

// makeGetRequest makes GET requests for any provider
func makeGetRequest(url string, timeout int, headers *map[string]string) (int, []byte, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Accept", "application/json")

	// Set all the custom headers
	if headers != nil {
		for key, value := range *headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		fields := e.Fields{"url": url}
		// If error string contains "context deadline exceeded" then it's a timeout error
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return 0, nil, e.Throw("httpTimeout", err.Error()).SetFields(fields)
		}
		return 0, nil, e.FromError(err).SetFields(fields)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	//contentType := resp.Header.Get("Content-Type")
	//if !strings.Contains(contentType, "application/json") {
	//	return resp.StatusCode, nil, e.Throwf(errUnhandled, "Unexpected content type: '%s'", contentType)
	//}

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, e.Throw(errUnhandled, err.Error())
	}

	return resp.StatusCode, bodyData, nil
}
