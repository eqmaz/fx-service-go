package ratecache

import (
	"sync"
	"time"
)

type RateCache struct {
	mu         sync.Mutex
	rates      map[string]float64
	expiry     time.Duration
	timestamps map[string]time.Time
}

var instance *RateCache
var once sync.Once

// GetInstance singleton pattern to get the RateCache instance
func GetInstance() *RateCache {
	once.Do(func() {
		instance = &RateCache{
			rates:      make(map[string]float64),
			timestamps: make(map[string]time.Time),
		}
	})
	return instance
}

// SetExpiry globally set the expiry time for cache entries
func (rc *RateCache) SetExpiry(seconds int) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.expiry = time.Duration(seconds) * time.Second
}

// Set saves a rate in the cache
func (rc *RateCache) Set(from, to string, rate float64) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	key := from + "_" + to
	rc.rates[key] = rate
	rc.timestamps[key] = time.Now()
}

// Get retrieves a rate from the cache. Returns nil if the rate is not found or expired
func (rc *RateCache) Get(from, to string) *float64 {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	key := from + "_" + to

	rate, exists := rc.rates[key]
	if !exists {
		return nil
	}

	if time.Since(rc.timestamps[key]) > rc.expiry {
		delete(rc.rates, key)
		delete(rc.timestamps, key)
		return nil
	}

	return &rate
}

// GetAll returns all rates in the cache
func (rc *RateCache) GetAll() map[string]float64 {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return rc.rates
}

// Clear removes all rates from the cache
func (rc *RateCache) Clear() {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.rates = make(map[string]float64)
	rc.timestamps = make(map[string]time.Time)
}
