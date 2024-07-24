// stats_test.go

package stats

import (
	"testing"
)

// TestGetInstance checks the singleton instance
func TestGetInstance(t *testing.T) {
	instance1 := GetInstance()
	instance2 := GetInstance()
	if instance1 != instance2 {
		t.Error("GetInstance should return the same instance")
	}
}

// TestIncHitCount checks incrementing the hit count
func TestIncHitCount(t *testing.T) {
	s := GetInstance()
	s.hitCount = 0
	s.IncHitCount()
	if s.hitCount != 1 {
		t.Errorf("Expected hit count to be 1, got %d", s.hitCount)
	}
}

// TestIncRequestCount checks incrementing the request count
func TestIncRequestCount(t *testing.T) {
	s := GetInstance()
	s.requestCount = 0
	s.IncRequestCount()
	if s.requestCount != 1 {
		t.Errorf("Expected request count to be 1, got %d", s.requestCount)
	}
}

// TestIncErrorCount checks incrementing the error count
func TestIncErrorCount(t *testing.T) {
	s := GetInstance()
	s.errorCount = 0
	s.IncErrorCount()
	if s.errorCount != 1 {
		t.Errorf("Expected error count to be 1, got %d", s.errorCount)
	}
}

// TestIncPath checks incrementing the path count
func TestIncPath(t *testing.T) {
	s := GetInstance()
	s.pathCount = make(map[string]uint64)
	s.IncPath("/test")
	if s.pathCount["/test"] != 1 {
		t.Errorf("Expected path count for '/test' to be 1, got %d", s.pathCount["/test"])
	}
}

// TestGetStats checks retrieving the statistics
func TestGetStats(t *testing.T) {
	s := GetInstance()
	s.hitCount = 5
	s.requestCount = 4
	s.errorCount = 3
	s.failCount = 2
	s.pathCount = map[string]uint64{"/test": 1}

	stats := s.GetStats()
	if stats["hitCount"].(uint64) != 5 {
		t.Errorf("Expected hit count to be 5, got %d", stats["hitCount"])
	}
	if stats["requestCount"].(uint64) != 4 {
		t.Errorf("Expected request count to be 4, got %d", stats["requestCount"])
	}
	if stats["errorCount"].(uint64) != 3 {
		t.Errorf("Expected error count to be 3, got %d", stats["errorCount"])
	}
	if stats["failCount"].(uint64) != 2 {
		t.Errorf("Expected fail count to be 2, got %d", stats["failCount"])
	}
	if stats["pathCount"].(map[string]uint64)["/test"] != 1 {
		t.Errorf("Expected path count for '/test' to be 1, got %d", stats["pathCount"].(map[string]uint64)["/test"])
	}
}
