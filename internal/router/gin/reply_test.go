package ginHandlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestReplyResult(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/testReplyResult", func(c *gin.Context) {
		replyResult(c, "test data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testReplyResult", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d but got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"result": "test data"}
	if !reflect.DeepEqual(response, expected) {
		t.Errorf("expected response %v but got %v", expected, response)
	}
}

func TestReplyError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/testReplyError", func(c *gin.Context) {
		replyError(c, http.StatusBadRequest, "error data")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/testReplyError", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d but got %d", http.StatusBadRequest, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := map[string]interface{}{"error": "error data"}
	if response["error"] != expected["error"] {
		t.Errorf("expected error response %v but got %v", expected["error"], response["error"])
	}
	// Check if there's an unexpected result key
	if _, exists := response["result"]; !exists {
		t.Errorf("nil result key should exist in response: %v", response)
	}
}
