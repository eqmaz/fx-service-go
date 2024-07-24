package stats

import (
	util "fx-service/pkg/helpers"
	"sync"
)

type Stats struct {
	mu           sync.Mutex
	hitCount     uint64            // Count of all requests, valid or not
	requestCount uint64            // Count of successful requests
	errorCount   uint64            // Count of app level error responses
	failCount    uint64            // Count of provider API request failures
	pathCount    map[string]uint64 // Detailed count of requests, by path
}

var instance *Stats
var once sync.Once

// GetInstance singleton pattern to get the RateCache instance
func GetInstance() *Stats {
	once.Do(func() {
		instance = &Stats{
			hitCount:     0,
			requestCount: 0,
			failCount:    0,
			pathCount:    make(map[string]uint64),
		}
	})
	return instance
}

// GetStats retrieves the current statistics in a thread-safe way
func (s *Stats) GetStats() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy of the pathCount map to avoid race conditions
	pathCountCopy := util.CloneMapShallow(s.pathCount)

	return map[string]interface{}{
		"hitCount":     s.hitCount,
		"requestCount": s.requestCount,
		"errorCount":   s.errorCount,
		"failCount":    s.failCount,
		"pathCount":    pathCountCopy,
	}
}

// IncHitCount increments the global hit count (valid requests or not)
func (s *Stats) IncHitCount() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hitCount++
}

// IncRequestCount increments the global request count (successful requests only)
func (s *Stats) IncRequestCount() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestCount++
}

// IncErrorCount increments the global error count (app level error responses)
func (s *Stats) IncErrorCount() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errorCount++
}

// IncPath increments the count of requests for a specific path
func (s *Stats) IncPath(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.pathCount[path]; !ok {
		s.pathCount[path] = 0
	}
	s.pathCount[path]++
}
