package ratecache

// Driver interface for future implementations
type Driver interface {
	Set(from, to string, rate float64)
	SetExpiry(seconds int)
	Get(from, to string) *float64
	Clear()
}

// SetDriver set a driver for the RateCache
func (rc *RateCache) SetDriver(driver Driver) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	instance = driver.(*RateCache)
}

// InMemoryDriver implements the Driver interface
type InMemoryDriver struct {
	RateCache
}

// NewInMemoryDriver initialize with in-memory driver
func NewInMemoryDriver() *InMemoryDriver {
	return &InMemoryDriver{*GetInstance()}
}
