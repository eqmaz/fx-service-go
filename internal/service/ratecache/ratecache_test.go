package ratecache

import (
	"testing"
	"time"
)

// TestGetInstance checks the singleton instance
func TestGetInstance(t *testing.T) {
	instance1 := GetInstance()
	instance2 := GetInstance()
	if instance1 != instance2 {
		t.Error("GetInstance should return the same instance")
	}
}

// TestSetExpiry checks if the expiry time is set correctly
func TestSetExpiry(t *testing.T) {
	rc := GetInstance()
	rc.SetExpiry(10)
	if rc.expiry != 10*time.Second {
		t.Errorf("Expected expiry time to be %v, got %v", 10*time.Second, rc.expiry)
	}
}

// TestSetAndGetRate checks setting and getting a rate
func TestSetAndGetRate(t *testing.T) {
	rc := GetInstance()
	rc.Clear()
	rc.Set("USD", "EUR", 0.85)
	rate := rc.Get("USD", "EUR")
	if rate == nil || *rate != 0.85 {
		t.Errorf("Expected rate to be 0.85, got %v", rate)
	}
}

// TestGetExpiredRate checks getting an expired rate
func TestGetExpiredRate(t *testing.T) {
	rc := GetInstance()
	rc.Clear()
	rc.SetExpiry(1) // 1 second expiry
	rc.Set("USD", "EUR", 0.85)
	time.Sleep(2 * time.Second)
	rate := rc.Get("USD", "EUR")
	if rate != nil {
		t.Error("Expected rate to be nil after expiry")
	}
}

// TestGetNonExistentRate checks getting a non-existent rate
func TestGetNonExistentRate(t *testing.T) {
	rc := GetInstance()
	rc.Clear()
	rate := rc.Get("USD", "GBP")
	if rate != nil {
		t.Error("Expected rate to be nil for non-existent rate")
	}
}

// TestGetAllRates checks getting all rates
func TestGetAllRates(t *testing.T) {
	rc := GetInstance()
	rc.Clear()
	rc.Set("USD", "EUR", 0.85)
	rc.Set("USD", "GBP", 0.75)
	rates := rc.GetAll()
	if len(rates) != 2 || rates["USD_EUR"] != 0.85 || rates["USD_GBP"] != 0.75 {
		t.Errorf("Expected rates to be {USD_EUR: 0.85, USD_GBP: 0.75}, got %v", rates)
	}
}

// TestClearRates checks clearing all rates
func TestClearRates(t *testing.T) {
	rc := GetInstance()
	rc.Set("USD", "EUR", 0.85)
	rc.Clear()
	rates := rc.GetAll()
	if len(rates) != 0 {
		t.Error("Expected all rates to be cleared")
	}
}
